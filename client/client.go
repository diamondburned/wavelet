// Package server half-assedly provides functions and methods for a
// Wavelet server to interact with the TUI.
package client

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/cipher"
	"github.com/perlin-network/noise/handshake"
	"github.com/perlin-network/noise/nat"
	"github.com/perlin-network/noise/skademlia"
	"github.com/perlin-network/wavelet"
	"github.com/perlin-network/wavelet/api"
	"github.com/perlin-network/wavelet/cmd/tui/tui/logger"
	"github.com/perlin-network/wavelet/internal/snappy"
	"github.com/perlin-network/wavelet/store"
	"github.com/perlin-network/wavelet/sys"
	"google.golang.org/grpc"
)

// Config contains the server configs
type Config struct {
	NAT      bool
	Host     string
	Port     uint
	Wallet   string
	Genesis  string
	APIPort  uint
	Peers    []string
	Database string
}

// Client provides a Wavelet client
type Client struct {
	*skademlia.Client

	logger *logger.Logger
	cfg    Config

	listener   net.Listener
	address    string
	tcpAddress *net.TCPAddr

	db store.KV

	Ledger *wavelet.Ledger
	Keys   *skademlia.Keypair

	History HistoryStore
}

// New creates a new server. `log' is needed.
func New(cfg Config, log *logger.Logger) (*Client, error) {
	c := &Client{
		logger: log,
		cfg:    cfg,
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, err
	}

	c.listener = l
	c.tcpAddress = l.Addr().(*net.TCPAddr)

	c.address = net.JoinHostPort(
		cfg.Host,
		strconv.Itoa(c.tcpAddress.Port),
	)

	if cfg.NAT {
		if len(cfg.Peers) > 1 {
			if err := nat.NewPMP().AddMapping(
				"tcp",
				uint16(c.tcpAddress.Port),
				uint16(c.tcpAddress.Port),
				30*time.Minute,
			); err != nil {
				return nil, err
			}
		}

		resp, err := http.Get("http://myexternalip.com/raw")
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		ip, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		c.address = net.JoinHostPort(
			string(ip),
			strconv.Itoa(c.tcpAddress.Port),
		)
	}

	c.logger.Level(logger.WithInfo("Listening for peers").
		F("address", c.address))

	// Read the keys from file or private key
	k, err := c.readKeys(cfg.Wallet)
	if err != nil {
		return nil, err
	}

	c.Keys = k

	// Set up a new S/Kademlia client
	c.Client = skademlia.NewClient(
		c.address, k,
		skademlia.WithC1(sys.SKademliaC1),
		skademlia.WithC2(sys.SKademliaC2),
		skademlia.WithDialOptions(
			grpc.WithDefaultCallOptions(grpc.UseCompressor(snappy.Name)),
		),
	)

	// Initiate events

	c.Client.SetCredentials(noise.NewCredentials(
		c.address, handshake.NewECDH(), cipher.NewAEAD(), c.Client.Protocol(),
	))

	c.Client.OnPeerJoin(func(conn *grpc.ClientConn, id *skademlia.ID) {
		pub := fmt.Sprintf("%x", id.PublicKey())
		c.History.add(pub, HistoryEntryAccount)

		c.logger.Level(logger.WithInfo("Peer joined").
			F("public_key", "%s", pub).
			F("address", id.Address()))
	})

	c.Client.OnPeerLeave(func(conn *grpc.ClientConn, id *skademlia.ID) {
		pub := fmt.Sprintf("%x", id.PublicKey())
		c.History.add(pub, HistoryEntryAccount)

		c.logger.Level(logger.WithInfo("Peer left").
			F("public_key", "%s", pub).
			F("address", id.Address()))
	})

	// Initialize the key-value store
	if cfg.Database != "" {
		kv, err := store.NewLevelDB(c.cfg.Database)
		if err != nil {
			// s.logger.Level(logger.WithError(err).
			// Wrap("Failed to create/open database").
			// F("location", s.cfg.Database))

			return nil, err
		}

		c.db = kv
	} else {
		// Make one in memory instead
		c.db = store.NewInmem()
	}

	// Create a new ledger
	c.Ledger = wavelet.NewLedger(c.db, c.Client, &cfg.Genesis)

	return c, nil
}

// Start runs the listeners. This does not block.
func (c *Client) Start() {
	go func() {
		server := c.Client.Listen()

		wavelet.RegisterWaveletServer(server, c.Ledger.Protocol())

		if err := server.Serve(c.listener); err != nil {
			c.logger.Level(logger.WithError(err).
				Wrap("An error occured in the server"))
		}
	}()

	for _, addr := range c.cfg.Peers {
		if _, err := c.Client.Dial(addr); err != nil {
			c.logger.Level(logger.WithError(err).
				Wrap("Error dialing").
				F("address", addr))
		}
	}

	if peers := c.Client.Bootstrap(); len(peers) > 0 {
		lvl := logger.WithSuccess(fmt.Sprintf(
			"Bootstrapped with %d peer(s)", len(peers),
		))

		for i, peer := range peers {
			_i := strconv.Itoa(i)

			lvl.F(_i+"_address", peer.Address())
			lvl.F(_i+"_pubkey", "%x", peer.PublicKey())
		}

		c.logger.Level(lvl)
	}

	if c.cfg.APIPort > 0 {
		go api.New().StartHTTP(int(c.cfg.APIPort), c.Client, c.Ledger, c.Keys)
	}
}

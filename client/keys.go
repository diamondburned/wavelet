package client

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/perlin-network/noise/edwards25519"
	"github.com/perlin-network/noise/skademlia"
	"github.com/perlin-network/wavelet/cmd/tui/tui/logger"
	"github.com/perlin-network/wavelet/sys"
)

func (c *Client) readKeys(wallet string) (*skademlia.Keypair, error) {
	keys, err := keysFromFile(wallet)
	if err != nil {
		if os.IsNotExist(err) { // maybe it's a private key?
			keys, err = keysFromPrivate(wallet)
			if err == nil {
				goto Worked
			}
		}

		c.logger.Level(logger.WithError(err).
			Wrap("Generating a new key"))

		keys, err = keysGenerate()
	}

	if err != nil {
		c.logger.Level(logger.WithError(err).
			Wrap("Cannot generate a new key"))

		return nil, err
	}

Worked:
	c.logger.Level(logger.WithSuccess("Wallet loaded.").
		F("privatekey", "%x", keys.PrivateKey()).
		F("publickey", "%x", keys.PublicKey()))

	return keys, err
}

func keysGenerate() (*skademlia.Keypair, error) {
	keys, err := skademlia.NewKeys(sys.SKademliaC1, sys.SKademliaC2)
	if err != nil {
		return nil, errors.New("failed to generate a new wallet")
	}

	return keys, nil
}

func keysFromPrivate(key string) (*skademlia.Keypair, error) {
	var privateKey edwards25519.PrivateKey

	n, err := hex.Decode(privateKey[:], []byte(key))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode the private key specified: %s", key)
	}

	if n != edwards25519.SizePrivateKey {
		return nil, fmt.Errorf(
			"private key %s is not of the right length (%d instead of %d)",
			key, n, edwards25519.SizePrivateKey)
	}

	keys, err := skademlia.LoadKeys(privateKey, sys.SKademliaC1, sys.SKademliaC2)
	if err != nil {
		return nil, fmt.Errorf(
			"the private key specified is invalid: %s", key)
	}

	return keys, nil
}

func keysFromFile(filename string) (*skademlia.Keypair, error) {
	privateKeyBuf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return keysFromPrivate(string(privateKeyBuf))
}

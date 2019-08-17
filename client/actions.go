package client

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	"github.com/perlin-network/wavelet"
	"github.com/perlin-network/wavelet/cmd/tui/tui/logger"
	"github.com/perlin-network/wavelet/sys"
	"github.com/pkg/errors"
)

// TODO(diamond): Sometime in the future, maybe when I'm 70, I would like all
// this to return structs instead of calling logger directly.

var ErrInsufficientBalance = errors.New("Insufficient balance")

var (
	ErrInvalidID       = errors.New("Invalid transaction/account ID")
	ErrInvalidIDLength = errors.New("Invalid transaction/account ID length")
)

// PlaceStake places a stake
func (c *Client) PlaceStake(amount int) error {
	var intBuf [8]byte
	var payload bytes.Buffer

	payload.WriteByte(sys.PlaceStake)

	// Write the amount
	binary.LittleEndian.PutUint64(intBuf[:8], uint64(amount))
	payload.Write(intBuf[:8])

	tx, err := s.sendTx(wavelet.NewTransaction(
		s.Keys, sys.TagStake, payload.Bytes(),
	))

	if err != nil {
		return err
	}

	s.logger.Level(logger.WithSuccess("Stake placed").
		F("id", "%x", tx.ID))

	return nil
}

// WithdrawStake withdraws the stake
func (c *Client) WithdrawStake(amount int) error {
	var intBuf [8]byte
	var payload bytes.Buffer

	payload.WriteByte(sys.WithdrawStake)

	// Write the amount
	binary.LittleEndian.PutUint64(intBuf[:8], uint64(amount))
	payload.Write(intBuf[:8])

	tx, err := s.sendTx(wavelet.NewTransaction(
		s.Keys, sys.TagStake, payload.Bytes(),
	))

	if err != nil {
		return err
	}

	s.logger.Level(logger.WithSuccess("Stake withdrew").
		F("id", "%x", tx.ID))

	return nil
}

// WithdrawReward withdraws the reward
func (c *Client) WithdrawReward(amount int) error {
	var intBuf [8]byte
	var payload bytes.Buffer

	payload.WriteByte(sys.WithdrawReward)

	// Write the amount
	binary.LittleEndian.PutUint64(intBuf[:8], uint64(amount))
	payload.Write(intBuf[:8])

	tx, err := s.sendTx(wavelet.NewTransaction(
		s.Keys, sys.TagStake, payload.Bytes(),
	))

	if err != nil {
		return err
	}

	s.logger.Level(logger.WithSuccess("Reward withdrew").
		F("id", "%x", tx.ID))

	return nil
}

func (c *Client) sendTx(tx wavelet.Transaction) (*wavelet.Transaction, error) {
	tx = wavelet.AttachSenderToTransaction(
		s.Keys, tx,
		s.Ledger.Graph().FindEligibleParents()...,
	)

	if err := s.Ledger.AddTransaction(tx); err != nil {
		if errors.Cause(err) != wavelet.ErrMissingParents {
			e := logger.WithError(err)
			e.Wrap("Failed to create transaction")
			e.F("tx_id", "%x", tx.ID)

			s.logger.Level(e)

			return nil, e
		}
	}

	// Add the ID into history
	s.History.add(hex.EncodeToString(tx.ID[:]), HistoryEntryTransaction)

	return &tx, nil
}

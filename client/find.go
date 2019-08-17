package client

import (
	"encoding/hex"

	"github.com/perlin-network/wavelet"
	"github.com/perlin-network/wavelet/cmd/tui/tui/logger"
	"github.com/perlin-network/wavelet/sys"
	"github.com/pkg/errors"
)

type Account struct {
	Balance    uint64
	Stake      uint64
	Nonce      uint64
	Reward     uint64
	IsContract bool
	NumPages   uint64
}

var ErrNotAccount = errors.New("Address is not account")

func (c *Client) FindAccount(addr []byte) (*Account, error) {
	if len(addr) != wavelet.SizeTransactionID && len(addr) != wavelet.SizeAccountID {
		return nil, ErrInvalidIDLength
	}

	snapshot := c.Ledger.Snapshot()

	// See if the ID given is an account ID
	var accountID wavelet.AccountID
	copy(accountID[:], addr)

	// Read account info
	balance, _ := wavelet.ReadAccountBalance(snapshot, accountID)
	stake, _ := wavelet.ReadAccountStake(snapshot, accountID)
	nonce, _ := wavelet.ReadAccountNonce(snapshot, accountID)
	reward, _ := wavelet.ReadAccountReward(snapshot, accountID)

	_, isContract := wavelet.ReadAccountContractCode(snapshot, accountID)
	numPages, _ := wavelet.ReadAccountContractNumPages(snapshot, accountID)

	// The ID is an account
	if !(balance > 0 || stake > 0 || nonce > 0 || isContract || numPages > 0) {
		return nil, ErrNotAccount
	}

	return &Account{
		balance, stake, nonce, reward, isContract, numPages,
	}, nil
}

type Transaction struct {
	*wavelet.Transaction
	Parents []string
	Sender  []byte
	Creator []byte
	Nonce   uint64
	Tag     sys.Tag
	Depth   uint64

	Seed              []byte
	SeedZeroPrefixLen byte
}

var ErrNotTransaction = errors.New("Address is not transaction")

func (c *Client) FindTransaction(addr []byte) (*Transaction, error) {
	if len(addr) != wavelet.SizeTransactionID {
		return nil, ErrInvalidIDLength
	}

	// ID was not an account, probably a transaction then
	var txID wavelet.TransactionID
	copy(txID[:], addr)

	tx := c.Ledger.Graph().FindTransaction(txID)
	if tx == nil {
		return nil, ErrNotTransaction
	}

	var parents = make([]string, 0, len(tx.ParentIDs))
	for _, parentID := range tx.ParentIDs {
		parents = append(parents, hex.EncodeToString(parentID[:]))
	}

	return &Transaction{
		tx, parents, tx.Sender[:], tx.Creator[:], tx.Nonce,
		tx.Tag, tx.Depth, tx.Seed[:], tx.SeedLen,
	}, nil
}

func (c *Client) Find(addr []byte) error {
	err := c.FindAccount(addr)
	if err == ErrNotAccount {
		return c.FindTransaction(addr)
	}

	return err
}

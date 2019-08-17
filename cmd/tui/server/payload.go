package server

import (
	"github.com/perlin-network/noise/skademlia"
	"github.com/perlin-network/wavelet"
	"github.com/perlin-network/wavelet/avl"
	"github.com/pkg/errors"
)

var (
	ErrInsufficientBalance = errors.New("Insufficient balance")
)

type Payload struct {
	wavelet.Transfer
	snapshot *avl.Tree
	keys     *skademlia.Keypair
}

func NewPayload(ledger *wavelet.Ledger, keys *skademlia.Keypair) *Payload {
	return &Payload{
		Transfer: wavelet.Transfer{},
		snapshot: ledger.Snapshot(),
	}
}

func (p *Payload) SetRecipient(r [wavelet.SizeAccountID]byte) {
	copy(p.Recipient[:], r[:])
}

func (p *Payload) SetAmount(a uint64) error {
	bal, _ := wavelet.ReadAccountBalance(p.snapshot)
}

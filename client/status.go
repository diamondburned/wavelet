package client

import (
	"encoding/hex"

	"github.com/perlin-network/wavelet"
	"github.com/perlin-network/wavelet/sys"
)

type Status struct {
	Difficulty byte
	Round      uint64
	RootID     []byte
	ID         []byte
	Height     uint64

	Balance uint64
	Stake   uint64
	Reward  uint64
	Nonce   uint64
	PeerIDs []string

	Transactions        int
	MissingTransactions int
	TransactionsInStore int
	AccountsInStore     uint64

	PreferredID    string
	PreferredVotes int
}

func (s *Server) Status() Status {
	preferredID := "N/A"

	if preferred := s.Ledger.Finalizer().Preferred(); preferred != nil {
		preferredID = hex.EncodeToString(preferred.ID[:])
	}

	count := s.Ledger.Finalizer().Progress()

	snapshot := s.Ledger.Snapshot()
	publicKey := s.Keys.PublicKey()

	accountsLen := wavelet.ReadAccountsLen(snapshot)

	balance, _ := wavelet.ReadAccountBalance(snapshot, publicKey)
	stake, _ := wavelet.ReadAccountStake(snapshot, publicKey)
	reward, _ := wavelet.ReadAccountReward(snapshot, publicKey)
	nonce, _ := wavelet.ReadAccountNonce(snapshot, publicKey)

	round := s.Ledger.Rounds().Latest()
	rootDepth := s.Ledger.Graph().RootDepth()

	peers := s.Client.ClosestPeerIDs()
	peerIDs := make([]string, 0, len(peers))

	for _, id := range peers {
		peerIDs = append(peerIDs, id.String())
	}

	// Add the root ID and self ID to the history
	s.History.add(hex.EncodeToString(round.End.ID[:]), HistoryEntryRoot)
	s.History.add(hex.EncodeToString(publicKey[:]), HistoryEntrySelf)

	return Status{
		Difficulty: round.ExpectedDifficulty(
			sys.MinDifficulty, sys.DifficultyScaleFactor),
		Round:   round.Index,
		RootID:  round.End.ID[:],
		ID:      publicKey[:],
		Height:  s.Ledger.Graph().Height(),
		Balance: balance,
		Stake:   stake,
		Reward:  reward,
		Nonce:   nonce,
		PeerIDs: peerIDs,

		Transactions:        s.Ledger.Graph().DepthLen(&rootDepth, nil),
		MissingTransactions: s.Ledger.Graph().MissingLen(),
		TransactionsInStore: s.Ledger.Graph().Len(),
		AccountsInStore:     accountsLen,

		PreferredID:    preferredID,
		PreferredVotes: count,
	}
}

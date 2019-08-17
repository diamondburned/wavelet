package client

import (
	"github.com/perlin-network/wavelet"
	"github.com/perlin-network/wavelet/cmd/tui/tui/logger"
	"github.com/perlin-network/wavelet/sys"
)

// gasLimit is optional
func (s *Server) Pay(recipient [wavelet.SizeAccountID]byte,
	amount, gasLimit uint64, additional []byte) (*wavelet.Transaction, error) {

	var payload wavelet.Transfer
	copy(payload.Recipient[:], recipient[:])

	payload.Amount = amount

	snapshot := s.Ledger.Snapshot()
	balance, _ := wavelet.ReadAccountBalance(snapshot, s.Keys.PublicKey())

	if balance < uint64(amount)+sys.TransactionFeeAmount {
		err := logger.WithError(ErrInsufficientBalance).
			Wrap("Can't pay")

		err.F("balance", "%d", balance)
		err.F("cost", "%d", amount)
		err.F("offset", "%d", int(balance)-int(amount))

		s.logger.Level(err)

		return nil, err
	}

	_, codeAvailable := wavelet.ReadAccountContractCode(
		snapshot, payload.Recipient,
	)

	if codeAvailable {
		// Set gas limit by default to the balance the user has.
		payload.GasLimit = balance - amount - sys.TransactionFeeAmount
		payload.FuncName = []byte("on_money_received")
	}

	tx, err := s.sendTx(wavelet.NewTransaction(
		s.Keys, sys.TagTransfer, payload.Marshal(),
	))

	if err != nil {
		return nil, err
	}

	s.logger.Level(logger.WithSuccess("Paid").
		F("tx_id", "%x", tx.ID))

	return tx, nil
}

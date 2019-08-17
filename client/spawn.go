package client

/*
// Spawn spawns a smart contract
func (c *Client) Spawn(file *os.File) error {
	code, err := ioutil.ReadFile(pathToContract)
	if err != nil {
		e := logger.WithError(err)
		e.Wrap("Failed to find the smart contract")
		e.F("path", pathToContract)

		s.logger.Level(e)
		return e
	}

	var buf [8]byte

	w := bytes.NewBuffer(nil)

	// Write a fake gas fee
	binary.LittleEndian.PutUint64(buf[:], 100000000)
	w.Write(buf[:])

	// Write the payload size
	binary.LittleEndian.PutUint64(buf[:], uint64(len(code)))
	w.Write(buf[:])

	w.Write(code) // Smart contract code.

	newTx := wavelet.NewTransaction(s.Keys, sys.TagContract, w.Bytes())

	tx, err := s.sendTx(newTx)
	if err != nil {
		return err
	}

	s.logger.Level(logger.WithSuccess("Smart contract deployed").
		F("id", "%x", tx.ID))

	return nil
}

*/

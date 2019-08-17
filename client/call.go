package client

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	"github.com/perlin-network/wavelet"
	"github.com/perlin-network/wavelet/cmd/tui/tui/logger"
)

// FunctionCall is the struct containing parameters to call a function.
type FunctionCall struct {
	Name   string
	Params [][]byte
}

func NewFunctionCall(name string, params ...[]byte) FunctionCall {
	return FunctionCall{
		Name:   name,
		Params: params,
	}
}

// Call calls a smart contract function
func (c *Client) Call(recipient [wavelet.SizeAccountID]byte,
	amount, gasLimit uint64, fn FunctionCall) (*wavelet.Transaction, error) {

	// Make an int buffer
	var intBuf = make([]byte, 8)

	// Make a payload buffer
	var payload bytes.Buffer

	// Write the function name length and name
	payload.Write(EncodeString(fn.Name))

	// Make a function parameters buffer
	var params bytes.Buffer
	for _, b := range fn.Params {
		params.Write(b)
	}

	// Write the parameters buffer into the payload
	buf := params.Bytes()
	binary.LittleEndian.PutUint32(intBuf[:4], uint32(len(buf)))
	payload.Write(intBuf[:4]) // write len
	payload.Write(buf)        // write body

	tx, err := s.Pay(recipient, amount, gasLimit, payload.Bytes())
	if err != nil {
		return nil, err
	}

	s.logger.Level(logger.WithSuccess("Called function "+fn.Name).
		F("tx_id", "%x", tx.ID))

	return tx, nil
}

// Jesus Christ why did I do this
// I could've just used a []byte and functions to convert them, instead of
// doing this stupid OOP pattern

// LOL fixed it. Functional FTW.

func DecodeHex(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func EncodeString(s string) []byte {
	return EncodeBytes([]byte(s))
}

func EncodeBytes(b []byte) []byte {
	var lenbuf = make([]byte, 4)
	binary.LittleEndian.PutUint32(lenbuf, uint32(len(b)))

	var buf bytes.Buffer
	buf.Write(lenbuf)
	buf.Write(b)

	return buf.Bytes()
}

func EncodeByte(u byte) []byte {
	return []byte{u}
}

func EncodeUint8(u uint8) []byte {
	return EncodeByte(byte(u))
}

func EncodeUint16(u uint16) []byte {
	var buf = make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, u)
	return buf
}

func EncodeUint32(u uint32) []byte {
	var buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, u)
	return buf
}

func EncodeUint64(u uint64) []byte {
	var buf = make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, u)
	return buf
}

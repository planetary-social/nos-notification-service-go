package domain

import (
	"encoding/hex"

	"github.com/boreq/errors"
)

type PublicKey struct {
	s string
}

func NewPublicKey(s string) (PublicKey, error) {
	decodeString, err := hex.DecodeString(s)
	if err != nil {
		return PublicKey{}, errors.Wrap(err, "error decoding hex")
	}

	s = hex.EncodeToString(decodeString)
	return PublicKey{s}, nil
}

func (k PublicKey) Hex() string {
	return k.s
}

func (k PublicKey) Bytes() []byte {
	b, err := hex.DecodeString(k.s)
	if err != nil {
		panic(err)
	}
	return b
}

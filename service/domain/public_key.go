package domain

import (
	"encoding/hex"

	"github.com/boreq/errors"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

type PublicKey struct {
	s string
}

func NewPublicKey(s string) (PublicKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return PublicKey{}, errors.Wrap(err, "error decoding hex")
	}

	if len(b) != secp256k1.PrivKeyBytesLen {
		return PublicKey{}, errors.New("invalid length")
	}

	s = hex.EncodeToString(b)
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

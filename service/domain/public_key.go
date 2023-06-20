package domain

type PublicKey struct {
	s string
}

func NewPublicKey(s string) (PublicKey, error) {
	// todo validate

	return PublicKey{s}, nil
}

func (k PublicKey) Bytes() []byte {
	return k.Bytes()
}

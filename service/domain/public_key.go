package domain

type PublicKey struct {
	s string
}

func NewPublicKey(s string) (PublicKey, error) {
	return PublicKey{s}, nil
}

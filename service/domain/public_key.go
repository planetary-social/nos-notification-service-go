package domain

type PublicKey struct {
	s string
}

func NewPublicKey(s string) (PublicKey, error) {
	// todo validate

	return PublicKey{s}, nil
}

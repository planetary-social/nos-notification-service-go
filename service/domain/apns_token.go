package domain

type APNSToken struct {
	s string
}

func (t APNSToken) String() string {
	return t.s
}

func NewAPNSToken(s string) (APNSToken, error) {
	// todo validate

	return APNSToken{s}, nil
}

package domain

type APNSToken struct {
	s string
}

func NewAPNSToken(s string) (APNSToken, error) {
	// todo validate

	return APNSToken{s}, nil
}

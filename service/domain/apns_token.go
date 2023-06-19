package domain

type APNSToken struct {
	s string
}

func NewAPNSToken(s string) (APNSToken, error) {
	return APNSToken{s}, nil
}

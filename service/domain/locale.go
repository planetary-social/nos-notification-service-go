package domain

type Locale struct {
	s string
}

func NewLocale(s string) (Locale, error) {
	// todo validate

	return Locale{s: s}, nil
}

func (l Locale) String() string {
	return l.s
}

package jwt

type ErrInvalidToken struct{}

func (e ErrInvalidToken) Error() string {
	return "invalid-token"
}

type ErrExpiredToken struct{}

func (e ErrExpiredToken) Error() string {
	return "expired-token"
}

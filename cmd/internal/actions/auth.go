package actions

import "time"

func (a *Actions) SendVerificationEmail(email string) error {
	token, err := a.jwt.Sign(TokenPayload{
		Email: email,
	}, JwtSessionToken, time.Hour)
	if err != nil {
		return err
	}

	err = a.mailer.SendVerificationEmail(token, email)
	if err != nil {
		return err
	}

	return nil
}

func (a *Actions) ConfirmAuth(token string) error {
	err := a.jwt.Validate(token, JwtSessionToken)
	if err != nil {
		return err
	}

	return nil
}

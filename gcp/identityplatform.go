package gcp

import (
	"context"

	fb "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/jarrodhroberson/ossgo/secrets"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"
)

type IdentityToolkitAccountSignUpError struct {
	Error struct {
		Code   int `json:"code"`
		Errors []struct {
			Domain  string `json:"domain"`
			Message string `json:"message"`
			Reason  string `json:"reason"`
		} `json:"errors"`
		Message string `json:"message"`
	} `json:"error"`
}

type SignUpWithEmailPasswordResponse struct {
	IdToken      string                            `json:"idToken"`
	Email        string                            `json:"email"`
	RefreshToken string                            `json:"refreshToken"`
	ExpiresIn    string                            `json:"expiresIn"`
	LocalId      string                            `json:"localId"`
	Error        IdentityToolkitAccountSignUpError `json:"error"`
}

type IdentityPlatform struct {
}

type Client struct {
	host   string
	apiKey string
}

func NewIdentityPlatformClient() *Client {
	return &Client{
		host:   "identitytoolkit.googleapis.com",
		apiKey: secrets.GetSecretValueAsString(context.Background(), "IDENTITY_PLATFORM_API_KEY"), //TODO: make this dynamic otherwise it will require server restarts to change it.
	}
}

func (ipc *Client) GetByEmailAddress(ctx context.Context, emailAddress string) (*auth.UserRecord, error) {
	app, err := fb.NewApp(context.Background(), nil)
	if err != nil {
		err = errorx.InitializationFailed.Wrap(err, "error initialising firebase app")
		log.Fatal().Err(err).Msg(err.Error())
	}
	authClient, err := app.Auth(ctx)
	if err != nil {
		err = errorx.InitializationFailed.Wrap(err, "error initialising firebase auth client")
		log.Fatal().Err(err).Msg(err.Error())
		return nil, err
	}
	u, err := authClient.GetUserByEmail(ctx, emailAddress)
	if err != nil {
		return nil, err
	}
	return u, err
}

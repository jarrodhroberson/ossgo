package gcp

import (
	"context"
	"fmt"
	"strconv"
	"time"

	fb "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"resty.dev/v3"

	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/jarrodhroberson/ossgo/secrets"
)

type NestedError struct {
	Domain  string `json:"domain"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

func (n NestedError) MarshalZerologObject(e *zerolog.Event) {
	e.Str("domain", n.Domain).Str("message", n.Message).Str("reason", n.Reason)
}

type NestedErrors []NestedError

func (n NestedErrors) MarshalZerologArray(e *zerolog.Array) {
	for _, ne := range n {
		e.Object(ne)
	}
}

/*
IdentityToolkitError
This is an example of the payload when there is an error from the IdentityToolkit

		<code>
	    {
		  "error": {
		    "code": 400,
		    "message": "EMAIL_EXISTS",
		    "errors": [
		      {
		        "message": "EMAIL_EXISTS",
		        "domain": "global",
		        "reason": "invalid"
		      }
		    ]
		  }
		}
	    </code>
*/
type IdentityToolkitError struct {
	Code    int           `json:"code"`
	Errors  []NestedError `json:"errors"`
	Message string        `json:"message"`
}

func (i IdentityToolkitError) MarshalZerologObject(e *zerolog.Event) {
	e.Str("code", strconv.Itoa(i.Code)).
		Str("message", i.Message)
}

func (i IdentityToolkitError) Error() string {
	return fmt.Sprintf("%s: %s", i.Code, i.Message)
}

type SignUpWithEmailPasswordResponse struct {
	IdToken      string               `json:"idToken"`
	Email        string               `json:"email"`
	RefreshToken string               `json:"refreshToken"`
	ExpiresIn    string               `json:"expiresIn"`
	LocalId      string               `json:"localId"`
	Error        IdentityToolkitError `json:"error"`
}

type SignInWithEmailPasswordResponse struct {
	IdToken      string               `json:"idToken"`
	Email        string               `json:"email"`
	RefreshToken string               `json:"refreshToken"`
	ExpiresIn    string               `json:"expiresIn"`
	LocalId      string               `json:"localId"`
	Registered   bool                 `json:"registered"`
	Error        IdentityToolkitError `json:"error"`
}

func (s SignInWithEmailPasswordResponse) MarshalZerologObject(e *zerolog.Event) {
	e.Stack().
		Str("idToken", s.IdToken).
		Str("email", s.Email).
		Str("refreshToken", s.RefreshToken).
		Str("expiresIn", s.ExpiresIn).
		Str("localId", s.LocalId).
		Bool("registered", s.Registered).
		Err(s.Error)
}

type authRequestBody struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
	TenantId          string `json:"tenantId"`
}

type SignUpRequestBody authRequestBody

type SignInRequestBody authRequestBody

type Client interface {
	CustomToken(ctx context.Context, uid string) (string, error)
	CustomTokenWithClaims(ctx context.Context, uid string, developerClaims map[string]interface{}) (string, error)
	SessionCookie(ctx context.Context, idToken string, expiresIn time.Duration) (string, error)
	VerifyIdToken(ctx context.Context, token string) (*auth.Token, error)
	VerifyIDTokenAndCheckRevoked(ctx context.Context, idToken string) (*auth.Token, error)
	VerifySessionCookie(ctx context.Context, sessionCookie string) (*auth.Token, error)
	VerifySessionCookieAndCheckRevoked(ctx context.Context, sessionCookie string) (*auth.Token, error)
	GetUserData(ctx context.Context, idToken string) (*auth.UserRecord, error)
	GetByEmailAddress(ctx context.Context, emailAddress string) (*auth.UserRecord, error)
	GetByLocalId(ctx context.Context, localId string) (*auth.UserRecord, error)
	SignUpWithEmailPassword(ctx context.Context, email string, password string) (*SignUpWithEmailPasswordResponse, error)
	SignInWithEmailPassword(ctx context.Context, email string, password string) (*SignInWithEmailPasswordResponse, error)
	EmailVerificationLink(ctx context.Context, email string) (string, error)
	PasswordResetLink(ctx context.Context, email string) (string, error)
	EmailSignInLink(ctx context.Context, email string, settings *auth.ActionCodeSettings) (string, error)
	RevokeRefreshTokens(ctx context.Context, uid string) error
}

func NewIdentityPlatformClient() Client {
	return &client{
		host:   "identitytoolkit.googleapis.com",
		apiKey: secrets.GetSecretValueAsString(context.Background(), "IDENTITY_PLATFORM_API_KEY"), //TODO: make this dynamic otherwise it will require server restarts to change it.
		app:    must.Must(fb.NewApp(context.Background(), nil)),
	}
}

type client struct {
	host   string
	apiKey string
	app    *fb.App
}

/*
RevokeRefreshTokens
This method revokes all refresh tokens for a specified user.
endpoint: N/A - this is a firebase SDK method
documentation: https://firebase.google.com/docs/reference/admin/go/reference/admin#authclient.RevokeRefreshTokens
*/
func (ipc *client) RevokeRefreshTokens(ctx context.Context, uid string) error {
	authClient, err := ipc.app.Auth(ctx)
	if err != nil {
		err = errorx.InitializationFailed.Wrap(err, "error initialising firebase auth client")

		return err
	}
	return authClient.RevokeRefreshTokens(ctx, uid)
}

/*
EmailSignInLink
This method generates an email sign-in link for a specified email address.
endpoint: N/A - this is a firebase SDK method
documentation: https://firebase.google.com/docs/reference/admin/go/reference/admin#authclient.EmailSignInLink
*/
func (ipc *client) EmailSignInLink(ctx context.Context, email string, settings *auth.ActionCodeSettings) (string, error) {
	authClient, err := ipc.app.Auth(ctx)
	if err != nil {
		err = errorx.InitializationFailed.Wrap(err, "error initialising firebase auth client")

		return "", err
	}
	return authClient.EmailSignInLink(ctx, email, settings)
}

/*
PasswordResetLink
This method generates a password reset link for a specified email address.
endpoint: N/A - this is a firebase SDK method
documentation: https://firebase.google.com/docs/reference/admin/go/reference/admin#authclient.PasswordResetLink
*/
func (ipc *client) PasswordResetLink(ctx context.Context, email string) (string, error) {
	authClient, err := ipc.app.Auth(ctx)
	if err != nil {
		err = errorx.InitializationFailed.Wrap(err, "error initialising firebase auth client")

		return "", err
	}
	return authClient.PasswordResetLink(ctx, email)
}

/*
EmailVerificationLink
This method generates an email verification link for a specified email address.
endpoint: N/A - this is a firebase SDK method
documentation: https://firebase.google.com/docs/reference/admin/go/reference/admin#authclient.EmailVerificationLink
*/
func (ipc *client) EmailVerificationLink(ctx context.Context, email string) (string, error) {
	authClient, err := ipc.app.Auth(ctx)
	if err != nil {
		err = errorx.InitializationFailed.Wrap(err, "error initialising firebase auth client")

		return "", err
	}
	return authClient.EmailVerificationLink(ctx, email)
}

/*
VerifySessionCookieAndCheckRevoked
This method verifies a session cookie and checks if it has been revoked.
endpoint: N/A - this is a firebase SDK method
documentation: https://firebase.google.com/docs/reference/admin/go/reference/admin#authclient.VerifySessionCookieAndCheckRevoked
*/
func (ipc *client) VerifySessionCookieAndCheckRevoked(ctx context.Context, sessionCookie string) (*auth.Token, error) {
	authClient, err := ipc.app.Auth(ctx)
	if err != nil {
		err = errorx.InitializationFailed.Wrap(err, "error initialising firebase auth client")

		return nil, err
	}
	return authClient.VerifySessionCookieAndCheckRevoked(ctx, sessionCookie)
}

/*
VerifySessionCookie
This method verifies a session cookie.
endpoint: N/A - this is a firebase SDK method
documentation: https://firebase.google.com/docs/reference/admin/go/reference/admin#authclient.VerifySessionCookie
*/
func (ipc *client) VerifySessionCookie(ctx context.Context, sessionCookie string) (*auth.Token, error) {
	authClient, err := ipc.app.Auth(ctx)
	if err != nil {
		err = errorx.InitializationFailed.Wrap(err, "error initialising firebase auth client")

		return nil, err
	}
	return authClient.VerifySessionCookieAndCheckRevoked(ctx, sessionCookie)
}

/*
VerifyIDTokenAndCheckRevoked
This method verifies an ID token and checks if it has been revoked.
endpoint: N/A - this is a firebase SDK method
documentation: https://firebase.google.com/docs/reference/admin/go/reference/admin#authclient.VerifyIDTokenAndCheckRevoked
*/
func (ipc *client) VerifyIDTokenAndCheckRevoked(ctx context.Context, idToken string) (*auth.Token, error) {
	authClient, err := ipc.app.Auth(ctx)
	if err != nil {
		err = errorx.InitializationFailed.Wrap(err, "error initialising firebase auth client")

		return nil, err
	}
	return authClient.VerifyIDTokenAndCheckRevoked(ctx, idToken)
}

/*
SessionCookie
This method creates a session cookie for a given ID token.
endpoint: N/A - this is a firebase SDK method
documentation: https://firebase.google.com/docs/reference/admin/go/reference/admin#authclient.SessionCookie
*/
func (ipc *client) SessionCookie(ctx context.Context, idToken string, expiresIn time.Duration) (string, error) {
	authClient, err := ipc.app.Auth(ctx)
	if err != nil {
		err = errorx.InitializationFailed.Wrap(err, "error initialising firebase auth client")

		return "", err
	}
	return authClient.SessionCookie(ctx, idToken, expiresIn)
}

/*
CustomToken
This method creates a custom token for a given user ID.
endpoint: N/A - this is a firebase SDK method
documentation: https://firebase.google.com/docs/reference/admin/go/reference/admin#authclient.CustomToken
*/
func (ipc *client) CustomToken(ctx context.Context, uid string) (string, error) {
	authClient, err := ipc.app.Auth(ctx)
	if err != nil {
		err = errorx.InitializationFailed.Wrap(err, "error initialising firebase auth client")

		return "", err
	}
	return authClient.CustomToken(ctx, uid)
}

/*
CustomTokenWithClaims
This method creates a custom token for a given user ID with additional developer claims.
endpoint: N/A - this is a firebase SDK method
documentation: https://firebase.google.com/docs/reference/admin/go/reference/admin#authclient.CustomTokenWithClaims
*/
func (ipc *client) CustomTokenWithClaims(ctx context.Context, uid string, developerClaims map[string]interface{}) (string, error) {
	authClient, err := ipc.app.Auth(ctx)
	if err != nil {
		err = errorx.InitializationFailed.Wrap(err, "error initialising firebase auth client")

		return "", err
	}
	return authClient.CustomTokenWithClaims(ctx, uid, developerClaims)
}

/*
VerifyIdToken
This method verifies the provided ID token against the Firebase Authentication service.
It returns a Token object if the token is valid, or an error if the token is invalid or if there is an issue communicating with the Firebase service.
endpoint: N/A - this is a firebase SDK method
*/
func (ipc *client) VerifyIdToken(ctx context.Context, token string) (*auth.Token, error) {
	authClient, err := ipc.app.Auth(ctx)
	if err != nil {
		err = errorx.InitializationFailed.Wrap(err, "error initialising firebase auth client")

		return nil, err
	}
	return authClient.VerifyIDToken(ctx, token)
}

/*
GetUserData
endpoint https://identitytoolkit.googleapis.com/v1/accounts:lookup?key=[API_KEY]
documentation https://cloud.google.com/identity-platform/docs/use-rest-api#section-get-account-info
common error codes: INVALID_ID_TOKEN, USER_NOT_FOUND
*/
func (ipc *client) GetUserData(ctx context.Context, idToken string) (*auth.UserRecord, error) {
	authClient, err := ipc.app.Auth(ctx)
	if err != nil {
		err = errorx.InitializationFailed.Wrap(err, "error initialising firebase auth client")

		return nil, err
	}

	token, err := authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}
	return authClient.GetUser(ctx, token.UID)
}

/*
GetByEmailAddress
endpoint https://identitytoolkit.googleapis.com/v1/accounts:lookup?key=[API_KEY]
documentation https://cloud.google.com/identity-platform/docs/use-rest-api#section-get-account-info
common error codes: INVALID_ID_TOKEN, USER_NOT_FOUND
*/
func (ipc *client) GetByEmailAddress(ctx context.Context, emailAddress string) (*auth.UserRecord, error) {
	authClient, err := ipc.app.Auth(ctx)
	if err != nil {
		err = errorx.InitializationFailed.Wrap(err, "error initialising firebase auth client")

		return nil, err
	}
	return authClient.GetUserByEmail(ctx, emailAddress)
}

/*
GetByLocalId
endpoint https://identitytoolkit.googleapis.com/v1/accounts:lookup?key=[API_KEY]
documentation https://cloud.google.com/identity-platform/docs/use-rest-api#section-get-account-info
common error codes: INVALID_ID_TOKEN, USER_NOT_FOUND
*/
func (ipc *client) GetByLocalId(ctx context.Context, localId string) (*auth.UserRecord, error) {
	authClient, err := ipc.app.Auth(ctx)
	if err != nil {
		err = errorx.InitializationFailed.Wrap(err, "error initialising firebase auth client")

		return nil, err
	}
	return authClient.GetUser(ctx, localId)
}

/*
SignUpWithEmailPassword
endpoint `https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=[API_KEY]`
documentation https://cloud.google.com/identity-platform/docs/use-rest-api#section-create-email-password
common error messages: EMAIL_EXISTS, OPERATION_NOT_ALLOWED, TOO_MANY_ATTEMPTS_TRY_LATER
*/
func (ipc *client) SignUpWithEmailPassword(ctx context.Context, email string, password string) (*SignUpWithEmailPasswordResponse, error) {
	requestBody := SignUpRequestBody{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
		TenantId:          "",
	}
	log.Debug().Msgf("SignUpWithEmailPassword: %s:%s", requestBody.Email, requestBody.Password)
	var responseBody SignUpWithEmailPasswordResponse
	var errorBody IdentityToolkitError

	httpClient := resty.New()
	httpClient.SetLogger(&ZerologResty{log: log.Logger})
	defer httpClient.Close()

	res, err := httpClient.R().
		EnableTrace().
		SetContentType("application/json").
		SetMethod("POST").
		SetURL("https://identitytoolkit.googleapis.com/v1/accounts:signUp").
		SetQueryParam("key", ipc.apiKey).
		SetBody(requestBody).
		SetResult(&responseBody).
		SetError(&errorBody).
		Send()
	if err != nil {
		log.Error().Err(err).Msgf("http.Client error %s", err.Error())
		return nil, err
	}
	if res.IsError() {
		log.Error().Err(err).Msgf("http.Client.Response.IsError %s", res.Err.Error())
		return nil, errorBody
	}
	return &responseBody, nil
}

/*
SignInWithEmailPassword
endpoint https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=[API_KEY]
documentation https://cloud.google.com/identity-platform/docs/use-rest-api#section-sign-in-email-password
common error codes: EMAIL_NOT_FOUND, INVALID_PASSWORD, USER_DISABLED
*/
func (ipc *client) SignInWithEmailPassword(ctx context.Context, email string, password string) (*SignInWithEmailPasswordResponse, error) {
	requestBody := SignInRequestBody{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
		TenantId:          "",
	}
	var responseBody SignInWithEmailPasswordResponse
	var errorBody IdentityToolkitError

	httpClient := resty.New()
	httpClient.SetLogger(&ZerologResty{log: log.Logger})
	defer httpClient.Close()

	res, err := httpClient.R().
		EnableTrace().
		SetContentType("application/json").
		SetMethod("POST").
		SetURL("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword").
		SetQueryParam("key", ipc.apiKey).
		SetBody(requestBody).
		SetResult(&responseBody).
		SetError(&errorBody).
		Send()

	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, errorBody
	}
	return &responseBody, nil
}

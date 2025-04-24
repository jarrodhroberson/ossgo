package crypt

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"

	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"

	"github.com/jarrodhroberson/ossgo/gcp"
)

var malformedCypherTextError = errorx.NewErrorBuilder(errorx.IllegalFormat).WithCause(fmt.Errorf("malformed cypher text")).Create()

// generateRandomBytes generates random bytes with entropy sourced from the
// current location.
func GenerateRandomBytes(w io.Writer, numBytes int32) ([]byte, error) {
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create kms client: %w", err) //TODO: replace with errorx error
	}
	defer func(client *kms.KeyManagementClient) {
		err := client.Close()
		if err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
		}
	}(client)

	req := &kmspb.GenerateRandomBytesRequest{
		Location:        must.Must(gcp.Region()),
		LengthBytes:     numBytes,
		ProtectionLevel: kmspb.ProtectionLevel_HSM,
	}
	result, err := client.GenerateRandomBytes(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err) //TODO: replace with errorx error
	}
	return result.GetData(), nil
}

// Encrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Encrypt(plaintext []byte, key []byte) (ciphertext []byte, err error) {
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Decrypt(ciphertext []byte, key []byte) (plaintext []byte, err error) {
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, malformedCypherTextError
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

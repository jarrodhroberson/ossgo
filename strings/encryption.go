package strings

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/joomcode/errorx"

	"github.com/jarrodhroberson/ossgo/gcp"
)

var malformed_cypher_text_error = errorx.NewErrorBuilder(errorx.IllegalFormat).WithCause(fmt.Errorf("malformed cypher text")).Create()

// generateRandomBytes generates random bytes with entropy sourced from the
// given location.
func generateRandomBytes(numBytes int32) (string, error) {
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return NUL, err
	}
	defer func(client *kms.KeyManagementClient) {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}(client)

	// Build the request.
	req := &kmspb.GenerateRandomBytesRequest{
		// Location := "projects/my-project/locations/us-east1"
		Location:        fmt.Sprintf("projects/%s/locations/%s", gcp.Must(gcp.ProjectId()), gcp.Must(gcp.Region())),
		LengthBytes:     numBytes,
		ProtectionLevel: kmspb.ProtectionLevel_HSM,
	}

	result, err := client.GenerateRandomBytes(ctx, req)
	if err != nil {
		return NUL, err
	}

	encodedData := base64.StdEncoding.EncodeToString(result.Data)

	return encodedData, nil
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
		return nil, malformed_cypher_text_error
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}
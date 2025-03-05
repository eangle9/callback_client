package callbackreceiver

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func VerifyRequestHash(sk string, payload []byte, t string, hash string) ([]byte, error) {
	mac := hmac.New(sha256.New, []byte(sk))

	if _, err := mac.Write([]byte(t)); err != nil {
		return nil, err
	}

	if _, err := mac.Write([]byte(".")); err != nil {
		return nil, err
	}

	if _, err := mac.Write([]byte(payload)); err != nil {
		return nil, err
	}

	expectedHash := hex.EncodeToString(mac.Sum(nil))

	if expectedHash != hash {
		return nil, fmt.Errorf("hash verification failed: expected %s, but got %s", expectedHash, hash)
	}

	return payload, nil
}

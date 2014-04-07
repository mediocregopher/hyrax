package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"

	"github.com/mediocregopher/hyrax/server/config"
	"github.com/mediocregopher/hyrax/types"
)

// Auth checks whether the given command is authorized given its secret as-is.
// It returns a boolean of the result, or an error if something went wrong
// checking
func Auth(cmd *types.Action) (bool, error) {
	if !config.UseGlobalAuth && !config.UseKeyAuth {
		return true, nil
	}

	cmdB := []byte(cmd.Command)
	keyB := []byte(cmd.StorageKey)
	idB := []byte(cmd.Id)

	if config.UseGlobalAuth {
		for _, secret := range GetGlobalSecrets() {
			if ok := checkSecret(secret, cmdB, keyB, idB, cmd.Secret); ok {
				return true, nil
			}
		}
	}

	if config.UseKeyAuth {
		for _, secret := range GetKeySecrets(cmd.StorageKey) {
			if ok := checkSecret(secret, cmdB, keyB, idB, cmd.Secret); ok {
				return true, nil
			}
		}
	}

	return false, nil
}

func checkSecret(secret, cmd, key, id []byte, cmdSecret string) bool {
	mac := hmac.New(sha1.New, secret)
	mac.Write(cmd)
	mac.Write(key)
	mac.Write(id)
	sum := hex.EncodeToString(mac.Sum(nil))
	return sum == cmdSecret
}

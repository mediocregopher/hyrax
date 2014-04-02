package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"

	"github.com/mediocregopher/hyrax/types"
)

// Auth checks whether the given command is authorized given its secret as-is.
// It returns a boolean of the result, or an error if something went wrong
// checking
func Auth(cmd *types.Action) (bool, error) {
	// If no global serets are set, then everything is allowed
	globalSecs := GetGlobalSecrets()
	if len(globalSecs) == 0 {
		return true, nil
	}

	for _, secret := range GetGlobalSecrets() {
		if ok := checkSecret(secret, cmd); ok {
			return true, nil
		}
	}

	//keySecrets, err := GetSecrets(cmd.StorageKey)
	//if err != nil {
	//	return false, err
	//}

	//for _, secret := range keySecrets {
	//	if ok := checkSecret(secret, cmd); ok {
	//		return true, nil
	//	}
	//}

	return false, nil
}

func checkSecret(secret []byte, cmd *types.Action) bool {
	mac := hmac.New(sha1.New, secret)
	mac.Write([]byte(cmd.Command))
	mac.Write([]byte(cmd.StorageKey))
	mac.Write([]byte(cmd.Id))
	sum := mac.Sum(nil)
	sumhex := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(sumhex, sum)
	return hmac.Equal(sumhex, []byte(cmd.Secret))
}

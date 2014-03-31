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
func Auth(cmd *types.ClientCommand) (bool, error) {
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

func checkSecret(secret []byte, cmd *types.ClientCommand) bool {
	mac := hmac.New(sha1.New, secret)
	mac.Write(cmd.Command)
	mac.Write(cmd.StorageKey)
	mac.Write(cmd.Id)
	sum := mac.Sum(nil)
	sumhex := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(sumhex, sum)
	return hmac.Equal(sumhex, cmd.Secret)
}

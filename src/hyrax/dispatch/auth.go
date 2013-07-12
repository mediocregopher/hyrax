package dispatch

import (
    "crypto/sha1"
    "encoding/hex"
    "hyrax/types"
)

var secretkeys []string

// SetSecretkeys sets the list of secret keys to a different list
func SetSecretKeys(keys []string) {
    secretkeys = keys
}

// GetSecretKeys returns the list of keys currently in use
func GetSecretKeys() []string {
    return secretkeys
}

// Given a command payload (which presumably has a secret set), checks
// whether that secret checks out for one of the secret keys
func CheckAuth(cmdP *types.Payload) bool {
    h := sha1.New()
    for i := range secretkeys {
        h.Write( []byte(cmdP.Domain+cmdP.Name+secretkeys[i]) )
        if hex.EncodeToString(h.Sum(nil)) == cmdP.Secret { return true }
        h.Reset()
    }
    return false
}

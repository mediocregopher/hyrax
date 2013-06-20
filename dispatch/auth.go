package dispatch

import (
    "crypto/sha1"
    "encoding/hex"
    "hyrax/types"
)

var secretkeys []string

func SetSecretKeys(keys []string) {
    secretkeys = keys
}

func GetSecretKeys() []string {
    return secretkeys
}

func CheckAuth(cmdP *types.Payload) bool {
    h := sha1.New()
    for i := range secretkeys {
        h.Write( []byte(cmdP.Domain+cmdP.Name+secretkeys[i]) )
        if hex.EncodeToString(h.Sum(nil)) == cmdP.Secret { return true }
        h.Reset()
    }
    return false
}

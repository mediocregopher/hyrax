package dispatch

import (
    "crypto/sha1"
    "encoding/hex"
    "hyrax/types"
    "bytes"
)

var secretkeys [][]byte

// SetSecretkeys sets the list of secret keys to a different list
func SetSecretKeys(keys [][]byte) {
    secretkeys = keys
}

// GetSecretKeys returns the list of keys currently in use
func GetSecretKeys() [][]byte {
    return secretkeys
}

// Given a command payload (which presumably has a secret set), checks
// whether that secret checks out for one of the secret keys
func CheckAuth(cmdP *types.Payload) bool {
    h := sha1.New()
    for i := range secretkeys {
        h.Write( authMsg(cmdP.Domain,cmdP.Name,secretkeys[i]) )
        sum := h.Sum(nil)
        sumencodedsize := hex.EncodedLen(len(sum))
        sumencoded := make([]byte,sumencodedsize)
        hex.Encode(sumencoded,sum)
        if bytes.Equal(sumencoded,cmdP.Secret) { return true }
        h.Reset()
    }
    return false
}

func authMsg(domain, name, secret []byte) []byte {
    dl := len(domain)
    dnl := dl + len(name)
    buf := make([]byte,dnl+len(secret))
    copy(buf,domain)
    copy(buf[dl:],name)
    copy(buf[dnl:],secret)
    return buf
}

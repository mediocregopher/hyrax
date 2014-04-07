package auth

var globalSecrets = [][]byte{}
var getSecretsCh = make(chan chan [][]byte)
var setSecretsCh = make(chan [][]byte)

type keySecretsCast struct {
	key     string
	secrets [][]byte
}

type keySecretsCall struct {
	key   string
	retCh chan [][]byte
}

var keySecrets = map[string]map[string]bool{}
var getKeySecretsCh = make(chan *keySecretsCall)
var setKeySecretsCh = make(chan *keySecretsCast)
var addKeySecretsCh = make(chan *keySecretsCast)
var remKeySecretsCh = make(chan *keySecretsCast)

func init() {
	go keeper()
}

func keeper() {
	for {
		select {
		case retCh := <-getSecretsCh:
			retCh <- globalSecrets
		case secrets := <-setSecretsCh:
			globalSecrets = secrets
		case call := <-getKeySecretsCh:
			if secrets, ok := keySecrets[call.key]; ok {
				ret := make([][]byte, 0, len(secrets))
				for secretStr := range secrets {
					ret = append(ret, []byte(secretStr))
				}
				call.retCh <- ret
			} else {
				call.retCh <- [][]byte{}
			}
		case cast := <-setKeySecretsCh:
			secrets := map[string]bool{}
			for _, secret := range cast.secrets {
				secrets[string(secret)] = true
			}
			keySecrets[cast.key] = secrets
		case cast := <-addKeySecretsCh:
			secrets, ok := keySecrets[cast.key]
			if !ok {
				secrets = map[string]bool{}
				keySecrets[cast.key] = secrets
			}
			for _, secret := range cast.secrets {
				secrets[string(secret)] = true
			}
		case cast := <-remKeySecretsCh:
			if secrets, ok := keySecrets[cast.key]; ok {
				for _, secret := range cast.secrets {
					delete(secrets, string(secret))
				}
				if len(secrets) == 0 {
					delete(keySecrets, cast.key)
				}
			}
		}
	}
}

// GetGlobalSecrets returns the list of currently active global secrets
func GetGlobalSecrets() [][]byte {
	retCh := make(chan [][]byte)
	getSecretsCh <- retCh
	return <-retCh
}

// Overwrites the current list of global secrets
func SetGlobalSecrets(secrets [][]byte) {
	setSecretsCh <- secrets
}

// Gets the list of currently active secrets for a specific key
func GetKeySecrets(key string) [][]byte {
	call := keySecretsCall{key, make(chan [][]byte)}
	getKeySecretsCh <- &call
	return <-call.retCh
}

// Overwrites the active list of secrets for a specific key
func SetKeySecrets(key string, secrets [][]byte) {
	setKeySecretsCh <- &keySecretsCast{key, secrets}
}

// Adds secrets to the list of active secrets for a specific key
func AddKeySecrets(key string, secrets [][]byte) {
	addKeySecretsCh <- &keySecretsCast{key, secrets}
}

// Removes secrets from the list of active secrets for a specific key
func RemKeySecrets(key string, secrets [][]byte) {
	remKeySecretsCh <- &keySecretsCast{key, secrets}
}

package auth

import (
	"bytes"
	"github.com/grooveshark/golib/gslog"
)

var globalSecrets = [][]byte{}

var getSecretsCh = make(chan chan [][]byte)
var addSecretCh = make(chan []byte)
var remSecretCh = make(chan []byte)

func init() {
	go keeper()
}

// keeper is a little loop which keeps track of the global secrets on each hyrax
// node so that they are readily available for all operations. Secrets for
// individual keys are kept on the node for that key
func keeper() {
	for {
		select {
		case retCh := <-getSecretsCh:
			retCh <- globalSecrets
		case secret := <-addSecretCh:
			globalSecrets = append(globalSecrets, secret)
		case secret := <-remSecretCh:
			newgs := make([][]byte, 0, len(globalSecrets))
			for i := range globalSecrets {
				if !bytes.Equal(globalSecrets[i], secret) {
					newgs = append(newgs, globalSecrets[i])
				}
			}
			globalSecrets = newgs
		}
	}
}

// TODO make sure global secrets make it to all other nodes
// TODO make a way to "set" global secrets, so one node can sync them to other
//      nodes

// AddGlobalSecret appends a secret to the list of global ones currently in use.
func AddGlobalSecret(s []byte) {
	gslog.Infof("Loading secret:", string(s))
	addSecretCh <- s
}

// RemGlobalSecret removes all instances of the given secret from the list of
// global secrets
func RemGlobalSecret(s []byte) {
	remSecretCh <- s
}

// GetGlobalSecrets returns the list of currently active global secrets
func GetGlobalSecrets() [][]byte {
	retCh := make(chan [][]byte)
	getSecretsCh <- retCh
	return <-retCh
}

//var secretns = []byte("sec")
//
//// AddSecret adds a secret to the set of valid secrets for a given key
//func AddSecret(key, s []byte) error {
//	secKey := storage.KeyMaker.Namespace(secretns, key)
//	addCmd := storage.CommandFactory.GenericSetAdd(secKey, s)
//	_, err := storage.RoutedCmd(key, addCmd)
//	return err
//}
//
//// RemSecret removes a secret from the set of valid secrets for a given key, if
//// it existed in it at all
//func RemSecret(key, s []byte) error {
//	secKey := storage.KeyMaker.Namespace(secretns, key)
//	remCmd := storage.CommandFactory.GenericSetRem(secKey, s)
//	_, err := storage.RoutedCmd(key, remCmd)
//	return err
//}
//
//// GetSecrets returns the set of valid secrets for a given key
//func GetSecrets(key []byte) ([][]byte, error) {
//	secKey := storage.KeyMaker.Namespace(secretns, key)
//	getCmd := storage.CommandFactory.GenericSetMembers(secKey)
//	r, err := storage.RoutedCmd(key, getCmd)
//	if err != nil {
//		return nil, err
//	}
//	return r.([][]byte), nil
//}

package core

import (
	"errors"
	"github.com/mediocregopher/hyrax/server/auth"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
)

// TODO make success condition not be nil

var wrongNumArgs = errors.New("wrong number of arguments")
var wrongArgType = errors.New("wrong argument type")

func argsToEndpoint(cmd *types.ClientCommand) (string, error) {
	if len(cmd.Args) != 1 {
		return "", wrongNumArgs
	} else if addr, ok := cmd.Args[0].(string); ok {
		return addr, nil
	} else {
		return "", wrongArgType
	}
}

// If another node calls ALISTENTOME it is insisting that we ask it for its
// local keychange events
func AListenToMe(
	_ stypes.Client, cmd *types.ClientCommand) (interface{}, error) {

	listenEndpoint, err := argsToEndpoint(cmd)
	if err != nil {
		return nil, err
	}

	return OK, PullFromLocalManager.EnsureClient(listenEndpoint)
}

// If another node calls AIGNOREME it is insisting that we stop caring about its
// local keychange events
func AIgnoreMe(_ stypes.Client, cmd *types.ClientCommand) (interface{}, error) {
	listenEndpoint, err := argsToEndpoint(cmd)
	if err != nil {
		return nil, err
	}
	return OK, PullFromLocalManager.CloseClient(listenEndpoint)
}

func argsToByteSlice(cmd *types.ClientCommand) ([]byte, error) {
	if len(cmd.Args) != 1 {
		return nil, wrongNumArgs
	} else if s, ok := cmd.Args[0].(string); ok {
		return []byte(s), nil
	} else {
		return nil, wrongArgType
	}
}

// AGlobalSecretAdd adds a global secret to every node in the mesh's global
// secret list
func AGlobalSecretAdd(
	_ stypes.Client,
	cmd *types.ClientCommand) (interface{}, error) {
	secret, err := argsToByteSlice(cmd)
	if err != nil {
		return nil, err
	}

	auth.AddGlobalSecret(secret)
	return OK, nil
}

// AGlobalSecretRem removes a global secret from every node in the mesh's global
// secret list
func AGlobalSecretRem(
	_ stypes.Client,
	cmd *types.ClientCommand) (interface{}, error) {
	secret, err := argsToByteSlice(cmd)
	if err != nil {
		return nil, err
	}

	auth.RemGlobalSecret(secret)
	return OK, nil
}

// AGlobalSecrets returns the list of currently active secrets on this node (and
// presumably every other node)
func AGlobalSecrets(
	_ stypes.Client,
	cmd *types.ClientCommand) (interface{}, error) {
	secretsB := auth.GetGlobalSecrets()
	secrets := make([]string, len(secretsB))
	for i := range secretsB {
		secrets[i] = string(secretsB[i])
	}
	return secrets, nil
}

// ASecretAdd adds a particular secret to an individual key
//func ASecretAdd(
//	_ stypes.Client,
//	cmd *types.ClientCommand) (interface{}, error) {
//	secret, err := argsToByteSlice(cmd)
//	if err != nil {
//		return nil, err
//	}
//
//	return OK, auth.AddSecret(cmd.StorageKey, secret)
//}

// ASecretRem removes a particular secret from an individual key
//func ASecretRem(
//	_ stypes.Client,
//	cmd *types.ClientCommand) (interface{}, error) {
//	secret, err := argsToByteSlice(cmd)
//	if err != nil {
//		return nil, err
//	}
//
//	return OK, auth.RemSecret(cmd.StorageKey, secret)
//}

// ASecrets returns all the particular secrets for an individual key
//func ASecrets(_ stypes.Client, cmd *types.ClientCommand) (interface{}, error) {
//	keyB, err := argsToByteSlice(cmd)
//	if err != nil {
//		return nil, err
//	}
//
//	secretsB, err := auth.GetSecrets(keyB)
//	if err != nil {
//		return nil, err
//	}
//
//	secrets := make([]string, len(secretsB))
//	for i := range secretsB {
//		secrets[i] = string(secretsB[i])
//	}
//
//	return secrets, nil
//}

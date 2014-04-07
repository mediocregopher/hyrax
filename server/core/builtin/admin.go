package builtin

import (
	"errors"

	"github.com/mediocregopher/hyrax/server/auth"
	"github.com/mediocregopher/hyrax/server/core/dist"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
)

var wrongNumArgs = errors.New("wrong number of arguments")
var wrongArgType = errors.New("wrong argument type")

func argsToEndpoint(cmd *types.Action) (*types.ListenEndpoint, error) {
	if len(cmd.Args) != 1 {
		return nil, wrongNumArgs
	} else if addr, ok := cmd.Args[0].(string); ok {
		return types.ListenEndpointFromString(addr)
	} else {
		return nil, wrongArgType
	}
}

// If another node calls ALISTENTOME it is insisting that we ask it for its
// local keychange events
func AListenToMe(
	_ stypes.Client, cmd *types.Action) (interface{}, error) {

	listenEndpoint, err := argsToEndpoint(cmd)
	if err != nil {
		return nil, err
	}

	return OK, dist.PullFromLocalManager.EnsureClient(listenEndpoint)
}

// If another node calls AIGNOREME it is insisting that we stop caring about its
// local keychange events
func AIgnoreMe(_ stypes.Client, cmd *types.Action) (interface{}, error) {
	listenEndpoint, err := argsToEndpoint(cmd)
	if err != nil {
		return nil, err
	}
	return OK, dist.PullFromLocalManager.CloseClient(listenEndpoint)
}

func argsToByteSliceSlice(cmd *types.Action) ([][]byte, error) {
	if len(cmd.Args) < 1 {
		return nil, wrongNumArgs
	}
	bss := make([][]byte, len(cmd.Args))
	for i := range cmd.Args {
		if s, ok := cmd.Args[i].(string); ok {
			bss[i] = []byte(s)
		} else {
			return nil, wrongArgType
		}
	}
	return bss, nil
}

// AGlobalSecrets returns the list of currently active global secrets on this
// node. The list of global secrets is set in the configuration file, and can be
// changed by issuing a reload
func AGlobalSecrets(_ stypes.Client, cmd *types.Action) (interface{}, error) {
	secretsB := auth.GetGlobalSecrets()
	secrets := make([]string, len(secretsB))
	for i := range secretsB {
		secrets[i] = string(secretsB[i])
	}
	return secrets, nil
}

// ASecretsSet overwrites the set of secrets for a specific key
func ASecretsSet(_ stypes.Client, cmd *types.Action) (interface{}, error) {
	secrets, err := argsToByteSliceSlice(cmd)
	if err != nil {
		return nil, err
	}

	auth.SetKeySecrets(cmd.StorageKey, secrets)
	return OK, nil
}

// ASecretsAdd adds secrets to an individual key
func ASecretsAdd(_ stypes.Client, cmd *types.Action) (interface{}, error) {
	secrets, err := argsToByteSliceSlice(cmd)
	if err != nil {
		return nil, err
	}

	auth.AddKeySecrets(cmd.StorageKey, secrets)
	return OK, nil
}

// ASecretsRem removes secrets from an individual key
func ASecretsRem(_ stypes.Client, cmd *types.Action) (interface{}, error) {
	secrets, err := argsToByteSliceSlice(cmd)
	if err != nil {
		return nil, err
	}

	auth.RemKeySecrets(cmd.StorageKey, secrets)
	return OK, nil
}

// ASecrets returns all secrets currently active for an individual key
func ASecrets(_ stypes.Client, cmd *types.Action) (interface{}, error) {
	secretsB := auth.GetKeySecrets(cmd.StorageKey)
	secrets := make([]string, len(secretsB))
	for i := range secretsB {
		secrets[i] = string(secretsB[i])
	}
	return secrets, nil
}

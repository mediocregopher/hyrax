package builtin

import (
	"errors"
	"github.com/mediocregopher/hyrax/server/auth"
	"github.com/mediocregopher/hyrax/server/dist"
	"github.com/mediocregopher/hyrax/server/storage-router"
	stypes "github.com/mediocregopher/hyrax/server/types"
	"github.com/mediocregopher/hyrax/types"
)

// TODO make success condition not be nil

var wrongNumArgs = errors.New("wrong number of arguments")
var wrongArgType = errors.New("wrong argument type")

func argsToAddr(cmd *types.ClientCommand) (*string, error) {
	if len(cmd.Args) != 1 {
		return nil, wrongNumArgs
	} else if addr, ok := cmd.Args[0].(string); ok {
		return &addr, nil
	} else {
		return nil, wrongArgType
	}
}

// ANodeAdd uses the first argument as a node address, and adds that address to
// the hyrax mesh this one is in.
func ANodeAdd(
	_ stypes.ClientId,
	cmd *types.ClientCommand) (interface{}, error) {
	node, err := argsToAddr(cmd)
	if err != nil {
		return nil, err
	}

	return nil, dist.AddNode(node)
}

// ANodeRem uses the first argument as a node address, and removes that node
// from the hyrax mesh, if it was in there to begin with
func ANodeRem(
	_ stypes.ClientId,
	cmd *types.ClientCommand) (interface{}, error) {
	node, err := argsToAddr(cmd)
	if err != nil {
		return nil, err
	}

	dist.RemNode(node)
	return nil, nil
}

// TODO command to list nodes

// ABucketSet sets the storage bucket at the index given as the first argument.
// The second and third arguments are the connection type and address, and the
// rest are extra that are passed through depending on the storage type
func ABucketSet(
	_ stypes.ClientId,
	cmd *types.ClientCommand) (interface{}, error) {
	if len(cmd.Args) >= 3 {
		return nil, wrongNumArgs
	}

	var ok bool
	var bucket int
	var conntype, addr string
	var extra []interface{}

	if bucket, ok = cmd.Args[0].(int); !ok {
		return nil, wrongArgType
	} else if conntype, ok = cmd.Args[1].(string); !ok {
		return nil, wrongArgType
	} else if addr, ok = cmd.Args[2].(string); !ok {
		return nil, wrongArgType
	}

	extra = cmd.Args[3:]
	return nil, router.SetBucket(bucket, conntype, addr, extra...)
}

// ABuckets returns the current bucket list
func ABuckets(_ stypes.ClientId, _ *types.ClientCommand) (interface{}, error) {
	return router.GetBuckets, nil
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
	_ stypes.ClientId,
	cmd *types.ClientCommand) (interface{}, error) {
	secret, err := argsToByteSlice(cmd)
	if err != nil {
		return nil, err
	}

	auth.AddGlobalSecret(secret)
	return nil, nil
}

// AGlobalSecretRem removes a global secret from every node in the mesh's global
// secret list
func AGlobalSecretRem(
	_ stypes.ClientId,
	cmd *types.ClientCommand) (interface{}, error) {
	secret, err := argsToByteSlice(cmd)
	if err != nil {
		return nil, err
	}

	auth.RemGlobalSecret(secret)
	return nil, nil
}

// AGlobalSecrets returns the list of currently active secrets on this node (and
// presumably every other node)
func AGlobalSecrets(
	_ stypes.ClientId,
	cmd *types.ClientCommand) (interface{}, error) {
	secretsB := auth.GetGlobalSecrets()
	secrets := make([]string, len(secretsB))
	for i := range secretsB {
		secrets[i] = string(secretsB[i])
	}
	return secrets, nil
}

// ASecretAdd adds a particular secret to an individual key
func ASecretAdd(
	_ stypes.ClientId,
	cmd *types.ClientCommand) (interface{}, error) {
	secret, err := argsToByteSlice(cmd)
	if err != nil {
		return nil, err
	}

	return nil, auth.AddSecret(cmd.StorageKey, secret)
}

// ASecretRem removes a particular secret from an individual key
func ASecretRem(
	_ stypes.ClientId,
	cmd *types.ClientCommand) (interface{}, error) {
	secret, err := argsToByteSlice(cmd)
	if err != nil {
		return nil, err
	}

	return nil, auth.RemSecret(cmd.StorageKey, secret)
}

// ASecrets returns all the particular secrets for an individual key
func ASecrets(
	_ stypes.ClientId,
	cmd *types.ClientCommand) (interface{}, error) {
	keyB, err := argsToByteSlice(cmd)
	if err != nil {
		return nil, err
	}

	secretsB, err := auth.GetSecrets(keyB)
	if err != nil {
		return nil, err
	}

	secrets := make([]string, len(secretsB))
	for i := range secretsB {
		secrets[i] = string(secretsB[i])
	}

	return secrets, nil
}

package main

import (
	"bytes"
	"hyrax-server/config"
	"hyrax-server/dispatch"
	"hyrax-server/net"
	"hyrax-server/storage"
	"strconv"
)

func main() {
	config.LoadConfig()

	keys := bytes.Split([]byte(config.GetStr("initial-secret-keys")), []byte{':'})
	dispatch.SetSecretKeys(keys)

	err := storage.RedisConnect()
	if err != nil {
		panic(err)
	}

	err = CleanupTransientData()
	if err != nil {
		panic(err)
	}

	port := strconv.Itoa(config.GetInt("port"))
	addr := ":" + port
	err = net.TcpListen(addr)
	if err != nil {
		panic(err)
	}

	select {}
}

// CleanupTransientData uses AllWildcards to go through and delete all keys
// containing data related to connections from the last instance of hyrax
// that existed on the redis instance
func CleanupTransientData() error {
	queries := storage.AllWildcards()
	for i := range queries {
		keysr, err := storage.CmdPretty(storage.KEYS, queries[i])
		if err != nil {
			return err
		}

		keys := keysr.([][]byte)
		for j := range keys {
			_, err := storage.CmdPretty(storage.DEL, keys[j])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

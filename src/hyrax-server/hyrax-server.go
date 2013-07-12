package main

import (
    "strings"
    "strconv"
    "hyrax/config"
    "hyrax/storage"
    "hyrax/dispatch"
    "hyrax/net"
)

func main() {
    config.LoadConfig()

    keys := strings.Split(config.GetStr("initial-secret-keys"),":")
    dispatch.SetSecretKeys(keys)

    err := storage.RedisConnect()
    if err != nil { panic(err); }

    err = CleanupTransientData()
    if err != nil { panic(err); }

    port := strconv.Itoa(config.GetInt("port"))
    addr := ":"+port
    err = net.TcpListen(addr)
    if err != nil { panic(err); }

    select {}
}

// CleanupTransientData uses AllWildcards to go through and delete all keys
// containing data related to connections from the last instance of hyrax
// that existed on the redis instance
func CleanupTransientData() error {
    queries := storage.AllWildcards()
    for i := range queries {
        keysr,err := storage.CmdPretty("KEYS",queries[i])
        if err != nil { return err }

        keys := keysr.([]string)
        for j := range keys {
            _,err := storage.CmdPretty("DEL",keys[j])
            if err != nil { return err }
        }
    }

    return nil
}

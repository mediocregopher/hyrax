package main

import (
    "fmt"
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

    port := strconv.Itoa(config.GetInt("port"))
    addr := ":"+port
    err = net.TcpListen(addr)
    if err != nil { panic(err); }

    //cmd(`{"command":"set","payload":{"domain":"a","id":"k1","values":["wut012345wut"],"secret":"dea83285cb755ddb47e2b24b68b5321f394e3641"}}`)
    //cmd(`{"command":"getrange","payload":{"domain":"a","id":"k1","values":["0","2"]}}`)

    select {}
}

func cmd(c string) {
    r,err := dispatch.DoCommand([]byte(c))
    fmt.Println(string(r),err)
}

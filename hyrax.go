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

    //BUG(mediocregopher): proper cleanup
    _,err = storage.CmdPretty("FLUSHALL")
    if err != nil { panic(err); }

    port := strconv.Itoa(config.GetInt("port"))
    addr := ":"+port
    err = net.TcpListen(addr)
    if err != nil { panic(err); }

    //cmd(`{"command":"amon","payload":{"domain":"a","id":"k1"}}`)
    //cmd(`{"command":"set","payload":{"domain":"a","id":"k1","secret":"dea83285cb755ddb47e2b24b68b5321f394e3641","values":["ohai"]}}`)

    select {}
}

func cmd(c string) {
    r,err := dispatch.DoCommand(0,[]byte(c))
    fmt.Println(string(r),err)
}

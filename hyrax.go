package main

import (
    "fmt"
    "hyrax/config"
    "hyrax/storage"
    "hyrax/dispatch"
)

func main() {
    config.LoadConfig()
    err := storage.RedisConnect()
    if err != nil { panic(err); }

    cmd(`{"command":"set","payload":[{"domain":"t2","id":"k1","values":["012345"]}]}`)
    cmd(`{"command":"getrange","payload":[{"domain":"t2","id":"k1","values":["0","2"]}]}`)
}

func cmd(c string) {
    r,err := dispatch.DoCommand([]byte(c))
    fmt.Println(string(r),err)
}
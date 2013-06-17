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

    cmd(`{"command":"mset","payload":[{"domain":"t1","id":"k1","values":["v1"]},
                                      {"domain":"t1","id":"k2","values":["v2"]}]}`)
    cmd(`{"command":"mget","payload":[{"domain":"t1","id":"k1"},
                                      {"domain":"t1","id":"k2"}]}`)
}

func cmd(c string) {
    r,err := dispatch.DoCommand([]byte(c))
    fmt.Println(string(r),err)
}

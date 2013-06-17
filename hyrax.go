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

    cmd(`{"command":"hgetall","payload":[{"domain":"t1","id":"m1"}]}`)
    cmd(`{"command":"hset","payload":[{"domain":"t1","id":"m1","values":["b","a"]}]}`)
    cmd(`{"command":"hgetall","payload":[{"domain":"t1","id":"m1"}]}`)
}

func cmd(c string) {
    r,err := dispatch.DoCommand([]byte(c))
    fmt.Println(string(r),err)
}

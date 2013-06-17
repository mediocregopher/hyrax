package main

import (
    //"fmt"
    "hyrax/config"
    "hyrax/storage"
)

func main() {
    config.LoadConfig()
    err := storage.RedisConnect()
    if err != nil { panic(err); }

}

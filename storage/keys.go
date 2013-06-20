package storage

import (
    "strings"
)

const SEP string = ":"

func createKey(pieces... string) string {
    return strings.Join(pieces,SEP)
}

func DirectKey(domain,id string) string {
    return createKey("direct",domain,id)
}

func MonKey(domain,id string) string {
    return createKey("mon",domain,id)
}

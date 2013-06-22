package storage

import (
    "strings"
    "hyrax/types"
)

const SEP string = ":"

func createKey(pieces... string) string {
    return strings.Join(pieces,SEP)
}

// DirectKey returns the key for redis that will be used as a key
// for commands that interact directly with redis
func DirectKey(domain,id string) string {
    return createKey("direct",domain,id)
}

// MonKey returns the key that will be used to store the set of
// connection ids that are monitoring a value
func MonKey(domain,id string) string {
    return createKey("mon",domain,id)
}

// ConnMonKey returns the key that will be used to store the set
// of values being monitored by a connection
func ConnMonKey(cid types.ConnId) string {
    return createKey("conn","mon",cid.Serialize())
}

// ConnMonVal returns the value that will be stored at a ConnMonKey
func ConnMonVal(domain, id string) string {
    //We use createKey cause it does what we want, even though
    //we're actually making the value that's going to be set
    return createKey(domain,id)
}

// DeconstructConnMonVal gets the domain and id from a ConnMonVal
func DeconstructConnMonVal(connmonval string) (string,string) {
    s := strings.Split(connmonval,SEP)
    return s[0],s[1]
}

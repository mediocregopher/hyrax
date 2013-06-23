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

// EkgKey returns the key that will be used to store the set
// of connection id/names being monitored by the ekg
func EkgKey(domain, id string) string {
    return createKey("ekg",domain,id)
}

// EkgVal returns the value that will be stored at an EkgKey
func EkgVal(cid types.ConnId, name string) string {
    return createKey(cid.Serialize(),name)
}

// DeconstructEkgVal returns the connection id and name being
// represented by the given EkgVal
func DeconstructEkgVal(ekgval string) (types.ConnId,string) {
    s := strings.Split(ekgval,SEP)

    cid,err := types.ConnIdDeserialize(s[0])
    if err != nil { panic(err) }

    return cid,s[1]
}

// ConnEkgKey returns the key that will be used to store the set
// of ekgs that a connection is hooked up to
func ConnEkgKey(cid types.ConnId) string {
    return createKey("conn","ekg",cid.Serialize())
}

// ConnEkgVal returns the value that will be stored at the ConnEkgKey
func ConnEkgVal(domain, id, name string) string {
    return createKey(domain,id,name)
}

// DeconstructConnEkgVal returns the domain, id, and name from a ConnEkgVal
func DeconstructConnEkgVal(connekgval string) (string,string,string) {
    s := strings.Split(connekgval,SEP)
    return s[0],s[1],s[2]
}

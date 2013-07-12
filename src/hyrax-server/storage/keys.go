package storage

import (
    "bytes"
    "hyrax-server/types"
)

//These are for use by this and other modules so we don't have to
//re-allocate them everytime they get used
var SEP = []byte{':'}
var CONN = []byte("conn")
var DIRECT = []byte("direct")
var WILDCARD = []byte{'*'}

var SADD = []byte("SADD")
var SREM = []byte("SREM")
var SMEMBERS = []byte("SMEMBERS")
var SCARD = []byte("SCARD")
var DEL = []byte("DEL")
var KEYS = []byte("KEYS")
var OK = []byte("OK")

var DISCONNECT = []byte("disconnect")
var MONPUSH = []byte("mon-push")


func createKey(pieces... []byte) []byte {
    return bytes.Join(pieces,SEP)
}

// DirectKey returns the key for redis that will be used as a key
// for commands that interact directly with redis
func DirectKey(domain,id []byte) []byte {
    return createKey(DIRECT,domain,id)
}

////////////////////////////////////////////////////////////////////////
// Mon
////////////////////////////////////////////////////////////////////////

var MON = []byte("mon")

// MonKey returns the key that will be used to store the set of
// connection ids that are monitoring a value
func MonKey(domain,id []byte) []byte {
    return createKey(MON,domain,id)
}

// ConnMonKey returns the key that will be used to store the set
// of values being monitored by a connection
func ConnMonKey(cid types.ConnId) []byte {
    return createKey(CONN,MON,cid.Serialize())
}

// ConnMonVal returns the value that will be stored at a ConnMonKey
func ConnMonVal(domain, id []byte) []byte {
    //We use createKey cause it does what we want, even though
    //we're actually making the value that's going to be set
    return createKey(domain,id)
}

// DeconstructConnMonVal gets the domain and id from a ConnMonVal
func DeconstructConnMonVal(connmonval []byte) ([]byte,[]byte) {
    b := bytes.SplitN(connmonval,SEP,2)
    return b[0],b[1]
}

// MonWildcards returns the list of wildcarded keys that will cover
// all Mon related data in redis
func MonWildcards() [][]byte {
    return [][]byte{ createKey(MON,WILDCARD),
                     createKey(CONN,MON,WILDCARD) }
}

////////////////////////////////////////////////////////////////////////
// EKG
////////////////////////////////////////////////////////////////////////

var EKG = []byte("ekg")

// EkgKey returns the key that will be used to store the set
// of connection id/names being monitored by the ekg
func EkgKey(domain, id []byte) []byte {
    return createKey(EKG,domain,id)
}

// EkgVal returns the value that will be stored at an EkgKey
func EkgVal(cid types.ConnId, name []byte) []byte {
    return createKey(cid.Serialize(),name)
}

// DeconstructEkgVal returns the connection id and name being
// represented by the given EkgVal
func DeconstructEkgVal(ekgval []byte) (types.ConnId,[]byte) {
    b := bytes.SplitN(ekgval,SEP,2)

    cid,err := types.ConnIdDeserialize(b[0])
    if err != nil { panic(err) }

    return cid,b[1]
}

// ConnEkgKey returns the key that will be used to store the set
// of ekgs that a connection is hooked up to
func ConnEkgKey(cid types.ConnId) []byte {
    return createKey(CONN,EKG,cid.Serialize())
}

// ConnEkgVal returns the value that will be stored at the ConnEkgKey
func ConnEkgVal(domain, id, name []byte) []byte {
    return createKey(domain,id,name)
}

// DeconstructConnEkgVal returns the domain, id, and name from a ConnEkgVal
func DeconstructConnEkgVal(connekgval []byte) ([]byte,[]byte,[]byte) {
    b := bytes.SplitN(connekgval,SEP,3)
    return b[0],b[1],b[2]
}

// EkgWildcards returns the list of wildcarded keys that will cover
// all Ekg related data in redis
func EkgWildcards() [][]byte {
    return [][]byte{ createKey(EKG,WILDCARD),
                     createKey(CONN,EKG,WILDCARD) }
}

////////////////////////////////////////////////////////////////////////
// Util
////////////////////////////////////////////////////////////////////////

// AllWildcards returns the list of wildcarded keys that will cover
// all "transient" data (aka, all data related to connections that would
// become invalid on a server restart)
func AllWildcards() [][]byte {
    mon := MonWildcards()
    ekg := EkgWildcards()
    ret := make([][]byte,0,len(mon)+len(ekg))

    for i := range mon {
        ret = append(ret,mon[i])
    }

    for i := range ekg {
        ret = append(ret,ekg[i])
    }

    return ret
}

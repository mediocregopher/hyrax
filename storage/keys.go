package storage

import (
    "strings"
    "hyrax/types"
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

func ConnMonKey(cid types.ConnId) string {
    return createKey("conn","mon",cid.Serialize())
}

func ConnMonVal(domain, id string) string {
    //We use createKey cause it does what we want, even though
    //we're actually making the value that's going to be set
    return createKey(domain,id)
}

//Get the domain and id from the connmonval
func DeconstructConnMonVal(connmonval string) (string,string) {
    s := strings.Split(connmonval,SEP)
    return s[0],s[1]
}

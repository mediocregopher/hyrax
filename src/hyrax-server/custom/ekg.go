package custom

import (
    "hyrax/types"
    "hyrax-server/storage"
    "errors"
)

// EAdd adds the connection's id (and name) to an ekg's set of things it's
// watching, and adds the ekg's information to the connection's set of
// ekgs its hooked up to
func EAdd(cid types.ConnId, pay *types.Payload) (interface{},error) {
    connekgkey := storage.ConnEkgKey(cid)
    connekgval := storage.ConnEkgVal(pay.Domain,pay.Id,pay.Name)
    _,err := storage.CmdPretty("SADD",connekgkey,connekgval)
    if err != nil { return nil,err }

    ekgkey := storage.EkgKey(pay.Domain,pay.Id)
    ekgval := storage.EkgVal(cid,pay.Name)
    _,err = storage.CmdPretty("SADD",ekgkey,ekgval)
    return "OK",err
}

// ERem removes the connection's id (and name) from an ekg's set of things
// it's watching, and removes the ekg's information from the connection's
// set of ekgs its hooked up to
func ERem(cid types.ConnId, pay *types.Payload) (interface{},error) {
    ekgkey := storage.EkgKey(pay.Domain,pay.Id)
    ekgval := storage.EkgVal(cid,pay.Name)
    _,err := storage.CmdPretty("SREM",ekgkey,ekgval)
    if err != nil { return nil,err }

    connekgkey := storage.ConnEkgKey(cid)
    connekgval := storage.ConnEkgVal(pay.Domain,pay.Id,pay.Name)
    _,err = storage.CmdPretty("SREM",connekgkey,connekgval)
    return "OK",err
}

// CleanConnEkg takes in a connection id and cleans up all of its
// ekgs, and the set which keeps track of those ekgs. It also
// sends out alerts for all the ekgs it's hooked up to, since
// this only gets called on a disconnect.
func CleanConnEkg(cid types.ConnId) error {
    connekgkey := storage.ConnEkgKey(cid)
    r,err := storage.CmdPretty("SMEMBERS",connekgkey)
    if err != nil { return err }

    ekgs := r.([]string)

    for i := range ekgs {
        domain,id,name := storage.DeconstructConnEkgVal(ekgs[i])
        ekgkey := storage.EkgKey(domain,id)
        ekgval := storage.EkgVal(cid,name)
        _,err = storage.CmdPretty("SREM",ekgkey,ekgval)
        if err != nil { return err }

        cmd := types.Command{
            Command: "disconnect",
            Payload: types.Payload{
                Domain: domain,
                Id:     id,
                Name:   name,
            },
        }
        MonMakeAlert(&cmd)
    }

    _,err = storage.CmdPretty("DEL",connekgkey)
    return err

}

// EMembers returns the list of names being monitored by an ekg
func EMembers(cid types.ConnId, pay *types.Payload) (interface{},error) {
    ekgkey := storage.EkgKey(pay.Domain,pay.Id)
    r,err := storage.CmdPretty("SMEMBERS",ekgkey)
    if err != nil { return nil,err }

    members := r.([]string)
    for i := range members {
        _,name := storage.DeconstructEkgVal(members[i])
        members[i] = name
    }

    return members,nil
}

// ECard returns the number of connection/name combinations being monitored
func ECard(cid types.ConnId, pay *types.Payload) (interface{},error) {
    ekgkey := storage.EkgKey(pay.Domain,pay.Id)
    return storage.CmdPretty("SCARD",ekgkey)
}

// EIsMember returns whether or not the given name is being monitored by the ekg
func EIsMember(cid types.ConnId, pay *types.Payload) (interface{},error) {

    if !(len(pay.Values) > 0) {
        return nil,errors.New("ERR wrong number of arguments for 'eismember' command")
    }

    ekgkey := storage.EkgKey(pay.Domain,pay.Id)
    r,err := storage.CmdPretty("SMEMBERS",ekgkey)
    if err != nil { return nil,err }

    members := r.([]string)
    for i := range members {
        _,name := storage.DeconstructEkgVal(members[i])
        if name == pay.Values[0] { return 1,nil }
    }

    return 0,nil
}

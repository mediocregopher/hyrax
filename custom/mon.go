package custom

import (
    "hyrax/types"
    "hyrax/storage"
)

var monCh chan *types.Command

func AMon(cid types.ConnId, pay *types.Payload) (interface{},error) {
    monkey := storage.MonKey(pay.Domain,pay.Id)

    _,err := storage.CmdPretty("SADD",monkey,cid)
    return "OK",err
}

func Mon(cid types.ConnId, pay *types.Payload) (interface{},error) {
    _,err := AMon(cid,pay)
    if err != nil { return nil,err }
    dirkey := storage.DirectKey(pay.Domain,pay.Id)
    return storage.CmdPretty("GET",dirkey)
}

func HMon(cid types.ConnId, pay *types.Payload) (interface{},error) {
    _,err := AMon(cid,pay)
    if err != nil { return nil,err }
    dirkey := storage.DirectKey(pay.Domain,pay.Id)

    r,err := storage.CmdPretty("HGETALL",dirkey)
    if err != nil { return nil,err }

    return storage.RedisListToMap(r.([]string))
}

func LMon(cid types.ConnId, pay *types.Payload) (interface{},error) {
    _,err := AMon(cid,pay)
    if err != nil { return nil,err }
    dirkey := storage.DirectKey(pay.Domain,pay.Id)
    return storage.CmdPretty("LRANGE",dirkey,0,-1)
}

func SMon(cid types.ConnId, pay *types.Payload) (interface{},error) {
    _,err := AMon(cid,pay)
    if err != nil { return nil,err }
    dirkey := storage.DirectKey(pay.Domain,pay.Id)
    return storage.CmdPretty("SMEMBERS",dirkey)
}

func ZMon(cid types.ConnId, pay *types.Payload) (interface{},error) {
    _,err := AMon(cid,pay)
    if err != nil { return nil,err }
    dirkey := storage.DirectKey(pay.Domain,pay.Id)

    r,err := storage.CmdPretty("SMEMBERS",dirkey)
    if err != nil { return nil,err }

    return storage.RedisListToIntMap(r.([]string))
}

//TODO emon

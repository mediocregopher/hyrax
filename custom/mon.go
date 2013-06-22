package custom

import (
    "hyrax/types"
    "hyrax/storage"
    "hyrax/router"
    "hyrax/parse"
    "strconv"
    "log"
    "errors"
)


func AMon(cid types.ConnId, pay *types.Payload) (interface{},error) {
    monkey := storage.MonKey(pay.Domain,pay.Id)
    connmonkey := storage.ConnMonKey(cid)
    connmonval := storage.ConnMonVal(pay.Domain,pay.Id)

    _,err := storage.CmdPretty("SADD",monkey,cid)
    if err != nil { return nil,err }

    _,err = storage.CmdPretty("SADD",connmonkey,connmonval)
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

func CleanConnMon(cid types.ConnId) error {
    connmonkey := storage.ConnMonKey(cid)
    r,err := storage.CmdPretty("SMEMBERS",connmonkey)
    if err != nil { return err }

    mons := r.([]string)

    for i := range mons {
        domain,id := storage.DeconstructConnMonVal(mons[i])
        monkey := storage.MonKey(domain,id)
        _,err = storage.CmdPretty("SREM",monkey,cid)
        if err != nil { return err }
    }

    _,err = storage.CmdPretty("DEL",connmonkey)
    return err
}

var monCh chan *types.Command

type MonPushPayload struct {
    Domain  string   `json:"domain"`
    Id      string   `json:"id"`
    Command string   `json:"command"`
    Values  []string `json:"values"`
}

func MonMakeAlert(cmd *types.Command) {
    monCh <- cmd
}

func MonHandleAlert(cmd *types.Command) error {
    monkey := storage.MonKey(cmd.Payload.Domain,cmd.Payload.Id)
    r,err := storage.CmdPretty("SMEMBERS",monkey)
    if err != nil { return err }
    idstrs := r.([]string)

    if len(idstrs) == 0 { return nil }
    var pay MonPushPayload
    pay.Domain = cmd.Payload.Domain
    pay.Id = cmd.Payload.Id
    pay.Command = cmd.Command
    pay.Values = cmd.Payload.Values

    for i := range idstrs {
        id,err := strconv.Atoi(idstrs[i])
        if err != nil {
            return errors.New(err.Error()+" when converting "+idstrs[i]+" to int")
        }

        msg,err := parse.EncodeMessage("mon-push",pay)
        if err != nil {
            return errors.New(err.Error()+" when encoding mon push message")
        }

        router.SendPushMessage(types.ConnId(id),msg)
    }

    return nil

}

func init() {
    monCh = make(chan *types.Command)

    for i:=0; i<10; i++ {
        go func(){
            for {
                cmd := <-monCh
                err := MonHandleAlert(cmd)
                if err != nil {
                    log.Printf("%s when calling MonHandleAlert(%v)\n",err.Error(),cmd)
                }
            }
        }()
    }

}

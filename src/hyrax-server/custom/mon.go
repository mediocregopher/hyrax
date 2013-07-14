package custom

import (
    types  "hyrax/types"
    stypes "hyrax-server/types"
    . "hyrax-server/storage"
    "hyrax-server/router"
    "log"
    "errors"
)


// MAdd adds the connection's id to the set of connections that
// are monitoring the domain/id in redis (so it can receive alerts)
// and adds the domain/id to the set of domain/ids that the
// connection is monitoring in redis (so it can clean up when the connection
// closes).
func MAdd(cid stypes.ConnId, pay *types.Payload) (interface{},error) {
    monkey := MonKey(pay.Domain,pay.Id)
    connmonkey := ConnMonKey(cid)
    connmonval := ConnMonVal(pay.Domain,pay.Id)

    _,err := CmdPretty(SADD,connmonkey,connmonval)
    if err != nil { return nil,err }

    _,err = CmdPretty(SADD,monkey,cid)
    return "OK",err
}

// MRem removes the connection's id form the set of connections that
// are monitoring the domain/id in redis, and removes the domain/id
// from the set of domain/ids that the connection is monitoring in
// redis
func MRem(cid stypes.ConnId, pay *types.Payload) (interface{},error) {
    monkey := MonKey(pay.Domain,pay.Id)
    connmonkey := ConnMonKey(cid)
    connmonval := ConnMonVal(pay.Domain,pay.Id)

    _,err := CmdPretty(SREM,connmonkey,connmonval)
    if err != nil { return nil,err }

    _,err = CmdPretty(SREM,monkey,cid)
    return "OK",err
}

// CleanConnMon takes in a connection id and cleans up all of its
// monitors, and the set which keeps track of those monitors
func CleanConnMon(cid stypes.ConnId) error {
    connmonkey := ConnMonKey(cid)
    r,err := CmdPretty(SMEMBERS,connmonkey)
    if err != nil { return err }

    mons := r.([][]byte)

    for i := range mons {
        domain,id := DeconstructConnMonVal(mons[i])
        monkey := MonKey(domain,id)
        _,err = CmdPretty(SREM,monkey,cid)
        if err != nil { return err }
    }

    _,err = CmdPretty(DEL,connmonkey)
    return err
}

var monCh chan *types.Command

// MonMakeAlert takes in a command that's being performed and sends
// out alerts to anyone monitoring that command
func MonMakeAlert(cmd *types.Command) {
    monCh <- cmd
}

// monPushPayload is the payload for push notifications. It is basically
// the standard payload object but without the secret, and with a command
// string field instead
type monPushPayload struct {
    Domain  []byte   `json:"domain"`
    Id      []byte   `json:"id"`
    Name    []byte   `json:"name,omitempty"`
    Command []byte   `json:"command"`
    Values  [][]byte `json:"values,omitempty"`
}


// monHandleAlert takes commands to be alerted and does the alert
func monHandleAlert(cmd *types.Command) error {

    var pay monPushPayload
    pay.Domain = cmd.Payload.Domain
    pay.Id = cmd.Payload.Id
    pay.Name = cmd.Payload.Name
    pay.Command = cmd.Command
    pay.Values = cmd.Payload.Values

    return MonDoAlert(&pay)

}

// MonDoAlert actually does the fetching of monitors on a value and
// and sends them alerts
func MonDoAlert(pay *monPushPayload) error {
    monkey := MonKey(pay.Domain,pay.Id)
    r,err := CmdPretty(SMEMBERS,monkey)
    if err != nil { return err }
    idstrs := r.([][]byte)

    if len(idstrs) == 0 {
        return nil
    }

    msg,err := types.EncodeMessage(MONPUSH,pay)
    if err != nil {
        return errors.New(err.Error()+" when encoding mon push message")
    }

    for i := range idstrs {
        id,err := stypes.ConnIdDeserialize(idstrs[i])
        if err != nil {
            return errors.New(err.Error()+" when converting "+string(idstrs[i])+" to int")
        }

        router.SendPushMessage(stypes.ConnId(id),msg)
    }

    return nil
}

// init creates a bunch of routines that will read in commands that require alerts
// and call monHandleAlert on them
func init() {
    monCh = make(chan *types.Command)

    for i:=0; i<10; i++ {
        go func(){
            for {
                cmd := <-monCh
                err := monHandleAlert(cmd)
                if err != nil {
                    log.Printf("%s when calling monHandleAlert(%v)\n",err.Error(),cmd)
                }
            }
        }()
    }

}

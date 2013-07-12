package custom

import (
    "hyrax/types"
    "hyrax-server/storage"
    "hyrax-server/router"
    "strconv"
    "log"
    "errors"
)


// MAdd adds the connection's id to the set of connections that
// are monitoring the domain/id in redis (so it can receive alerts)
// and adds the domain/id to the set of domain/ids that the
// connection is monitoring in redis (so it can clean up when the connection
// closes).
func MAdd(cid types.ConnId, pay *types.Payload) (interface{},error) {
    monkey := storage.MonKey(pay.Domain,pay.Id)
    connmonkey := storage.ConnMonKey(cid)
    connmonval := storage.ConnMonVal(pay.Domain,pay.Id)

    _,err := storage.CmdPretty("SADD",connmonkey,connmonval)
    if err != nil { return nil,err }

    _,err = storage.CmdPretty("SADD",monkey,cid)
    return "OK",err
}

// MRem removes the connection's id form the set of connections that
// are monitoring the domain/id in redis, and removes the domain/id
// from the set of domain/ids that the connection is monitoring in
// redis
func MRem(cid types.ConnId, pay *types.Payload) (interface{},error) {
    monkey := storage.MonKey(pay.Domain,pay.Id)
    connmonkey := storage.ConnMonKey(cid)
    connmonval := storage.ConnMonVal(pay.Domain,pay.Id)

    _,err := storage.CmdPretty("SREM",connmonkey,connmonval)
    if err != nil { return nil,err }

    _,err = storage.CmdPretty("SREM",monkey,cid)
    return "OK",err
}

// CleanConnMon takes in a connection id and cleans up all of its
// monitors, and the set which keeps track of those monitors
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

// MonMakeAlert takes in a command that's being performed and sends
// out alerts to anyone monitoring that command
func MonMakeAlert(cmd *types.Command) {
    monCh <- cmd
}

// monHandleAlert takes commands to be alerted and does the alert
func monHandleAlert(cmd *types.Command) error {

    var pay types.MonPushPayload
    pay.Domain = cmd.Payload.Domain
    pay.Id = cmd.Payload.Id
    pay.Name = cmd.Payload.Name
    pay.Command = cmd.Command
    pay.Values = cmd.Payload.Values

    return MonDoAlert(&pay)

}

// MonDoAlert actually does the fetching of monitors on a value and
// and sends them alerts
func MonDoAlert(pay *types.MonPushPayload) error {
    monkey := storage.MonKey(pay.Domain,pay.Id)
    r,err := storage.CmdPretty("SMEMBERS",monkey)
    if err != nil { return err }
    idstrs := r.([]string)

    for i := range idstrs {
        id,err := strconv.Atoi(idstrs[i])
        if err != nil {
            return errors.New(err.Error()+" when converting "+idstrs[i]+" to int")
        }

        msg,err := types.EncodeMessage("mon-push",pay)
        if err != nil {
            return errors.New(err.Error()+" when encoding mon push message")
        }

        router.SendPushMessage(types.ConnId(id),msg)
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

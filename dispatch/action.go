package dispatch

import (
    "hyrax/storage"
    "hyrax/types"
    "hyrax/parse"
    "hyrax/custom"
    "errors"
)

func DoCommand(cid types.ConnId, rawJson []byte) ([]byte,error) {
    cmd,err := parse.DecodeCommand(rawJson)
    if err != nil {
        return parse.EncodeError(err.Error())
    }

    ret,err := doCommandWrap(cid,cmd)
    if err != nil {
        return parse.EncodeError(err.Error())
    }

    return parse.EncodeMessage(cmd.Command,ret)
}

func doCommandWrap(cid types.ConnId, cmd *types.Command) (interface{},error) {
    pay := &cmd.Payload

    if !CommandExists(&cmd.Command) {
        return nil,errors.New("Unsupported command")
    }

    if CommandModifies(&cmd.Command) {
        if !CheckAuth(pay) {
            return nil,errors.New("cannot authenticate with key "+pay.Secret)
        }
        custom.MonMakeAlert(cmd)
    }

    if pay.Id == "" {
        return nil,errors.New("missing key id")
    }

    if CommandIsCustom(&cmd.Command) {
        return doCustomCommand(cid,cmd)
    }

    numArgs := len(pay.Values)+1

    args := make([]interface{},0,numArgs)
    strKey := storage.DirectKey(pay.Domain,pay.Id)
    args = append(args,strKey)
    for j:=0; j<len(pay.Values); j++ {
        args = append(args,pay.Values[j])
    }

    r,err := storage.Cmd(cmd.Command,args)
    if err != nil { return nil,err }

    if CommandReturnsMap(&cmd.Command) {
        return storage.RedisListToMap(r.([]string))
    }

    return r,nil
}

func doCustomCommand(cid types.ConnId, cmd *types.Command) (interface{},error) {
    f,ok := customCommandMap[cmd.Command]
    if !ok { return nil,errors.New("Command in main map not listed in custom map") }

    return f(cid,&cmd.Payload)
}

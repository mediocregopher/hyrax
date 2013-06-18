package dispatch

import (
    "hyrax/storage"
    "errors"
)

func DoCommand(rawJson []byte) ([]byte,error) {
    cmd,err := DecodeCommand(rawJson)
    if err != nil {
        return EncodeError(err.Error())
    }

    ret,err := doCommandWrap(cmd)
    if err != nil {
        return EncodeError(err.Error())
    }

    return EncodeMessage(cmd.Command,ret)
}

func doCommandWrap(cmd *Command) (interface{},error) {
    pay := &cmd.Payload

    if !CommandExists(&cmd.Command) {
        return nil,errors.New("Unsupported command")
    }

    if CommandModifies(&cmd.Command) {
        if !CheckAuth(pay) {
            return nil,errors.New("cannot authenticate with key "+pay.Secret)
        }
    }

    if pay.Id == "" {
        return nil,errors.New("missing key id")
    }

    numArgs := len(pay.Values)+1

    args := make([]interface{},0,numArgs)
    strKey := storage.CreateKey(pay.Domain,pay.Id)
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

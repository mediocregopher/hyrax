package dispatch

import (
    "hyrax/storage"
    "hyrax/types"
    "hyrax/parse"
    "hyrax/custom"
    "hyrax/router"
    "errors"
)

// DoCommand takes in raw bytes that presumably have json data, decodes them,
// and performs whatever commands are needed based on them. It then returns
// raw bytes that contain the json return, either a return message or an error.
// The second return value, error, is only used if encoding the return value fails
// for some reason. Any actual errors are returned json encoded in the first return
// parameter.
func DoCommand(cid types.ConnId, rawJson []byte) ([]byte,error) {
    cmd,err := parse.DecodeCommand(rawJson)
    if err != nil {
        return parse.EncodeError("",err.Error())
    }

    ret,err := doCommandWrap(cid,cmd)
    if err != nil {
        return parse.EncodeError(cmd.Command,err.Error())
    }

    return parse.EncodeMessage(cmd.Command,ret)
}

func doCommandWrap(cid types.ConnId, cmd *types.Command) (interface{},error) {
    pay := &cmd.Payload
    cinfo,cexists := GetCommandInfo(&cmd.Command)

    if cexists {
        return nil,errors.New("Unsupported command")
    }

    if cinfo.Modifies {
        if !CheckAuth(pay) {
            return nil,errors.New("cannot authenticate with key "+pay.Secret)
        }
        if !cinfo.IsQuiet {
            custom.MonMakeAlert(cmd)
        }
    }

    if pay.Id == "" {
        return nil,errors.New("missing key id")
    }

    if cinfo.IsQuiet {
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

    if cinfo.ReturnsMap {
        return storage.RedisListToMap(r.([]string))
    }

    return r,nil
}

// doCustomCommand dispatches commands that don't go directly to redis, and instead
// are handled elsewhere
func doCustomCommand(cid types.ConnId, cmd *types.Command) (interface{},error) {
    f,ok := customCommandMap[cmd.Command]
    if !ok { return nil,errors.New("Command in main map not listed in custom map") }

    return f(cid,&cmd.Payload)
}

// DoCleanup takes in a connection id which is now defunct and cleans up any data it
// may have accumulated during its life (entry in router map, monitors, etc...)
func DoCleanup(cid types.ConnId) error {
    router.CleanId(cid)
    err := custom.CleanConnMon(cid)
    if err != nil { return err }
    return custom.CleanConnEkg(cid)
}

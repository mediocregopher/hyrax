package dispatch

import (
    types  "hyrax/types"
    stypes "hyrax-server/types"
    "hyrax-server/storage"
    "hyrax-server/custom"
    "hyrax-server/router"
    "errors"
    "fmt"
)

func errBytes (e error) []byte {
    return []byte(e.Error())
}

// DoCommand takes in raw bytes that presumably have json data, decodes them,
// and performs whatever commands are needed based on them. It then returns
// raw bytes that contain the json return, either a return message or an error.
// The second return value, error, is only used if encoding the return value fails
// for some reason. Any actual errors are returned json encoded in the first return
// parameter.
func DoCommand(cid stypes.ConnId, rawJson []byte) ([]byte,error) {

    if rawJson[0] == '{' {
        cmd,err := types.DecodeCommand(rawJson)
        if err != nil {
            return types.EncodeError(nil,err)
        }

        ret,err := doCommandWrap(cid,cmd)
        if err != nil {
            return types.EncodeError(cmd.Command,err)
        }

        return types.EncodeMessage(cmd.Command,ret)
    } else if rawJson[0] == '[' {
        cmds,err := types.DecodeCommandPackage(rawJson)
        if err != nil {
            return types.EncodeError(nil,err)
        }

        rets := make([][]byte,len(cmds))
        for i := range cmds {
            ret,err := doCommandWrap(cid,cmds[i])
            if err != nil {
                rets[i],_ = types.EncodeError(cmds[i].Command,err)
            } else {
                rets[i],_ = types.EncodeMessage(cmds[i].Command,ret)
            }
        }

        return types.EncodeMessagePackage(rets)
    }

    return types.EncodeError(nil,errors.New("Unknown command format"))
}

func doCommandWrap(cid stypes.ConnId, cmd *types.Command) (interface{},error) {
    pay := &cmd.Payload
    cinfo,cexists := GetCommandInfo(cmd.Command)

    if !cexists {
        return nil,fmt.Errorf("Unsupported command: %s",cmd.Command)
    }

    if cinfo.Modifies {
        if !CheckAuth(pay) {
            return nil,errors.New("cannot authenticate with key "+string(pay.Secret))
        }
        if !cmd.Quiet {
            custom.MonMakeAlert(cmd)
        }
    }

    if len(pay.Id) == 0 {
        return nil,errors.New("missing key id")
    }

    if cinfo.IsCustom {
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
        return storage.RedisListToMap(r.([][]byte))
    }

    return r,nil
}

// doCustomCommand dispatches commands that don't go directly to redis, and instead
// are handled elsewhere
func doCustomCommand(cid stypes.ConnId, cmd *types.Command) (interface{},error) {
    f,ok := GetCustomCommandFunc(cmd.Command)
    if !ok { return nil,errors.New("Command in main map not listed in custom map") }

    return f(cid,&cmd.Payload)
}

// DoCleanup takes in a connection id which is now defunct and cleans up any data it
// may have accumulated during its life (entry in router map, monitors, etc...)
func DoCleanup(cid stypes.ConnId) error {
    router.CleanId(cid)
    err := custom.CleanConnMon(cid)
    if err != nil { return err }
    return custom.CleanConnEkg(cid)
}

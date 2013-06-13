package config

import (
    "flag"
    "strconv"
    "fmt"
    "os"
)

var configInt = map[string]int{}
var configStr = map[string]string{}

func GetInt(name string) int {
    val,ok := configInt[name]
    if !ok { panic("attempted to access non-int-parameter "+name) }
    return val
}

func GetStr(name string) string {
    val,ok := configStr[name]
    if !ok { panic("attempted to access non-str-parameter "+name) }
    return val
}

func LoadConfig() {

    //Load cli into its own set of config maps
    cliConfigInt := map[string]*int{}
    cliConfigStr := map[string]*string{}
    for name,param := range params {
        if param.Type == INT {
            dummy := 0
            cliConfigInt[name] = &dummy
            flag.IntVar(cliConfigInt[name],name,param.Default.(int),param.Description)
        } else {
            dummy := ""
            cliConfigStr[name] = &dummy
            flag.StringVar(cliConfigStr[name],name,param.Default.(string),param.Description)
        }
    }

    //Some extra cli args
    dumpExample := flag.Bool("example",false,"Dump example configuration to stdout and exit")
    configFile  := flag.String("config","","Configuration file to load, empty means don't load any file and only use command-line args")

    flag.Parse()

    //If the flag to dump example config is set to true, do that
    if *dumpExample {
        fmt.Print(dumpExampleConfig())
        os.Exit(0)
    }

    //If config file is specified, load the string map from it and load the values into
    //global config
    if *configFile != "" {
        configFileMap,err := readConfig(*configFile)
        if err != nil { panic(err) }

        for name,val := range configFileMap {
            if param,ok := params[name]; ok {
                if param.Type == INT {
                    valint,err := strconv.Atoi(val)
                    if err != nil {
                        panic("field "+name+" in "+*configFile+" cannot be read as a number")
                    }
                    configInt[name] = valint
                } else {
                    configStr[name] = val
                }
            }
        }
    }

    //Now we look through each param. If it's set on the command-line (not to the default)
    //we set that in the global config maps. If it's also not set in the conf we set it
    //to the param's default. If it is set in the conf then it's already set in the global
    //configs by the previous section
    for name,param := range params {
        if param.Type == INT {
            cliVal := *cliConfigInt[name]
            _,confSet := configInt[name]
            if cliVal != param.Default {
                configInt[name] = cliVal
            } else if !confSet {
                configInt[name] = param.Default.(int)
            }
        } else {
            cliVal := *cliConfigStr[name]
            _,confSet := configStr[name]
            if cliVal != param.Default {
                configStr[name] = cliVal
            } else if !confSet {
                configStr[name] = param.Default.(string)
            }
        }
    }

}

package main

import (
	"fmt"
	"github.com/mediocregopher/flagconfig"
	"github.com/mediocregopher/hyrax/client"
	"github.com/mediocregopher/hyrax/types"
)

func printError(err error) {
	fmt.Println("ERR:", err)
}

func main() {
	fc := flagconfig.New("hyrax-cli")
	fc.DisallowConfig()
	fc.StrParam("format", "protocol format to use", "json")
	fc.StrParam("conn-type", "connection type to use", "tcp")
	fc.StrParam("addr", "address or socket location to connect to", "127.0.0.1:2379")
	fc.FlagParam("hold", "hold onto the connection even after the result has been returned, and output push messages as they come in", false)
	fc.RequiredStrParam("cmd", "cmd to execute")
	fc.StrParam("key", "key the command is executed against, if any", "")
	fc.StrParams("arg", "argument to command")
	fc.StrParam("id", "id of the client issuing command, if any", "")
	fc.StrParam("secret-key", "", "secret key used to construct hmac and validate command")

	if err := fc.Parse(); err != nil {
		fmt.Println(err)
		return
	}

	format := fc.GetStr("format")
	conntype := fc.GetStr("conn-type")
	addr := fc.GetStr("addr")
	hold := fc.GetFlag("hold")
	var pushCh chan *types.ClientCommand
	if hold {
		pushCh = make(chan *types.ClientCommand)
	}

	la := types.NewListenAddr(conntype, format, addr)
	c, err := client.NewClient(la, pushCh)
	if err != nil {
		printError(err)
		return
	}

	cmd := []byte(fc.GetStr("cmd"))
	keyB := []byte(fc.GetStr("key"))
	id := []byte(fc.GetStr("id"))
	secretKey := []byte(fc.GetStr("secret-key"))

	argsStrs := fc.GetStrs("arg")
	args := make([][]byte, len(argsStrs))
	for i := range argsStrs {
		args[i] = []byte(argsStrs[i])
	}

	ccmd := client.CreateClientCommand(cmd, keyB, id, secretKey, args...)

	ret, err := c.Cmd(ccmd)
	if err != nil {
		printError(err)
		return
	}

	fmt.Println(ret)

	if hold {
		for pushed := range pushCh {
			fmt.Printf(
				"PUSH from '%s' %s %s %v\n",
				string(pushed.Id),
				string(pushed.Command),
				string(pushed.StorageKey),
				pushed.Args,
			)
		}
	}

	c.Close()
}

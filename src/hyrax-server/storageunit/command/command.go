package command

// CommandRet is returned from a Command in the RetCh. It's really just a tuple
// around the return value and an error
type CommandRet struct {
	Ret interface{}
	Err error
}

// Command is sent to a StorageUnitConn, and contains all data necessary to
// complete a call and return any data from it.
type Command struct {
	Cmd   []byte
	Args  []interface{}
	RetCh chan *CommandRet
}

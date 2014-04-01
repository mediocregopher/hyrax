package translate

import (
	"encoding/json"

	. "github.com/mediocregopher/hyrax/types"
)

// JsonTranslator can encode/decode all messages required by hyrax
// servers/clients to communicate, and implements to the Translator interface
type JsonTranslator struct{}

func (j *JsonTranslator) ToClientCommand(b []byte) (*ClientCommand, error) {
	c := &ClientCommand{}
	err := json.Unmarshal(b, c)
	return c, err
}

func (j *JsonTranslator) FromClientCommand(c *ClientCommand) ([]byte, error) {
	return json.Marshal(c)
}

func (j *JsonTranslator) ToClientReturn(b []byte) (*ClientReturn, error) {
	c := &ClientReturn{}
	err := json.Unmarshal(b, c)
	return c, err
}

func (j *JsonTranslator) FromClientReturn(c *ClientReturn) ([]byte, error) {
	return json.Marshal(c)
}

package translate

import (
	. "github.com/mediocregopher/hyrax/types"
	"github.com/mediocregopher/gojson"
)

// JsonTranslator can encode/decode all messages required by hyrax
// servers/clients to communicate, and implements to the Translator interface
type JsonTranslator struct{}

func (j *JsonTranslator) ToClientCommand(b []byte) (*ClientCommand, error) {
	c := &ClientCommand{}
	err := gojson.Unmarshal(b, c)
	return c, err
}

func (j *JsonTranslator) FromClientCommand(c *ClientCommand) ([]byte, error) {
	return gojson.Marshal(c)
}

func (j *JsonTranslator) ToClientReturn(b []byte) (*ClientReturn, error) {
	c := &ClientReturn{}
	err := gojson.Unmarshal(b, c)
	return c, err
}

func (j *JsonTranslator) FromClientReturn(c *ClientReturn) ([]byte, error) {
	return gojson.Marshal(c)
}

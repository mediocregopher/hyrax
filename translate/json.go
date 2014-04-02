package translate

import (
	"encoding/json"

	. "github.com/mediocregopher/hyrax/types"
)

// JsonTranslator can encode/decode all messages required by hyrax
// servers/clients to communicate, and implements to the Translator interface
type JsonTranslator struct{}

func (j *JsonTranslator) ToAction(b []byte) (*Action, error) {
	a := &Action{}
	err := json.Unmarshal(b, a)
	return a, err
}

func (j *JsonTranslator) FromAction(a *Action) ([]byte, error) {
	return json.Marshal(a)
}

func (j *JsonTranslator) ToActionReturn(b []byte) (*ActionReturn, error) {
	ar := &ActionReturn{}
	err := json.Unmarshal(b, ar)
	return ar, err
}

func (j *JsonTranslator) FromActionReturn(ar *ActionReturn) ([]byte, error) {
	return json.Marshal(ar)
}

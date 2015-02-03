package jsv

import (
	"encoding/json"
)

var (
	JsonCodec Codec = &jsonCodec{}
)

type Codec interface {
	Marshal(i interface{}) ([]byte, error)
	Unmarshal(b []byte, i interface{}) error
}

type jsonCodec struct{}

func (this *jsonCodec) Marshal(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}
func (this *jsonCodec) Unmarshal(b []byte, i interface{}) error {
	return json.Unmarshal(b, &i)
}

package network

import (
	"reflect"
	"consensus_layer/serializer"
	"fmt"
)

func marshalBinary(v interface{}) ([]byte, error) {
	s := serializer.NewSerializer()
	extension := func(v interface{}) error {
		switch t := v.(type) {
		case MessageType:
			return s.WriteBytes([]byte{byte(t)})
		default:
			rv := reflect.Indirect(reflect.ValueOf(v))
			return fmt.Errorf("wrong type: %s", rv.Type().String())
		}
	}
	s.Extension = extension
	err := s.Serialize(v)
	return s.Bytes(), err
}

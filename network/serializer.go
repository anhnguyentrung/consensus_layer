package network

import (
	"reflect"
	"consensus_layer/serializer"
	"fmt"
	"consensus_layer/crypto"
	"consensus_layer/blockchain"
)

func MarshalBinary(v interface{}) ([]byte, error) {
	s := serializer.NewSerializer()
	extension := func(v interface{}) error {
		switch t := v.(type) {
		case MessageType:
			return s.WriteBytes([]byte{byte(t)})
		case blockchain.SHA256Type:
			return s.WriteBytes(t[:])
		case crypto.Signature:
			if len(t.Data) != 65 {
				return fmt.Errorf("length of signature data is not 65 bytes")
			}
			return s.WriteBytes(t.Data)
		case crypto.PublicKey:
			if len(t.Data) != 33 {
				return fmt.Errorf("length of public key is not 33 bytes")
			}
			return s.WriteBytes(t.Data)
		default:
			rv := reflect.Indirect(reflect.ValueOf(v))
			return fmt.Errorf("wrong type: %s", rv.Type().String())
		}
	}
	s.Extension = extension
	err := s.Serialize(v)
	return s.Bytes(), err
}

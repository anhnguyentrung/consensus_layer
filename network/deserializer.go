package network

import (
	"reflect"
	"consensus_layer/serializer"
	"fmt"
	"bufio"
	"io"
	"encoding/binary"
	"consensus_layer/crypto"
)

const SHA256TypeSize = 32
const PublicKeySize  = 33
const SignatureSize  = 65

func unmarshalBinary(buf []byte, v interface{}) error {
	d := serializer.NewDeserializer(buf)
	extension := func(v interface{}) error {
		rv := reflect.Indirect(reflect.ValueOf(v))
		switch v.(type) {
		case *MessageType:
			bytes, err := d.ReadBytes(1)
			if err != nil {
				return err
			}
			rv.SetUint(uint64(bytes[0]))
			return nil
		case *SHA256Type:
			bytes, err := d.ReadBytes(SHA256TypeSize)
			if err != nil {
				return err
			}
			sha256 := SHA256Type{}
			copy(sha256[:], bytes)
			rv.Set(reflect.ValueOf(sha256))
			return nil
		case *crypto.PublicKey:
			bytes, err := d.ReadBytes(PublicKeySize)
			if err != nil {
				return err
			}
			publicKey := crypto.PublicKey{ Data: bytes}
			rv.Set(reflect.ValueOf(publicKey))
			return nil
		case *crypto.Signature:
			bytes, err := d.ReadBytes(SignatureSize)
			if err != nil {
				return err
			}
			signature := crypto.Signature{ Data: bytes}
			rv.Set(reflect.ValueOf(signature))
			return nil
		default:
			rv := reflect.Indirect(reflect.ValueOf(v))
			return fmt.Errorf("wrong type: %s", rv.Type().String())
		}
	}
	d.Extension = extension
	return d.Deserialize(v)
}

func unmarshalBinaryMessage(reader *bufio.Reader, message *Message) error {
	typeBuf := make([]byte, 1, 1)
	_, err := io.ReadFull(reader, typeBuf)
	if err != nil {
		return err
	}
	messageType := MessageType(typeBuf[0])
	lenBuf := make([]byte, 4, 4)
	_, err = io.ReadFull(reader, lenBuf)
	if err != nil {
		return err
	}
	length := binary.BigEndian.Uint32(lenBuf)
	messageData := make([]byte, length, length)
	n, err := io.ReadFull(reader, messageData)
	if uint32(n) != length {
		return fmt.Errorf("wrong length")
	}
	message.Header.Type = messageType
	message.Header.Length = length
	message.Payload = messageData
	return nil
}

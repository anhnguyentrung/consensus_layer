package network

import (
	"reflect"
	"consensus_layer/serializer"
	"fmt"
	"bufio"
	"io"
	"encoding/binary"
)

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
	//fmt.Println("size of recevied message: ", length)
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

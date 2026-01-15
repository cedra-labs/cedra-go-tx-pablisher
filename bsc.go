package cedra

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/spf13/cast"
)

type BCSEncoder struct {
	buf bytes.Buffer
}

func NewBCSEncoder() *BCSEncoder {
	return &BCSEncoder{}
}

func (bcs *BCSEncoder) WriteRawBytes(value []byte) {
	bcs.buf.Write(value)
}

func (bcs *BCSEncoder) GetBytes() []byte {
	return bcs.buf.Bytes()
}

func (bcs *BCSEncoder) Reset() {
	bcs.buf.Reset()
}

func (bcs *BCSEncoder) EncodeEnum(value uint64) {
	for value >= 0x80 {
		bcs.buf.WriteByte(byte(value&0x7F) | 0x80)
		value >>= 7
	}

	bcs.buf.WriteByte(byte(value & 0x7F))
}

func (bcs *BCSEncoder) EncodeString(value string) {
	byteValue := []byte(value)
	length := len(byteValue)
	bcs.EncodeEnum(cast.ToUint64(length))
	bcs.buf.Write(byteValue)
}

func (bcs *BCSEncoder) EnncodeBytes(data []byte) {
	bcs.EncodeEnum(cast.ToUint64(len(data)))
	bcs.WriteRawBytes(data)
}

func EnncodeToBCSBytes(data []byte) []byte {
	bcs := NewBCSEncoder()
	defer bcs.buf.Reset()

	bcs.EncodeEnum(cast.ToUint64(len(data)))
	bcs.WriteRawBytes(data)
	buff := bcs.GetBytes()

	return buff
}

func EncodeToBCSString(value string) []byte {
	bcs := NewBCSEncoder()
	defer bcs.buf.Reset()

	byteValue := []byte(value)
	length := len(byteValue)
	bcs.EncodeEnum(cast.ToUint64(length))
	bcs.buf.Write(byteValue)

	return bcs.GetBytes()
}

func EncodeUintToBCS[Type uint8 | uint16 | uint32 | uint64](value Type) []byte {
	switch v := any(value).(type) {
	case uint8:
		return []byte{cast.ToUint8(v)}
	case uint16:
		buf := make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, v)

		return buf
	case uint32:
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, v)

		return buf
	case uint64:
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, v)

		return buf
	}

	panic(errors.New("EncodeUintToBCS: invalid received type"))
}

func EncodeIntToBCS[Type int8 | int16 | int32 | int64](value Type) []byte {
	switch v := any(value).(type) {
	case int8:
		return []byte{cast.ToUint8(v)}
	case int16:
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, cast.ToUint16(v))

		return buf
	case int32:
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, cast.ToUint32(v))

		return buf
	case int64:
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, cast.ToUint64(v))

		return buf
	}

	panic(errors.New("EncodeIntToBCS: invalid received type"))
}

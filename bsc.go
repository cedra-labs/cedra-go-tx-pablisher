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

func (bcs *BCSEncoder) EncodeHexBytes(value []byte) {

}

func (bcs *BCSEncoder) WriteRawBytes(value []byte) {
	bcs.buf.Write(value)
}

func (bcs *BCSEncoder) GetBytes() []byte {
	return bcs.buf.Bytes()
}

func (bcs *BCSEncoder) Uleb128(val uint8) {
	for val>>7 != 0 {
		bcs.buf.WriteByte((val & 0x7F) | 0x80)
		val >>= 7
	}

	bcs.buf.WriteByte(val & 0x7F)
}

func EncodeToBCSString(value string, enc *BCSEncoder) []byte {
	byteValue := []byte(value)
	length := len(byteValue)
	enc.Uleb128(cast.ToUint8(length))
	enc.buf.Write(byteValue)

	return enc.GetBytes()
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

type Uint interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}

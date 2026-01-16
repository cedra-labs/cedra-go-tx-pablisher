package cedra

import (
	"bytes"
	"encoding/binary"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// BCSEncoder provides Binary Canonical Serialization (BCS) encoding functionality.
// BCS is a deterministic serialization format used by the Cedra blockchain.
type BCSEncoder struct {
	buf bytes.Buffer
}

// NewBCSEncoder creates a new BCS encoder instance.
func NewBCSEncoder() *BCSEncoder {
	return &BCSEncoder{}
}

// WriteRawBytes appends raw bytes to the encoder buffer without any length prefix.
func (bcs *BCSEncoder) WriteRawBytes(value []byte) {
	bcs.buf.Write(value)
}

// GetBytes returns a copy of the encoded bytes from the buffer.
func (bcs *BCSEncoder) GetBytes() []byte {
	return bcs.buf.Bytes()
}

// Reset clears the encoder buffer, allowing it to be reused.
func (bcs *BCSEncoder) Reset() {
	bcs.buf.Reset()
}

// EncodeEnum encodes a uint64 value using variable-length encoding (ULEB128).
// This is used for encoding enum variants and length values.
func (bcs *BCSEncoder) EncodeEnum(value uint64) {
	for value >= 0x80 {
		bcs.buf.WriteByte(byte(value&0x7F) | 0x80)
		value >>= 7
	}

	bcs.buf.WriteByte(byte(value & 0x7F))
}

// EncodeString encodes a string value with its length prefix.
// The length is encoded as a ULEB128-encoded uint64, followed by the string bytes.
func (bcs *BCSEncoder) EncodeString(value string) {
	byteValue := []byte(value)
	length := len(byteValue)
	bcs.EncodeEnum(cast.ToUint64(length))
	bcs.buf.Write(byteValue)
}

// EncodeBytes encodes a byte slice with its length prefix.
// The length is encoded as a ULEB128-encoded uint64, followed by the bytes.
func (bcs *BCSEncoder) EncodeBytes(data []byte) {
	bcs.EncodeEnum(cast.ToUint64(len(data)))
	bcs.WriteRawBytes(data)
}

// GetMessageLen returns the current length of the encoded message in bytes.
func (bcs *BCSEncoder) GetMessageLen() int {
	return bcs.buf.Len()
}

// SetMessageLen prepends a message length prefix to the encoded buffer.
// This is useful for encoding messages that need a length header.
func (bcs *BCSEncoder) SetMessageLen(msgLen uint8) {
	newBuf := append(EncodeUintToBCS(msgLen), bcs.buf.Bytes()...)
	bcs.buf.Reset()
	bcs.WriteRawBytes(newBuf)
}

// EncodeToBCSBytes encodes a byte slice to BCS format with length prefix.
// This is a convenience function that creates a new encoder, encodes the data, and returns the result.
func EncodeToBCSBytes(data []byte) []byte {
	bcs := NewBCSEncoder()
	bcs.EncodeEnum(cast.ToUint64(len(data)))
	bcs.WriteRawBytes(data)

	return bcs.GetBytes()
}

// EncodeToBCSString encodes a string to BCS format with length prefix.
// This is a convenience function that creates a new encoder, encodes the string, and returns the result.
func EncodeToBCSString(value string) []byte {
	bcs := NewBCSEncoder()
	byteValue := []byte(value)
	bcs.EncodeEnum(cast.ToUint64(len(byteValue)))
	bcs.buf.Write(byteValue)

	return bcs.GetBytes()
}

// EncodeUintToBCS encodes an unsigned integer to BCS format using little-endian byte order.
// Supports uint8, uint16, uint32, and uint64 types.
// Panics if an unsupported type is provided.
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

// EncodeIntToBCS encodes a signed integer to BCS format using big-endian byte order.
// Supports int8, int16, int32, and int64 types.
// Panics if an unsupported type is provided.
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

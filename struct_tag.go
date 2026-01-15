package cedra

import (
	"encoding/hex"
	"errors"
	"strings"
)

const (
	structTagVariant   = 7
	structTagTyArgsLen = 0
	tagSeparator       = "::"
)

type StructTag struct {
	Address [32]byte
	Module  string
	Name    string
}

func NewStringStructTag(tag string) (StructTag, error) {
	parts := strings.Split(tag, tagSeparator)
	if len(parts) != 3 {
		return StructTag{}, errors.New("") // TODO:
	}
	bytes, err := hex.DecodeString(strings.TrimPrefix(parts[0], "0x"))
	if err != nil {
		return StructTag{}, err // TODO:
	}

	buf := [32]byte{}
	copy((buf)[32-len(bytes):], bytes)

	return StructTag{
		Address: buf,
		Module:  parts[1],
		Name:    parts[2],
	}, nil
}

func (st StructTag) ToBCSBytes() []byte {
	bcs := NewBCSEncoder()
	defer bcs.buf.Reset()
	// Encode FaFaAddress
	bcs.EncodeEnum(structTagVariant)
	bcs.WriteRawBytes(st.Address[:])
	bcs.EncodeString(st.Module)
	bcs.EncodeString(st.Name)
	bcs.EncodeEnum(structTagTyArgsLen)

	encodedTag := bcs.GetBytes()

	return encodedTag
}

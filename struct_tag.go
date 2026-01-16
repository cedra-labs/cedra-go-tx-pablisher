package cedra

import (
	"encoding/hex"
	"strings"

	"github.com/pkg/errors"
)

const (
	// structTagVariant is the variant identifier for struct tags.
	structTagVariant = 7
	// structTagTyArgsLen is the length of type arguments (currently always 0).
	structTagTyArgsLen = 0
	// tagSeparator is the separator used in struct tag strings (e.g., "address::module::name").
	tagSeparator = "::"
)

// StructTag represents a type identifier in the Cedra blockchain.
// It consists of an address, module name, and type name.
type StructTag struct {
	// Address is the 32-byte address of the module.
	Address [32]byte
	// Module is the name of the module.
	Module string
	// Name is the name of the type.
	Name string
}

// NewStringStructTag parses a struct tag string in the format "address::module::name"
// and creates a StructTag instance. The address can optionally include the "0x" prefix.
// Returns an error if the tag format is invalid.
func NewStringStructTag(tag string) (StructTag, error) {
	parts := strings.Split(tag, tagSeparator)
	if len(parts) != 3 {
		return StructTag{}, errors.New("can't create new struct tag: invalid struct tag")
	}
	bytes, err := hex.DecodeString(strings.TrimPrefix(parts[0], "0x"))
	if err != nil {
		return StructTag{}, errors.Wrap(err, "can't create new struct tag: invalid module address")
	}

	buf := [32]byte{}
	copy((buf)[32-len(bytes):], bytes)

	return StructTag{
		Address: buf,
		Module:  parts[1],
		Name:    parts[2],
	}, nil
}

// ToBCSBytes encodes the struct tag into Binary Canonical Serialization (BCS) format.
// Returns the serialized byte representation of the struct tag.
func (st *StructTag) ToBCSBytes() []byte {
	bcs := NewBCSEncoder()
	bcs.EncodeEnum(structTagVariant)
	bcs.WriteRawBytes(st.Address[:])
	bcs.EncodeString(st.Module)
	bcs.EncodeString(st.Name)
	bcs.EncodeEnum(structTagTyArgsLen)

	return bcs.GetBytes()
}

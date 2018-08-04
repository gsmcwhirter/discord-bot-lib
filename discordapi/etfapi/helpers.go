package etfapi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
)

func writeAtom(b io.Writer, val []byte) error {
	// assumes the Atom identifier byte has already been written
	size, err := IntToInt16Slice(len(val))
	if err != nil {
		return errors.Wrap(err, "couldn't marshal size")
	}

	_, err = b.Write(size)
	if err != nil {
		return errors.Wrap(err, "could not write size")
	}

	_, err = b.Write(val)
	if err != nil {
		return errors.Wrap(err, "could not write value")
	}

	return nil
}

func writeLength32(b io.Writer, n int) error {
	size, err := IntToInt32Slice(n)
	if err != nil {
		return errors.Wrap(err, "could not marshal length")
	}

	_, err = b.Write(size)
	return errors.Wrap(err, "could not write length")
}

func marshalInterface(code ETFCode, val interface{}) ([]byte, error) {
	// var data []byte
	var err error

	b := bytes.Buffer{}
	b.WriteByte(byte(code))

	switch code {
	case Map:
		v, ok := val.([]Element)
		if !ok {
			return nil, errors.Wrap(ErrBadMarshalData, "not a list of elements")
		}

		if len(v)%2 != 0 {
			return nil, errors.Wrap(ErrBadMarshalData, "bad parity on map list")
		}

		err = writeLength32(&b, len(v)/2)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't marshal map length")
		}

		for i := 0; i < len(v); i += 2 {

			if !v[i].Code.IsStringish() {
				return nil, errors.Wrap(ErrBadMarshalData, "bad map key")
			}

			_, err = v[i].WriteTo(&b)
			if err != nil {
				return nil, errors.Wrap(err, "couldn't marshal map key")
			}

			_, err = v[i+1].WriteTo(&b)
			if err != nil {
				return nil, errors.Wrap(err, "couldn't marshal map value")
			}
		}
	case List:
		v, ok := val.([]Element)
		if !ok {
			return nil, errors.Wrap(ErrBadMarshalData, "not a list of elements")
		}

		err = writeLength32(&b, len(v))
		if err != nil {
			return nil, errors.Wrap(err, "couldn't marshal list length")
		}

		for _, e := range v {
			_, err = e.WriteTo(&b)
			if err != nil {
				return nil, errors.Wrap(err, "couldn't marshal list value")
			}
		}

		err = b.WriteByte(byte(EmptyList))
		if err != nil {
			return nil, errors.Wrap(err, "couldn't write trailing list byte")
		}

	case Binary:
		v, ok := val.([]byte)
		if !ok {
			return nil, errors.Wrap(ErrBadMarshalData, "not a byte slice")
		}

		err = writeLength32(&b, len(v))
		if err != nil {
			return nil, errors.Wrap(err, "couldn't marshal binary length")
		}

		_, err = b.Write(v)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't marshal binary value")
		}
	case String:
		fallthrough
	case Atom:
		v, ok := val.([]byte)
		if !ok {
			return nil, errors.Wrap(ErrBadMarshalData, "not a byte slice")
		}

		err = writeAtom(&b, v)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't marshal Atom value")
		}
	case Int32:
		v, ok := val.([]byte)
		if !ok {
			return nil, errors.Wrap(ErrBadMarshalData, "not a byte slice")
		}

		if len(v) != 4 {
			return nil, errors.Wrap(ErrBadMarshalData, "not a int32 byte slice")
		}

		_, err = b.Write(v)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't marshal Int32 value")
		}

	case Int8:
		v, ok := val.([]byte)
		if !ok {
			return nil, errors.Wrap(ErrBadMarshalData, "not a byte slice")
		}

		if len(v) != 1 {
			return nil, errors.Wrap(ErrBadMarshalData, "not a int8 byte slice")
		}

		_, err = b.Write(v)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't marshal Int8 value")
		}

	default:
		return nil, ErrBadMarshalData
	}

	return b.Bytes(), nil
}

func unmarshalSlice(raw []byte, numElements int) (uint32, []Element, error) {
	var size int
	var idx uint32
	var deltaIdx uint32
	var err error

	e := make([]Element, numElements)

	for i := 0; i < numElements; i++ {
		e[i].Code = ETFCode(raw[idx])
		idx++
		switch e[i].Code {
		case Map:
			size, err = Int32SliceToInt(raw[idx : idx+4])
			if err != nil {
				return 0, nil, errors.Wrap(err, "could not read map length")
			}
			idx += 4

			deltaIdx, e[i].Vals, err = unmarshalSlice(raw[idx:], size*2)
			if err != nil {
				return 0, nil, errors.Wrap(err, "could not unmarshal map")
			}
			idx += deltaIdx
		case String:
			fallthrough
		case Atom:
			size, err = Int16SliceToInt(raw[idx : idx+2])
			if err != nil {
				return 0, nil, errors.Wrap(err, "could not read atom/string length")
			}
			idx += 2
			e[i].Val = raw[idx : idx+uint32(size)]
			idx += uint32(size)
		case List:
			size, err = Int32SliceToInt(raw[idx : idx+4])
			if err != nil {
				return 0, nil, errors.Wrap(err, "coult not read list length")
			}
			idx += 4
			deltaIdx, e[i].Vals, err = unmarshalSlice(raw[idx:], size)
			if err != nil {
				return 0, nil, err
			}
			idx += deltaIdx

			if raw[idx] != byte(EmptyList) {
				return 0, nil, ErrBadPayload
			}
			idx++
		case Binary:
			size, err = Int32SliceToInt(raw[idx : idx+4])
			if err != nil {
				return 0, nil, errors.Wrap(err, "could not read binary length")
			}
			idx += 4
			e[i].Val = raw[idx : idx+uint32(size)]
			idx += uint32(size)
		case Int32:
			e[i].Val = raw[idx : idx+4]
			idx += 4
		case Int8: // small int
			e[i].Val = raw[idx : idx+1]
			idx++
		case EmptyList:
		case SmallBig:
			size = int(raw[idx])
			idx++

			e[i].Val = raw[idx : idx+uint32(size)+1]
			idx += uint32(size) + 1
		// case LargeBig:
		// 	size, err = Int32SliceToInt(raw[idx : idx+4])
		// 	if err != nil {
		// 		return 0, nil, errors.Wrap(err, "could not read largebig length")
		// 	}
		// 	idx += 4
		// 	e[i].Val = raw[idx : idx+uint32(size)]
		// 	idx += uint32(size)
		default:
			return 0, nil, errors.Wrap(ErrBadFieldType, fmt.Sprintf("type=%v", e[i].Code))
		}
	}

	return idx, e, nil
}

// IntToInt8Slice TODOC
func IntToInt8Slice(v int) ([]byte, error) {
	if v < 0 || v > 255 {
		return nil, ErrOutOfBounds
	}

	return []byte{byte(v)}, nil
}

// Int8SliceToInt TODOC
func Int8SliceToInt(v []byte) (int, error) {
	if len(v) != 1 {
		return 0, ErrOutOfBounds
	}

	return int(v[0]), nil
}

// IntToInt16Slice TODOC
func IntToInt16Slice(v int) ([]byte, error) {
	if v < 0 || v >= (1<<16) {
		return nil, ErrOutOfBounds
	}

	size := make([]byte, 2)
	binary.BigEndian.PutUint16(size, uint16(v))

	return size, nil
}

// Int16SliceToInt TODOC
func Int16SliceToInt(v []byte) (int, error) {
	if len(v) != 2 {
		return 0, ErrOutOfBounds
	}

	return int(binary.BigEndian.Uint16(v)), nil
}

// IntToInt32Slice TODOC
func IntToInt32Slice(v int) ([]byte, error) {
	if v < 0 || v >= (1<<32) {
		return nil, ErrOutOfBounds
	}

	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(v))

	return size, nil
}

// Int32SliceToInt TODOC
func Int32SliceToInt(v []byte) (int, error) {
	if len(v) != 4 {
		return 0, ErrOutOfBounds
	}

	return int(binary.BigEndian.Uint32(v)), nil
}

// IntNSliceToInt64 TODOC
func IntNSliceToInt64(v []byte) (int64, error) {
	var newV []byte

	if len(v) > 8 {
		return 0, ErrOutOfBounds
	}

	if len(v) < 8 {
		newV = make([]byte, 8)
		copy(newV[8-len(v):], v)
	} else {
		newV = v
	}

	return int64(binary.LittleEndian.Uint64(newV)), nil
}

// ElementMapToElementSlice TODOC
func ElementMapToElementSlice(m map[string]Element) ([]Element, error) {
	e := make([]Element, 0, len(m)*2)
	for k, v := range m {
		el, err := NewBinaryElement([]byte(k))
		if err != nil {
			return nil, errors.Wrap(err, "could not create Element for key")
		}

		e = append(e, el)
		e = append(e, v)
	}

	return e, nil
}

// MapAndIDFromElement TODOC
func MapAndIDFromElement(e Element) (eMap map[string]Element, id snowflake.Snowflake, err error) {
	eMap, err = e.ToMap()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("could not inflate element to map: %v", e))
		return
	}

	id, err = SnowflakeFromElement(eMap["id"])
	err = errors.Wrap(err, "could not get id snowflake.Snowflake")
	return
}

// SnowflakeFromElement TODOC
func SnowflakeFromElement(e Element) (s snowflake.Snowflake, err error) {
	temp, err := e.ToInt64()
	if err != nil {
		err = errors.Wrap(err, "could not unmarshal snowflake.Snowflake")
	}
	s = snowflake.Snowflake(temp)
	return
}

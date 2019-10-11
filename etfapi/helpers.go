package etfapi

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/gsmcwhirter/go-util/v5/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v12/snowflake"
)

func writeLength16(b io.Writer, n int) error {
	// assumes the Atom identifier byte has already been written
	size, err := intToInt16Slice(n)
	if err != nil {
		return errors.Wrap(err, "couldn't marshal length")
	}

	_, err = b.Write(size)
	return errors.Wrap(err, "could not write length")
}

func writeLength32(b io.Writer, n int) error {
	size, err := intToInt32Slice(n)
	if err != nil {
		return errors.Wrap(err, "could not marshal length")
	}

	_, err = b.Write(size)
	return errors.Wrap(err, "could not write length")
}

func marshalMapTo(b io.Writer, v []Element) error {
	var err error
	if len(v)%2 != 0 {
		return errors.Wrap(ErrBadMarshalData, "bad parity on map list")
	}

	err = writeLength32(b, len(v)/2)
	if err != nil {
		return errors.Wrap(err, "couldn't marshal map length")
	}

	for i := 0; i < len(v); i += 2 {
		if !v[i].Code.IsStringish() {
			return errors.Wrap(ErrBadMarshalData, "bad map key")
		}

		err = v[i].MarshalTo(b)
		if err != nil {
			return errors.Wrap(err, "couldn't marshal map key")
		}

		err = v[i+1].MarshalTo(b)
		if err != nil {
			return errors.Wrap(err, "couldn't marshal map value")
		}
	}

	return nil
}

func marshalListTo(b io.Writer, v []Element) error {
	err := writeLength32(b, len(v))
	if err != nil {
		return errors.Wrap(err, "couldn't marshal list length")
	}

	for _, e := range v {
		err = e.MarshalTo(b)
		if err != nil {
			return errors.Wrap(err, "couldn't marshal list value")
		}
	}

	_, err = b.Write([]byte{byte(EmptyList)})
	return errors.Wrap(err, "couldn't write trailing list byte")
}

func marshalBinaryTo(b io.Writer, v []byte) error {
	err := writeLength32(b, len(v))
	if err != nil {
		return errors.Wrap(err, "couldn't marshal binary length")
	}

	_, err = b.Write(v)
	return errors.Wrap(err, "couldn't marshal binary value")
}

// for Atom, String
func marshalStringTo(b io.Writer, v []byte) error {
	err := writeLength16(b, len(v))
	if err != nil {
		return errors.Wrap(err, "couldn't marshal string length")
	}

	_, err = b.Write(v)
	return errors.Wrap(err, "couldn't marshal string value")
}

// for SmallBig, LargeBig
func marshalInt64To(b io.Writer, v []byte) error {
	var err error

	if len(v) != 9 {
		return errors.Wrap(ErrBadMarshalData, "not a int64 byte slice")
	}

	_, err = b.Write([]byte{byte(len(v) - 1)})
	if err != nil {
		return errors.Wrap(err, "couldn't marshal Int64 size")
	}

	_, err = b.Write(v)
	return errors.Wrap(err, "couldn't marshal Int64 value")
}

func marshalInt32To(b io.Writer, v []byte) error {
	var err error

	if len(v) != 4 {
		return errors.Wrap(ErrBadMarshalData, "not a int32 byte slice")
	}

	_, err = b.Write(v)
	return errors.Wrap(err, "couldn't marshal Int32 value")
}

func marshalInt8To(b io.Writer, v []byte) error {
	var err error

	if len(v) != 1 {
		return errors.Wrap(ErrBadMarshalData, "not a int8 byte slice")
	}

	_, err = b.Write(v)
	return errors.Wrap(err, "couldn't marshal Int8 value")
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
			size, err = int32SliceToInt(raw[idx : idx+4])
			if err != nil {
				return 0, nil, errors.Wrap(err, "could not read map length")
			}
			idx += 4

			deltaIdx, e[i].Vals, err = unmarshalSlice(raw[idx:], size*2)
			if err != nil {
				return 0, nil, errors.Wrap(err, "could not unmarshal map")
			}
			idx += deltaIdx
		case Atom, String:
			size, err = int16SliceToInt(raw[idx : idx+2])
			if err != nil {
				return 0, nil, errors.Wrap(err, "could not read atom/string length")
			}
			idx += 2
			e[i].Val = raw[idx : idx+uint32(size)]
			idx += uint32(size)
		case List:
			size, err = int32SliceToInt(raw[idx : idx+4])
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
			size, err = int32SliceToInt(raw[idx : idx+4])
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

func intToInt8Slice(v int) ([]byte, error) {
	if v < 0 || v > 255 {
		return nil, ErrOutOfBounds
	}

	return []byte{byte(v)}, nil
}

func int8SliceToInt(v []byte) (int, error) {
	if len(v) != 1 {
		return 0, ErrOutOfBounds
	}

	return int(v[0]), nil
}

func intToInt16Slice(v int) ([]byte, error) {
	if v < 0 || v >= (1<<16) {
		return nil, ErrOutOfBounds
	}

	size := make([]byte, 2)
	binary.BigEndian.PutUint16(size, uint16(v))

	return size, nil
}

func int16SliceToInt(v []byte) (int, error) {
	if len(v) != 2 {
		return 0, ErrOutOfBounds
	}

	return int(binary.BigEndian.Uint16(v)), nil
}

func intToInt32Slice(v int) ([]byte, error) {
	if v < 0 || v >= (1<<32) {
		return nil, ErrOutOfBounds
	}

	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(v))

	return size, nil
}

func int64ToInt64Slice(v int64) ([]byte, error) {

	data := make([]byte, 9)
	if v < 0 {
		v = -v
		data[0] = 1
	}
	binary.LittleEndian.PutUint64(data[1:], uint64(v))

	return data, nil
}

func int32SliceToInt(v []byte) (int, error) {
	if len(v) != 4 {
		return 0, ErrOutOfBounds
	}

	return int(binary.BigEndian.Uint32(v)), nil
}

func intNSliceToInt64(v []byte) (int64, error) {
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

// ElementMapToElementSlice converts an string->Element map into a slice of Elements (kv pairs)
func ElementMapToElementSlice(m map[string]Element) ([]Element, error) {
	e := make([]Element, 0, len(m)*2)
	for k, v := range m {
		el, err := NewBinaryElement([]byte(k))
		if err != nil {
			return nil, errors.Wrap(err, "could not create Element for key")
		}

		e = append(e, el, v)
	}

	return e, nil
}

// MapAndIDFromElement converts a Map element into a string->Element map and attempts to extract
// an id Snowflake from the "id" field
func MapAndIDFromElement(e Element) (map[string]Element, snowflake.Snowflake, error) {
	eMap, err := e.ToMap()
	if err != nil {
		return eMap, 0, errors.Wrap(err, fmt.Sprintf("could not inflate element to map: %v", e))
	}

	id, err := SnowflakeFromElement(eMap["id"])
	return eMap, id, errors.Wrap(err, "could not get id snowflake.Snowflake")
}

// SnowflakeFromElement converts a number-like Element into a Snowflake
func SnowflakeFromElement(e Element) (snowflake.Snowflake, error) {
	temp, err := e.ToInt64()
	s := snowflake.Snowflake(temp)
	return s, errors.Wrap(err, "could not unmarshal snowflake.Snowflake")
}

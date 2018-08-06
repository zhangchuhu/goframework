package gojce

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	_ "io"
	"reflect"
)

type jce_input_stream struct {
	//bs     []byte
	hd     head_data
	reader *bytes.Reader
}

func (j *jce_input_stream) readHead(hd *head_data) (int, error) {
	if j.reader.Len() <= 0 {
		return 0, fmt.Errorf("read file to end")
	}

	b, err := j.reader.ReadByte()
	if err != nil {
		return 0, err
	}

	hd.ty = byte((b & 0xF))
	hd.tag = (int)((b & (0xF << 4)) >> 4)
	if hd.tag == 15 {
		b, err = j.reader.ReadByte()
		if err != nil {
			return 0, err
		}
		hd.tag = int(b)
		return 2, nil
	}
	return 1, nil
}

func (j *jce_input_stream) peekHead(hd *head_data) (int, error) {
	l, err := j.readHead(hd)
	if err != nil {
		return 0, err
	}

	for i := 0; i < l; i++ {
		j.reader.UnreadByte()
	}
	return l, nil
}

func (j *jce_input_stream) skip(n int) error {
	_, err := j.reader.Seek(int64(n), 1)
	return err
}

func (j *jce_input_stream) skipToTag(tag int) (bool, error) {
	for {
		l, err := j.peekHead(&j.hd)
		if err != nil {
			return false, err
		}

		if tag <= j.hd.tag || j.hd.ty == STRUCT_END {
			return tag == j.hd.tag, nil
		}
		if err = j.skip(l); err != nil {
			return false, err
		}

		if err = j.skipFieldFrom(j.hd.ty); err != nil {
			return false, err
		}
	}
	return false, nil
}

func (j *jce_input_stream) skipField() error {
	_, err := j.readHead(&j.hd)
	if err != nil {
		return err
	}

	err = j.skipFieldFrom(j.hd.ty)
	if err != nil {
		return err
	}

	return nil
}

func (j *jce_input_stream) skipToStructEnd() error {
	var hd head_data // fix: skipFieldFrom() changes hd.ty, might break the loop earlier, using another head_data
	for hd.ty != STRUCT_END {
		_, err := j.readHead(&hd)
		if err != nil {
			return err
		}

		err = j.skipFieldFrom(hd.ty)
		if err != nil {
			return err
		}
	}
	return nil
}

func (j *jce_input_stream) skipFieldFrom(ty byte) error {
	switch ty {
	case BYTE:
		j.skip(1)
		break
	case SHORT:
		j.skip(2)
		break
	case INT:
		j.skip(4)
		break
	case LONG:
		j.skip(8)
		break
	case FLOAT:
		j.skip(4)
		break
	case DOUBLE:
		j.skip(8)
		break
	case STRING1:
		b, err := j.reader.ReadByte()
		if err != nil {
			return err
		}

		l := int(b)
		if l < 0 {
			l += 256
		}
		j.skip(l)
		break
	case STRING4:
		var l int32
		err := binary.Read(j.reader, binary.BigEndian, &l) //get from net endian
		if err != nil {
			return err
		}
		j.skip(int(l))
		break
	case MAP:
		sizeinf, err := j.readInt(0, true)
		if err != nil {
			return err
		} else if sizeinf == nil {
			return fmt.Errorf("read int nil")
		}
		size := sizeinf.(int)
		for i := 0; i < size*2; i++ {
			j.skipField()
		}
		break
	case LIST:
		sizeinf, err := j.readInt(0, true)
		if err != nil {
			return err
		} else if sizeinf == nil {
			return fmt.Errorf("read int nil")
		}
		size := sizeinf.(int)
		for i := 0; i < size; i++ {
			j.skipField()
		}
		break
	case SIMPLE_LIST:
		j.readHead(&j.hd)
		if j.hd.ty != BYTE {
			return fmt.Errorf("skipField with invalid type, type value: %d, hd.type: %d", ty, j.hd.ty)
		}
		sizeinf, err := j.readInt(0, true)
		if err != nil {
			return err
		} else if sizeinf == nil {
			return fmt.Errorf("read int nil")
		}
		size := sizeinf.(int)
		j.skip(size)
		break
	case STRUCT_BEGIN:
		err := j.skipToStructEnd()
		if err != nil {
			return err
		}
		break
	case STRUCT_END:
		break
	case ZERO_TAG:
		break
	default:
		return fmt.Errorf("invalid type.")
	}
	return nil
}

// func (j *jce_input_stream) ReadBool(tag int, required bool) (bool, error) {
// 	c, err := j.ReadByte(tag, required)
// 	return c != 0, err
// }

// func (j *jce_input_stream) ReadByte(tag int, required bool) (byte, error) {
// 	var c byte
// 	if ok, err := j.skipToTag(tag); ok {
// 		hd := new(head_data)
// 		j.readHead(hd)
// 		switch hd.ty {
// 		case ZERO_TAG:
// 			c = 0
// 			break
// 		case BYTE:
// 			c, err = j.reader.ReadByte()
// 			if err != nil {
// 				return c, err
// 			}
// 			break
// 		default:
// 			return c, fmt.Errorf("type mismatch.")
// 		}
// 	} else if required {
// 		return c, fmt.Errorf("require field not exist.")
// 	}
// 	return c, nil
// }

func (j *jce_input_stream) readByte(tag int, required bool) (interface{}, error) {
	var n byte
	var err error
	if ok, _ := j.skipToTag(tag); ok {
		j.readHead(&j.hd)
		switch j.hd.ty {
		case ZERO_TAG:
			n = 0
			break
		case BYTE:
			n, err = j.reader.ReadByte()
			if err != nil {
				return n, err
			}
			break
		default:
			return n, fmt.Errorf("type mismatch.")
		}
		return n, nil
	} else if required {
		return nil, fmt.Errorf("require field not exist.")
	}
	return nil, nil
}

//read short
func (j *jce_input_stream) readInt16(tag int, required bool) (interface{}, error) {
	var n int16
	if ok, _ := j.skipToTag(tag); ok {
		j.readHead(&j.hd)
		switch j.hd.ty {
		case ZERO_TAG:
			n = 0
			break
		case BYTE:
			var c int8
			binary.Read(j.reader, binary.BigEndian, &c)
			n = int16(c)
			break
		case SHORT:
			binary.Read(j.reader, binary.BigEndian, &n)
			break
		default:
			return n, fmt.Errorf("type mismatch.")
		}
		return n, nil
	} else if required {
		return nil, fmt.Errorf("require field not exist.")
	}
	return nil, nil
}

//read ushort
func (j *jce_input_stream) readUint16(tag int, required bool) (interface{}, error) {
	var n uint16
	if ok, _ := j.skipToTag(tag); ok {
		j.readHead(&j.hd)
		switch j.hd.ty {
		case ZERO_TAG:
			n = 0
			break
		case BYTE:
			c, err := j.reader.ReadByte()
			if err != nil {
				return nil, err
			}
			n = uint16(c)
			break
		case SHORT:
			binary.Read(j.reader, binary.BigEndian, &n)
			break
		default:
			return nil, fmt.Errorf("type mismatch.")
		}
		return n, nil
	} else if required {
		return nil, fmt.Errorf("require field not exist.")
	}
	return nil, nil
}

func (j *jce_input_stream) readInt(tag int, required bool) (interface{}, error) {
	var n int
	if ok, _ := j.skipToTag(tag); ok {
		j.readHead(&j.hd)
		switch j.hd.ty {
		case ZERO_TAG:
			n = 0
			break
		case BYTE:
			var c int8
			binary.Read(j.reader, binary.BigEndian, &c)
			n = int(c)
			break
		case SHORT:
			var c int16
			binary.Read(j.reader, binary.BigEndian, &c)
			n = int(c)
			break
		case INT:
			var c int32
			binary.Read(j.reader, binary.BigEndian, &c)
			n = int(c)
			break
		case LONG:
			var c int64
			binary.Read(j.reader, binary.BigEndian, &c)
			n = int(c)
			break
		default:
			return nil, fmt.Errorf("type mismatch.")
		}
		return n, nil
	} else if required {
		return nil, fmt.Errorf("require field not exist.")
	}
	return nil, nil
}

func (j *jce_input_stream) readUint(tag int, required bool) (interface{}, error) {
	n, err := j.readInt64(tag, required)
	if n != nil {
		return uint(n.(int64)), err
	}
	return n, err
}

func (j *jce_input_stream) readUint32(tag int, required bool) (interface{}, error) {
	n, err := j.readUint64(tag, required)
	if n != nil {
		return uint32(n.(uint64)), err
	}
	return n, err
}

func (j *jce_input_stream) readInt32(tag int, required bool) (interface{}, error) {
	n, err := j.readInt64(tag, required)
	if n != nil {
		return int32(n.(int64)), err
	}
	return n, err
}

//read long
func (j *jce_input_stream) readInt64(tag int, required bool) (interface{}, error) {
	var n int64
	if ok, _ := j.skipToTag(tag); ok {
		j.readHead(&j.hd)
		switch j.hd.ty {
		case ZERO_TAG:
			n = 0
			break
		case BYTE:
			var c int8
			binary.Read(j.reader, binary.BigEndian, &c)
			n = int64(c)
			break
		case SHORT:
			var c int16
			binary.Read(j.reader, binary.BigEndian, &c)
			n = int64(c)
			break
		case INT:
			var c int32
			binary.Read(j.reader, binary.BigEndian, &c)
			n = int64(c)
			break
		case LONG:
			binary.Read(j.reader, binary.BigEndian, &n)
			break
		default:
			return nil, fmt.Errorf("type mismatch.")
		}
		return n, nil
	} else if required {
		return nil, fmt.Errorf("require field not exist.")
	}
	return nil, nil
}

//read ulong
func (j *jce_input_stream) readUint64(tag int, required bool) (interface{}, error) {
	var n uint64
	if ok, _ := j.skipToTag(tag); ok {
		j.readHead(&j.hd)
		switch j.hd.ty {
		case ZERO_TAG:
			n = 0
			break
		case BYTE:
			c, err := j.reader.ReadByte()
			if err != nil {
				return nil, err
			}
			n = uint64(c)
			break
		case SHORT:
			var c uint16
			binary.Read(j.reader, binary.BigEndian, &c)
			n = uint64(c)
			break
		case INT:
			var c uint32
			binary.Read(j.reader, binary.BigEndian, &c)
			n = uint64(c)
			break
		case LONG:
			binary.Read(j.reader, binary.BigEndian, &n)
			break
		default:
			return nil, fmt.Errorf("type mismatch.")
		}
		return n, nil
	} else if required {
		return nil, fmt.Errorf("require field not exist.")
	}
	return nil, nil
}

func (j *jce_input_stream) readFloat32(tag int, required bool) (interface{}, error) {
	var n float32
	if ok, _ := j.skipToTag(tag); ok {
		j.readHead(&j.hd)
		switch j.hd.ty {
		case ZERO_TAG:
			n = 0
			break
		case FLOAT:
			binary.Read(j.reader, binary.BigEndian, &n)
			break
		default:
			return n, fmt.Errorf("type mismatch.")
		}
		return n, nil
	} else if required {
		return nil, fmt.Errorf("require field not exist.")
	}
	return nil, nil
}

func (j *jce_input_stream) readFloat64(tag int, required bool) (interface{}, error) {
	var n float64
	if ok, _ := j.skipToTag(tag); ok {
		j.readHead(&j.hd)
		switch j.hd.ty {
		case ZERO_TAG:
			n = 0
			break
		case FLOAT:
			var c float32
			binary.Read(j.reader, binary.BigEndian, &c)
			n = float64(c)
			break
		case DOUBLE:
			binary.Read(j.reader, binary.BigEndian, &n)
			break
		default:
			return nil, fmt.Errorf("type mismatch.")
		}
		return n, nil
	} else if required {
		return nil, fmt.Errorf("require field not exist.")
	}
	return nil, nil
}

func (j *jce_input_stream) readByteString(tag int, required bool) (interface{}, error) {
	var n string
	if ok, _ := j.skipToTag(tag); ok {
		j.readHead(&j.hd)
		switch j.hd.ty {
		case STRING1:
			b, err := j.reader.ReadByte()
			if err != nil {
				return nil, err
			}

			l := int(b)
			if l < 0 {
				l += 256
			}

			str := make([]byte, l)
			j.reader.Read(str)
			n = hex.EncodeToString(str)
			break
		case STRING4:
			var l int32
			binary.Read(j.reader, binary.BigEndian, &l)

			if l > JCE_MAX_STRING_LENGTH || l < 0 {
				return nil, fmt.Errorf("string too long: %s", l)
			}

			str := make([]byte, l)
			j.reader.Read(str)
			n = hex.EncodeToString(str)
			break
		default:
			return n, fmt.Errorf("type mismatch.")
		}
		return n, nil
	} else if required {
		return nil, fmt.Errorf("require field not exist.")
	}
	return nil, nil
}

func (j *jce_input_stream) readString(tag int, required bool) (interface{}, error) {
	var n string
	if ok, _ := j.skipToTag(tag); ok {
		j.readHead(&j.hd)
		switch j.hd.ty {
		case STRING1:
			b, err := j.reader.ReadByte()
			if err != nil {
				return nil, err
			}

			l := int(b)
			if l < 0 {
				l += 256
			}

			str := make([]byte, l)
			j.reader.Read(str)
			n = string(str)
			break
		case STRING4:
			var l int32
			binary.Read(j.reader, binary.BigEndian, &l)

			if l > JCE_MAX_STRING_LENGTH || l < 0 {
				return nil, fmt.Errorf("string too long: %s", l)
			}

			str := make([]byte, l)
			j.reader.Read(str)
			n = string(str)
			break
		default:
			return nil, fmt.Errorf("type mismatch.")
		}
		return n, nil
	} else if required {
		return nil, fmt.Errorf("require field not exist.")
	}
	return nil, nil
}

func (j *jce_input_stream) readJceStruct(ty reflect.Type, tag int, required bool) (interface{}, error) {
	var obj reflect.Value
	//var inf JceStruct
	modelType := reflect.TypeOf((*JceStruct)(nil)).Elem()
	obj = reflect.New(ty)
	if !obj.Type().Implements(modelType) {
		return nil, fmt.Errorf("wrong struct type")
	}

	if ok, _ := j.skipToTag(tag); ok {
		j.readHead(&j.hd)
		if j.hd.ty != STRUCT_BEGIN {
			return nil, fmt.Errorf("type mismatch.")
		}

		method := obj.MethodByName("ReadFrom")

		rs := method.Call([]reflect.Value{reflect.ValueOf(j)})[0]
		if !rs.IsNil() { //err != nil
			return nil, fmt.Errorf("call readfrom method error")
		}
		j.skipToStructEnd()
	} else if required {
		return nil, fmt.Errorf("require field not exist.")
	}
	return obj.Elem().Interface(), nil
}

func (j *jce_input_stream) readSlice(ty reflect.Type, tag int, required bool) (interface{}, error) {
	if ok, _ := j.skipToTag(tag); ok {
		if !(ty.Kind() == reflect.Array || ty.Kind() == reflect.Slice) {
			return nil, fmt.Errorf("read item is not array or slice: %v", ty.Kind())
		}

		elem_type := ty.Elem()

		j.readHead(&j.hd)
		switch j.hd.ty {
		case LIST:
			sizeinf, err := j.readInt(0, true)
			if err != nil {
				return nil, err
			} else if sizeinf == nil {
				return nil, fmt.Errorf("read int nil")
			}
			size := sizeinf.(int)
			if size < 0 {
				return nil, fmt.Errorf("size invalid: %d", size)
			}

			slice := reflect.MakeSlice(ty, size, size)
			for i := 0; i < size; i++ {
				v := slice.Index(i)
				val, err := j.Read(elem_type, 0, true)
				if err != nil {
					return nil, err
				}

				v.Set(reflect.ValueOf(val))
			}
			return slice.Interface(), nil
		case SIMPLE_LIST:
			hh := new(head_data)
			j.readHead(hh)

			if hh.ty == ZERO_TAG {
				slice := reflect.MakeSlice(elem_type, 0, 0)
				return slice.Interface(), nil
			}
			if hh.ty != BYTE {
				return nil, fmt.Errorf("type mismatch, tag: %d, type: %d, %d", tag, j.hd.ty, hh.ty)
			}

			sizeinf, err := j.readInt(0, true)
			if err != nil {
				return nil, err
			} else if sizeinf == nil {
				return nil, fmt.Errorf("read int nil")
			}
			size := sizeinf.(int)
			if size < 0 {
				return nil, fmt.Errorf("size invalid: tag: %d, type: %d, size: %d", tag, j.hd.ty, size)
			}

			bs := make([]byte, size)
			j.reader.Read(bs)
			slice := reflect.MakeSlice(ty, size, size)
			for i := 0; i < size; i++ {
				v := slice.Index(i)
				v.Set(reflect.ValueOf(bs[i]))
			}
			return slice.Interface(), nil
		default:
			return nil, fmt.Errorf("type mismatch.")
		}
	} else if required {
		return nil, fmt.Errorf("require field not exist.")
	}
	return nil, nil
}

func (j *jce_input_stream) readMap(ty reflect.Type, tag int, required bool) (interface{}, error) {
	if ok, _ := j.skipToTag(tag); ok {
		if ty.Kind() != reflect.Map {
			return nil, fmt.Errorf("read item is not map")
		}

		key_type := ty.Key()
		val_type := ty.Elem()
		m := reflect.MakeMap(ty)

		j.readHead(&j.hd)
		switch j.hd.ty {
		case MAP:
			sizeinf, err := j.readInt(0, true)
			if err != nil {
				return nil, err
			} else if sizeinf == nil {
				return nil, fmt.Errorf("read int nil")
			}
			size := sizeinf.(int)
			if size < 0 {
				return nil, fmt.Errorf("size invalid: %d", size)
			}

			for i := 0; i < size; i++ {
				mk, err := j.Read(key_type, 0, true)
				if err != nil {
					return nil, err
				}
				mv, err := j.Read(val_type, 1, true)
				if err != nil {
					return nil, err
				}

				m.SetMapIndex(reflect.ValueOf(mk), reflect.ValueOf(mv))
			}
			return m.Interface(), nil
		default:
			return nil, fmt.Errorf("type mismatch.")
		}
	} else if required {
		return nil, fmt.Errorf("require field not exist.")
	}
	return nil, nil
}

func (j *jce_input_stream) Read(ty reflect.Type, tag int, required bool) (interface{}, error) {
	switch ty.Kind() {
	case reflect.Bool:
		c, err := j.readByte(tag, required)
		if c != nil {
			return c.(byte) != 0, err
		} else {
			return nil, err
		}
	case reflect.Uint8: //byte
		return j.readByte(tag, required)
	case reflect.Uint16:
		return j.readUint16(tag, required)
	case reflect.Uint32:
		return j.readUint32(tag, required)
	case reflect.Uint64:
		return j.readUint64(tag, required)
	case reflect.Uint:
		return j.readUint(tag, required)
	case reflect.Int8: //TODO: check
		return j.readByte(tag, required)
	case reflect.Int16:
		return j.readInt16(tag, required)
	case reflect.Int32:
		return j.readInt32(tag, required)
	case reflect.Int64:
		return j.readInt64(tag, required)
	case reflect.Int:
		return j.readInt(tag, required)
	case reflect.Float32:
		return j.readFloat32(tag, required)
	case reflect.Float64:
		return j.readFloat64(tag, required)
	case reflect.String:
		return j.readString(tag, required)
	case reflect.Map:
		return j.readMap(ty, tag, required)
	case reflect.Struct:
		return j.readJceStruct(ty, tag, required)
	case reflect.Array:
		return j.readSlice(ty, tag, required)
	case reflect.Slice:
		return j.readSlice(ty, tag, required)
	default:
		return nil, fmt.Errorf("type mismatch.")
	}
}

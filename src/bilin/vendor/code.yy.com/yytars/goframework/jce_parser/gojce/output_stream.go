package gojce

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"reflect"
)

type jce_output_stream struct {
	buf *bytes.Buffer
}

func (j *jce_output_stream) ToBytes() []byte {
	return j.buf.Bytes()
}

func (j *jce_output_stream) GetLength() int {
	return j.buf.Len()
}

func (j *jce_output_stream) writeHead(ty byte, tag int) error {
	if tag < 15 {
		b := byte((byte(tag) << 4) | ty)
		return j.buf.WriteByte(b)
	} else if tag < 256 {
		b := byte((15 << 4) | ty)
		if err := j.buf.WriteByte(b); err != nil {
			return err
		}
		return j.buf.WriteByte(byte(tag))
	} else {
		return fmt.Errorf("tag is too large: %d", tag)
	}

}

func (j *jce_output_stream) writeBool(b bool, tag int) error {
	var by byte
	if b {
		by = 0x01
	} else {
		by = 0
	}
	return j.writeByte(by, tag)
}

func (j *jce_output_stream) writeByte(by byte, tag int) error {
	var err error
	if by == 0 {
		if err = j.writeHead(byte(ZERO_TAG), tag); err != nil {
			return err
		}
	} else {
		if err = j.writeHead(byte(BYTE), tag); err != nil {
			return err
		}

		if err = binary.Write(j.buf, binary.BigEndian, by); err != nil {
			return err
		}
	}
	return nil
}

func (j *jce_output_stream) writeInt16(n int16, tag int) error {
	var err error
	if n >= math.MinInt8 && n <= math.MaxInt8 {
		if err = j.writeByte(byte(n), tag); err != nil {
			return err
		}
	} else {
		if err = j.writeHead(byte(SHORT), tag); err != nil {
			return err
		}

		if err = binary.Write(j.buf, binary.BigEndian, n); err != nil {
			return err
		}
	}
	return nil
}

func (j *jce_output_stream) writeUint16(n uint16, tag int) error {
	return j.writeInt16(int16(n), tag)
}

func (j *jce_output_stream) writeInt32(n int32, tag int) error {
	var err error
	if n >= math.MinInt16 && n <= math.MaxInt16 {
		if err = j.writeInt16(int16(n), tag); err != nil {
			return err
		}
	} else {
		if err = j.writeHead(byte(INT), tag); err != nil {
			return err
		}

		if err = binary.Write(j.buf, binary.BigEndian, n); err != nil {
			return err
		}
	}
	return nil
}

func (j *jce_output_stream) writeUint32(n uint32, tag int) error {
	return j.writeUint64(uint64(n), tag)
}

func (j *jce_output_stream) writeInt64(n int64, tag int) error {
	var err error
	if n >= math.MinInt32 && n <= math.MaxInt32 {
		if err = j.writeInt32(int32(n), tag); err != nil {
			return err
		}
	} else {
		if err = j.writeHead(byte(LONG), tag); err != nil {
			return err
		}

		if err = binary.Write(j.buf, binary.BigEndian, n); err != nil {
			return err
		}
	}
	return nil
}

func (j *jce_output_stream) writeUint64(n uint64, tag int) error {
	return j.writeInt64(int64(n), tag)
}

func (j *jce_output_stream) writeFloat32(n float32, tag int) error {
	var err error
	if err = j.writeHead(byte(FLOAT), tag); err != nil {
		return err
	}

	if err = binary.Write(j.buf, binary.BigEndian, n); err != nil {
		return err
	}
	return nil
}

func (j *jce_output_stream) writeFloat64(n float64, tag int) error {
	var err error
	if err = j.writeHead(byte(DOUBLE), tag); err != nil {
		return err
	}

	if err = binary.Write(j.buf, binary.BigEndian, n); err != nil {
		return err
	}
	return nil
}

func (j *jce_output_stream) writeStringByte(s string, tag int) error {
	bs, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	if len(bs) > 255 {
		if err = j.writeHead(byte(STRING4), tag); err != nil {
			return err
		}

		if err = binary.Write(j.buf, binary.BigEndian, int32(len(bs))); err != nil {
			return err
		}

		if _, err = j.buf.Write(bs); err != nil {
			return err
		}
	} else {
		if err = j.writeHead(byte(STRING1), tag); err != nil {
			return err
		}

		if err = binary.Write(j.buf, binary.BigEndian, byte(len(bs))); err != nil {
			return err
		}

		if _, err = j.buf.Write(bs); err != nil {
			return err
		}
	}
	return nil
}

func (j *jce_output_stream) writeString(s string, tag int) error {
	var err error
	bs := []byte(s)
	if len(bs) > 255 {
		if err = j.writeHead(byte(STRING4), tag); err != nil {
			return err
		}

		if err = binary.Write(j.buf, binary.BigEndian, int32(len(bs))); err != nil {
			return err
		}

		if _, err = j.buf.Write(bs); err != nil {
			return err
		}
	} else {
		if err = j.writeHead(byte(STRING1), tag); err != nil {
			return err
		}

		if err = binary.Write(j.buf, binary.BigEndian, byte(len(bs))); err != nil {
			return err
		}

		if _, err = j.buf.Write(bs); err != nil {
			return err
		}
	}
	return nil
}

func (j *jce_output_stream) writeMap(v reflect.Value, tag int) error {
	var err error
	ty := v.Type().Elem()
	ele := v.Elem()
	if !v.IsValid() {
		return fmt.Errorf("val is not valid")
	}

	if ty.Kind() != reflect.Map {
		return fmt.Errorf("mismatch type:%v", ty.Kind())
	}

	if err = j.writeHead(byte(MAP), tag); err != nil {
		return err
	}

	if err = j.writeInt32(int32(ele.Len()), 0); err != nil {
		return err
	}

	for _, key := range ele.MapKeys() {
		key_ptr := reflect.New(key.Type())
		key_ptr.Elem().Set(key)
		if err = j.Write(key_ptr, 0); err != nil {
			return err
		}

		v_ptr := reflect.New(ele.MapIndex(key).Type())
		v_ptr.Elem().Set(ele.MapIndex(key))
		if err = j.Write(v_ptr, 1); err != nil {
			return err
		}
	}
	return nil
}

func (j *jce_output_stream) writeSlice(v reflect.Value, tag int) error {
	var err error
	ty := v.Type().Elem()
	ele := v.Elem()
	if !v.IsValid() {
		return fmt.Errorf("val is not valid")
	}

	if !(ty.Kind() == reflect.Array || ty.Kind() == reflect.Slice) {
		return fmt.Errorf("mismatch type: %v", ty.Kind())
	}

	switch ele.Interface().(type) {
	case []byte:
		if err = j.writeHead(byte(SIMPLE_LIST), tag); err != nil {
			return err
		}

		if err = j.writeHead(byte(BYTE), 0); err != nil {
			return err
		}

		if err = j.writeInt32(int32(ele.Len()), 0); err != nil {
			return err
		}

		if _, err = j.buf.Write(ele.Interface().([]byte)); err != nil {
			return err
		}
	default:
		if err = j.writeHead(byte(LIST), tag); err != nil {
			return err
		}

		if err = j.writeInt32(int32(ele.Len()), 0); err != nil {
			return err
		}

		for i := 0; i < ele.Len(); i++ {
			val := reflect.New(ele.Index(i).Type())
			val.Elem().Set(ele.Index(i))
			if err = j.Write(val, 0); err != nil {
				return err
			}
		}
	}

	return nil
}

func (j *jce_output_stream) writeJceStruct(v reflect.Value, tag int) error {
	var err error
	st := v.Interface().(JceStruct)
	if err = j.writeHead(byte(STRUCT_BEGIN), tag); err != nil {
		return err
	}
	if err = st.WriteTo(j); err != nil {
		return err
	}
	if err = j.writeHead(byte(STRUCT_END), tag); err != nil {
		return err
	}
	return nil
}

func (j *jce_output_stream) Write(v reflect.Value, tag int) error {
	kind := v.Kind()
	if kind == reflect.Ptr {
		kind = v.Elem().Kind()
		switch kind {
		case reflect.Bool:
			return j.writeBool(v.Elem().Interface().(bool), tag)
		case reflect.Uint8:
			return j.writeByte(byte(v.Elem().Interface().(uint8)), tag)
		case reflect.Uint16:
			return j.writeUint16(v.Elem().Interface().(uint16), tag)
		case reflect.Uint32:
			return j.writeUint32(v.Elem().Interface().(uint32), tag)
		case reflect.Uint64:
			return j.writeUint64(v.Elem().Interface().(uint64), tag)
		case reflect.Uint:
			return j.writeUint32(uint32(v.Elem().Interface().(uint)), tag)
		case reflect.Int8:
			return j.writeByte(byte(v.Elem().Interface().(int8)), tag)
		case reflect.Int16:
			return j.writeInt16(v.Elem().Interface().(int16), tag)
		case reflect.Int32:
			return j.writeInt32(v.Elem().Interface().(int32), tag)
		case reflect.Int64:
			return j.writeInt64(v.Elem().Interface().(int64), tag)
		case reflect.Int:
			return j.writeInt32(int32(v.Elem().Interface().(int)), tag)
		case reflect.Float32:
			return j.writeFloat32(v.Elem().Interface().(float32), tag)
		case reflect.Float64:
			return j.writeFloat64(v.Elem().Interface().(float64), tag)
		case reflect.String:
			return j.writeString(v.Elem().Interface().(string), tag)
		case reflect.Map:
			return j.writeMap(v, tag)
		case reflect.Struct:
			return j.writeJceStruct(v, tag)
		case reflect.Array:
			return j.writeSlice(v, tag)
		case reflect.Slice:
			return j.writeSlice(v, tag)
		default:
			return fmt.Errorf("type mismatch.")
		}
	} else {
		switch kind {
		case reflect.Bool:
			return j.writeBool(v.Interface().(bool), tag)
		case reflect.Uint8:
			return j.writeByte(byte(v.Interface().(uint8)), tag)
		case reflect.Uint16:
			return j.writeUint16(v.Interface().(uint16), tag)
		case reflect.Uint32:
			return j.writeUint32(v.Interface().(uint32), tag)
		case reflect.Uint64:
			return j.writeUint64(v.Interface().(uint64), tag)
		case reflect.Uint:
			return j.writeUint32(uint32(v.Interface().(uint)), tag)
		case reflect.Int8:
			return j.writeByte(byte(v.Interface().(int8)), tag)
		case reflect.Int16:
			return j.writeInt16(v.Interface().(int16), tag)
		case reflect.Int32:
			return j.writeInt32(v.Interface().(int32), tag)
		case reflect.Int64:
			return j.writeInt64(v.Interface().(int64), tag)
		case reflect.Int:
			return j.writeInt32(int32(v.Interface().(int)), tag)
		case reflect.Float32:
			return j.writeFloat32(v.Interface().(float32), tag)
		case reflect.Float64:
			return j.writeFloat64(v.Interface().(float64), tag)
		case reflect.String:
			return j.writeString(v.Interface().(string), tag)
		default:
			return fmt.Errorf("type mismatch. %v", kind)
		}
	}
}

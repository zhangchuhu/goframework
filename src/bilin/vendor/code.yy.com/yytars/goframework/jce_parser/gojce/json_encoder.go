package gojce

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

const (
	nestClassArray = iota
	nestClassStruct
)

type jce_json_encoder struct {
	buf     *bytes.Buffer
	hasPrev bool
}

func (j *jce_json_encoder) setPrev(has bool) {
	j.hasPrev = has
}

func (j *jce_json_encoder) stringEscape(s string) string {
	//s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, `"`, `\"`, -1)
	return s
}

func (j *jce_json_encoder) printComma() {
	if j.hasPrev {
		j.buf.WriteByte(',')
	}
}

func (j *jce_json_encoder) printField(fieldname string) {
	if fieldname != "" {
		j.buf.WriteString(fmt.Sprintf("\"%s\":", fieldname))
	}
}

func (j *jce_json_encoder) addFieldString(fieldName, content string) {
	j.printComma()
	j.printField(fieldName)
	j.buf.WriteString(fmt.Sprintf("\"%s\"", j.stringEscape(content)))
	j.setPrev(true)
}

func (j *jce_json_encoder) addFieldPlain(fieldName, content string) {
	j.printComma()
	j.printField(fieldName)
	j.buf.WriteString(fmt.Sprintf("%s", content))
	j.setPrev(true)
}

func (j *jce_json_encoder) addFieldComplex(fieldName string, content reflect.Value) {
	j.printComma()
	j.printField(fieldName)
	j.setPrev(false)
	j.EncodeJSON(content, "")
	j.setPrev(true)
}

func (j *jce_json_encoder) beginNest(fieldName string, class int) {
	j.printComma()
	j.printField(fieldName)
	switch class {
	case nestClassArray:
		j.buf.WriteRune('[')
	case nestClassStruct:
		j.buf.WriteRune('{')
	}
	j.setPrev(false)
}

func (j *jce_json_encoder) endNest(class int) {
	switch class {
	case nestClassArray:
		j.buf.WriteRune(']')
	case nestClassStruct:
		j.buf.WriteRune('}')
	}
	j.setPrev(true)
}

func (j *jce_json_encoder) appendContent(content string) {
	j.buf.WriteString(content)
}

func (j *jce_json_encoder) encodeMap(v reflect.Value, fieldName string) error {
	ty := v.Type().Elem()
	ele := v.Elem()
	if !v.IsValid() {
		return fmt.Errorf("val is not valid")
	}

	if ty.Kind() != reflect.Map {
		return fmt.Errorf("mismatch type:%v", ty.Kind())
	}

	if len(ele.MapKeys()) == 0 {
		j.addFieldPlain(fieldName, "{}")
		return nil
	}

	j.beginNest(fieldName, nestClassStruct)
	for _, key := range ele.MapKeys() {
		key_ptr := reflect.New(key.Type())
		key_ptr.Elem().Set(key)
		v_ptr := reflect.New(ele.MapIndex(key).Type())
		v_ptr.Elem().Set(ele.MapIndex(key))

		j.addFieldComplex(fmt.Sprintf("%v", key_ptr.Elem()), v_ptr)
	}
	j.endNest(nestClassStruct)
	return nil
}

func (j *jce_json_encoder) encodeJceStruct(v reflect.Value, fieldName string) error {
	st := v.Interface().(JceJsonSupporter)
	j.beginNest(fieldName, nestClassStruct)
	st.WriteJson(j)
	//WriteJson will append '}' rune to the end, so don't need to call endNest().
	return nil
}

func (j *jce_json_encoder) encodeSlice(v reflect.Value, fieldName string) error {
	var err error
	ty := v.Type().Elem()
	ele := v.Elem()
	if !v.IsValid() {
		return fmt.Errorf("val is not valid")
	}

	if !(ty.Kind() == reflect.Array || ty.Kind() == reflect.Slice) {
		return fmt.Errorf("mismatch type: %v", ty.Kind())
	}

	if ele.Len() == 0 {
		j.addFieldPlain(fieldName, "[]")
		return nil
	}

	j.beginNest(fieldName, nestClassArray)
	for i := 0; i < ele.Len(); i++ {
		val := reflect.New(ele.Index(i).Type())
		val.Elem().Set(ele.Index(i))
		if err = j.EncodeJSON(val, ""); err != nil {
			return err
		}
	}
	j.endNest(nestClassArray)
	return nil
}

func (j *jce_json_encoder) encodeString(v reflect.Value, fieldName string) error {
	j.addFieldString(fieldName, fmt.Sprintf("%v", v.Elem()))
	return nil
}

func (j *jce_json_encoder) EncodeJSON(v reflect.Value, fieldName string) error {
	kind := v.Kind()
	if kind == reflect.Ptr {
		kind = v.Elem().Kind()
		switch kind {
		case reflect.Map:
			return j.encodeMap(v, fieldName)
		case reflect.Struct:
			return j.encodeJceStruct(v, fieldName)
		case reflect.Array:
			return j.encodeSlice(v, fieldName)
		case reflect.Slice:
			return j.encodeSlice(v, fieldName)
		case reflect.String:
			return j.encodeString(v, fieldName)
		default:
			j.addFieldPlain(fieldName, fmt.Sprintf("%v", v.Elem()))
		}
	}

	return nil
}

func (j *jce_json_encoder) ToBytes() []byte {
	j.buf.WriteByte('}')
	return j.buf.Bytes()
}

package gojce

import (
	"bytes"
	"fmt"
	"reflect"
)

type jce_displayer struct {
	buf   *bytes.Buffer
	level int
}

func (j *jce_displayer) appendLine(fieldName, content string) {
	j.printField(fieldName)
	j.buf.WriteString(content)
	j.buf.WriteString("\n")
}

func (j *jce_displayer) append(fieldName, content string) {
	j.printField(fieldName)
	j.buf.WriteString(content)
}

func (j *jce_displayer) printField(fieldname string) {
	j.appendLevelTab()
	if fieldname != "" {
		j.buf.WriteString(fieldname + ": ")
	}
}

func (j *jce_displayer) appendLevelTab() {
	for i := 0; i < j.level; i++ {
		j.buf.WriteString("\t")
	}
}

func (j *jce_displayer) displayMap(v reflect.Value, fieldName string) error {
	ty := v.Type().Elem()
	ele := v.Elem()
	if !v.IsValid() {
		return fmt.Errorf("val is not valid")
	}

	if ty.Kind() != reflect.Map {
		return fmt.Errorf("mismatch type:%v", ty.Kind())
	}

	if len(ele.MapKeys()) == 0 {
		j.appendLine(fieldName, "0, {}")
		return nil
	}

	jd1 := NewDisplayer(j.buf, j.level+1)
	jd2 := NewDisplayer(j.buf, j.level+2)

	j.appendLine(fieldName, fmt.Sprintf("%d, {", len(ele.MapKeys())))
	i := 0
	for _, key := range ele.MapKeys() {
		key_ptr := reflect.New(key.Type())
		key_ptr.Elem().Set(key)
		v_ptr := reflect.New(ele.MapIndex(key).Type())
		v_ptr.Elem().Set(ele.MapIndex(key))

		jd1.Display(reflect.ValueOf("("), "")
		jd2.Display(key_ptr, fmt.Sprintf("key[%d]", i))
		jd2.Display(v_ptr, fmt.Sprintf("value[%d]", i))
		jd1.Display(reflect.ValueOf(")"), "")
		i++
	}
	j.appendLine("", "}")
	return nil
}

func (j *jce_displayer) displayJceStruct(v reflect.Value, fieldName string) error {
	st := v.Interface().(JceStruct)
	jd1 := NewDisplayer(j.buf, j.level+1)
	j.appendLine(fieldName, "{")
	st.Display(jd1)
	j.appendLine("", "}")
	return nil
}

func (j *jce_displayer) displaySlice(v reflect.Value, fieldName string) error {
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
		j.appendLine(fieldName, "{}")
		return nil
	}

	jd1 := NewDisplayer(j.buf, j.level+1)
	j.appendLine(fieldName, "{")
	for i := 0; i < ele.Len(); i++ {
		val := reflect.New(ele.Index(i).Type())
		val.Elem().Set(ele.Index(i))
		if err = jd1.Display(val, fmt.Sprintf("[%d]", i)); err != nil {
			return err
		}
	}
	j.appendLine("", "}")
	return nil
}

func (j *jce_displayer) displayString(v reflect.Value, fieldName string) error {
	j.appendLine(fieldName, fmt.Sprintf("\"%v\"", v.Elem()))
	return nil
}

func (j *jce_displayer) Display(v reflect.Value, fieldName string) error {
	kind := v.Kind()
	if kind == reflect.Ptr {
		kind = v.Elem().Kind()
		switch kind {
		case reflect.Map:
			return j.displayMap(v, fieldName)
		case reflect.Struct:
			return j.displayJceStruct(v, fieldName)
		case reflect.Array:
			return j.displaySlice(v, fieldName)
		case reflect.Slice:
			return j.displaySlice(v, fieldName)
		case reflect.String:
			return j.displayString(v, fieldName)
		default:
			j.appendLine(fieldName, fmt.Sprintf("%v", v.Elem()))
		}
	}

	return nil
}

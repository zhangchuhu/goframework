package gojce

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type jce_json_decoder struct {
	jsonBytes []byte
}

func (j *jce_json_decoder) DecodeJSON(rawVal reflect.Value) error {
	var json_obj interface{}
	err := json.Unmarshal(j.jsonBytes, &json_obj)
	if err != nil {
		return err
	}
	rawData := json_obj.(map[string]interface{})
	return j.DecodeRawData(rawData, rawVal)
}

func (j *jce_json_decoder) DecodeRawData(rawData map[string]interface{}, rawVal reflect.Value) error {
	var val reflect.Value
	kind := rawVal.Kind()
	switch kind {
	case reflect.Ptr:
		val = rawVal.Elem()
		if val.Kind() != reflect.Struct {
			return fmt.Errorf("Incompatible Type : %v : Looking For Struct", kind)
		}
	default:
		return fmt.Errorf("Incompatible Type : %v", kind)
	}

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		fieldName := tag.Get("json")

		if fieldName == "" {
			if strings.HasPrefix(typeField.Name, "M_") {
				fieldName = string([]byte(typeField.Name)[2:])
			} else {
				return fmt.Errorf("wrong field name[%s], use jce2go generate code", typeField.Name)
			}

			// if valueField.Kind() == reflect.Struct {
			// 	// We have a struct that may have indivdual tags. Process separately
			// 	d.DecodePath(m, valueField)
			// 	continue
			// } else if valueField.Kind() == reflect.Ptr && reflect.TypeOf(valueField).Kind() == reflect.Struct {
			// 	// We have a pointer to a struct
			// 	if valueField.IsNil() {
			// 		// Create the object since it doesn't exist
			// 		valueField.Set(reflect.New(valueField.Type().Elem()))
			// 		decoded, _ = d.DecodePath(m, valueField.Elem())
			// 		if decoded == false {
			// 			// If nothing was decoded for this object return the pointer to nil
			// 			valueField.Set(reflect.NewAt(valueField.Type().Elem(), nil))
			// 		}
			// 		continue
			// 	}

			// 	d.DecodePath(m, valueField.Elem())
			// 	continue
			// }
		}

		if data, ok := rawData[fieldName]; ok {
			err := j.decode(data, valueField.Addr())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (j *jce_json_decoder) decode(data interface{}, v reflect.Value) error {
	kind := v.Kind()
	//fmt.Printf("type=%v, data=%v, v=%v, kind=%v\n", reflect.TypeOf(data), data, v.Type(), v.Elem().Kind())
	if kind == reflect.Ptr {
		kind = v.Elem().Kind()
		switch kind {
		case reflect.Bool:
			return j.decodeBool(data, v)
		case reflect.Uint8:
			return j.decodeUint(data, v)
		case reflect.Uint16:
			return j.decodeUint(data, v)
		case reflect.Uint32:
			return j.decodeUint(data, v)
		case reflect.Uint64:
			return j.decodeUint(data, v)
		case reflect.Uint:
			return j.decodeUint(data, v)
		case reflect.Int8:
			return j.decodeInt(data, v)
		case reflect.Int16:
			return j.decodeInt(data, v)
		case reflect.Int32:
			return j.decodeInt(data, v)
		case reflect.Int64:
			return j.decodeInt(data, v)
		case reflect.Int:
			return j.decodeInt(data, v)
		case reflect.Float32:
			return j.decodeFloat(data, v)
		case reflect.Float64:
			return j.decodeFloat(data, v)
		case reflect.String:
			return j.decodeString(data, v)
		case reflect.Map:
			return j.decodeMap(data, v)
		case reflect.Struct:
			return j.decodeJceStruct(data, v)
		case reflect.Array:
			return j.decodeSlice(data, v)
		case reflect.Slice:
			return j.decodeSlice(data, v)
		default:
			return fmt.Errorf("type mismatch.")
		}
	}
	return fmt.Errorf("type mismatch.")
}

func (j *jce_json_decoder) decodeBool(data interface{}, v reflect.Value) error {
	if ret, ok := data.(bool); ok {
		if v.Elem().IsValid() && v.Elem().CanSet() {
			v.Elem().SetBool(ret)
		}
	}

	return nil
}

func (j *jce_json_decoder) decodeUint(data interface{}, v reflect.Value) error {
	if ret, ok := data.(float64); ok {
		if v.Elem().IsValid() && v.Elem().CanSet() {
			v.Elem().SetUint(uint64(ret))
		}
	}
	return nil
}

func (j *jce_json_decoder) decodeInt(data interface{}, v reflect.Value) error {
	if ret, ok := data.(float64); ok {
		if v.Elem().IsValid() && v.Elem().CanSet() {
			v.Elem().SetInt(int64(ret))
		}
	}
	return nil
}

func (j *jce_json_decoder) decodeFloat(data interface{}, v reflect.Value) error {
	if ret, ok := data.(float64); ok {
		if v.Elem().IsValid() && v.Elem().CanSet() {
			v.Elem().SetFloat(ret)
		}
	}
	return nil
}

func (j *jce_json_decoder) decodeString(data interface{}, v reflect.Value) error {
	if ret, ok := data.(string); ok {
		if v.Elem().IsValid() && v.Elem().CanSet() {
			v.Elem().SetString(ret)
		}
	}
	return nil
}

func (j *jce_json_decoder) decodeMap(data interface{}, v reflect.Value) error {
	if ret, ok := data.(map[string]interface{}); ok {
		if v.Elem().IsValid() && v.Elem().CanSet() {
			elem_type := v.Elem().Type()
			key_type := elem_type.Key()
			val_type := elem_type.Elem()
			m := reflect.MakeMap(elem_type)

			for mk, mv := range ret {
				var key_val reflect.Value
				var_val := reflect.New(val_type)
				switch key_type.Kind() {
				case reflect.String:
					key_val = reflect.New(key_type)
					key_val.Elem().SetString(mk)
				case reflect.Float32, reflect.Float64:
					key_val = reflect.New(reflect.TypeOf(float64(0)))
					f, err := strconv.ParseFloat(mk, 64)
					if err != nil {
						return err
					}
					key_val.Elem().SetFloat(f)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					key_val = reflect.New(key_type)
					u, err := strconv.ParseUint(mk, 10, 64)
					if err != nil {
						return err
					}
					key_val.Elem().SetUint(u)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					key_val = reflect.New(key_type)
					i, err := strconv.ParseInt(mk, 10, 64)
					if err != nil {
						return err
					}
					key_val.Elem().SetInt(i)
				default:
					return fmt.Errorf("unknown map key type:%v", key_type.Kind())
				}

				j.decode(mv, var_val)
				m.SetMapIndex(key_val.Elem(), var_val.Elem())
			}
			v.Elem().Set(m)
		}
	}
	return nil
}

func (j *jce_json_decoder) decodeJceStruct(data interface{}, v reflect.Value) error {
	if ret, ok := data.(map[string]interface{}); ok {
		if v.Elem().IsValid() && v.Elem().CanSet() {
			return j.DecodeRawData(ret, v)
		}
	}
	return nil
}

func (j *jce_json_decoder) decodeSlice(data interface{}, v reflect.Value) error {
	if ret, ok := data.([]interface{}); ok {
		if v.Elem().IsValid() && v.Elem().CanSet() {
			size := len(ret)
			elem_type := v.Elem().Type()
			slice := reflect.MakeSlice(elem_type, size, size)
			for i := 0; i < size; i++ {
				elem_val := slice.Index(i)
				if err := j.decode(ret[i], elem_val.Addr()); err != nil {
					return err
				}
			}
			v.Elem().Set(slice)
		}
	}
	return nil
}

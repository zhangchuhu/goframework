// **********************************************************************
// This file was generated by a TAF parser!
// TAF version 1.6.0 by WSRD Tencent.
// Generated from `BusF.jce'
// **********************************************************************

package taf
import "reflect"
import "code.yy.com/yytars/goframework/jce_parser/gojce"

type BusCommuData struct {
    CommuKey string
    ErrorInfo string
    C2sMmapName string
    C2sMmapSize int64
    C2sFifoName string
    S2cMmapName string
    S2cMmapSize int64
    S2cFifoName string
}

func (_obj *BusCommuData) resetDefault() {
    _obj.CommuKey = ""
    _obj.ErrorInfo = ""
    _obj.C2sMmapName = ""
    _obj.C2sMmapSize = 0
    _obj.C2sFifoName = ""
    _obj.S2cMmapName = ""
    _obj.S2cMmapSize = 0
    _obj.S2cFifoName = ""
}

func (_obj *BusCommuData) WriteTo(_os gojce.JceOutputStream) error {
    var _err error
    if _err = _os.Write(reflect.ValueOf(&_obj.CommuKey), 0); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.ErrorInfo), 1); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.C2sMmapName), 2); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.C2sMmapSize), 3); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.C2sFifoName), 4); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.S2cMmapName), 5); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.S2cMmapSize), 6); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.S2cFifoName), 7); _err != nil {
        return _err
    }
    return nil
}

func (_obj *BusCommuData) ReadFrom(_is gojce.JceInputStream) error {
    var _err error
    var _i interface{}
    _obj.resetDefault()
    _i, _err = _is.Read(reflect.TypeOf(_obj.CommuKey), 0, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.CommuKey = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.ErrorInfo), 1, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.ErrorInfo = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.C2sMmapName), 2, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.C2sMmapName = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.C2sMmapSize), 3, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.C2sMmapSize = _i.(int64)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.C2sFifoName), 4, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.C2sFifoName = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.S2cMmapName), 5, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.S2cMmapName = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.S2cMmapSize), 6, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.S2cMmapSize = _i.(int64)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.S2cFifoName), 7, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.S2cFifoName = _i.(string)
    }
    return nil
}

func (_obj *BusCommuData) Display(_ds gojce.JceDisplayer) {
    _ds.Display(reflect.ValueOf(&_obj.CommuKey), "CommuKey")
    _ds.Display(reflect.ValueOf(&_obj.ErrorInfo), "ErrorInfo")
    _ds.Display(reflect.ValueOf(&_obj.C2sMmapName), "c2sMmapName")
    _ds.Display(reflect.ValueOf(&_obj.C2sMmapSize), "c2sMmapSize")
    _ds.Display(reflect.ValueOf(&_obj.C2sFifoName), "c2sFifoName")
    _ds.Display(reflect.ValueOf(&_obj.S2cMmapName), "s2cMmapName")
    _ds.Display(reflect.ValueOf(&_obj.S2cMmapSize), "s2cMmapSize")
    _ds.Display(reflect.ValueOf(&_obj.S2cFifoName), "s2cFifoName")
}

func (_obj *BusCommuData) WriteJson(_en gojce.JceJsonEncoder) ([]byte, error) {
    var _err error
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.CommuKey), "CommuKey")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.ErrorInfo), "ErrorInfo")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.C2sMmapName), "c2sMmapName")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.C2sMmapSize), "c2sMmapSize")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.C2sFifoName), "c2sFifoName")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.S2cMmapName), "s2cMmapName")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.S2cMmapSize), "s2cMmapSize")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.S2cFifoName), "s2cFifoName")
    if _err != nil {
        return nil, _err
    }
    return _en.ToBytes(), nil
}

func (_obj *BusCommuData) ReadJson(_de gojce.JceJsonDecoder) error {
    return _de.DecodeJSON(reflect.ValueOf(_obj))
}


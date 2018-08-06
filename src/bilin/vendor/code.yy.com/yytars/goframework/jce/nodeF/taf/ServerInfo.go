// **********************************************************************
// This file was generated by a TAF parser!
// TAF version 1.6.0 by WSRD Tencent.
// Generated from `NodeF.jce'
// **********************************************************************

package taf
import "reflect"
import "code.yy.com/yytars/goframework/jce_parser/gojce"

type ServerInfo struct {
    Application string
    ServerName string
    Pid int32
    Adapter string
    ModuleType string
    Container string
}

func (_obj *ServerInfo) resetDefault() {
    _obj.Application = ""
    _obj.ServerName = ""
    _obj.Pid = 0
    _obj.Adapter = ""
    _obj.ModuleType = "taf"
    _obj.Container = ""
}

func (_obj *ServerInfo) WriteTo(_os gojce.JceOutputStream) error {
    var _err error
    if _err = _os.Write(reflect.ValueOf(&_obj.Application), 0); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.ServerName), 1); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.Pid), 2); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.Adapter), 3); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.ModuleType), 4); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.Container), 5); _err != nil {
        return _err
    }
    return nil
}

func (_obj *ServerInfo) ReadFrom(_is gojce.JceInputStream) error {
    var _err error
    var _i interface{}
    _obj.resetDefault()
    _i, _err = _is.Read(reflect.TypeOf(_obj.Application), 0, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.Application = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.ServerName), 1, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.ServerName = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.Pid), 2, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.Pid = _i.(int32)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.Adapter), 3, false)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.Adapter = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.ModuleType), 4, false)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.ModuleType = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.Container), 5, false)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.Container = _i.(string)
    }
    return nil
}

func (_obj *ServerInfo) Display(_ds gojce.JceDisplayer) {
    _ds.Display(reflect.ValueOf(&_obj.Application), "application")
    _ds.Display(reflect.ValueOf(&_obj.ServerName), "serverName")
    _ds.Display(reflect.ValueOf(&_obj.Pid), "pid")
    _ds.Display(reflect.ValueOf(&_obj.Adapter), "adapter")
    _ds.Display(reflect.ValueOf(&_obj.ModuleType), "moduleType")
    _ds.Display(reflect.ValueOf(&_obj.Container), "container")
}

func (_obj *ServerInfo) WriteJson(_en gojce.JceJsonEncoder) ([]byte, error) {
    var _err error
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.Application), "application")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.ServerName), "serverName")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.Pid), "pid")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.Adapter), "adapter")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.ModuleType), "moduleType")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.Container), "container")
    if _err != nil {
        return nil, _err
    }
    return _en.ToBytes(), nil
}

func (_obj *ServerInfo) ReadJson(_de gojce.JceJsonDecoder) error {
    return _de.DecodeJSON(reflect.ValueOf(_obj))
}


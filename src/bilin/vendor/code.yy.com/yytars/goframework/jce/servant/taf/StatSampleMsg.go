// **********************************************************************
// This file was generated by a TAF parser!
// TAF version 1.6.0 by WSRD Tencent.
// Generated from `StatF.jce'
// **********************************************************************

package taf
import "reflect"
import "code.yy.com/yytars/goframework/jce_parser/gojce"

type StatSampleMsg struct {
    Unid string
    MasterName string
    SlaveName string
    InterfaceName string
    MasterIp string
    SlaveIp string
    Depth int32
    Width int32
    ParentWidth int32
}

func (_obj *StatSampleMsg) resetDefault() {
    _obj.Unid = ""
    _obj.MasterName = ""
    _obj.SlaveName = ""
    _obj.InterfaceName = ""
    _obj.MasterIp = ""
    _obj.SlaveIp = ""
    _obj.Depth = 0
    _obj.Width = 0
    _obj.ParentWidth = 0
}

func (_obj *StatSampleMsg) WriteTo(_os gojce.JceOutputStream) error {
    var _err error
    if _err = _os.Write(reflect.ValueOf(&_obj.Unid), 0); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.MasterName), 1); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.SlaveName), 2); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.InterfaceName), 3); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.MasterIp), 4); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.SlaveIp), 5); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.Depth), 6); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.Width), 7); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.ParentWidth), 8); _err != nil {
        return _err
    }
    return nil
}

func (_obj *StatSampleMsg) ReadFrom(_is gojce.JceInputStream) error {
    var _err error
    var _i interface{}
    _obj.resetDefault()
    _i, _err = _is.Read(reflect.TypeOf(_obj.Unid), 0, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.Unid = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.MasterName), 1, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.MasterName = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.SlaveName), 2, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.SlaveName = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.InterfaceName), 3, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.InterfaceName = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.MasterIp), 4, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.MasterIp = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.SlaveIp), 5, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.SlaveIp = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.Depth), 6, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.Depth = _i.(int32)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.Width), 7, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.Width = _i.(int32)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.ParentWidth), 8, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.ParentWidth = _i.(int32)
    }
    return nil
}

func (_obj *StatSampleMsg) Display(_ds gojce.JceDisplayer) {
    _ds.Display(reflect.ValueOf(&_obj.Unid), "unid")
    _ds.Display(reflect.ValueOf(&_obj.MasterName), "masterName")
    _ds.Display(reflect.ValueOf(&_obj.SlaveName), "slaveName")
    _ds.Display(reflect.ValueOf(&_obj.InterfaceName), "interfaceName")
    _ds.Display(reflect.ValueOf(&_obj.MasterIp), "masterIp")
    _ds.Display(reflect.ValueOf(&_obj.SlaveIp), "slaveIp")
    _ds.Display(reflect.ValueOf(&_obj.Depth), "depth")
    _ds.Display(reflect.ValueOf(&_obj.Width), "width")
    _ds.Display(reflect.ValueOf(&_obj.ParentWidth), "parentWidth")
}

func (_obj *StatSampleMsg) WriteJson(_en gojce.JceJsonEncoder) ([]byte, error) {
    var _err error
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.Unid), "unid")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.MasterName), "masterName")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.SlaveName), "slaveName")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.InterfaceName), "interfaceName")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.MasterIp), "masterIp")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.SlaveIp), "slaveIp")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.Depth), "depth")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.Width), "width")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.ParentWidth), "parentWidth")
    if _err != nil {
        return nil, _err
    }
    return _en.ToBytes(), nil
}

func (_obj *StatSampleMsg) ReadJson(_de gojce.JceJsonDecoder) error {
    return _de.DecodeJSON(reflect.ValueOf(_obj))
}

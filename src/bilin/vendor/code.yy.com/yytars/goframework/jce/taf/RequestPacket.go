// **********************************************************************
// This file was generated by a TAF parser!
// TAF version 1.6.0 by WSRD Tencent.
// Generated from `RequestF.jce'
// **********************************************************************

package taf
import "reflect"
import "code.yy.com/yytars/goframework/jce_parser/gojce"

type RequestPacket struct {
    IVersion int16
    CPacketType byte
    IMessageType int32
    IRequestId int32
    SServantName string
    SFuncName string
    SBuffer []byte
    ITimeout int32
    Context map[string]string
    Status map[string]string
}

func (_obj *RequestPacket) resetDefault() {
    _obj.IVersion = 0
    _obj.CPacketType = 0
    _obj.IMessageType = 0
    _obj.IRequestId = 0
    _obj.SServantName = ""
    _obj.SFuncName = ""
    _obj.ITimeout = 0
}

func (_obj *RequestPacket) WriteTo(_os gojce.JceOutputStream) error {
    var _err error
    if _err = _os.Write(reflect.ValueOf(&_obj.IVersion), 1); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.CPacketType), 2); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.IMessageType), 3); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.IRequestId), 4); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.SServantName), 5); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.SFuncName), 6); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.SBuffer), 7); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.ITimeout), 8); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.Context), 9); _err != nil {
        return _err
    }
    if _err = _os.Write(reflect.ValueOf(&_obj.Status), 10); _err != nil {
        return _err
    }
    return nil
}

func (_obj *RequestPacket) ReadFrom(_is gojce.JceInputStream) error {
    var _err error
    var _i interface{}
    _obj.resetDefault()
    _i, _err = _is.Read(reflect.TypeOf(_obj.IVersion), 1, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.IVersion = _i.(int16)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.CPacketType), 2, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.CPacketType = _i.(byte)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.IMessageType), 3, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.IMessageType = _i.(int32)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.IRequestId), 4, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.IRequestId = _i.(int32)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.SServantName), 5, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.SServantName = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.SFuncName), 6, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.SFuncName = _i.(string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.SBuffer), 7, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.SBuffer = _i.([]byte)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.ITimeout), 8, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.ITimeout = _i.(int32)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.Context), 9, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.Context = _i.(map[string]string)
    }
    _i, _err = _is.Read(reflect.TypeOf(_obj.Status), 10, true)
    if _err != nil {
        return _err
    }
    if _i != nil {
        _obj.Status = _i.(map[string]string)
    }
    return nil
}

func (_obj *RequestPacket) Display(_ds gojce.JceDisplayer) {
    _ds.Display(reflect.ValueOf(&_obj.IVersion), "iVersion")
    _ds.Display(reflect.ValueOf(&_obj.CPacketType), "cPacketType")
    _ds.Display(reflect.ValueOf(&_obj.IMessageType), "iMessageType")
    _ds.Display(reflect.ValueOf(&_obj.IRequestId), "iRequestId")
    _ds.Display(reflect.ValueOf(&_obj.SServantName), "sServantName")
    _ds.Display(reflect.ValueOf(&_obj.SFuncName), "sFuncName")
    _ds.Display(reflect.ValueOf(&_obj.SBuffer), "sBuffer")
    _ds.Display(reflect.ValueOf(&_obj.ITimeout), "iTimeout")
    _ds.Display(reflect.ValueOf(&_obj.Context), "context")
    _ds.Display(reflect.ValueOf(&_obj.Status), "status")
}

func (_obj *RequestPacket) WriteJson(_en gojce.JceJsonEncoder) ([]byte, error) {
    var _err error
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.IVersion), "iVersion")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.CPacketType), "cPacketType")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.IMessageType), "iMessageType")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.IRequestId), "iRequestId")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.SServantName), "sServantName")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.SFuncName), "sFuncName")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.SBuffer), "sBuffer")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.ITimeout), "iTimeout")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.Context), "context")
    if _err != nil {
        return nil, _err
    }
    _err = _en.EncodeJSON(reflect.ValueOf(&_obj.Status), "status")
    if _err != nil {
        return nil, _err
    }
    return _en.ToBytes(), nil
}

func (_obj *RequestPacket) ReadJson(_de gojce.JceJsonDecoder) error {
    return _de.DecodeJSON(reflect.ValueOf(_obj))
}


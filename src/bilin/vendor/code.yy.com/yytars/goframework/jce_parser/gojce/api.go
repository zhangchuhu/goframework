package gojce

import (
	"bytes"
	"reflect"
)

//jce数据流解码接口
type JceInputStream interface {
	//读取一个特定的jce
	Read(ty reflect.Type, tag int, required bool) (interface{}, error)
}

//jce数据流编码接口
type JceOutputStream interface {
	Write(v reflect.Value, tag int) error
	ToBytes() []byte
	GetLength() int
}

//jce内容格式化接口
type JceDisplayer interface {
	//格式化打印
	Display(v reflect.Value, name string) error
}

//jce到json的编码接口
type JceJsonEncoder interface {
	EncodeJSON(v reflect.Value, name string) error
	ToBytes() []byte
}

//json到jce的解码接口
type JceJsonDecoder interface {
	DecodeJSON(v reflect.Value) error
}

//jce结构接口
type JceStruct interface {
	// 将jce数据通过编码器编码
	WriteTo(os JceOutputStream) error
	// 通过解析器解析jce数据
	ReadFrom(is JceInputStream) error
	// 格式化jce内容
	Display(ds JceDisplayer)
}

//jce结构对json支持接口
type JceJsonSupporter interface {
	WriteJson(JceJsonEncoder) ([]byte, error)
	ReadJson(JceJsonDecoder) error
}

//创建一个新的jce内容格式化结构
func NewDisplayer(buf *bytes.Buffer, level int) JceDisplayer {
	ds := new(jce_displayer)
	ds.buf = buf
	ds.level = level
	return ds
}

//创建一个jce数据流解析器
func NewInputStream(bs []byte) JceInputStream {
	is := new(jce_input_stream)
	// is.bs = make([]byte, len(bs))
	// copy(is.bs, bs)
	is.reader = bytes.NewReader(bs)
	return is
}

//创建一个jce数据流解析器，带位移
func NewInputStreamWithPos(bs []byte, pos int) JceInputStream {
	is := new(jce_input_stream)
	// is.bs = make([]byte, len(bs[pos:]))
	// copy(is.bs, bs[pos:])
	is.reader = bytes.NewReader(bs[pos:])
	return is
}

//创建一个jce数据流编码器
func NewOutputStream() JceOutputStream {
	os := new(jce_output_stream)
	os.buf = bytes.NewBuffer(nil)
	return os
}

//创建一个jce到JSON的编码器
func NewJceJsonEncoder() JceJsonEncoder {
	encoder := new(jce_json_encoder)
	encoder.buf = bytes.NewBuffer(nil)
	encoder.buf.WriteByte('{')
	return encoder
}

//创建一个JSON到jce的解码器
func NewJceJsonDecoder(bs []byte) JceJsonDecoder {
	decoder := new(jce_json_decoder)
	decoder.jsonBytes = bs
	return decoder
}

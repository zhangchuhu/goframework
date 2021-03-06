// Autogenerated by Thrift Compiler (0.11.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package bilin_msg_filter

import (
	"bytes"
	"reflect"
	"context"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = context.Background
var _ = reflect.DeepEqual
var _ = bytes.Equal

type MsgFilter interface {
  // Parameters:
  //  - Msg
  CheckMsg(ctx context.Context, msg string) (r int32, err error)
}

type MsgFilterClient struct {
  c thrift.TClient
}

// Deprecated: Use NewMsgFilter instead
func NewMsgFilterClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *MsgFilterClient {
  return &MsgFilterClient{
    c: thrift.NewTStandardClient(f.GetProtocol(t), f.GetProtocol(t)),
  }
}

// Deprecated: Use NewMsgFilter instead
func NewMsgFilterClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *MsgFilterClient {
  return &MsgFilterClient{
    c: thrift.NewTStandardClient(iprot, oprot),
  }
}

func NewMsgFilterClient(c thrift.TClient) *MsgFilterClient {
  return &MsgFilterClient{
    c: c,
  }
}

// Parameters:
//  - Msg
func (p *MsgFilterClient) CheckMsg(ctx context.Context, msg string) (r int32, err error) {
  var _args0 MsgFilterCheckMsgArgs
  _args0.Msg = msg
  var _result1 MsgFilterCheckMsgResult
  if err = p.c.Call(ctx, "check_msg", &_args0, &_result1); err != nil {
    return
  }
  return _result1.GetSuccess(), nil
}

type MsgFilterProcessor struct {
  processorMap map[string]thrift.TProcessorFunction
  handler MsgFilter
}

func (p *MsgFilterProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
  p.processorMap[key] = processor
}

func (p *MsgFilterProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
  processor, ok = p.processorMap[key]
  return processor, ok
}

func (p *MsgFilterProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
  return p.processorMap
}

func NewMsgFilterProcessor(handler MsgFilter) *MsgFilterProcessor {

  self2 := &MsgFilterProcessor{handler:handler, processorMap:make(map[string]thrift.TProcessorFunction)}
  self2.processorMap["check_msg"] = &msgFilterProcessorCheckMsg{handler:handler}
return self2
}

func (p *MsgFilterProcessor) Process(ctx context.Context, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  name, _, seqId, err := iprot.ReadMessageBegin()
  if err != nil { return false, err }
  if processor, ok := p.GetProcessorFunction(name); ok {
    return processor.Process(ctx, seqId, iprot, oprot)
  }
  iprot.Skip(thrift.STRUCT)
  iprot.ReadMessageEnd()
  x3 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function " + name)
  oprot.WriteMessageBegin(name, thrift.EXCEPTION, seqId)
  x3.Write(oprot)
  oprot.WriteMessageEnd()
  oprot.Flush()
  return false, x3

}

type msgFilterProcessorCheckMsg struct {
  handler MsgFilter
}

func (p *msgFilterProcessorCheckMsg) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := MsgFilterCheckMsgArgs{}
  if err = args.Read(iprot); err != nil {
    iprot.ReadMessageEnd()
    x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
    oprot.WriteMessageBegin("check_msg", thrift.EXCEPTION, seqId)
    x.Write(oprot)
    oprot.WriteMessageEnd()
    oprot.Flush()
    return false, err
  }

  iprot.ReadMessageEnd()
  result := MsgFilterCheckMsgResult{}
var retval int32
  var err2 error
  if retval, err2 = p.handler.CheckMsg(ctx, args.Msg); err2 != nil {
    x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing check_msg: " + err2.Error())
    oprot.WriteMessageBegin("check_msg", thrift.EXCEPTION, seqId)
    x.Write(oprot)
    oprot.WriteMessageEnd()
    oprot.Flush()
    return true, err2
  } else {
    result.Success = &retval
}
  if err2 = oprot.WriteMessageBegin("check_msg", thrift.REPLY, seqId); err2 != nil {
    err = err2
  }
  if err2 = result.Write(oprot); err == nil && err2 != nil {
    err = err2
  }
  if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
    err = err2
  }
  if err2 = oprot.Flush(); err == nil && err2 != nil {
    err = err2
  }
  if err != nil {
    return
  }
  return true, err
}


// HELPER FUNCTIONS AND STRUCTURES

// Attributes:
//  - Msg
type MsgFilterCheckMsgArgs struct {
  Msg string `thrift:"msg,1" db:"msg" json:"msg"`
}

func NewMsgFilterCheckMsgArgs() *MsgFilterCheckMsgArgs {
  return &MsgFilterCheckMsgArgs{}
}


func (p *MsgFilterCheckMsgArgs) GetMsg() string {
  return p.Msg
}
func (p *MsgFilterCheckMsgArgs) Read(iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *MsgFilterCheckMsgArgs)  ReadField1(iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.Msg = v
}
  return nil
}

func (p *MsgFilterCheckMsgArgs) Write(oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin("check_msg_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *MsgFilterCheckMsgArgs) writeField1(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("msg", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:msg: ", p), err) }
  if err := oprot.WriteString(string(p.Msg)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.msg (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:msg: ", p), err) }
  return err
}

func (p *MsgFilterCheckMsgArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("MsgFilterCheckMsgArgs(%+v)", *p)
}

// Attributes:
//  - Success
type MsgFilterCheckMsgResult struct {
  Success *int32 `thrift:"success,0" db:"success" json:"success,omitempty"`
}

func NewMsgFilterCheckMsgResult() *MsgFilterCheckMsgResult {
  return &MsgFilterCheckMsgResult{}
}

var MsgFilterCheckMsgResult_Success_DEFAULT int32
func (p *MsgFilterCheckMsgResult) GetSuccess() int32 {
  if !p.IsSetSuccess() {
    return MsgFilterCheckMsgResult_Success_DEFAULT
  }
return *p.Success
}
func (p *MsgFilterCheckMsgResult) IsSetSuccess() bool {
  return p.Success != nil
}

func (p *MsgFilterCheckMsgResult) Read(iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 0:
      if fieldTypeId == thrift.I32 {
        if err := p.ReadField0(iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *MsgFilterCheckMsgResult)  ReadField0(iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI32(); err != nil {
  return thrift.PrependError("error reading field 0: ", err)
} else {
  p.Success = &v
}
  return nil
}

func (p *MsgFilterCheckMsgResult) Write(oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin("check_msg_result"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField0(oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *MsgFilterCheckMsgResult) writeField0(oprot thrift.TProtocol) (err error) {
  if p.IsSetSuccess() {
    if err := oprot.WriteFieldBegin("success", thrift.I32, 0); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err) }
    if err := oprot.WriteI32(int32(*p.Success)); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T.success (0) field write error: ", p), err) }
    if err := oprot.WriteFieldEnd(); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err) }
  }
  return err
}

func (p *MsgFilterCheckMsgResult) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("MsgFilterCheckMsgResult(%+v)", *p)
}



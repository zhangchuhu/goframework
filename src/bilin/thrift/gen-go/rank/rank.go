// Autogenerated by Thrift Compiler (0.11.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package rank

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

// Attributes:
//  - UID
//  - Value
//  - Rank
type TRank struct {
  UID int64 `thrift:"uid,1" db:"uid" json:"uid"`
  Value int64 `thrift:"value,2" db:"value" json:"value"`
  Rank int64 `thrift:"rank,3" db:"rank" json:"rank"`
}

func NewTRank() *TRank {
  return &TRank{}
}


func (p *TRank) GetUID() int64 {
  return p.UID
}

func (p *TRank) GetValue() int64 {
  return p.Value
}

func (p *TRank) GetRank() int64 {
  return p.Rank
}
func (p *TRank) Read(iprot thrift.TProtocol) error {
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
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField1(iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField2(iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.I64 {
        if err := p.ReadField3(iprot); err != nil {
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

func (p *TRank)  ReadField1(iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.UID = v
}
  return nil
}

func (p *TRank)  ReadField2(iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Value = v
}
  return nil
}

func (p *TRank)  ReadField3(iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI64(); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.Rank = v
}
  return nil
}

func (p *TRank) Write(oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin("TRank"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(oprot); err != nil { return err }
    if err := p.writeField2(oprot); err != nil { return err }
    if err := p.writeField3(oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *TRank) writeField1(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("uid", thrift.I64, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:uid: ", p), err) }
  if err := oprot.WriteI64(int64(p.UID)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.uid (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:uid: ", p), err) }
  return err
}

func (p *TRank) writeField2(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("value", thrift.I64, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:value: ", p), err) }
  if err := oprot.WriteI64(int64(p.Value)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.value (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:value: ", p), err) }
  return err
}

func (p *TRank) writeField3(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("rank", thrift.I64, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:rank: ", p), err) }
  if err := oprot.WriteI64(int64(p.Rank)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.rank (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:rank: ", p), err) }
  return err
}

func (p *TRank) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("TRank(%+v)", *p)
}

type TRankService interface {
  // Parameters:
  //  - Code
  //  - Rtype
  //  - Ctype
  //  - Latest
  //  - Size
  //  - TimeParam
  QueryRank(ctx context.Context, code string, rtype string, ctype string, latest bool, size int32, timeParam int32) (r []*TRank, err error)
}

type TRankServiceClient struct {
  c thrift.TClient
}

// Deprecated: Use NewTRankService instead
func NewTRankServiceClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *TRankServiceClient {
  return &TRankServiceClient{
    c: thrift.NewTStandardClient(f.GetProtocol(t), f.GetProtocol(t)),
  }
}

// Deprecated: Use NewTRankService instead
func NewTRankServiceClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *TRankServiceClient {
  return &TRankServiceClient{
    c: thrift.NewTStandardClient(iprot, oprot),
  }
}

func NewTRankServiceClient(c thrift.TClient) *TRankServiceClient {
  return &TRankServiceClient{
    c: c,
  }
}

// Parameters:
//  - Code
//  - Rtype
//  - Ctype
//  - Latest
//  - Size
//  - TimeParam
func (p *TRankServiceClient) QueryRank(ctx context.Context, code string, rtype string, ctype string, latest bool, size int32, timeParam int32) (r []*TRank, err error) {
  var _args0 TRankServiceQueryRankArgs
  _args0.Code = code
  _args0.Rtype = rtype
  _args0.Ctype = ctype
  _args0.Latest = latest
  _args0.Size = size
  _args0.TimeParam = timeParam
  var _result1 TRankServiceQueryRankResult
  if err = p.c.Call(ctx, "queryRank", &_args0, &_result1); err != nil {
    return
  }
  return _result1.GetSuccess(), nil
}

type TRankServiceProcessor struct {
  processorMap map[string]thrift.TProcessorFunction
  handler TRankService
}

func (p *TRankServiceProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
  p.processorMap[key] = processor
}

func (p *TRankServiceProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
  processor, ok = p.processorMap[key]
  return processor, ok
}

func (p *TRankServiceProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
  return p.processorMap
}

func NewTRankServiceProcessor(handler TRankService) *TRankServiceProcessor {

  self2 := &TRankServiceProcessor{handler:handler, processorMap:make(map[string]thrift.TProcessorFunction)}
  self2.processorMap["queryRank"] = &tRankServiceProcessorQueryRank{handler:handler}
return self2
}

func (p *TRankServiceProcessor) Process(ctx context.Context, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
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

type tRankServiceProcessorQueryRank struct {
  handler TRankService
}

func (p *tRankServiceProcessorQueryRank) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  args := TRankServiceQueryRankArgs{}
  if err = args.Read(iprot); err != nil {
    iprot.ReadMessageEnd()
    x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
    oprot.WriteMessageBegin("queryRank", thrift.EXCEPTION, seqId)
    x.Write(oprot)
    oprot.WriteMessageEnd()
    oprot.Flush()
    return false, err
  }

  iprot.ReadMessageEnd()
  result := TRankServiceQueryRankResult{}
var retval []*TRank
  var err2 error
  if retval, err2 = p.handler.QueryRank(ctx, args.Code, args.Rtype, args.Ctype, args.Latest, args.Size, args.TimeParam); err2 != nil {
    x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing queryRank: " + err2.Error())
    oprot.WriteMessageBegin("queryRank", thrift.EXCEPTION, seqId)
    x.Write(oprot)
    oprot.WriteMessageEnd()
    oprot.Flush()
    return true, err2
  } else {
    result.Success = retval
}
  if err2 = oprot.WriteMessageBegin("queryRank", thrift.REPLY, seqId); err2 != nil {
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
//  - Code
//  - Rtype
//  - Ctype
//  - Latest
//  - Size
//  - TimeParam
type TRankServiceQueryRankArgs struct {
  Code string `thrift:"code,1" db:"code" json:"code"`
  Rtype string `thrift:"rtype,2" db:"rtype" json:"rtype"`
  Ctype string `thrift:"ctype,3" db:"ctype" json:"ctype"`
  Latest bool `thrift:"latest,4" db:"latest" json:"latest"`
  Size int32 `thrift:"size,5" db:"size" json:"size"`
  TimeParam int32 `thrift:"timeParam,6" db:"timeParam" json:"timeParam"`
}

func NewTRankServiceQueryRankArgs() *TRankServiceQueryRankArgs {
  return &TRankServiceQueryRankArgs{}
}


func (p *TRankServiceQueryRankArgs) GetCode() string {
  return p.Code
}

func (p *TRankServiceQueryRankArgs) GetRtype() string {
  return p.Rtype
}

func (p *TRankServiceQueryRankArgs) GetCtype() string {
  return p.Ctype
}

func (p *TRankServiceQueryRankArgs) GetLatest() bool {
  return p.Latest
}

func (p *TRankServiceQueryRankArgs) GetSize() int32 {
  return p.Size
}

func (p *TRankServiceQueryRankArgs) GetTimeParam() int32 {
  return p.TimeParam
}
func (p *TRankServiceQueryRankArgs) Read(iprot thrift.TProtocol) error {
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
    case 2:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField2(iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(fieldTypeId); err != nil {
          return err
        }
      }
    case 3:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField3(iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(fieldTypeId); err != nil {
          return err
        }
      }
    case 4:
      if fieldTypeId == thrift.BOOL {
        if err := p.ReadField4(iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(fieldTypeId); err != nil {
          return err
        }
      }
    case 5:
      if fieldTypeId == thrift.I32 {
        if err := p.ReadField5(iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(fieldTypeId); err != nil {
          return err
        }
      }
    case 6:
      if fieldTypeId == thrift.I32 {
        if err := p.ReadField6(iprot); err != nil {
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

func (p *TRankServiceQueryRankArgs)  ReadField1(iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.Code = v
}
  return nil
}

func (p *TRankServiceQueryRankArgs)  ReadField2(iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Rtype = v
}
  return nil
}

func (p *TRankServiceQueryRankArgs)  ReadField3(iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(); err != nil {
  return thrift.PrependError("error reading field 3: ", err)
} else {
  p.Ctype = v
}
  return nil
}

func (p *TRankServiceQueryRankArgs)  ReadField4(iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBool(); err != nil {
  return thrift.PrependError("error reading field 4: ", err)
} else {
  p.Latest = v
}
  return nil
}

func (p *TRankServiceQueryRankArgs)  ReadField5(iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI32(); err != nil {
  return thrift.PrependError("error reading field 5: ", err)
} else {
  p.Size = v
}
  return nil
}

func (p *TRankServiceQueryRankArgs)  ReadField6(iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI32(); err != nil {
  return thrift.PrependError("error reading field 6: ", err)
} else {
  p.TimeParam = v
}
  return nil
}

func (p *TRankServiceQueryRankArgs) Write(oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin("queryRank_args"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(oprot); err != nil { return err }
    if err := p.writeField2(oprot); err != nil { return err }
    if err := p.writeField3(oprot); err != nil { return err }
    if err := p.writeField4(oprot); err != nil { return err }
    if err := p.writeField5(oprot); err != nil { return err }
    if err := p.writeField6(oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *TRankServiceQueryRankArgs) writeField1(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("code", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:code: ", p), err) }
  if err := oprot.WriteString(string(p.Code)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.code (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:code: ", p), err) }
  return err
}

func (p *TRankServiceQueryRankArgs) writeField2(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("rtype", thrift.STRING, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:rtype: ", p), err) }
  if err := oprot.WriteString(string(p.Rtype)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.rtype (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:rtype: ", p), err) }
  return err
}

func (p *TRankServiceQueryRankArgs) writeField3(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("ctype", thrift.STRING, 3); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:ctype: ", p), err) }
  if err := oprot.WriteString(string(p.Ctype)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.ctype (3) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 3:ctype: ", p), err) }
  return err
}

func (p *TRankServiceQueryRankArgs) writeField4(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("latest", thrift.BOOL, 4); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:latest: ", p), err) }
  if err := oprot.WriteBool(bool(p.Latest)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.latest (4) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 4:latest: ", p), err) }
  return err
}

func (p *TRankServiceQueryRankArgs) writeField5(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("size", thrift.I32, 5); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:size: ", p), err) }
  if err := oprot.WriteI32(int32(p.Size)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.size (5) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 5:size: ", p), err) }
  return err
}

func (p *TRankServiceQueryRankArgs) writeField6(oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin("timeParam", thrift.I32, 6); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:timeParam: ", p), err) }
  if err := oprot.WriteI32(int32(p.TimeParam)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.timeParam (6) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 6:timeParam: ", p), err) }
  return err
}

func (p *TRankServiceQueryRankArgs) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("TRankServiceQueryRankArgs(%+v)", *p)
}

// Attributes:
//  - Success
type TRankServiceQueryRankResult struct {
  Success []*TRank `thrift:"success,0" db:"success" json:"success,omitempty"`
}

func NewTRankServiceQueryRankResult() *TRankServiceQueryRankResult {
  return &TRankServiceQueryRankResult{}
}

var TRankServiceQueryRankResult_Success_DEFAULT []*TRank

func (p *TRankServiceQueryRankResult) GetSuccess() []*TRank {
  return p.Success
}
func (p *TRankServiceQueryRankResult) IsSetSuccess() bool {
  return p.Success != nil
}

func (p *TRankServiceQueryRankResult) Read(iprot thrift.TProtocol) error {
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
      if fieldTypeId == thrift.LIST {
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

func (p *TRankServiceQueryRankResult)  ReadField0(iprot thrift.TProtocol) error {
  _, size, err := iprot.ReadListBegin()
  if err != nil {
    return thrift.PrependError("error reading list begin: ", err)
  }
  tSlice := make([]*TRank, 0, size)
  p.Success =  tSlice
  for i := 0; i < size; i ++ {
    _elem4 := &TRank{}
    if err := _elem4.Read(iprot); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem4), err)
    }
    p.Success = append(p.Success, _elem4)
  }
  if err := iprot.ReadListEnd(); err != nil {
    return thrift.PrependError("error reading list end: ", err)
  }
  return nil
}

func (p *TRankServiceQueryRankResult) Write(oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin("queryRank_result"); err != nil {
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

func (p *TRankServiceQueryRankResult) writeField0(oprot thrift.TProtocol) (err error) {
  if p.IsSetSuccess() {
    if err := oprot.WriteFieldBegin("success", thrift.LIST, 0); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err) }
    if err := oprot.WriteListBegin(thrift.STRUCT, len(p.Success)); err != nil {
      return thrift.PrependError("error writing list begin: ", err)
    }
    for _, v := range p.Success {
      if err := v.Write(oprot); err != nil {
        return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
      }
    }
    if err := oprot.WriteListEnd(); err != nil {
      return thrift.PrependError("error writing list end: ", err)
    }
    if err := oprot.WriteFieldEnd(); err != nil {
      return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err) }
  }
  return err
}

func (p *TRankServiceQueryRankResult) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("TRankServiceQueryRankResult(%+v)", *p)
}


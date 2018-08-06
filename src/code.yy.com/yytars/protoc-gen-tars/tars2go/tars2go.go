// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2015 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Package grpc outputs gRPC service descriptions in Go code.
// It runs as a plugin for the Go protocol buffer compiler plugin.
// It is linked in to protoc-gen-go.
package tars2go

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	pb "code.yy.com/yytars/protoc-gen-tars/descriptor"
	"code.yy.com/yytars/protoc-gen-tars/generator"
)

// Paths for packages used by code generated in this file,
// relative to the import_prefix of the generator.Generator.
const (
	contextPkgPath = "context"
)

func init() {
	generator.RegisterPlugin(new(tars))
}

// grpc is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for gRPC support.
type tars struct {
	gen *generator.Generator
}

// Name returns the name of this plugin, "grpc".
func (g *tars) Name() string {
	return "tars2go"
}

// The names for packages imported in the generated code.
// They may vary from the final path component of the import path
// if the name is used by other packages.
var (
	contextPkg string
	grpcPkg    string
)

// Init initializes the plugin.
func (g *tars) Init(gen *generator.Generator) {
	g.gen = gen
	contextPkg = generator.RegisterUniquePackageName("context", nil)
	grpcPkg = generator.RegisterUniquePackageName("grpc", nil)
}

// Given a type name defined in a .proto, return its object.
// Also record that we're using it, to guarantee the associated import.
func (g *tars) objectNamed(name string) generator.Object {
	g.gen.RecordTypeUse(name)
	return g.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (g *tars) typeName(str string) string {
	return g.gen.TypeName(g.objectNamed(str))
}

// P forwards to g.gen.P.
func (g *tars) P(args ...interface{}) { g.gen.P(args...) }

// Generate generates code for the services in the given file.
func (g *tars) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	g.P("// Reference imports to suppress errors if they are not otherwise used.")
	g.P("var _ ", contextPkg, ".Context")
	g.P()
	for i, service := range file.FileDescriptorProto.Service {
		g.generateService(file, service, i)
	}
}

// GenerateImports generates the import declaration for this file.
func (g *tars) GenerateImports(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	g.P("import (")
	g.P(contextPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, contextPkgPath)))
	//g.P(grpcPkg, " ", strconv.Quote(path.Join(g.gen.ImportPrefix, grpcPkgPath)))
	g.P("\"","code.yy.com/yytars/goframework/tars/servant","\"")
	g.P("\"","code.yy.com/yytars/goframework/tars/servant/model","\"")
	g.P("\"","code.yy.com/yytars/goframework/tars/servant/protocol","\"")
	g.P("\"","errors","\"")
	g.P(")")
	g.P()
}

func (g *tars)generateSetServant(servName string)  {
	g.P("func (c *", unexport(servName), "Client) ", "SetServant(s model.PbServant)", "{")
	//g.P("func(c *helloClient) SetServant(s model.Servant){")
	g.P("c.s = s")
	g.P("}")
}
// reservedClientName records whether a client name is reserved on the client side.
var reservedClientName = map[string]bool{
// TODO: do we need any in gRPC?
}

func unexport(s string) string { return strings.ToLower(s[:1]) + s[1:] }

func (g *tars)generateDispatcher(servName,serverType string,service *pb.ServiceDescriptorProto)  {
	dispacthername := unexport(servName)+"Dispatcher"

	g.P("type ",dispacthername," struct {")
	g.P("}")
	g.P()

	g.P("func New",servName,"Dispatcher() ","servant.PbDispatcher  {")
	g.P("return &",dispacthername,"{}")
	g.P("}")
	g.P()

	g.P("func(_obj *",dispacthername,")Dispatch(ctx context.Context, _val interface{}, req *pbtaf.RequestPacket) (*pbtaf.ResponsePacket,error){")
	g.P("var pbbuf []byte")
	g.P("_imp := _val.(",serverType,")")

	g.P("switch req.SFuncName {")
	for _,method :=range service.Method{
		g.P("case ",`"`,method.GetName(),`":`)
		g.P("var req_ ",g.typeName(method.GetInputType()))
		g.P("if err := proto.Unmarshal(req.SBuffer,&req_);err != nil{")
		g.P("return nil,err")
		g.P("}")
		g.P()

		origMethName := method.GetName()
		methName := generator.CamelCase(origMethName)
		if reservedClientName[methName] {
			methName += "_"
		}
		g.P("_ret,err := _imp.",methName,"(ctx,&req_)")
		g.P("if err != nil{")
		g.P("return nil,err")
		g.P("}")
		g.P()

		g.P("if pbbuf,err = proto.Marshal(_ret);err != nil{")
		g.P("return nil,err")
		g.P("}")
		g.P()

	}
	g.P("default:")
	g.P("return nil,errors.New(\"unknow func\")")
	g.P("}")

	//g.P("var status map[string]string")
	g.P("return &pbtaf.ResponsePacket{")
	g.P("IVersion:     1,")
	g.P("IRequestId:   req.IRequestId,")
	g.P("SBuffer:      pbbuf,")
	//g.P("Status:       status,")
	g.P("Context:      req.Context,")
	g.P("}",",nil")

	g.P("}")
}
// generateService generates all the code for the named service.
func (g *tars) generateService(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto, index int) {
	path := fmt.Sprintf("6,%d", index) // 6 means service.

	origServName := service.GetName()
	fullServName := origServName
	if pkg := file.GetPackage(); pkg != "" {
		fullServName = pkg + "." + fullServName
	}
	servName := generator.CamelCase(origServName)

	g.P()
	g.P("// Client API for ", servName, " service")
	g.P()

	// Client interface.
	g.P("type ", servName, "Client interface {")
	//g.P("servant.HashCodeSetter")
	for i, method := range service.Method {
		g.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) // 2 means method in a service.
		g.P(g.generateClientSignature(servName, method))
	}
	g.P("}")
	g.P()

	// Client structure.
	g.P("type ", unexport(servName), "Client struct {")
	g.P("s ", "model.PbServant")
	//g.P("servant.ConsistenHashInfo")
	g.P("}")
	g.P()

	// NewClient factory.
	g.P("func New", servName, "Client(objname string, comm *servant.Communicator) ", servName, "Client {")
	g.P("if comm == nil || objname == \"\"{")
	g.P("return nil")
	g.P("}")
	g.P("return &", unexport(servName), "Client{s : comm.GetServantProxy(objname)}")
	g.P("}")
	g.P()

	var methodIndex, streamIndex int
	serviceDescVar := "_" + servName + "_serviceDesc"
	// Client method implementations.
	for _, method := range service.Method {
		var descExpr string
		if !method.GetServerStreaming() && !method.GetClientStreaming() {
			// Unary RPC method
			descExpr = fmt.Sprintf("&%s.Methods[%d]", serviceDescVar, methodIndex)
			methodIndex++
		} else {
			// Streaming RPC method
			descExpr = fmt.Sprintf("&%s.Streams[%d]", serviceDescVar, streamIndex)
			streamIndex++
		}
		g.generateClientMethod(servName, fullServName, serviceDescVar, method, descExpr)
	}
	// TODO will remove
	g.P()
	//g.generateSetServant(servName)

	g.P("// Server API for ", servName, " service")
	g.P()

	// Server interface.
	serverType := servName + "Server"
	g.P("type ", serverType, " interface {")
	for i, method := range service.Method {
		g.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) // 2 means method in a service.
		g.P(g.generateServerSignature(servName, method))
	}
	g.P("}")
	g.P()

	g.generateDispatcher(servName,serverType,service)

	g.P()
}

// generateClientSignature returns the client-side signature for a method.
func (g *tars) generateClientSignature(servName string, method *pb.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)
	if reservedClientName[methName] {
		methName += "_"
	}
	reqArg := ", in *" + g.typeName(method.GetInputType())
	if method.GetClientStreaming() {
		reqArg = ""
	}
	respName := "*" + g.typeName(method.GetOutputType())
	if method.GetServerStreaming() || method.GetClientStreaming() {
		respName = servName + "_" + generator.CamelCase(origMethName) + "Client"
	}
	return fmt.Sprintf("%s(ctx %s.Context%s, opts ...map[string]string) (%s, error)", methName, contextPkg, reqArg, respName)
}

func (g *tars) generateClientMethod(servName, fullServName, serviceDescVar string, method *pb.MethodDescriptorProto, descExpr string) {
	//sname := fmt.Sprintf("/%s/%s", fullServName, method.GetName())
	//methName := generator.CamelCase(method.GetName())
	//inType := g.typeName(method.GetInputType())
	outType := g.typeName(method.GetOutputType())

	g.P("func (c *", unexport(servName), "Client) ", g.generateClientSignature(servName, method), "{")

	g.P("var (")
	//g.P("_status map[string]string")
	//g.P("_context map[string]string = make(map[string]string, 1)")
	g.P("reply ",outType)
	g.P(")")
	g.P()

	g.P("pbbuf,err := proto.Marshal(in)")
	g.P("if err != nil {")
	g.P("return nil,err")
	g.P("}")
	g.P()

	//g.P("if c.ConsistentHashEnable{")
	//g.P(`_context["consisthashkey"] = c.HashCode`)
	//g.P("}")
	//g.P()

	g.P("_resp,err := ",`c.s.Pb_invoke(ctx,0,"`,method.GetName(),`", pbbuf)`)
	g.P("if err != nil {")
	g.P("return nil,err")
	g.P("}")
	g.P()

	g.P("if err = proto.Unmarshal(_resp.SBuffer,&reply);err != nil{")
	g.P("return nil,err")
	g.P("}")
	g.P("return &reply,nil")
	g.P("}")
	return
}

// generateServerSignature returns the server-side signature for a method.
func (g *tars) generateServerSignature(servName string, method *pb.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)
	if reservedClientName[methName] {
		methName += "_"
	}

	var reqArgs []string
	ret := "error"
	if !method.GetServerStreaming() && !method.GetClientStreaming() {
		reqArgs = append(reqArgs, contextPkg+".Context")
		ret = "(*" + g.typeName(method.GetOutputType()) + ", error)"
	}
	if !method.GetClientStreaming() {
		reqArgs = append(reqArgs, "*"+g.typeName(method.GetInputType()))
	}
	if method.GetServerStreaming() || method.GetClientStreaming() {
		reqArgs = append(reqArgs, servName+"_"+generator.CamelCase(origMethName)+"Server")
	}

	return methName + "(" + strings.Join(reqArgs, ", ") + ") " + ret
}

//func (g *tars) generateServerMethod(servName string, method *pb.MethodDescriptorProto) string {
//	methName := generator.CamelCase(method.GetName())
//	hname := fmt.Sprintf("_%s_%s_Handler", servName, methName)
//	inType := g.typeName(method.GetInputType())
//	outType := g.typeName(method.GetOutputType())
//
//	if !method.GetServerStreaming() && !method.GetClientStreaming() {
//		g.P("func ", hname, "(srv interface{}, ctx ", contextPkg, ".Context, dec func(interface{}) error) (interface{}, error) {")
//		g.P("in := new(", inType, ")")
//		g.P("if err := dec(in); err != nil { return nil, err }")
//		g.P("out, err := srv.(", servName, "Server).", methName, "(ctx, in)")
//		g.P("if err != nil { return nil, err }")
//		g.P("return out, nil")
//		g.P("}")
//		g.P()
//		return hname
//	}
//	streamType := unexport(servName) + methName + "Server"
//	g.P("func ", hname, "(srv interface{}, stream ", grpcPkg, ".ServerStream) error {")
//	if !method.GetClientStreaming() {
//		g.P("m := new(", inType, ")")
//		g.P("if err := stream.RecvMsg(m); err != nil { return err }")
//		g.P("return srv.(", servName, "Server).", methName, "(m, &", streamType, "{stream})")
//	} else {
//		g.P("return srv.(", servName, "Server).", methName, "(&", streamType, "{stream})")
//	}
//	g.P("}")
//	g.P()
//
//	genSend := method.GetServerStreaming()
//	genSendAndClose := !method.GetServerStreaming()
//	genRecv := method.GetClientStreaming()
//
//	// Stream auxiliary types and methods.
//	g.P("type ", servName, "_", methName, "Server interface {")
//	if genSend {
//		g.P("Send(*", outType, ") error")
//	}
//	if genSendAndClose {
//		g.P("SendAndClose(*", outType, ") error")
//	}
//	if genRecv {
//		g.P("Recv() (*", inType, ", error)")
//	}
//	g.P(grpcPkg, ".ServerStream")
//	g.P("}")
//	g.P()
//
//	g.P("type ", streamType, " struct {")
//	g.P(grpcPkg, ".ServerStream")
//	g.P("}")
//	g.P()
//
//	if genSend {
//		g.P("func (x *", streamType, ") Send(m *", outType, ") error {")
//		g.P("return x.ServerStream.SendMsg(m)")
//		g.P("}")
//		g.P()
//	}
//	if genSendAndClose {
//		g.P("func (x *", streamType, ") SendAndClose(m *", outType, ") error {")
//		g.P("return x.ServerStream.SendMsg(m)")
//		g.P("}")
//		g.P()
//	}
//	if genRecv {
//		g.P("func (x *", streamType, ") Recv() (*", inType, ", error) {")
//		g.P("m := new(", inType, ")")
//		g.P("if err := x.ServerStream.RecvMsg(m); err != nil { return nil, err }")
//		g.P("return m, nil")
//		g.P("}")
//		g.P()
//	}
//
//	return hname
//}

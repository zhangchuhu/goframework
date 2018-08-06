// Autogenerated by Thrift Compiler (0.11.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package main

import (
        "context"
        "flag"
        "fmt"
        "math"
        "net"
        "net/url"
        "os"
        "strconv"
        "strings"
        "git.apache.org/thrift.git/lib/go/thrift"
	"common"
        "user"
)

var _ = common.GoUnusedProtection__

func Usage() {
  fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
  flag.PrintDefaults()
  fmt.Fprintln(os.Stderr, "\nFunctions:")
  fmt.Fprintln(os.Stderr, "  ComRet internalHttpForbidUser(string userIds, string bilinIds)")
  fmt.Fprintln(os.Stderr, "  GetUserForConServerRet getUserForConServer(i64 userId, i64 fromUserId, i64 groupId, bool ifPush, string requestType)")
  fmt.Fprintln(os.Stderr, "  QueryAttentionMeCountRet QueryAttentionMeCount( userList)")
  fmt.Fprintln(os.Stderr, "  i32 ping()")
  fmt.Fprintln(os.Stderr)
  os.Exit(0)
}

func main() {
  flag.Usage = Usage
  var host string
  var port int
  var protocol string
  var urlString string
  var framed bool
  var useHttp bool
  var parsedUrl *url.URL
  var trans thrift.TTransport
  _ = strconv.Atoi
  _ = math.Abs
  flag.Usage = Usage
  flag.StringVar(&host, "h", "localhost", "Specify host and port")
  flag.IntVar(&port, "p", 9090, "Specify port")
  flag.StringVar(&protocol, "P", "binary", "Specify the protocol (binary, compact, simplejson, json)")
  flag.StringVar(&urlString, "u", "", "Specify the url")
  flag.BoolVar(&framed, "framed", false, "Use framed transport")
  flag.BoolVar(&useHttp, "http", false, "Use http")
  flag.Parse()
  
  if len(urlString) > 0 {
    var err error
    parsedUrl, err = url.Parse(urlString)
    if err != nil {
      fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
      flag.Usage()
    }
    host = parsedUrl.Host
    useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http"
  } else if useHttp {
    _, err := url.Parse(fmt.Sprint("http://", host, ":", port))
    if err != nil {
      fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
      flag.Usage()
    }
  }
  
  cmd := flag.Arg(0)
  var err error
  if useHttp {
    trans, err = thrift.NewTHttpClient(parsedUrl.String())
  } else {
    portStr := fmt.Sprint(port)
    if strings.Contains(host, ":") {
           host, portStr, err = net.SplitHostPort(host)
           if err != nil {
                   fmt.Fprintln(os.Stderr, "error with host:", err)
                   os.Exit(1)
           }
    }
    trans, err = thrift.NewTSocket(net.JoinHostPort(host, portStr))
    if err != nil {
      fmt.Fprintln(os.Stderr, "error resolving address:", err)
      os.Exit(1)
    }
    if framed {
      trans = thrift.NewTFramedTransport(trans)
    }
  }
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error creating transport", err)
    os.Exit(1)
  }
  defer trans.Close()
  var protocolFactory thrift.TProtocolFactory
  switch protocol {
  case "compact":
    protocolFactory = thrift.NewTCompactProtocolFactory()
    break
  case "simplejson":
    protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
    break
  case "json":
    protocolFactory = thrift.NewTJSONProtocolFactory()
    break
  case "binary", "":
    protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid protocol specified: ", protocol)
    Usage()
    os.Exit(1)
  }
  iprot := protocolFactory.GetProtocol(trans)
  oprot := protocolFactory.GetProtocol(trans)
  client := user.NewUserServiceClient(thrift.NewTStandardClient(iprot, oprot))
  if err := trans.Open(); err != nil {
    fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
    os.Exit(1)
  }
  
  switch cmd {
  case "internalHttpForbidUser":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "InternalHttpForbidUser requires 2 args")
      flag.Usage()
    }
    argvalue0 := flag.Arg(1)
    value0 := argvalue0
    argvalue1 := flag.Arg(2)
    value1 := argvalue1
    fmt.Print(client.InternalHttpForbidUser(context.Background(), value0, value1))
    fmt.Print("\n")
    break
  case "getUserForConServer":
    if flag.NArg() - 1 != 5 {
      fmt.Fprintln(os.Stderr, "GetUserForConServer requires 5 args")
      flag.Usage()
    }
    argvalue0, err12 := (strconv.ParseInt(flag.Arg(1), 10, 64))
    if err12 != nil {
      Usage()
      return
    }
    value0 := argvalue0
    argvalue1, err13 := (strconv.ParseInt(flag.Arg(2), 10, 64))
    if err13 != nil {
      Usage()
      return
    }
    value1 := argvalue1
    argvalue2, err14 := (strconv.ParseInt(flag.Arg(3), 10, 64))
    if err14 != nil {
      Usage()
      return
    }
    value2 := argvalue2
    argvalue3 := flag.Arg(4) == "true"
    value3 := argvalue3
    argvalue4 := flag.Arg(5)
    value4 := argvalue4
    fmt.Print(client.GetUserForConServer(context.Background(), value0, value1, value2, value3, value4))
    fmt.Print("\n")
    break
  case "QueryAttentionMeCount":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "QueryAttentionMeCount requires 1 args")
      flag.Usage()
    }
    arg17 := flag.Arg(1)
    mbTrans18 := thrift.NewTMemoryBufferLen(len(arg17))
    defer mbTrans18.Close()
    _, err19 := mbTrans18.WriteString(arg17)
    if err19 != nil { 
      Usage()
      return
    }
    factory20 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt21 := factory20.GetProtocol(mbTrans18)
    containerStruct0 := user.NewUserServiceQueryAttentionMeCountArgs()
    err22 := containerStruct0.ReadField1(jsProt21)
    if err22 != nil {
      Usage()
      return
    }
    argvalue0 := containerStruct0.UserList
    value0 := argvalue0
    fmt.Print(client.QueryAttentionMeCount(context.Background(), value0))
    fmt.Print("\n")
    break
  case "ping":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "Ping requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.Ping(context.Background()))
    fmt.Print("\n")
    break
  case "":
    Usage()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
  }
}
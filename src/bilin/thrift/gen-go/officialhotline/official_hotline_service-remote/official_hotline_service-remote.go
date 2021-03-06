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
        "bilin/thrift/gen-go/common"
        "officialhotline"
)

var _ = common.GoUnusedProtection__

func Usage() {
  fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
  flag.PrintDefaults()
  fmt.Fprintln(os.Stderr, "\nFunctions:")
  fmt.Fprintln(os.Stderr, "  OfficialHotlineRet switchCallback(i32 sid, i64 oldUserId, i64 newUserId, string msg, i32 resultCode)")
  fmt.Fprintln(os.Stderr, "  OfficialHotlineRet onOfficialMike(i32 sid, i64 userId)")
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
  client := officialhotline.NewOfficialHotlineServiceClient(thrift.NewTStandardClient(iprot, oprot))
  if err := trans.Open(); err != nil {
    fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
    os.Exit(1)
  }
  
  switch cmd {
  case "switchCallback":
    if flag.NArg() - 1 != 5 {
      fmt.Fprintln(os.Stderr, "SwitchCallback requires 5 args")
      flag.Usage()
    }
    tmp0, err5 := (strconv.Atoi(flag.Arg(1)))
    if err5 != nil {
      Usage()
      return
    }
    argvalue0 := int32(tmp0)
    value0 := argvalue0
    argvalue1, err6 := (strconv.ParseInt(flag.Arg(2), 10, 64))
    if err6 != nil {
      Usage()
      return
    }
    value1 := argvalue1
    argvalue2, err7 := (strconv.ParseInt(flag.Arg(3), 10, 64))
    if err7 != nil {
      Usage()
      return
    }
    value2 := argvalue2
    argvalue3 := flag.Arg(4)
    value3 := argvalue3
    tmp4, err9 := (strconv.Atoi(flag.Arg(5)))
    if err9 != nil {
      Usage()
      return
    }
    argvalue4 := int32(tmp4)
    value4 := argvalue4
    fmt.Print(client.SwitchCallback(context.Background(), value0, value1, value2, value3, value4))
    fmt.Print("\n")
    break
  case "onOfficialMike":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "OnOfficialMike requires 2 args")
      flag.Usage()
    }
    tmp0, err10 := (strconv.Atoi(flag.Arg(1)))
    if err10 != nil {
      Usage()
      return
    }
    argvalue0 := int32(tmp0)
    value0 := argvalue0
    argvalue1, err11 := (strconv.ParseInt(flag.Arg(2), 10, 64))
    if err11 != nil {
      Usage()
      return
    }
    value1 := argvalue1
    fmt.Print(client.OnOfficialMike(context.Background(), value0, value1))
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

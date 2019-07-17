// Autogenerated by Thrift Compiler (0.9.3)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package main

import (
	"ds_user_service"
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"math"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nFunctions:")
	fmt.Fprintln(os.Stderr, "  PublisherSegments getPublisherSegments(StorageId storage_id, Context context)")
	fmt.Fprintln(os.Stderr, "  Conversions getConversions(StorageId storage_id, ConversionId conversion_id, Context context)")
	fmt.Fprintln(os.Stderr, "  AdvertisersData getAdvertisersData(MarketId market_id, Context context)")
	fmt.Fprintln(os.Stderr, "  CookieState getCookieState(CookieStateId id, PlatformHash platform_hash, Context context)")
	fmt.Fprintln(os.Stderr, "  void setDeviceData(DeviceId id, DeviceData data, Context context)")
	fmt.Fprintln(os.Stderr, "  void delDeviceData(DeviceId id, Context context)")
	fmt.Fprintln(os.Stderr, "  void delDeviceDataByKey(DeviceId id, DeviceDataKeys keys, Context context)")
	fmt.Fprintln(os.Stderr, "  DeviceData getDeviceData(DeviceId id, Context context)")
	fmt.Fprintln(os.Stderr, "  bool isOptedOut(DeviceId id, Context context)")
	fmt.Fprintln(os.Stderr, "  CrossDeviceData getCrossDeviceData(StorageId storage_id, Context context)")
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
	var parsedUrl url.URL
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
		parsedUrl, err := url.Parse(urlString)
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
	client := ds_user_service.NewUserDataServiceClientFactory(trans, protocolFactory)
	if err := trans.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
		os.Exit(1)
	}

	switch cmd {
	case "getPublisherSegments":
		if flag.NArg()-1 != 2 {
			fmt.Fprintln(os.Stderr, "GetPublisherSegments requires 2 args")
			flag.Usage()
		}
		argvalue0 := flag.Arg(1)
		value0 := ds_user_service.StorageId(argvalue0)
		arg56 := flag.Arg(2)
		mbTrans57 := thrift.NewTMemoryBufferLen(len(arg56))
		defer mbTrans57.Close()
		_, err58 := mbTrans57.WriteString(arg56)
		if err58 != nil {
			Usage()
			return
		}
		factory59 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt60 := factory59.GetProtocol(mbTrans57)
		containerStruct1 := ds_user_service.NewUserDataServiceGetPublisherSegmentsArgs()
		err61 := containerStruct1.ReadField2(jsProt60)
		if err61 != nil {
			Usage()
			return
		}
		argvalue1 := containerStruct1.Context
		value1 := ds_user_service.Context(argvalue1)
		fmt.Print(client.GetPublisherSegments(value0, value1))
		fmt.Print("\n")
		break
	case "getConversions":
		if flag.NArg()-1 != 3 {
			fmt.Fprintln(os.Stderr, "GetConversions requires 3 args")
			flag.Usage()
		}
		argvalue0 := flag.Arg(1)
		value0 := ds_user_service.StorageId(argvalue0)
		argvalue1 := flag.Arg(2)
		value1 := ds_user_service.ConversionId(argvalue1)
		arg64 := flag.Arg(3)
		mbTrans65 := thrift.NewTMemoryBufferLen(len(arg64))
		defer mbTrans65.Close()
		_, err66 := mbTrans65.WriteString(arg64)
		if err66 != nil {
			Usage()
			return
		}
		factory67 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt68 := factory67.GetProtocol(mbTrans65)
		containerStruct2 := ds_user_service.NewUserDataServiceGetConversionsArgs()
		err69 := containerStruct2.ReadField3(jsProt68)
		if err69 != nil {
			Usage()
			return
		}
		argvalue2 := containerStruct2.Context
		value2 := ds_user_service.Context(argvalue2)
		fmt.Print(client.GetConversions(value0, value1, value2))
		fmt.Print("\n")
		break
	case "getAdvertisersData":
		if flag.NArg()-1 != 2 {
			fmt.Fprintln(os.Stderr, "GetAdvertisersData requires 2 args")
			flag.Usage()
		}
		argvalue0 := flag.Arg(1)
		value0 := ds_user_service.MarketId(argvalue0)
		arg71 := flag.Arg(2)
		mbTrans72 := thrift.NewTMemoryBufferLen(len(arg71))
		defer mbTrans72.Close()
		_, err73 := mbTrans72.WriteString(arg71)
		if err73 != nil {
			Usage()
			return
		}
		factory74 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt75 := factory74.GetProtocol(mbTrans72)
		containerStruct1 := ds_user_service.NewUserDataServiceGetAdvertisersDataArgs()
		err76 := containerStruct1.ReadField2(jsProt75)
		if err76 != nil {
			Usage()
			return
		}
		argvalue1 := containerStruct1.Context
		value1 := ds_user_service.Context(argvalue1)
		fmt.Print(client.GetAdvertisersData(value0, value1))
		fmt.Print("\n")
		break
	case "getCookieState":
		if flag.NArg()-1 != 3 {
			fmt.Fprintln(os.Stderr, "GetCookieState requires 3 args")
			flag.Usage()
		}
		argvalue0 := flag.Arg(1)
		value0 := ds_user_service.CookieStateId(argvalue0)
		argvalue1 := flag.Arg(2)
		value1 := ds_user_service.PlatformHash(argvalue1)
		arg79 := flag.Arg(3)
		mbTrans80 := thrift.NewTMemoryBufferLen(len(arg79))
		defer mbTrans80.Close()
		_, err81 := mbTrans80.WriteString(arg79)
		if err81 != nil {
			Usage()
			return
		}
		factory82 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt83 := factory82.GetProtocol(mbTrans80)
		containerStruct2 := ds_user_service.NewUserDataServiceGetCookieStateArgs()
		err84 := containerStruct2.ReadField3(jsProt83)
		if err84 != nil {
			Usage()
			return
		}
		argvalue2 := containerStruct2.Context
		value2 := ds_user_service.Context(argvalue2)
		fmt.Print(client.GetCookieState(value0, value1, value2))
		fmt.Print("\n")
		break
	case "setDeviceData":
		if flag.NArg()-1 != 3 {
			fmt.Fprintln(os.Stderr, "SetDeviceData requires 3 args")
			flag.Usage()
		}
		argvalue0 := flag.Arg(1)
		value0 := ds_user_service.DeviceId(argvalue0)
		arg86 := flag.Arg(2)
		mbTrans87 := thrift.NewTMemoryBufferLen(len(arg86))
		defer mbTrans87.Close()
		_, err88 := mbTrans87.WriteString(arg86)
		if err88 != nil {
			Usage()
			return
		}
		factory89 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt90 := factory89.GetProtocol(mbTrans87)
		containerStruct1 := ds_user_service.NewUserDataServiceSetDeviceDataArgs()
		err91 := containerStruct1.ReadField2(jsProt90)
		if err91 != nil {
			Usage()
			return
		}
		argvalue1 := containerStruct1.Data
		value1 := ds_user_service.DeviceData(argvalue1)
		arg92 := flag.Arg(3)
		mbTrans93 := thrift.NewTMemoryBufferLen(len(arg92))
		defer mbTrans93.Close()
		_, err94 := mbTrans93.WriteString(arg92)
		if err94 != nil {
			Usage()
			return
		}
		factory95 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt96 := factory95.GetProtocol(mbTrans93)
		containerStruct2 := ds_user_service.NewUserDataServiceSetDeviceDataArgs()
		err97 := containerStruct2.ReadField3(jsProt96)
		if err97 != nil {
			Usage()
			return
		}
		argvalue2 := containerStruct2.Context
		value2 := ds_user_service.Context(argvalue2)
		fmt.Print(client.SetDeviceData(value0, value1, value2))
		fmt.Print("\n")
		break
	case "delDeviceData":
		if flag.NArg()-1 != 2 {
			fmt.Fprintln(os.Stderr, "DelDeviceData requires 2 args")
			flag.Usage()
		}
		argvalue0 := flag.Arg(1)
		value0 := ds_user_service.DeviceId(argvalue0)
		arg99 := flag.Arg(2)
		mbTrans100 := thrift.NewTMemoryBufferLen(len(arg99))
		defer mbTrans100.Close()
		_, err101 := mbTrans100.WriteString(arg99)
		if err101 != nil {
			Usage()
			return
		}
		factory102 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt103 := factory102.GetProtocol(mbTrans100)
		containerStruct1 := ds_user_service.NewUserDataServiceDelDeviceDataArgs()
		err104 := containerStruct1.ReadField2(jsProt103)
		if err104 != nil {
			Usage()
			return
		}
		argvalue1 := containerStruct1.Context
		value1 := ds_user_service.Context(argvalue1)
		fmt.Print(client.DelDeviceData(value0, value1))
		fmt.Print("\n")
		break
	case "delDeviceDataByKey":
		if flag.NArg()-1 != 3 {
			fmt.Fprintln(os.Stderr, "DelDeviceDataByKey requires 3 args")
			flag.Usage()
		}
		argvalue0 := flag.Arg(1)
		value0 := ds_user_service.DeviceId(argvalue0)
		arg106 := flag.Arg(2)
		mbTrans107 := thrift.NewTMemoryBufferLen(len(arg106))
		defer mbTrans107.Close()
		_, err108 := mbTrans107.WriteString(arg106)
		if err108 != nil {
			Usage()
			return
		}
		factory109 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt110 := factory109.GetProtocol(mbTrans107)
		containerStruct1 := ds_user_service.NewUserDataServiceDelDeviceDataByKeyArgs()
		err111 := containerStruct1.ReadField2(jsProt110)
		if err111 != nil {
			Usage()
			return
		}
		argvalue1 := containerStruct1.Keys
		value1 := ds_user_service.DeviceDataKeys(argvalue1)
		arg112 := flag.Arg(3)
		mbTrans113 := thrift.NewTMemoryBufferLen(len(arg112))
		defer mbTrans113.Close()
		_, err114 := mbTrans113.WriteString(arg112)
		if err114 != nil {
			Usage()
			return
		}
		factory115 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt116 := factory115.GetProtocol(mbTrans113)
		containerStruct2 := ds_user_service.NewUserDataServiceDelDeviceDataByKeyArgs()
		err117 := containerStruct2.ReadField3(jsProt116)
		if err117 != nil {
			Usage()
			return
		}
		argvalue2 := containerStruct2.Context
		value2 := ds_user_service.Context(argvalue2)
		fmt.Print(client.DelDeviceDataByKey(value0, value1, value2))
		fmt.Print("\n")
		break
	case "getDeviceData":
		if flag.NArg()-1 != 2 {
			fmt.Fprintln(os.Stderr, "GetDeviceData requires 2 args")
			flag.Usage()
		}
		argvalue0 := flag.Arg(1)
		value0 := ds_user_service.DeviceId(argvalue0)
		arg119 := flag.Arg(2)
		mbTrans120 := thrift.NewTMemoryBufferLen(len(arg119))
		defer mbTrans120.Close()
		_, err121 := mbTrans120.WriteString(arg119)
		if err121 != nil {
			Usage()
			return
		}
		factory122 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt123 := factory122.GetProtocol(mbTrans120)
		containerStruct1 := ds_user_service.NewUserDataServiceGetDeviceDataArgs()
		err124 := containerStruct1.ReadField2(jsProt123)
		if err124 != nil {
			Usage()
			return
		}
		argvalue1 := containerStruct1.Context
		value1 := ds_user_service.Context(argvalue1)
		fmt.Print(client.GetDeviceData(value0, value1))
		fmt.Print("\n")
		break
	case "isOptedOut":
		if flag.NArg()-1 != 2 {
			fmt.Fprintln(os.Stderr, "IsOptedOut requires 2 args")
			flag.Usage()
		}
		argvalue0 := flag.Arg(1)
		value0 := ds_user_service.DeviceId(argvalue0)
		arg126 := flag.Arg(2)
		mbTrans127 := thrift.NewTMemoryBufferLen(len(arg126))
		defer mbTrans127.Close()
		_, err128 := mbTrans127.WriteString(arg126)
		if err128 != nil {
			Usage()
			return
		}
		factory129 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt130 := factory129.GetProtocol(mbTrans127)
		containerStruct1 := ds_user_service.NewUserDataServiceIsOptedOutArgs()
		err131 := containerStruct1.ReadField2(jsProt130)
		if err131 != nil {
			Usage()
			return
		}
		argvalue1 := containerStruct1.Context
		value1 := ds_user_service.Context(argvalue1)
		fmt.Print(client.IsOptedOut(value0, value1))
		fmt.Print("\n")
		break
	case "getCrossDeviceData":
		if flag.NArg()-1 != 2 {
			fmt.Fprintln(os.Stderr, "GetCrossDeviceData requires 2 args")
			flag.Usage()
		}
		argvalue0 := flag.Arg(1)
		value0 := ds_user_service.StorageId(argvalue0)
		arg133 := flag.Arg(2)
		mbTrans134 := thrift.NewTMemoryBufferLen(len(arg133))
		defer mbTrans134.Close()
		_, err135 := mbTrans134.WriteString(arg133)
		if err135 != nil {
			Usage()
			return
		}
		factory136 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt137 := factory136.GetProtocol(mbTrans134)
		containerStruct1 := ds_user_service.NewUserDataServiceGetCrossDeviceDataArgs()
		err138 := containerStruct1.ReadField2(jsProt137)
		if err138 != nil {
			Usage()
			return
		}
		argvalue1 := containerStruct1.Context
		value1 := ds_user_service.Context(argvalue1)
		fmt.Print(client.GetCrossDeviceData(value0, value1))
		fmt.Print("\n")
		break
	case "":
		Usage()
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
	}
}

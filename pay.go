package pay

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func checkRemote() {
	conn, err := net.Dial("ip:icmp", urlPay)
	if err != nil {
		logs.Critical("支付平台无法访问，修改host文件 将jxpay.com指向指定地址")
		return
	}
	add := conn.RemoteAddr()
	logs.Info(add.String())
}

// 注册函数
func Register(req *RegisterRequest) (*RegisterResponse, error) {
	//checkRemote()
	data, _ := json.Marshal(req)
	resp, err := http.Post(urlPay+":8091"+apiRegister,
		"application/json",
		strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	d := new(RegisterResponse)
	err = json.Unmarshal(body, d)
	logs.Error(err)
	return d, nil
}

// 支付函数 .
func DoPay(r *DoPayRequest) (interface{}, error) {
	if r.Money < 0.01 {
		return "支付金额不能小于0.01", errors.New("支付金额不能小于0.01")
	}
	req, _ := json.Marshal(r)
	resp, err := http.Post(urlPay+":8091"+apiDoPay,
		"application/json",
		strings.NewReader(string(req)))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	d := new(CommonResponse)
	if err = json.Unmarshal(body, d); err != nil {
		logs.Error(err)
		return "", err
	}
	if d.State == "failed" {
		return d.Message, errors.New(d.Message)
	}
	return d.Data, nil
}

// GetIPAddr 获取本机内网地址
func GetIPAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// 支付函数 .
func DoOutPay(r *DoOutPayRequest) (interface{}, error) {
	req, _ := json.Marshal(r)
	resp, err := http.Post(urlPay+":8091"+apiDoOutPay,
		"application/json",
		strings.NewReader(string(req)))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	d := new(CommonResponse)
	if err = json.Unmarshal(body, d); err != nil {
		logs.Error(err)
		return "", err
	}
	if d.State == "failed" {
		return d.Message, errors.New(d.Message)
	}
	return d.Data, nil
}

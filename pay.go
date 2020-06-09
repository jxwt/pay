package pay

import (
	"encoding/json"
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
	checkRemote()
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
	req, _ := json.Marshal(r)
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
	d := new(CommonResponse)
	if err = json.Unmarshal(body, d); err != nil {
		logs.Error(err)
		return "", err
	}
	return d.Message, nil
}

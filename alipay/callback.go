package alipay

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"net/http"
	"sort"
	"strings"
)

func AliWebCallback(w http.ResponseWriter, r *http.Request) (*AliWebPayResult, error) {
	var m = make(map[string]string)
	var signSlice []string
	r.ParseForm()
	for k, v := range r.Form {
		// k不会有多个值的情况
		m[k] = v[0]
		if k == "sign" || k == "sign_type" {
			continue
		}
		signSlice = append(signSlice, fmt.Sprintf("%s=%s", k, v[0]))
	}

	sort.Strings(signSlice)
	signData := strings.Join(signSlice, "&")
	if m["sign_type"] != "RSA" {
		//错误日志
		logs.Error("签名类型未知")
	}

	DefaultAliWebClient().CheckSign(signData, m["sign"])

	var aliPay AliWebPayResult
	err := MapStringToStruct(m, &aliPay)
	if err != nil {
		w.Write([]byte("error"))
		logs.Error(err)
	}

	w.Write([]byte("success"))
	return &aliPay, nil
}

func MapStringToStruct(m map[string]string, i interface{}) error {
	bin, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bin, i)
	if err != nil {
		return err
	}
	return nil
}

// 支付宝app支付回调
func AliAppCallback(w http.ResponseWriter, r *http.Request, getPublicKey func(string) *rsa.PublicKey) (*AliWebPayResult, error) {
	var result string
	defer func() {
		w.Write([]byte(result))
	}()

	var m = make(map[string]string)
	var signSlice []string
	r.ParseForm()
	aliAppClient := new(AliAppClient)
	for k, v := range r.Form {
		m[k] = v[0]
		if k == "sign" || k == "sign_type" {
			continue
		}
		if k == "out_trade_no" {
			aliAppClient.PublicKey = getPublicKey(v[0])
		}
		signSlice = append(signSlice, fmt.Sprintf("%s=%s", k, v[0]))
	}
	sort.Strings(signSlice)
	signData := strings.Join(signSlice, "&")
	if m["sign_type"] != "RSA2" {
		result = "error"
		logs.Warning("签名类型未知")
	}
	if aliAppClient.PublicKey == nil {
		logs.Error("支付宝public key 失效")
		return nil, errors.New("未知支付信息")
	}

	aliAppClient.CheckSign(signData, m["sign"])

	mByte, err := json.Marshal(m)
	logs.Warning(string(mByte))
	if err != nil {
		result = "error"
		logs.Warning(err)
	}

	var aliPay AliWebPayResult
	err = json.Unmarshal(mByte, &aliPay)
	if err != nil {
		result = "error"
		logs.Error("m is %v, err is %v", m, err)
	}
	result = "success"
	return &aliPay, nil
}

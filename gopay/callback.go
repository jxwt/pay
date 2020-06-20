package gopay

import (
	"crypto/rsa"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"encoding/json"
	"github.com/jxwt/pay/gopay/client"
	"github.com/jxwt/pay/gopay/common"
	"github.com/jxwt/pay/gopay/util"
)

func AliWebCallback(w http.ResponseWriter, r *http.Request) (*common.AliWebPayResult, error) {
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

	client.DefaultAliWebClient().CheckSign(signData, m["sign"])

	var aliPay common.AliWebPayResult
	err := util.MapStringToStruct(m, &aliPay)
	if err != nil {
		w.Write([]byte("error"))
		logs.Error(err)
	}

	w.Write([]byte("success"))
	return &aliPay, nil
}

// 支付宝app支付回调
func AliAppCallback(w http.ResponseWriter, r *http.Request, getPublicKey func(string) *rsa.PublicKey) (*common.AliWebPayResult, error) {
	var result string
	defer func() {
		w.Write([]byte(result))
	}()

	var m = make(map[string]string)
	var signSlice []string
	r.ParseForm()
	aliAppClient := new(client.AliAppClient)
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

	var aliPay common.AliWebPayResult
	err = json.Unmarshal(mByte, &aliPay)
	if err != nil {
		result = "error"
		logs.Error("m is %v, err is %v", m, err)
	}
	result = "success"
	return &aliPay, nil
}

// WeChatCallback 微信支付
func WeChatWebCallback(w http.ResponseWriter, r *http.Request) (*common.WeChatPayResult, error) {
	var returnCode = "FAIL"
	var returnMsg = ""
	defer func() {
		formatStr := `<xml><return_code><![CDATA[%s]]></return_code>
                  <return_msg>![CDATA[%s]]</return_msg></xml>`
		returnBody := fmt.Sprintf(formatStr, returnCode, returnMsg)
		w.Write([]byte(returnBody))
	}()
	var reXML common.WeChatPayResult
	//body := cb.Ctx.Input.RequestBody
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//log.Error(string(body))
		returnCode = "FAIL"
		returnMsg = "Bodyerror"
		logs.Error(err)
	}
	err = xml.Unmarshal(body, &reXML)
	if err != nil {
		//log.Error(err, string(body))
		returnMsg = "参数错误"
		returnCode = "FAIL"
		logs.Error(err)
	}

	if reXML.ReturnCode != "SUCCESS" {
		//log.Error(reXML)
		returnCode = "FAIL"
		return &reXML, errors.New(reXML.ReturnCode)
	}
	m := util.XmlToMap(body)

	var signData []string
	for k, v := range m {
		if k == "sign" {
			continue
		}
		signData = append(signData, fmt.Sprintf("%v=%v", k, v))
	}

	key := client.DefaultWechatAppClient().Key

	mySign, err := client.WechatGenSign(key, m)
	if err != nil {
		return &reXML, err
	}

	if mySign != m["sign"] {
		logs.Error(errors.New("签名交易错误"))
	}

	returnCode = "SUCCESS"
	return &reXML, nil
}

func WeChatAppCallback(w http.ResponseWriter, r *http.Request, callback func(string) string) (*common.WeChatPayResult, error) {
	var returnCode = "FAIL"
	var returnMsg = ""
	defer func() {
		formatStr := `<xml><return_code><![CDATA[%s]]></return_code>
                  <return_msg>![CDATA[%s]]</return_msg></xml>`
		returnBody := fmt.Sprintf(formatStr, returnCode, returnMsg)
		w.Write([]byte(returnBody))
	}()
	var reXML common.WeChatPayResult
	//body := cb.Ctx.Input.RequestBody
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//log.Error(string(body))
		returnCode = "FAIL"
		returnMsg = "Bodyerror"
		logs.Error(err)
	}
	err = xml.Unmarshal(body, &reXML)
	if err != nil {
		//log.Error(err, string(body))
		returnMsg = "参数错误"
		returnCode = "FAIL"
		logs.Error(err)
	}

	if reXML.ReturnCode != "SUCCESS" {
		//log.Error(reXML)
		returnCode = "FAIL"
		return &reXML, errors.New(reXML.ReturnCode)
	}

	m := util.XmlToMap(body)
	wxClient := new(client.WechatAppClient)
	for k, v := range m {
		if k == "sign" {
			continue
		}
		if k == "out_trade_no" {
			wxClient.Key = callback(v)
		}
	}
	mySign, err := client.WechatGenSign(wxClient.Key, m)
	if err != nil {
		return &reXML, err
	}
	if mySign != m["sign"] {
		return nil, errors.New("签名交易错误")
	}

	returnCode = "SUCCESS"
	return &reXML, nil
}

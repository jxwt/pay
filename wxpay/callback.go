package wxpay

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"net/http"
)

func WeChatAppCallback(w http.ResponseWriter,body []byte, callback func(string) string) (*WeChatPayResult, error) {
	var returnCode = "FAIL"
	var returnMsg = ""
	defer func() {
		formatStr := `<xml><return_code><![CDATA[%s]]></return_code>
                  <return_msg>![CDATA[%s]]</return_msg></xml>`
		returnBody := fmt.Sprintf(formatStr, returnCode, returnMsg)
		w.Write([]byte(returnBody))
	}()
	var reXML WeChatPayResult
	err := xml.Unmarshal(body, &reXML)
	if err != nil {
		logs.Warning(err)
		returnMsg = "参数错误"
		returnCode = "FAIL"
	}

	if reXML.ReturnCode != "SUCCESS" {
		logs.Error(reXML)
		returnCode = "FAIL"
		return &reXML, errors.New(reXML.ReturnCode)
	}

	m := XmlToMap(body)
	wxClient := new(WxClient)
	var signData []string
	for k, v := range m {
		if k == "sign" {
			continue
		}
		if k == "out_trade_no" {
			wxClient.PayKey = callback(v)
		}
		signData = append(signData, fmt.Sprintf("%v=%v", k, v))
	}

	mySign, err := WechatGenSign(wxClient.PayKey, m)
	if err != nil {
		return &reXML, err
	}

	if mySign != m["sign"] {
		logs.Error("签名交易错误")
	}

	returnCode = "SUCCESS"
	return &reXML, nil
}

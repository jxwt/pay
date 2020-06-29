package wxpay

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
)

func WeChatAppCallback(w http.ResponseWriter, r *http.Request, callback func(string) string) (*WeChatPayResult, error) {
	var returnCode = "FAIL"
	var returnMsg = ""
	defer func() {
		formatStr := `<xml><return_code><![CDATA[%s]]></return_code>
                  <return_msg>![CDATA[%s]]</return_msg></xml>`
		returnBody := fmt.Sprintf(formatStr, returnCode, returnMsg)
		w.Write([]byte(returnBody))
	}()
	var reXML WeChatPayResult
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

	m := XmlToMap(body)
	wxClient := new(WxClient)
	var signData []string
	for k, v := range m {
		if k == "sign" {
			continue
		}
		if k == "out_trade_no" {
			wxClient.Key = callback(v)
		}
		signData = append(signData, fmt.Sprintf("%v=%v", k, v))
	}

	mySign, err := WechatGenSign(wxClient.Key, m)
	if err != nil {
		return &reXML, err
	}

	if mySign != m["sign"] {
		logs.Error("签名交易错误")
	}

	returnCode = "SUCCESS"
	return &reXML, nil
}

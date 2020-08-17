package alipay

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/axgle/mahonia"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"sort"
	"strings"
	"time"
)

//支付宝h5支付
var DefaultAliWapClient *AliWapClient

type AliWapClient struct {
	SellerID   string //合作者ID
	AppID      string // 应用ID
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func InitAliWapClient(c *AliWapClient) {
	DefaultAliWapClient = c
}

// DefaultAliWapClient
func GetDefaultAliWapClient() *AliWapClient {
	return DefaultAliWapClient
}

func (i *AliWapClient) Pay(charge *Charge) (map[string]string, error) {
	return nil, nil
}

func (i *AliWapClient) PayToClient(charge *Charge) (map[string]string, error) {
	return map[string]string{}, errors.New("暂未开发该功能")
}

func (i *AliWapClient) MakePayMap(method string, charge *Charge, rsaType string) (map[string]string, error) {
	var m = make(map[string]string)
	var bizContent = make(map[string]string)
	m["app_id"] = i.AppID
	m["method"] = "alipay.trade.wap.pay"
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["notify_url"] = charge.CallbackURL
	m["sign_type"] = "RSA2"
	// bizContent["subject"] = TruncatedText(charge.Describe, 32)
	//bizContent["quit_url"]="www.baidu.com"

	bizContent["subject"] = "pay"
	bizContent["out_trade_no"] = charge.TradeNum
	//bizContent["quit_url"] = ""
	bizContent["product_code"] = "QUICK_WAP_WAY"
	bizContent["total_amount"] = AliyunMoneyFeeToString(charge.MoneyFee)
	if charge.BuyerId != "" {
		bizContent["buyer_id"] = charge.BuyerId
	}
	bizContentJson, err := json.Marshal(bizContent)
	if err != nil {
		return map[string]string{}, errors.New("json.Marshal: " + err.Error())
	}
	m["biz_content"] = string(bizContentJson)
	if rsaType == "RSA2" {
		m["sign"] = i.GenSign(m)
	} else {
		m["sign"] = i.GenSignRsa1(m)
	}

	logs.Warning(m)
	return m, nil
}

func (i *AliWapClient) MakeH5PayMap(method string, charge *Charge, rsaType string) (string, error) {
	var m = make(map[string]string)
	var bizContent = make(map[string]interface{})
	m["app_id"] = i.AppID
	m["method"] = "alipay.trade.wap.pay"
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["notify_url"] = charge.CallbackURL
	m["sign_type"] = "RSA2"
	if charge.AuthToken != "" {
		m["app_auth_token"] = charge.AuthToken
	}
	bizContent["subject"] = TruncatedText(charge.Describe, 32)
	bizContent["out_trade_no"] = charge.TradeNum
	//bizContent["quit_url"] = ""
	bizContent["product_code"] = "QUICK_WAP_WAY"
	bizContent["total_amount"] = AliyunMoneyFeeToString(charge.MoneyFee)
	if charge.BuyerId != "" {
		bizContent["buyer_id"] = charge.BuyerId
	}
	if charge.ExtendParam != "" {
		p := new(ExtendParam)

		err := json.Unmarshal([]byte(charge.ExtendParam), p)
		if err != nil {
			return "", errors.New("ali pay extend param error")
		}
		bizContent["extend_params"] = p
	}
	bizContentJson, err := json.Marshal(bizContent)
	if err != nil {
		return "", errors.New("json.Marshal: " + err.Error())
	}
	m["biz_content"] = string(bizContentJson)
	if rsaType == "RSA2" {
		m["sign"] = i.GenSign(m)
	} else {
		m["sign"] = i.GenSignRsa1(m)
	}
	logs.Warning(m)
	// 转form表单
	buf := bytes.NewBufferString("")
	for k, v := range m {
		buf.WriteString(fmt.Sprintf(`<input type='hidden' name='%s' value='%s'>`, k, strings.Replace(v, "'", "&apos;", -1)))
	}
	formatStr :=
		`<html>
	<meta http-equiv=Content-Type content="text/html;charset=utf-8">
	<body>
		<form id='paysubmit' name='paysubmit' action='%s' method = 'GET'>
		%s
		<input type='submit' value='ok' style='display:none'>
		</form>
		<script>
		(function(){
			document.forms['paysubmit'].submit();
		})();
		</script>
	</body>
	</html>`
	return fmt.Sprintf(formatStr, "https://openapi.alipay.com/gateway.do?charset=utf-8", buf.String()), nil
}

func (i *AliWapClient) ToPay(charge *Charge) (string, error) {
	payMap, err := i.MakePayMap("alipay.trade.wap.pay", charge, "RSA2")
	if err != nil {
		return "", err
	}
	return i.SendToAlipay(payMap, "post")
}

// ToH5Pay 支付宝h5支付,返回请求参数
func (i *AliWapClient) ToH5Pay(charge *Charge) (string, error) {
	formData, err := i.MakeH5PayMap("alipay.trade.wap.pay", charge, "RSA2")
	if err != nil {
		return "", err
	}
	// fmt.Println(formData)
	return formData, nil
}

func (i *AliWapClient) SendToAlipay(m map[string]string, method string) (string, error) {
	req := httplib.Get("https://openapi.alipay.com/gateway.do")
	if method == "post" {
		req = httplib.Post("https://openapi.alipay.com/gateway.do")
	}
	for k, v := range m {
		req.Param(k, v)
	}
	resp, err := req.Response()
	if err != nil {
		println(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	a:=ConvertToString(string(body),"gbk", "utf-8")
	fmt.Println(a)
	return string(body), err
}


func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

// GenSign 产生签名
func (i *AliWapClient) GenSign(m map[string]string) string {
	var data []string

	for k, v := range m {
		if v != "" && k != "sign" {
			data = append(data, fmt.Sprintf(`%s=%s`, k, v))
		}
	}
	sort.Strings(data)
	signData := strings.Join(data, "&")

	s := sha256.New()
	_, err := s.Write([]byte(signData))
	if err != nil {
		panic(err)
	}
	hashByte := s.Sum(nil)
	signByte, err := i.PrivateKey.Sign(rand.Reader, hashByte, crypto.SHA256)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(signByte)
}

// GenSign 产生签名
func (i *AliWapClient) GenSignRsa1(m map[string]string) string {
	var data []string

	for k, v := range m {
		if v != "" && k != "sign" {
			data = append(data, fmt.Sprintf(`%s=%s`, k, v))
		}
	}
	sort.Strings(data)
	signData := strings.Join(data, "&")

	s := sha1.New()
	_, err := s.Write([]byte(signData))
	if err != nil {
		panic(err)
	}
	hashByte := s.Sum(nil)
	signByte, err := i.PrivateKey.Sign(rand.Reader, hashByte, crypto.SHA1)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(signByte)
}

// NewEncoderToString 将带中文的[]byte 转GB18030字符串
func (c *AliWapClient) NewEncoderToString(req []byte) string {
	reader := bytes.NewReader(req)
	out := transform.NewReader(reader, simplifiedchinese.GB18030.NewEncoder())
	ret, _ := ioutil.ReadAll(out)
	return string(ret)
}

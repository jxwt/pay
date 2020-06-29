package alipay

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	PayTypeAli = 2 // 支付宝支付方式
)

var DefaultAliAppClient *AliAppClient

type AliAppClient struct {
	SellerID   string //合作者ID
	AppID      string // 应用ID
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}


func (i *AliAppClient) MakePayMap(method string, charge *Charge, rsaType string) (map[string]string, error) {
	var m = make(map[string]string)
	var bizContent = make(map[string]string)
	m["app_id"] = i.AppID
	m["method"] = method
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["notify_url"] = charge.CallbackURL
	m["sign_type"] = rsaType
	bizContent["subject"] = TruncatedText(charge.Describe, 32)
	bizContent["out_trade_no"] = charge.TradeNum
	//bizContent["product_code"] = "QUICK_MSECURITY_PAY"
	bizContent["total_amount"] = AliyunMoneyFeeToString(charge.MoneyFee)
	if charge.BuyerId != "" {
		bizContent["buyer_id"] = charge.BuyerId
	}
	if charge.ExtendParam != "" {
		bizContent["extend_params"] = charge.ExtendParam
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

func (i *AliAppClient) MakeRefund(method string, bizContent *AliRefundRequest, rsaType string) (map[string]string, error) {
	var m = make(map[string]string)
	m["app_id"] = i.AppID
	m["method"] = method
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["sign_type"] = rsaType

	bizContentJson, err := json.Marshal(bizContent)
	if err != nil {
		return map[string]string{}, errors.New("json.Marshal: " + err.Error())
	}
	m["biz_content"] = i.NewEncoderToString(bizContentJson)
	if rsaType == "RSA2" {
		m["sign"] = i.GenSign(m)
	} else {
		m["sign"] = i.GenSignRsa1(m)
	}

	logs.Warning(m)
	return m, nil
}

// MakeToaccountTransfer 单比转账请求
func (i *AliAppClient) MakeToaccountTransfer(method string, req *ToaccountTransferRequest, rsaType string) (map[string]string, error) {
	var m = make(map[string]string)
	m["app_id"] = i.AppID
	m["method"] = method
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["sign_type"] = rsaType
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return map[string]string{}, errors.New("json.Marshal: " + err.Error())
	}
	m["biz_content"] = i.NewEncoderToString(reqJSON)
	if rsaType == "RSA2" {
		m["sign"] = i.GenSign(m)
	} else {
		m["sign"] = i.GenSignRsa1(m)
	}
	logs.Warning(m)
	return m, nil
}

func (i *AliAppClient) ToPay(charge *Charge) (string, error) {
	payMap, err := i.MakePayMap("alipay.trade.apps.pay", charge, "RSA")
	if err != nil {
		return "", err
	}
	return i.SendToAlipay(payMap, "post")
}

// 支付宝退款
func (i *AliAppClient) Refund(refund *AliRefundRequest) (*AliRefundResponse, error) {
	payMap, err := i.MakeRefund("alipay.trade.refund", refund, "RSA2")
	if err != nil {
		return nil, err
	}
	response, err := i.SendToAlipay(payMap, "post")
	if err != nil || response == "" {
		return nil, err
	}
	result := new(AliRefundResponse)
	err = json.Unmarshal([]byte(response), result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// ToaccountTransfer 单笔转账到支付宝账户
func (i *AliAppClient) ToaccountTransfer(req *ToaccountTransferRequest) (*ToaccountTransferResponse, error) {
	reqMap, err := i.MakeToaccountTransfer("alipay.fund.trans.toaccount.transfer", req, "RSA")
	if err != nil {
		return nil, err
	}
	response, err := i.SendToAlipay(reqMap, "post")
	if err != nil || response == "" {
		return nil, err
	}
	result := new(ToaccountTransferResponse)
	err = json.Unmarshal([]byte(response), result)
	if err != nil {
		return result, err
	}
	return result, nil
}

/**
获取APP支付的链接码
*/
func (i *AliAppClient) AppPay(charge *Charge) (string, error) {
	payMap, err := i.MakePayMap("alipay.trade.app.pay", charge, "RSA2")
	if err != nil {
		return "", err
	}
	return i.ToURL(payMap), nil
}

// ToH5Pay 支付宝h5支付,返回请求参数
func (i *AliAppClient) ToH5Pay(charge *Charge) (string, error) {
	formData, err := i.MakeH5PayMap(charge, "RSA2")
	if err != nil {
		return "", err
	}
	// fmt.Println(formData)
	return formData, nil
}

func (i *AliAppClient) CreateOrder(charge *Charge) (string, error) {
	payMap, err := i.MakePayMap("alipay.trade.create", charge, "RSA2")
	if err != nil {
		return "", err
	}
	return i.SendToAlipay(payMap, "get")
}

func (i *AliAppClient) Login(code string) (string, error) {
	var m = make(map[string]string)
	m["app_id"] = i.AppID
	m["method"] = "alipay.systemsService.oauth.token"
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["sign_type"] = "RSA2"
	m["grant_type"] = "authorization_code"
	m["code"] = code
	m["sign"] = i.GenSign(m)
	return i.SendToAlipay(m, "post")
}

func (i *AliAppClient) GetLoginUserInfo(authToken string) (string, error) {
	var m = make(map[string]string)
	m["app_id"] = i.AppID
	m["method"] = "alipay.user.info.share"
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["sign_type"] = "RSA2"
	m["auth_token"] = authToken
	m["version"] = "1.0"
	m["sign"] = i.GenSign(m)
	return i.SendToAlipay(m, "post")
}

func (i *AliAppClient) GetAppLoginParams(targetId string) string {
	m := map[string]string{
		"apiname":    "com.alipay.account.auth",
		"app_id":     i.AppID,
		"app_name":   "mc",
		"auth_type":  "AUTHACCOUNT",
		"biz_type":   "openservice",
		"method":     "alipay.open.auth.sdk.code.get",
		"pid":        i.SellerID,
		"product_id": "APP_FAST_LOGIN",
		"scope":      "kuaijie",
		"target_id":  targetId,
		"sign_type":  "RSA2",
	}
	m["sign"] = i.GenSign(m)
	var data []string
	for k, v := range m {
		if v != "" {
			data = append(data, fmt.Sprintf(`%s=%s`, k, v))
		}
	}
	sort.Strings(data)
	return strings.Join(data, "&")
}

func (i *AliAppClient) SendToAlipay(m map[string]string, method string) (string, error) {
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
	return string(body), err
}



// 退款查询
func (i *AliAppClient) QueryRefund(outTradeNo string) (*AliRefundResponse, error) {
	var m = make(map[string]string)
	m["method"] = "alipay.trade.fastpay.refund.query"
	m["app_id"] = i.AppID
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["sign_type"] = "RSA"
	bizContent := map[string]string{"out_trade_no": outTradeNo, "out_request_no": "AliPay" + outTradeNo}
	bizContentJson, err := json.Marshal(bizContent)
	if err != nil {
		return nil, errors.New("json.Marshal: " + err.Error())
	}
	m["biz_content"] = string(bizContentJson)
	sign := i.GenSignRsa1(m)
	m["sign"] = sign

	resp, err := i.SendToAlipay(m, "post")
	if err != nil {
		return nil, err
	}
	result := new(AliRefundResponse)
	err = json.Unmarshal([]byte(resp), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 订单查询
func (i *AliAppClient) QueryOrder(outTradeNo string) (*AliWebAppQueryResult, error) {
	var m = make(map[string]string)
	m["method"] = "alipay.trade.query"
	m["app_id"] = i.AppID
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["sign_type"] = "RSA2"
	bizContent := map[string]string{"out_trade_no": outTradeNo}
	bizContentJson, err := json.Marshal(bizContent)
	if err != nil {
		return &AliWebAppQueryResult{}, errors.New("json.Marshal: " + err.Error())
	}
	m["biz_content"] = i.NewEncoderToString(bizContentJson)
	m["sign"] = i.GenSign(m)
	response, err := i.SendToAlipay(m, "post")
	if err != nil || response == "" {
		return nil, err
	}
	result := new(AliWebAppQueryResult)
	err = json.Unmarshal([]byte(response), result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (i *AliAppClient) AliPreCreate(preCreate Charge) (PreCreateResult, error) {
	preCreateResult := new(PreCreateResponse)
	payMap, err := i.MakePayMap("alipay.trade.precreate", &preCreate, "RSA2")
	if err != nil {
		return preCreateResult.PreCreateResult, errors.New("json.Marshal: " + err.Error())
	}
	response, err := i.SendToAlipay(payMap, "get")
	if err != nil || response == "" {
		return preCreateResult.PreCreateResult, err
	}
	err = json.Unmarshal([]byte(response), preCreateResult)
	if err != nil {
		return preCreateResult.PreCreateResult, err
	} else if preCreateResult.PreCreateResult.SubMsg != "" {
		return preCreateResult.PreCreateResult, errors.New("支付宝下单失败")
	}
	return preCreateResult.PreCreateResult, nil
}

// GenSign 产生签名
func (i *AliAppClient) GenSignRsa1(m map[string]string) string {
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

// GenSign 产生签名
func (i *AliAppClient) GenSign(m map[string]string) string {
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

// CheckSign 检测签名
func (i *AliAppClient) CheckSign(signData, sign string) {
	signByte, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		panic(err)
	}
	s := sha256.New()
	_, err = s.Write([]byte(signData))
	if err != nil {
		//panic(err)
		logs.Warning(err)
	}
	hash := s.Sum(nil)
	logs.Warning(i.PublicKey)
	err = rsa.VerifyPKCS1v15(i.PublicKey, crypto.SHA256, hash, signByte)
	if err != nil {
		logs.Warning(err)
	}
}

// ToURL
func (i *AliAppClient) ToURL(m map[string]string) string {
	var buf []string
	for k, v := range m {
		buf = append(buf, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}
	return strings.Join(buf, "&")
}

// NewEncoderToString 将带中文的[]byte 转GB18030字符串
func (i *AliAppClient) NewEncoderToString(req []byte) string {
	reader := bytes.NewReader(req)
	out := transform.NewReader(reader, simplifiedchinese.GB18030.NewEncoder())
	ret, _ := ioutil.ReadAll(out)
	return string(ret)
}

// ParsePKCS1PrivateKey 解析私钥
func ParsePKCS1PrivateKey(data []byte) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, errors.New("private key error")
	}

	keyTemp, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key = keyTemp.(*rsa.PrivateKey)
	return key, err
}

// ParsePKCS1PublicKey 解析公钥
func ParsePKCS1PublicKey(data []byte) (key *rsa.PublicKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, errors.New("public key error")
	}

	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key error")
	}

	return key, err
}

// FormatPublicKey 格式化公钥
func FormatPublicKey(raw string) (result []byte) {
	return formatKey(raw, "-----BEGIN PUBLIC KEY-----", "-----END PUBLIC KEY-----")
}

// FormatPrivateKey 格式化私钥
func FormatPrivateKey(raw string) (result []byte) {
	return formatKey(raw, "-----BEGIN RSA PRIVATE KEY-----", "-----END RSA PRIVATE KEY-----")
}

func formatKey(raw, prefix, suffix string) (result []byte) {
	if raw == "" {
		return nil
	}
	raw = strings.Replace(raw, prefix, "", 1)
	raw = strings.Replace(raw, suffix, "", 1)
	raw = strings.Replace(raw, " ", "", -1)
	raw = strings.Replace(raw, "\n", "", -1)
	raw = strings.Replace(raw, "\r", "", -1)
	raw = strings.Replace(raw, "\t", "", -1)

	var ll = 64
	var sl = len(raw)
	var c = sl / ll
	if sl%ll > 0 {
		c = c + 1
	}

	var buf bytes.Buffer
	buf.WriteString(prefix + "\n")
	for i := 0; i < c; i++ {
		var b = i * ll
		var e = b + ll
		if e > sl {
			buf.WriteString(raw[b:])
		} else {
			buf.WriteString(raw[b:e])
		}
		buf.WriteString("\n")
	}
	buf.WriteString(suffix)
	return buf.Bytes()
}

// AliTradePay 支付宝统一收单
func (i *AliAppClient) AliTradePay(aliTradePay *AliTradePayRequest) (*AliTradePayResponse, error) {
	payMap, err := i.MakeTradePay("alipay.trade.pay", aliTradePay, "RSA2")
	if err != nil {
		return nil, err
	}
	response, err := i.SendToAlipay(payMap, "post")
	if err != nil || response == "" {
		return nil, err
	}
	result := new(AliTradePayResponse)
	err = json.Unmarshal([]byte(response), result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// MakeTradePay 创建支付宝统一收单请求
func (i *AliAppClient) MakeTradePay(method string, bizContent *AliTradePayRequest, rsaType string) (map[string]string, error) {
	var m = make(map[string]string)
	m["app_id"] = i.AppID
	m["method"] = method
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["sign_type"] = rsaType

	bizContentJson, err := json.Marshal(bizContent)
	if err != nil {
		return map[string]string{}, errors.New("json.Marshal: " + err.Error())
	}
	m["biz_content"] = i.NewEncoderToString(bizContentJson)
	if rsaType == "RSA2" {
		m["sign"] = i.GenSign(m)
	} else {
		m["sign"] = i.GenSignRsa1(m)
	}

	logs.Warning(m)
	return m, nil
}

// AliTradeCancel 支付撤单
func (i *AliAppClient) AliTradeCancel(aliTradePay *AliTradeCancelRequest) (*AliTradeCancelResponse, error) {
	payMap, err := i.MakeTradeCancel("alipay.trade.cancel", aliTradePay, "RSA2")
	if err != nil {
		return nil, err
	}
	response, err := i.SendToAlipay(payMap, "post")
	if err != nil || response == "" {
		return nil, err
	}
	result := new(AliTradeCancelResponse)
	err = json.Unmarshal([]byte(response), result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// MakeTradePay 创建支付宝统一收单请求
func (i *AliAppClient) MakeTradeCancel(method string, bizContent *AliTradeCancelRequest, rsaType string) (map[string]string, error) {
	var m = make(map[string]string)
	m["app_id"] = i.AppID
	m["method"] = method
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["sign_type"] = rsaType

	bizContentJson, err := json.Marshal(bizContent)
	if err != nil {
		return map[string]string{}, errors.New("json.Marshal: " + err.Error())
	}
	m["biz_content"] = i.NewEncoderToString(bizContentJson)
	if rsaType == "RSA2" {
		m["sign"] = i.GenSign(m)
	} else {
		m["sign"] = i.GenSignRsa1(m)
	}

	logs.Warning(m)
	return m, nil
}


func (i *AliAppClient) MakeH5PayMap(charge *Charge, rsaType string) (string, error) {
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
	// bizContent["subject"] = TruncatedText(charge.Describe, 32)
	bizContent["subject"] = "test"
	bizContent["out_trade_no"] = charge.TradeNum
	//bizContent["quit_url"] = ""
	bizContent["product_code"] = "QUICK_WAP_WAY"
	bizContent["total_amount"] = AliyunMoneyFeeToString(charge.MoneyFee)
	if charge.BuyerId != "" {
		bizContent["buyer_id"] = charge.BuyerId
	}
	if charge.ExtendParam != "" {
		bizContent["extend_params"] = charge.ExtendParam
	}
	//if charge.IndustryRefluxInfo != nil {
	//	d, _ := json.Marshal(charge.IndustryRefluxInfo)
	//	extern := &ExtendParam{
	//		SysServiceProviderId: "2088521066336121",
	//		IndustryRefluxInfo:   string(d),
	//	}
	//
	//}
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

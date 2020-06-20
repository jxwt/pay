package client

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
	"github.com/jxwt/pay/gopay/common"
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

func InitAliAppClient(c *AliAppClient) {
	DefaultAliAppClient = c
}

// DefaultAliAppClient 得到默认支付宝app客户端
func GetDefaultAliAppClient() *AliAppClient {
	return DefaultAliAppClient
}

func (this *AliAppClient) MakePayMap(method string, charge *common.Charge, rsaType string) (map[string]string, error) {
	var m = make(map[string]string)
	var bizContent = make(map[string]string)
	m["app_id"] = this.AppID
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
	bizContentJson, err := json.Marshal(bizContent)
	if err != nil {
		return map[string]string{}, errors.New("json.Marshal: " + err.Error())
	}
	m["biz_content"] = string(bizContentJson)
	if rsaType == "RSA2" {
		m["sign"] = this.GenSign(m)
	} else {
		m["sign"] = this.GenSignRsa1(m)
	}

	logs.Warning(m)
	return m, nil
}

func (this *AliAppClient) MakeRefund(method string, bizContent *common.AliRefundRequest, rsaType string) (map[string]string, error) {
	var m = make(map[string]string)
	m["app_id"] = this.AppID
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
	m["biz_content"] = this.NewEncoderToString(bizContentJson)
	if rsaType == "RSA2" {
		m["sign"] = this.GenSign(m)
	} else {
		m["sign"] = this.GenSignRsa1(m)
	}

	logs.Warning(m)
	return m, nil
}

// MakeToaccountTransfer 单比转账请求
func (c *AliAppClient) MakeToaccountTransfer(method string, req *common.ToaccountTransferRequest, rsaType string) (map[string]string, error) {
	var m = make(map[string]string)
	m["app_id"] = c.AppID
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
	m["biz_content"] = c.NewEncoderToString(reqJSON)
	if rsaType == "RSA2" {
		m["sign"] = c.GenSign(m)
	} else {
		m["sign"] = c.GenSignRsa1(m)
	}
	logs.Warning(m)
	return m, nil
}

func (this *AliAppClient) Pay(charge *common.Charge) (map[string]string, error) {
	return nil, nil
}
func (this *AliAppClient) ToPay(charge *common.Charge) (string, error) {
	payMap, err := this.MakePayMap("alipay.trade.apps.pay", charge, "RSA")
	if err != nil {
		return "", err
	}
	return this.SendToAlipay(payMap, "post")
}

// 支付宝退款
func (this *AliAppClient) Refund(refund *common.AliRefundRequest) (*common.AliRefundResponse, error) {
	payMap, err := this.MakeRefund("alipay.trade.refund", refund, "RSA2")
	if err != nil {
		return nil, err
	}
	response, err := this.SendToAlipay(payMap, "post")
	if err != nil || response == "" {
		return nil, err
	}
	result := new(common.AliRefundResponse)
	err = json.Unmarshal([]byte(response), result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// ToaccountTransfer 单笔转账到支付宝账户
func (c *AliAppClient) ToaccountTransfer(req *common.ToaccountTransferRequest) (*common.ToaccountTransferResponse, error) {
	reqMap, err := c.MakeToaccountTransfer("alipay.fund.trans.toaccount.transfer", req, "RSA")
	if err != nil {
		return nil, err
	}
	response, err := c.SendToAlipay(reqMap, "post")
	if err != nil || response == "" {
		return nil, err
	}
	result := new(common.ToaccountTransferResponse)
	err = json.Unmarshal([]byte(response), result)
	if err != nil {
		return result, err
	}
	return result, nil
}

/**
获取APP支付的链接码
*/
func (this *AliAppClient) AppPay(charge *common.Charge) (string, error) {
	payMap, err := this.MakePayMap("alipay.trade.app.pay", charge, "RSA2")
	if err != nil {
		return "", err
	}
	return this.ToURL(payMap), nil
}

func (this *AliAppClient) CreateOrder(charge *common.Charge) (string, error) {
	payMap, err := this.MakePayMap("alipay.trade.create", charge, "RSA2")
	if err != nil {
		return "", err
	}
	return this.SendToAlipay(payMap, "get")
}

func (this *AliAppClient) Login(code string) (string, error) {
	var m = make(map[string]string)
	m["app_id"] = this.AppID
	m["method"] = "alipay.systemsService.oauth.token"
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["sign_type"] = "RSA2"
	m["grant_type"] = "authorization_code"
	m["code"] = code
	m["sign"] = this.GenSign(m)
	fmt.Println(m)
	return this.SendToAlipay(m, "post")
}

func (this *AliAppClient) GetLoginUserInfo(authToken string) (string, error) {
	var m = make(map[string]string)
	m["app_id"] = this.AppID
	m["method"] = "alipay.user.info.share"
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["sign_type"] = "RSA2"
	m["auth_token"] = authToken
	m["version"] = "1.0"
	m["sign"] = this.GenSign(m)
	return this.SendToAlipay(m, "post")
}

func (this *AliAppClient) GetAppLoginParams(targetId string) string {
	m := map[string]string{
		"apiname":    "com.alipay.account.auth",
		"app_id":     this.AppID,
		"app_name":   "mc",
		"auth_type":  "AUTHACCOUNT",
		"biz_type":   "openservice",
		"method":     "alipay.open.auth.sdk.code.get",
		"pid":        this.SellerID,
		"product_id": "APP_FAST_LOGIN",
		"scope":      "kuaijie",
		"target_id":  targetId,
		"sign_type":  "RSA2",
	}
	m["sign"] = this.GenSign(m)
	var data []string
	for k, v := range m {
		if v != "" {
			data = append(data, fmt.Sprintf(`%s=%s`, k, v))
		}
	}
	sort.Strings(data)
	return strings.Join(data, "&")
}

func (this *AliAppClient) SendToAlipay(m map[string]string, method string) (string, error) {
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

func (this *AliAppClient) PayToClient(charge *common.Charge) (map[string]string, error) {

	return map[string]string{}, errors.New("暂未开发该功能")
}

// 退款查询
func (this *AliAppClient) QueryRefund(outTradeNo string) (*common.AliRefundResponse, error) {
	var m = make(map[string]string)
	m["method"] = "alipay.trade.fastpay.refund.query"
	m["app_id"] = this.AppID
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
	sign := this.GenSignRsa1(m)
	m["sign"] = sign

	resp, err := this.SendToAlipay(m, "post")
	if err != nil {
		return nil, err
	}
	result := new(common.AliRefundResponse)
	err = json.Unmarshal([]byte(resp), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 订单查询
func (this *AliAppClient) QueryOrder(outTradeNo string) (*common.AliWebAppQueryResult, error) {
	var m = make(map[string]string)
	m["method"] = "alipay.trade.query"
	m["app_id"] = this.AppID
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["sign_type"] = "RSA2"
	bizContent := map[string]string{"out_trade_no": outTradeNo}
	bizContentJson, err := json.Marshal(bizContent)
	if err != nil {
		return &common.AliWebAppQueryResult{}, errors.New("json.Marshal: " + err.Error())
	}
	m["biz_content"] = this.NewEncoderToString(bizContentJson)
	m["sign"] = this.GenSign(m)
	response, err := this.SendToAlipay(m, "post")
	if err != nil || response == "" {
		return nil, err
	}
	result := new(common.AliWebAppQueryResult)
	err = json.Unmarshal([]byte(response), result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (this *AliAppClient) AliPreCreate(preCreate common.Charge) (common.PreCreateResult, error) {
	preCreateResult := new(common.PreCreateResponse)
	payMap, err := this.MakePayMap("alipay.trade.precreate", &preCreate, "RSA2")
	if err != nil {
		return preCreateResult.PreCreateResult, errors.New("json.Marshal: " + err.Error())
	}
	response, err := this.SendToAlipay(payMap, "get")
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
func (this *AliAppClient) GenSignRsa1(m map[string]string) string {
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
	signByte, err := this.PrivateKey.Sign(rand.Reader, hashByte, crypto.SHA1)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(signByte)
}

// GenSign 产生签名
func (this *AliAppClient) GenSign(m map[string]string) string {
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
	signByte, err := this.PrivateKey.Sign(rand.Reader, hashByte, crypto.SHA256)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(signByte)
}

// CheckSign 检测签名
func (this *AliAppClient) CheckSign(signData, sign string) {
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
	logs.Warning(this.PublicKey)
	err = rsa.VerifyPKCS1v15(this.PublicKey, crypto.SHA256, hash, signByte)
	if err != nil {
		logs.Warning(err)
	}
}

// ToURL
func (this *AliAppClient) ToURL(m map[string]string) string {
	var buf []string
	for k, v := range m {
		buf = append(buf, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}
	return strings.Join(buf, "&")
}

// NewEncoderToString 将带中文的[]byte 转GB18030字符串
func (c *AliAppClient) NewEncoderToString(req []byte) string {
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
func (c *AliAppClient) AliTradePay(aliTradePay *common.AliTradePayRequest) (*common.AliTradePayResponse, error) {
	payMap, err := c.MakeTradePay("alipay.trade.pay", aliTradePay, "RSA2")
	if err != nil {
		return nil, err
	}
	response, err := c.SendToAlipay(payMap, "post")
	if err != nil || response == "" {
		return nil, err
	}
	result := new(common.AliTradePayResponse)
	err = json.Unmarshal([]byte(response), result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// MakeTradePay 创建支付宝统一收单请求
func (c *AliAppClient) MakeTradePay(method string, bizContent *common.AliTradePayRequest, rsaType string) (map[string]string, error) {
	var m = make(map[string]string)
	m["app_id"] = c.AppID
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
	m["biz_content"] = c.NewEncoderToString(bizContentJson)
	if rsaType == "RSA2" {
		m["sign"] = c.GenSign(m)
	} else {
		m["sign"] = c.GenSignRsa1(m)
	}

	logs.Warning(m)
	return m, nil
}

// AliTradeCancel 支付撤单
func (c *AliAppClient) AliTradeCancel(aliTradePay *common.AliTradeCancelRequest) (*common.AliTradeCancelResponse, error) {
	payMap, err := c.MakeTradeCancel("alipay.trade.cancel", aliTradePay, "RSA2")
	if err != nil {
		return nil, err
	}
	response, err := c.SendToAlipay(payMap, "post")
	if err != nil || response == "" {
		return nil, err
	}
	result := new(common.AliTradeCancelResponse)
	err = json.Unmarshal([]byte(response), result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// MakeTradePay 创建支付宝统一收单请求
func (c *AliAppClient) MakeTradeCancel(method string, bizContent *common.AliTradeCancelRequest, rsaType string) (map[string]string, error) {
	var m = make(map[string]string)
	m["app_id"] = c.AppID
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
	m["biz_content"] = c.NewEncoderToString(bizContentJson)
	if rsaType == "RSA2" {
		m["sign"] = c.GenSign(m)
	} else {
		m["sign"] = c.GenSignRsa1(m)
	}

	logs.Warning(m)
	return m, nil
}

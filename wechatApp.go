package pay

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jxwt/tools"
	"strconv"
	"strings"
	"time"
)

const (
	PayTypeWx = 1 // 微信支付方式
)

const (
	WxPayRefundURL = "https://api.mch.weixin.qq.com/secapi/pay/refund" // WxPayRefundURL 微信退款url
)

var defaultWechatAppClient *WechatAppClient

func InitWxAppClient(c *WechatAppClient) {
	defaultWechatAppClient = c
}

// DefaultWechatAppClient 默认微信app客户端
func DefaultWechatAppClient() *WechatAppClient {
	return defaultWechatAppClient
}

// WechatAppClient 微信app支付
type WechatAppClient struct {
	AppID       string       // 公众账号ID
	MchID       string       // 商户号ID
	Key         string       // 密钥
	PrivateKey  []byte       // 私钥文件内容
	PublicKey   []byte       // 公钥文件内容
	httpsClient *HTTPSClient // 双向证书链接
}

type WxPayTool struct {
	AppID       string
	MchID       string
	Key         string
	PrivateKey  []byte
	CallbackURL string

	CertPEM string // cert证书
	KeyPEM  string // 密钥证书
}

// Pay 支付
func (this *WechatAppClient) Pay(charge *Charge) (map[string]string, error) {
	result, err := this.WxUnifiedOrder(charge, "APP")
	if err != nil {
		return map[string]string{}, errors.New("WechatApp.sign: " + err.Error())
	}
	var c = make(map[string]string)
	c["appid"] = this.AppID
	c["partnerid"] = this.MchID
	c["prepayid"] = result.PrepayID
	c["package"] = "Sign=WXPay"
	c["noncestr"] = RandomStr()
	c["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
	sign2, err := WechatGenSign(this.Key, c)
	if err != nil {
		return map[string]string{}, errors.New("WechatApp.paySign: " + err.Error())
	}
	c["paySign"] = strings.ToUpper(sign2)
	return c, nil
}

// 支付到用户的微信账号
func (this *WechatAppClient) PayToClient(charge *Charge) (map[string]string, error) {
	return WachatCompanyChange(this.AppID, this.MchID, this.Key, this.httpsClient, charge)
}

// QueryOrder 查询订单
func (this *WechatAppClient) QueryOrder(tradeNum string) (WeChatQueryResult, error) {
	var m = make(map[string]string)
	m["appid"] = this.AppID
	m["mch_id"] = this.MchID
	m["out_trade_no"] = tradeNum
	m["nonce_str"] = RandomStr()

	sign, err := WechatGenSign(this.Key, m)
	if err != nil {
		return WeChatQueryResult{}, err
	}

	m["sign"] = sign

	return PostWechat("https://api.mch.weixin.qq.com/pay/orderquery", m, nil)
}

func (this *WechatAppClient) WxUnifiedOrder(charge *Charge, tradeType string) (WeChatQueryResult, error) {
	result := new(WeChatQueryResult)
	var m = make(map[string]string)
	m["appid"] = this.AppID
	m["mch_id"] = this.MchID
	m["nonce_str"] = RandomStr()
	m["body"] = TruncatedText(charge.Describe, 32)
	m["out_trade_no"] = charge.TradeNum
	m["total_fee"] = WechatMoneyFeeToString(charge.MoneyFee)
	m["spbill_create_ip"] = tools.GetLocalAddr()
	m["notify_url"] = charge.CallbackURL
	m["trade_type"] = tradeType
	m["sign_type"] = "MD5"
	sign, err := WechatGenSign(this.Key, m)
	if err != nil {
		return *result, errors.New("WechatApp.sign: " + err.Error())
	}
	m["sign"] = sign
	*result, err = PostWechat("https://api.mch.weixin.qq.com/pay/unifiedorder", m, nil)
	if err != nil {
		return *result, err
	}
	return *result, err
}

// WechatCallBackSuccessRes 构建微信回调成功返回
func WechatCallBackSuccessRes() string {
	var m = make(map[string]string)
	m["return_code"] = "SUCCESS"
	m["return_msg"] = "OK"
	buf := bytes.NewBufferString("")
	for k, v := range m {
		buf.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k))
	}
	xmlStr := fmt.Sprintf("<xml>%s</xml>", buf.String())
	return xmlStr
}

// MicroPay 微信付款码支付
// OutRefundNo 为后端自定义的随机字符串（尽量唯一） 与 商户退款单号（确保唯一性）
// TotalFee 订单的金额
// AuthCode 用户的授权码(条形码)
func (i *WechatAppClient) MicroPay(req *MicroPayRequest) (*WeChatQueryResult, error) {
	var m = make(map[string]string)
	m["appid"] = i.AppID
	m["mch_id"] = i.MchID
	m["nonce_str"] = RandomStr()
	m["body"] = req.Remark
	m["out_trade_no"] = req.OutTradeNo
	m["total_fee"] = strconv.Itoa(req.TotalFee)
	m["spbill_create_ip"] = tools.GetLocalAddr()
	m["auth_code"] = req.AuthCode

	sign, err := WechatGenSign(i.Key, m)
	if err != nil {
		return nil, errors.New("MicroPay.sign: " + err.Error())
	}

	m["sign"] = sign

	xmlRe, err := PostWechat("https://api.mch.weixin.qq.com/pay/micropay", m, nil)
	if err != nil {
		return &xmlRe, err
	}
	return &xmlRe, nil
}

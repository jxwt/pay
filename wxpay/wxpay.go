package wxpay

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jxwt/pay"
	"github.com/jxwt/tools"
	"strconv"
	"strings"
	"time"
)

type WxClient struct {
	AppID       string
	MchID       string
	SecretKey   string  // 登录秘钥
	PayKey  	string 	// 支付加签秘钥
	CallbackURL string
	SubMchId    string

	CertPEM string // cert证书
	KeyPEM  string // 密钥证书

	httpsClient *pay.HTTPSClient // 双向证书链接
	KeyPemNo    string

	Ciphertext string // 敏感信息加密使用的证书
	SerialNo   string // 敏感信息加密使用的证书号
}

func InitWxClient(AppID string, MchID string, SecretKey string, PayKey string, CallbackURL string, subMchId ...string) *WxClient {
	c := &WxClient{
		AppID:       AppID,
		MchID:       MchID,
		SecretKey:   SecretKey,
		PayKey:  	PayKey,
		CallbackURL: CallbackURL,
		httpsClient: nil,
	}
	if len(subMchId) > 0 {
		c.SubMchId = subMchId[0]
	}

	return c
}

// app支付
func (i *WxClient) AppPay(charge *Charge) (map[string]string, error) {
	result, err := i.WxUnifiedOrder(charge, "APP")
	if err != nil {
		return map[string]string{}, errors.New("wx app pay" + err.Error())
	}
	var c = make(map[string]string)
	c["appid"] = i.AppID
	c["partnerid"] = i.MchID
	c["prepayid"] = result.PrepayID
	c["package"] = "Sign=WXPay"
	c["noncestr"] = RandomStr()
	c["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
	//c["sign_type"] = "MD5"
	sign2, err := WechatGenSign(i.PayKey, c)
	if err != nil {
		return map[string]string{}, errors.New("wx app pay" + err.Error())
	}
	c["paySign"] = strings.ToUpper(sign2)
	return c, nil
}

// H5支付
func (i *WxClient) H5Pay(charge *Charge) (map[string]string, error) {
	result, err := i.WxUnifiedOrder(charge, "MWEB")
	if err != nil {
		return map[string]string{}, errors.New("wx app pay" + err.Error())
	}
	var c = make(map[string]string)
	c["appId"] = i.AppID
	c["timeStamp"] = fmt.Sprintf("%d", time.Now().Unix())
	c["nonceStr"] = RandomStr()
	c["package"] = fmt.Sprintf("prepay_id=%s", result.PrepayID)
	c["signType"] = "MD5"
	sign2, err := WechatGenSign(i.PayKey, c)
	if err != nil {
		return map[string]string{}, errors.New("WechatH5: " + err.Error())
	}
	c["paySign"] = sign2
	c["mweb_url"] = result.MwebUrl
	//delete(c, "appId")
	return c, nil
}

// 小程序 或者公众号支付
func (i *WxClient) MiniPay(charge *Charge) (map[string]string, error) {
	result, err := i.WxUnifiedOrder(charge, "JSAPI")
	if err != nil {
		return map[string]string{}, errors.New("wx app pay" + err.Error())
	}
	var c = make(map[string]string)
	c["appId"] = i.AppID
	c["timeStamp"] = fmt.Sprintf("%d", time.Now().Unix())
	c["nonceStr"] = RandomStr()
	c["package"] = fmt.Sprintf("prepay_id=%s", result.PrepayID)
	c["signType"] = "MD5"
	sign2, err := WechatGenSign(i.PayKey, c)
	if err != nil {
		return map[string]string{}, errors.New("WechatWeb: " + err.Error())
	}
	c["paySign"] = sign2
	logs.Warning("Pay res:", c)
	return c, nil
}

// 生成支付二维码信息
func (i *WxClient) WxNative(charge *Charge) (WeChatQueryResult, error) {
	result, err := i.WxUnifiedOrder(charge, "NATIVE")
	if err != nil {
		return result, err
	}
	return result, nil
}

func (i *WxClient) WxUnifiedOrder(charge *Charge, tradeType string) (WeChatQueryResult, error) {
	result := new(WeChatQueryResult)
	var m = make(map[string]string)
	m["appid"] = i.AppID
	m["mch_id"] = i.MchID
	m["nonce_str"] = RandomStr()
	m["body"] = TruncatedText(charge.Describe, 32)
	m["out_trade_no"] = charge.TradeNum
	m["total_fee"] = WechatMoneyFeeToString(charge.MoneyFee)
	m["spbill_create_ip"] = tools.GetLocalAddr()
	m["notify_url"] = i.CallbackURL
	m["trade_type"] = tradeType
	m["sign_type"] = "MD5"
	if i.SubMchId != "" {
		m["sub_mch_id"] = i.SubMchId
	}
	if charge.OpenID != "" {
		m["openid"] = charge.OpenID
	}
	if tradeType == "NWEB" {
		m["scene_info"] = charge.SceneInfo
	}
	// app H5支付需要附加app场景信息
	if charge.BundleId != "" {
		m["scene_info"] = fmt.Sprintf(`{"h5_info": {"type":"%s","app_name": "%s","bundle_id": "%s"}`, charge.AppType, charge.AppType, charge.BundleId)
	} else if charge.PackageName != "" {
		m["scene_info"] = fmt.Sprintf(`{"h5_info": {"type":"%s","app_name": "%s","package_name": "%s"}`, charge.AppType, charge.AppType, charge.PackageName)
	}
	sign, err := WechatGenSign(i.PayKey, m)
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

// QueryOrder 查询订单
func (i *WxClient) QueryOrder(tradeNum string) (WeChatQueryResult, error) {
	var m = make(map[string]string)
	m["appid"] = i.AppID
	m["mch_id"] = i.MchID
	m["out_trade_no"] = tradeNum
	m["nonce_str"] = RandomStr()

	sign, err := WechatGenSign(i.PayKey, m)
	if err != nil {
		return WeChatQueryResult{}, err
	}

	m["sign"] = sign

	return PostWechat("https://api.mch.weixin.qq.com/pay/orderquery", m, nil)
}

// MicroPay 微信付款码支付
// OutRefundNo 为后端自定义的随机字符串（尽量唯一） 与 商户退款单号（确保唯一性）
// TotalFee 订单的金额
// AuthCode 用户的授权码(条形码)
func (i *WxClient) MicroPay(req *MicroPayRequest) (*WeChatQueryResult, error) {
	var m = make(map[string]string)
	m["appid"] = i.AppID
	m["mch_id"] = i.MchID
	m["nonce_str"] = RandomStr()
	m["body"] = req.Remark
	m["out_trade_no"] = req.OutTradeNo
	m["total_fee"] = strconv.Itoa(req.TotalFee)
	m["spbill_create_ip"] = tools.GetLocalAddr()
	m["auth_code"] = req.AuthCode

	sign, err := WechatGenSign(i.PayKey, m)
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

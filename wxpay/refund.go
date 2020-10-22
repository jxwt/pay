package wxpay

import (
	"bytes"
	"crypto/tls"
	"errors"
	"github.com/jxwt/pay"
	"github.com/jxwt/tools"
	"log"
	"net/http"
	"strings"
)

// WithCert 附着商户证书
func (i *WxClient) WithCert(certFile, keyFile string) error {
	certByte := FormatCeritficate(certFile)
	keyByte := FormatPrivateKey(keyFile)
	i.WithCertBytes(certByte, keyByte)
	return nil
}

func (i *WxClient) WithCertBytes(cert, key []byte) error {
	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return err
	}
	conf := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}
	trans := &http.Transport{
		TLSClientConfig: conf,
	}

	httpsClient := http.Client{
		Transport: trans,
	}
	i.httpsClient = &pay.HTTPSClient{
		httpsClient,
	}
	return nil
}

// PayRefund 微信退款
// outRefundNo 为后端自定义的随机字符串（尽量唯一） 与 商户退款单号（确保唯一性）
// OutTradeNo 需要退款的微信订单号
// refundDesc 退款理由
// totalFee,refundFee 订单的金额,与退款的金额
func (i *WxClient) PayRefund(payRefundReq *PayRefundRequest) (*WeChatQueryResult, error) {
	if err := i.WithCert(i.CertPEM, i.KeyPEM); err != nil {
		log.Printf("PayRefund err:%v\n", err)
		return nil, err
	}
	m := make(map[string]string)
	m["appid"] = i.AppID
	if i.SubMchId != "" {
		m["sub_mch_id"] = i.SubMchId
	}
	m["mch_id"] = i.MchID
	m["nonce_str"] = RandomStr()
	m["out_trade_no"] = payRefundReq.OutTradeNo
	m["out_refund_no"] = payRefundReq.OutRefundNo
	m["total_fee"] = WechatMoneyFeeToString(payRefundReq.TotalFee)
	m["refund_fee"] = WechatMoneyFeeToString(payRefundReq.RefundFee)
	m["refund_desc"] = payRefundReq.RefundDesc
	sign, err := WechatGenSign(i.PayKey, m)
	if err != nil {
		return nil, errors.New("wx refund sign err " + err.Error())
	}
	m["sign"] = sign

	// 发起退款申请
	result, err := PostWechat("https://api.mch.weixin.qq.com/secapi/pay/refund", m, i.httpsClient)
	if err != nil {
		log.Printf("clientPost.Post error: %v", err)
		return nil, err
	}
	return &result, nil
}

// PayReverse 撤销订单
func (i *WxClient) PayReverse(tradeNum string) (*WeChatQueryResult, error) {
	if err := i.WithCert(i.CertPEM, i.KeyPEM); err != nil {
		log.Printf("PayRefund err:%v\n", err)
		return nil, err
	}
	m := make(map[string]string)
	m["appid"] = i.AppID
	m["mch_id"] = i.MchID
	m["nonce_str"] = RandomStr()
	m["out_trade_no"] = tradeNum
	sign, err := WechatGenSign(i.PayKey, m)
	if err != nil {
		return nil, errors.New("wx refund sign err " + err.Error())
	}
	m["sign"] = sign
	// 发起退款申请
	result, err := PostWechat("https://api.mch.weixin.qq.com/secapi/pay/reverse", m, i.httpsClient)
	if err != nil {
		log.Printf("clientPost.Post error: %v", err)
		return nil, err
	}
	return &result, nil
}

// FormatPrivateKey 格式化私钥
func FormatPrivateKey(raw string) (result []byte) {
	return formatKey(raw, "-----BEGIN PRIVATE KEY-----", "-----END PRIVATE KEY-----")
}

// FormatCeritficate 格式化cert
func FormatCeritficate(raw string) (result []byte) {
	return formatKey(raw, "-----BEGIN CERTIFICATE-----", "-----END CERTIFICATE-----")
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

//企业付款，成功返回自定义订单号，微信订单号，true，失败返回错误信息，false
func (i *WxClient) Transfer(payRefundReq *PayRefundRequest) error {
	if err := i.WithCert(i.CertPEM, i.KeyPEM); err != nil {
		log.Printf("PayRefund err:%v\n", err)
		return err
	}
	m := make(map[string]string)
	m["mch_appid"] = i.AppID
	m["mch_id"] = i.MchID
	m["open_id"] = payRefundReq.OpenId
	m["nonce_str"] = RandomStr()
	m["partner_trade_no"] = payRefundReq.OutRefundNo
	m["amount"] = WechatMoneyFeeToString(payRefundReq.RefundFee)
	m["check_name"] = "NO_CHECK"
	m["desc"] = payRefundReq.RefundDesc
	m["spbill_create_ip"] = tools.GetLocalAddr()
	sign, err := WechatGenSign(i.PayKey, m)
	if err != nil {
		return errors.New("wx refund sign err " + err.Error())
	}
	m["sign"] = sign

	// 发起退款申请
	result, err := PostWechat("https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers", m, i.httpsClient)
	if err != nil {
		log.Printf("clientPost.Post error: %v", err)
		return err
	}

	if result.ReturnCode == "SUCCESS" && result.ResultCode == "SUCCESS" {
		return nil
	}
	return errors.New(result.ReturnMsg)
}

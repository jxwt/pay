package wxpay

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/axgle/mahonia"
	"github.com/jxwt/pay"
	"github.com/jxwt/tools"
	"github.com/shopspring/decimal"
	"io"
	"io/ioutil"
	"sort"
	"strings"
	"time"
)

// Charge 支付参数
type Charge struct {
	TradeNum    string  `json:"tradeNum,omitempty"`
	Origin      string  `json:"origin,omitempty"`
	UserID      string  `json:"userId,omitempty"`
	PayMethod   int64   `json:"payMethod,omitempty"`
	MoneyFee    float64 `json:"MoneyFee,omitempty"`
	CallbackURL string  `json:"callbackURL,omitempty"`
	ReturnURL   string  `json:"returnURL,omitempty"`
	ShowURL     string  `json:"showURL,omitempty"`
	Describe    string  `json:"describe,omitempty"`
	OpenID      string  `json:"openid,omitempty"`
	CheckName   bool    `json:"check_name,omitempty"`
	ReUserName  string  `json:"re_user_name,omitempty"`
	BuyerId     string  `json:"buyerId,omitempty"`

	SceneInfo string `json:"omitempty"` //h5支付使用

	AppType     string // app场景名称
	AppName     string // app名称
	BundleId    string // ios传
	PackageName string // android传
}

//PayCallback 支付返回
type PayCallback struct {
	Origin      string `json:"origin"`
	TradeNum    string `json:"trade_num"`
	OrderNum    string `json:"order_num"`
	CallBackURL string `json:"callback_url"`
	Status      int64  `json:"static"`
}

// CallbackReturn 回调业务代码时的参数
type CallbackReturn struct {
	IsSucceed     bool   `json:"isSucceed"`
	OrderNum      string `json:"orderNum"`
	TradeNum      string `json:"tradeNum"`
	UserID        string `json:"userID"`
	MoneyFee      int64  `json:"moneyFee"`
	Sign          string `json:"sign"`
	ThirdDiscount int64  `json:"thirdDiscount"`
}

// BaseResult 支付结果
type BaseResult struct {
	IsSucceed     bool   // 是否交易成功
	TradeNum      string // 交易流水号
	MoneyFee      int64  // 支付金额
	TradeTime     string // 交易时间
	ContractNum   string // 交易单号
	UserInfo      string // 支付账号信息(有可能有，有可能没有)
	ThirdDiscount int64  // 第三方优惠
}

//RandomStr 获取一个随机字符串
func RandomStr() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// 微信企业付款到零钱
func WachatCompanyChange(mchAppid, mchid, key string, conn *pay.HTTPSClient, charge *Charge) (map[string]string, error) {
	var m = make(map[string]string)
	m["mch_appid"] = mchAppid
	m["mchid"] = mchid
	m["nonce_str"] = RandomStr()
	m["partner_trade_no"] = charge.TradeNum
	m["openid"] = charge.OpenID
	m["amount"] = WechatMoneyFeeToString(charge.MoneyFee)
	m["spbill_create_ip"] = tools.GetLocalAddr()
	m["desc"] = TruncatedText(charge.Describe, 32)

	// 是否验证用户名称
	if charge.CheckName {
		m["check_name"] = "FORCE_CHECK"
		m["re_user_name"] = charge.ReUserName
	} else {
		m["check_name"] = "NO_CHECK"
	}

	sign, err := WechatGenSign(key, m)
	if err != nil {
		return map[string]string{}, err
	}
	m["sign"] = sign

	// 转出xml结构
	result, err := PostWechat("https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers", m, conn)
	if err != nil {
		return map[string]string{}, err
	}

	return struct2Map(result)
}

func WechatGenSign(key string, m map[string]string) (string, error) {
	var signData []string
	for k, v := range m {
		if v != "" && k != "sign" && k != "key" {
			signData = append(signData, fmt.Sprintf("%s=%s", k, v))
		}
	}

	sort.Strings(signData)
	signStr := strings.Join(signData, "&")
	signStr = signStr + "&key=" + key

	c := md5.New()
	_, err := c.Write([]byte(signStr))
	if err != nil {
		return "", errors.New("WechatGenSign md5.Write: " + err.Error())
	}
	signByte := c.Sum(nil)
	if err != nil {
		return "", errors.New("WechatGenSign md5.Sum: " + err.Error())
	}
	return strings.ToUpper(fmt.Sprintf("%x", signByte)), nil
}

func TruncatedText(data string, length int) string {
	data = FilterTheSpecialSymbol(data)
	if len([]rune(data)) > length {
		return string([]rune(data)[:length-1])
	}
	return data
}

//过滤特殊符号
func FilterTheSpecialSymbol(data string) string {
	// 定义转换规则
	specialSymbol := func(r rune) rune {
		if r == '`' || r == '!' || r == '$' ||
			r == '^' || r == '(' || r == ')' || r == '=' ||
			r == ':' || r == ';' ||
			r == ',' || r == '\\' || r == '[' || r == '.' || r == '<' ||
			r == '>' || r == '/' || r == '?' || r == '~' || r == '！' || r == '@' || r == '#' ||
			r == '￥' || r == '…' || r == '*' || r == '（' || r == '）' || r == '—' ||
			r == '|' || r == '{' || r == '}' || r == '【' || r == '】' || r == '‘' || r == '；' ||
			r == '：' || r == '”' || r == '“' || r == '\'' || r == '。' || r == '，' ||
			r == '、' || r == '？' || r == '%' || r == '+' || r == '_' || r == ']' || r == '"' || r == '&' {
			return ' '
		}
		return r
	}
	data = strings.Map(specialSymbol, data)
	return strings.Replace(data, "\n", " ", -1)
}

//对微信下订单或者查订单
func PostWechat(url string, data map[string]string, h *pay.HTTPSClient) (WeChatQueryResult, error) {
	var xmlRe WeChatQueryResult

	hc := new(pay.HTTPSClient)
	var re []byte
	var err error
	if h != nil {
		resp, err := h.Post(url, "application/xml; charset=utf-8", XmlEncode(data))
		if err != nil {
			return xmlRe, errors.New("HTTPSC.PostData: " + err.Error())
		}
		defer resp.Body.Close()
		re, err = ioutil.ReadAll(resp.Body)
	} else {
		hc = pay.HTTPSC
		buf := bytes.NewBufferString("")
		for k, v := range data {
			buf.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k))
		}
		xmlStr := fmt.Sprintf("<xml>%s</xml>", buf.String())
		logs.Warning(xmlStr)
		re, err = hc.PostData(url, "text/xml:charset=UTF-8", xmlStr)
		if err != nil {
			return xmlRe, errors.New("HTTPSC.PostData: " + err.Error())
		}
	}

	err = xml.Unmarshal(re, &xmlRe)
	if err != nil {
		logs.Error("get body:", string(re))
		return xmlRe, errors.New("xml.Unmarshal: " + err.Error())
	}

	if xmlRe.ReturnCode != "SUCCESS" {
		// 通信失败
		return xmlRe, errors.New("xmlRe.ReturnMsg: " + xmlRe.ReturnMsg)
	}

	if xmlRe.ResultCode != "SUCCESS" {
		// 业务结果失败
		return xmlRe, errors.New("xmlRe.ErrCodeDes: " + xmlRe.ErrCodeDes)
	}
	return xmlRe, nil
}

func XmlEncode(params map[string]string) io.Reader {
	var buf bytes.Buffer
	decoder := mahonia.NewDecoder("utf-8")
	if decoder == nil {
		fmt.Println("编码不存在!")
	}
	buf.WriteString(`<xml>`)
	for k, v := range params {
		buf.WriteString(`<`)
		buf.WriteString(k)
		buf.WriteString(`>`)
		buf.WriteString(v)
		buf.WriteString(`</`)
		buf.WriteString(k)
		buf.WriteString(`>`)
	}
	buf.WriteString(`</xml>`)
	fmt.Println(buf.String())
	return &buf
}

// 微信金额浮点转字符串
func WechatMoneyFeeToString(moneyFee float64) string {
	aDecimal := decimal.NewFromFloat(moneyFee)
	bDecimal := decimal.NewFromFloat(100)
	return aDecimal.Mul(bDecimal).Truncate(0).String()
}

func struct2Map(obj interface{}) (map[string]string, error) {

	j2 := make(map[string]string)

	j1, err := json.Marshal(obj)
	if err != nil {
		return j2, err
	}

	err2 := json.Unmarshal(j1, &j2)
	return j2, err2
}

func XmlToMap(xmlData []byte) map[string]string {
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	m := make(map[string]string)
	var token xml.Token
	var err error
	var k string
	for token, err = decoder.Token(); err == nil; token, err = decoder.Token() {
		if v, ok := token.(xml.StartElement); ok {
			k = v.Name.Local
			continue
		}
		if v, ok := token.(xml.CharData); ok {
			data := string(v.Copy())
			if strings.TrimSpace(data) == "" {
				continue
			}
			m[k] = data
		}
	}

	if err != nil && err != io.EOF {
		panic(err)
	}
	return m
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

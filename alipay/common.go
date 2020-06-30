package alipay

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/jxwt/pay"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

type ExtendParam struct {
	SysServiceProviderId string `json:"sys_service_provider_id"`
	IndustryRefluxInfo   string `json:"industry_reflux_info"`
}

type SceneData struct {
	LicensePlate   string `json:"license_plate"`
	StartTime      string `json:"start_time"`
	ParkingLotName string `json:"parking_lot_name"`
	CityCode       string `json:"city_code"`
	ParkingLotId   string `json:"parking_lot_id"`
}

type IndustryRefluxInfo struct {
	Channel   string    `json:"channel"`
	SceneCode string    `json:"scene_code"`
	SceneData SceneData `json:"scene_data"`
}

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
	SceneType   string  `json:"omitempty"` //h5支付使用

	ExtendParam string
	IndustryRefluxInfo *IndustryRefluxInfo
}

// RandomStr 获取一个随机字符串
func RandomStr() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// 对支付宝者查订单
func GetAlipay(url string) (AliWebQueryResult, error) {
	var xmlRe AliWebQueryResult

	re, err := pay.HTTPSC.GetData(url)
	if err != nil {
		return xmlRe, errors.New("HTTPSC.PostData: " + err.Error())
	}
	err = xml.Unmarshal(re, &xmlRe)
	if err != nil {
		return xmlRe, errors.New("xml.Unmarshal: " + err.Error())
	}
	return xmlRe, nil
}

//对支付宝者查订单
func GetAlipayApp(urls string) (AliWebAppQueryResult, error) {
	var aliPay AliWebAppQueryResult

	re, err := pay.HTTPSC.GetData(urls)
	if err != nil {
		return aliPay, errors.New("HTTPSC.PostData: " + err.Error())
	}

	err = json.Unmarshal(re, &aliPay)
	if err != nil {
		panic(fmt.Sprintf("re is %v, err is %v", re, err))
	}

	return aliPay, nil
}

// 支付宝金额转字符串
func AliyunMoneyFeeToString(moneyFee float64) string {
	return decimal.NewFromFloat(moneyFee).Truncate(2).String()
}

// ToURL
func ToURL(payUrl string, m map[string]string) string {
	var buf []string
	for k, v := range m {
		buf = append(buf, fmt.Sprintf("%s=%s", k, v))
	}
	return fmt.Sprintf("%s?%s", payUrl, strings.Join(buf, "&"))
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

// TradeRelationBindRequest 分账关系绑定请求
type TradeRelationBindRequest struct {
	ReceiverLists []ReceiverList `json:"receiver_list"`
	OutRequestNo  string         `json:"out_request_no"` // 外部请求号 32 唯一
}

// ReceiverList 分账列表
type ReceiverList struct {
	Type    string `json:"type"`    // 分账方类型
	Account string `json:"account"` // 分账方帐号
	Name    string `json:"name"`    // 分账方全称
	Memo    string `json:"memo"`    // 分账关系描述
}

// EncryptedDataStruct .
type EncryptedDataStruct struct {
	Response string `json:"response"`
	Sign     string `json:"sign"`
}

// AliParsePhoneNumberResponse 解析手机号返回
type AliParsePhoneNumberResponse struct {
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	Mobile string `json:"mobile"`
}
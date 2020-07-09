package wxpay

// WeChatResult 微信支付返回
type WeChatReResult struct {
	PrepayID string `xml:"prepay_id" json:"prepay_id,omitempty"`
	CodeURL  string `xml:"code_url" json:"code_url,omitempty"`
	MwebUrl  string `xml:"mweb_url" json:"mweb_url,omitempty"`
}

// WechatBaseResult 基本信息
type WechatBaseResult struct {
	ReturnCode string `xml:"return_code" json:"return_code,omitempty"`
	ReturnMsg  string `xml:"return_msg" json:"return_msg,omitempty"`
}

// WechatReturnData 返回通用数据
type WechatReturnData struct {
	AppID      string `xml:"appid,omitempty" json:"appid,omitempty"`
	MchID      string `xml:"mch_id,omitempty" json:"mch_id,omitempty"`
	MchAppid   string `xml:"mch_appid,omitempty" json:"mch_appid,omitempty"`
	DeviceInfo string `xml:"device_info,omitempty" json:"device_info,omitempty"`
	NonceStr   string `xml:"nonce_str,omitempty" json:"nonce_str,omitempty"`
	Sign       string `xml:"sign,omitempty" json:"sign,omitempty"`
	ResultCode string `xml:"result_code,omitempty" json:"result_code,omitempty"`
	ErrCode    string `xml:"err_code,omitempty" json:"err_code,omitempty"`
	ErrCodeDes string `xml:"err_code_des,omitempty" json:"err_code_des,omitempty"`
}

// WechatResultData 结果通用数据
type WechatResultData struct {
	OpenID         string `xml:"openid,omitempty" json:"openid,omitempty"`
	IsSubscribe    string `xml:"is_subscribe,omitempty" json:"is_subscribe,omitempty"`
	TradeType      string `xml:"trade_type,omitempty" json:"trade_type,omitempty"`
	BankType       string `xml:"bank_type,omitempty" json:"bank_type,omitempty"`
	FeeType        string `xml:"fee_type,omitempty" json:"fee_type,omitempty"`
	TotalFee       string `xml:"total_fee,omitempty" json:"total_fee,omitempty"`
	CashFeeType    string `xml:"cash_fee_type,omitempty" json:"cash_fee_type,omitempty"`
	CashFee        string `xml:"cash_fee,omitempty" json:"cash_fee,omitempty"`
	TransactionID  string `xml:"transaction_id,omitempty" json:"transaction_id,omitempty"`
	OutTradeNO     string `xml:"out_trade_no,omitempty" json:"out_trade_no,omitempty"`
	Attach         string `xml:"attach,omitempty" json:"attach,omitempty"`
	TimeEnd        string `xml:"time_end,omitempty" json:"time_end,omitempty"`
	PartnerTradeNo string `xml:"partner_trade_no,omitempty" json:"partner_trade_no,omitempty"`
	PaymentNo      string `xml:"payment_no,omitempty" json:"payment_no,omitempty"`
	PaymentTime    string `xml:"payment_time,omitempty" json:"payment_time,omitempty"`
	DetailId       string `xml:"detail_id,omitempty" json:"detail_id,omitempty"`
}

type WeChatPayResult struct {
	WechatBaseResult
	WechatReturnData
	WechatResultData
}

//type WeChatPayResult struct {
//	Appid         string `json:"appid"`
//	BankType      string `json:"bank_type"`
//	CashFee       string `json:"cash_fee"`
//	FeeType       string `json:"fee_type"`
//	IsSubscribe   string `json:"is_subscribe"`
//	MchID         string `json:"mch_id"`
//	NonceStr      string `json:"nonce_str"`
//	Openid        string `json:"openid"`
//	OutTradeNo    string `json:"out_trade_no"`
//	ResultCode    string `json:"result_code"`
//	ReturnCode    string `json:"return_code"`
//	Sign          string `json:"sign"`
//	TimeEnd       string `json:"time_end"`
//	TotalFee      string `json:"total_fee"`
//	TradeType     string `json:"trade_type"`
//	TransactionID string `json:"transaction_id"`
//}
type WeChatQueryResult struct {
	WechatBaseResult
	WeChatReResult
	WechatReturnData
	WechatResultData
	PayRefundResponse
	TradeState     string `xml:"trade_state" json:"trade_state,omitempty"`
	TradeStateDesc string `xml:"trade_state_desc" json:"trade_state_desc,omitempty"`
}

// WxPayRefundResponse 微信退款请求返回
type PayRefundResponse struct {
	OutRefundNo         string `xml:"out_refund_no"`
	RefundID            string `xml:"refund_id"`
	RefundFee           int    `xml:"refund_fee"`
	SettlementRefundFee int    `xml:"settlement_refund_fee"`
	SettlementTotalFee  int    `xml:"settlement_total_fee"`
	CashRefundFee       int    `xml:"cash_refund_fee"`
	CouponRefundFee     int    `xml:"coupon_refund_fee"`
	CouponRefundCount   int    `xml:"coupon_refund_count"`
}

// PayRefundRequest 外部调用的退款请求
type PayRefundRequest struct {
	OutRefundNo string // 商户退款单号（确保唯一性）
	//TransactionID string // 需要退款的微信订单号
	RefundDesc string  // 退款理由
	TotalFee   float64 // 订单的金额
	RefundFee  float64 // 退款的金额
	OutTradeNo string  // 商户自定义单号（需要退款的单号）
	OpenId     string
}

// MicroPayRequest 付款码支付请求
type MicroPayRequest struct {
	OutTradeNo string `json:"out_trade_no"` // 商户订单号
	TotalFee   int    `json:"total_fee"`    // 单位分 总金额
	AuthCode   string `json:"auth_code"`    // 授权码（条形码）
	Remark     string `json:"remark"`       // 备注
}

// MicroPayResponse 付款码支付返回
type MicroPayResponse struct {
	ReturnCode         string `xml:"return_code"`
	ReturnMsg          string `xml:"return_msg"`
	ResultCode         string `xml:"result_code"`
	ErrCode            string `xml:"err_code"`
	ErrCodeDes         string `xml:"err_code_des"`
	AppID              string `xml:"appid"`
	MchID              string `xml:"mch_id"`
	NonceStr           string `xml:"nonce_str"`
	Sign               string `xml:"sign"`
	OpenID             string `xml:"open_id"`
	IsSubscribe        string `xml:"is_subscribe"` // 用户是否关注公众账号
	TradeType          string `xml:"trade_type"`   // MICROPAY 付款码支付
	BankType           string `xml:"bank_type"`    // 银行类型
	FeeType            string `xml:"fee_type"`     // 符合ISO 4217标准的三位字母代码，默认人民币：CNY
	TotalFee           int    `xml:"total_fee"`
	SettlementTotalFee int    `xml:"settlement_total_fee"` // 订单总金额，单位为分
	CouponFee          int    `xml:"coupon_fee"`           // “代金券”金额<=订单金额，订单金额-“代金券”金额=现金支付金额
	CashFeeType        string `xml:"cash_fee_type"`
	CashFee            int    `xml:"cash_fee"` // 订单现金支付金额
	TransactionID      string `xml:"transaction_id"`
	OutTradeNo         string `xml:"out_trade_no"`
	Attach             string `xml:"attach"`
	TimeEnd            string `xml:"time_end"` // 订单生成时间
}

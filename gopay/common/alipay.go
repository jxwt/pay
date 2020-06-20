package common

// AliWebPayResult 支付宝支付结果回调
type AliWebPayResult struct {
	AppID          string `json:"app_id"`
	AuthAppID      string `json:"auth_app_id"`
	BuyerID        string `json:"buyer_id"`
	BuyerLogonID   string `json:"buyer_logon_id"`
	BuyerPayAmount string `json:"buyer_pay_amount"`
	Charset        string `json:"charset"`
	FundBillList   string `json:"fund_bill_list"`
	GmtCreate      string `json:"gmt_create"`
	GmtPayment     string `json:"gmt_payment"`
	InvoiceAmount  string `json:"invoice_amount"`
	NotifyID       string `json:"notify_id"`
	NotifyTime     string `json:"notify_time"`
	NotifyType     string `json:"notify_type"`
	OutTradeNo     string `json:"out_trade_no"`
	PointAmount    string `json:"point_amount"`
	ReceiptAmount  string `json:"receipt_amount"`
	SellerEmail    string `json:"seller_email"`
	SellerID       string `json:"seller_id"`
	Sign           string `json:"sign"`
	SignType       string `json:"sign_type"`
	Subject        string `json:"subject"`
	TotalAmount    string `json:"total_amount"`
	TradeNo        string `json:"trade_no"`
	TradeStatus    string `json:"trade_status"`
	Version        string `json:"version"`
}

type FundBill struct {
	FundChannel string `json:"fundChannel"`
	Amount      string `json:"amount"`
}
type AliQueryResult struct {
	AppID          string `json:"app_id"`
	AuthAppID      string `json:"auth_app_id"`
	BuyerID        string `json:"buyer_id"`
	BuyerLogonID   string `json:"buyer_logon_id"`
	BuyerPayAmount string `json:"buyer_pay_amount"`
	Charset        string `json:"charset"`
	FundBillList   string `json:"fund_bill_list"`
	GmtCreate      string `json:"gmt_create"`
	GmtPayment     string `json:"gmt_payment"`
	InvoiceAmount  string `json:"invoice_amount"`
	NotifyID       string `json:"notify_id"`
	NotifyTime     string `json:"notify_time"`
	NotifyType     string `json:"notify_type"`
	OutTradeNo     string `json:"out_trade_no"`
	PointAmount    string `json:"point_amount"`
	ReceiptAmount  string `json:"receipt_amount"`
	SellerEmail    string `json:"seller_email"`
	SellerID       string `json:"seller_id"`
	Sign           string `json:"sign"`
	SignType       string `json:"sign_type"`
	Subject        string `json:"subject"`
	TotalAmount    string `json:"total_amount"`
	TradeNo        string `json:"trade_no"`
	TradeStatus    string `json:"trade_status"`
	Version        string `json:"version"`
}

type AliWebQueryResult struct {
	IsSuccess string `xml:"is_success"`
	ErrorMsg  string `xml:"error"`
	SignType  string `xml:"sign_type"`
	Sign      string `xml:"sign"`
	Response  struct {
		Trade struct {
			BuyerEmail          string `xml:"buyer_email"`
			BuyerId             string `xml:"buyer_id"`
			SellerID            string `xml:"seller_id"`
			TradeStatus         string `xml:"trade_status"`
			IsTotalFeeAdjust    string `xml:"is_total_fee_adjust"`
			OutTradeNum         string `xml:"out_trade_no"`
			Subject             string `xml:"subject"`
			FlagTradeLocked     string `xml:"flag_trade_locked"`
			Body                string `xml:"body"`
			GmtCreate           string `xml:"gmt_create"`
			GmtPayment          string `xml:"gmt_payment"`
			GmtLastModifiedTime string `xml:"gmt_last_modified_time"`
			SellerEmail         string `xml:"seller_email"`
			TotalFee            string `xml:"total_fee"`
			TradeNum            string `xml:"trade_no"`
		} `xml:"trade"`
	} `xml:"response"`
}

type AliWebAppQueryResult struct {
	AlipayTradeQueryResponse struct {
		Code                string `json:"code"`
		Msg                 string `json:"msg"`
		SubCode             string `json:"sub_code"`
		SubMsg              string `json:"sub_msg"`
		TradeNo             string `json:"trade_no"`
		OutTradeNo          string `json:"out_trade_no"`
		OpenId              string `json:"open_id"`
		BuyerLogonId        string `json:"buyer_logon_id"`
		TradeStatus         string `json:"trade_status"`
		TotalAmount         string `json:"total_amount"`
		ReceiptAmount       string `json:"receipt_amount"`
		BuyerPayAmount      string `json:"buyer_pay_amount"`
		PointAmount         string `json:"point_amount"`
		InvoiceAmount       string `json:"invoice_amount"`
		SendPayDate         string `json:"send_pay_date"`
		AlipayStoreId       string `json:"alipay_store_id"`
		StoreId             string `json:"store_id"`
		TerminalId          string `json:"terminal_id"`
		StoreName           string `json:"store_name"`
		BuyerUserId         string `json:"buyer_user_id"`
		DiscountGoodsDetail string `json:"discount_goods_detail"`
		IndustrySepcDetail  string `json:"industry_sepc_detail"`
	} `json:"alipay_trade_query_response"`
	Sign string `json:"sign"`
}

type AliPayResponse struct {
	Code    string `json:"code"`     //网关返回码
	Msg     string `json:"msg"`      //网关返回码描述
	SubCode string `json:"sub_code"` //业务返回码
	SubMsg  string `json:"sub_msg"`  //业务返回码描述
	Sign    string `json:"sign"`     //签名
}

//线下收单预创建请求参数
type PreCreateRequest struct {
	//必填参数
	OutTradeNo  string  `json:"out_trade_no"` //商户订单号,64个字符以内、只能包含字母、数字、下划线；需保证在商户端不重复
	TotalAmount float64 `json:"total_amount"` //订单总金额，单位为元，精确到小数点后两位
	Subject     string  `json:"subject"`      //订单标题
	//可选参数
	SellerId           string       `json:"seller_id"`           //卖家支付宝用户ID
	DiscountableAmount float64      `json:"discountable_amount"` //可打折金额. 参与优惠计算的金额，单位为元，精确到小数点后两位
	GoodsDetail        []GoodDetail `json:"goods_detail"`        //订单包含的商品列表信息.json格式. 其它说明详见：“商品明细说明”
	Body               string       `json:"body"`                //对商品的描述
	ProductCode        string       `json:"product_code"`        //销售产品码。
	OperatorId         string       `json:"operator_id"`         //商户操作员编号
	StoreId            string       `json:"store_id"`            //商户门店编号
	//懒得写了,下面还有,有需要再加
}
type PreCreateResponse struct {
	PreCreateResult PreCreateResult `json:"alipay_trade_precreate_response"`
	Sign            string          `json:"sign"`
}

type PreCreateResult struct {
	AliPayResponse
	OutTradeNo string `json:"out_trade_no"`
	QrCode     string `json:"qr_code"`
}
type GoodDetail struct {
	//必填参数
	GoodsId   string  `json:"goods_id"`   //商品的编号
	GoodsName string  `json:"goods_name"` //商品名称
	Quantity  int     `json:"quantity"`   //商品数量
	Price     float64 `json:"price"`      //商品单价，单位为元
	//可选参数
	GoodsCategory  string `json:"goods_category"`  //商品类目
	CategoriesTree string `json:"categories_tree"` //商品类目树
	Body           string `json:"body"`            //商品描述信息
	ShowUrl        string `json:"show_url"`        //商品展示URL
}

// 退款请求参数
type AliRefundRequest struct {
	AppAuthToken string `json:"-"`                      // 可选
	OutTradeNo   string `json:"out_trade_no,omitempty"` // 与 TradeNo 二选一
	TradeNo      string `json:"trade_no,omitempty"`     // 与 OutTradeNo 二选一
	RefundAmount string `json:"refund_amount"`          // 必须 需要退款的金额，该金额不能大于订单金额,单位为元，支持两位小数
	RefundReason string `json:"refund_reason"`          // 可选 退款的原因说明
	OutRequestNo string `json:"out_request_no"`         // 可选 标识一次退款请求，同一笔交易多次退款需要保证唯一，如需部分退款，则此参数必传。
	OperatorId   string `json:"operator_id"`            // 可选 商户的操作员编号
	StoreId      string `json:"store_id"`               // 可选 商户的门店编号
	TerminalId   string `json:"terminal_id"`            // 可选 商户的终端编号
}
type RefundDetailItem struct {
	FundChannel string `json:"fund_channel"` // 交易使用的资金渠道，详见 支付渠道列表
	Amount      string `json:"amount"`       // 该支付工具类型所使用的金额
	RealAmount  string `json:"real_amount"`  // 渠道实际付款金额
}

// 退款返回参数
type AliRefundResponse struct {
	AliPayTradeRefund struct {
		Code                 string              `json:"code"`
		Msg                  string              `json:"msg"`
		SubCode              string              `json:"sub_code"`
		SubMsg               string              `json:"sub_msg"`
		TradeNo              string              `json:"trade_no"`                          // 支付宝交易号
		OutTradeNo           string              `json:"out_trade_no"`                      // 商户订单号
		BuyerLogonId         string              `json:"buyer_logon_id"`                    // 用户的登录id
		BuyerUserId          string              `json:"buyer_user_id"`                     // 买家在支付宝的用户id
		FundChange           string              `json:"fund_change"`                       // 本次退款是否发生了资金变化
		RefundFee            string              `json:"refund_fee"`                        // 退款总金额
		GmtRefundPay         string              `json:"gmt_refund_pay"`                    // 退款支付时间
		StoreName            string              `json:"store_name"`                        // 交易在支付时候的门店名称
		RefundDetailItemList []*RefundDetailItem `json:"refund_detail_item_list,omitempty"` // 退款使用的资金渠道
	} `json:"alipay_trade_refund_response"`
	Sign string `json:"sign"`
}

// ToaccountTransferRequest 单笔转账请求
type ToaccountTransferRequest struct {
	// 必填
	OutBizNo     string `json:"out_biz_no"`    // 商户转账唯一订单号
	PayeeType    string `json:"payee_type"`    // 收款方账户类型 ALIPAY_LOGONID：支付宝登录号
	PayeeAccount string `json:"payee_account"` // 收款方账户
	Amount       string `json:"amount"`        // 转账金额字符串
	// 选填
	PayerShowName string `json:"payer_show_name"` // 付款方姓名
	PayeeRealName string `json:"payee_real_name"` // 收款方真实姓名
	Remark        string `json:"remark"`          // 备注
}

// ToaccountTransferResponse 单笔转账返回
type ToaccountTransferResponse struct {
	AlipayFundTransToaccountTransferResponse struct {
		Code     string `json:"code"`
		Msg      string `json:"msg"`
		OutBizNo string `json:"out_biz_no"`
		OrderID  string `json:"order_id"`
		PayDate  string `json:"pay_date"`
	} `json:"alipay_fund_trans_toaccount_transfer_response"`
	Sign string `json:"sign"`
}

// AliTradePayRequest 支付宝统一收单请求
type AliTradePayRequest struct {
	OutTradeNo    string `json:"out_trade_no"`   // 商户订单号
	Scene         string `json:"scene"`          // 支付场景 条码支付，取值：bar_code 声波支付，取值：wave_code
	AuthCode      string `json:"auth_code"`      // 支付授权码 25~30开头的长度为16~24位的数字
	Subject       string `json:"subject"`        // 订单标题
	TotalAmount   string `json:"total_amount"`   // 订单总金额
	TransCurrency string `json:"trans_currency"` // 人民币：CNY
}

// AliTradePayResponse 支付宝统一收单返回
type AliTradePayResponse struct {
	AlipayTradePayResponse struct {
		Code            string  `json:"code"`
		Msg             string  `json:"msg"`
		TradeNo         string  `json:"trade_no"`
		OutTradeNo      string  `json:"out_trade_no"`
		BuyerLogonID    string  `json:"buyer_logon_id"`
		SettleAmount    string  `json:"settle_amount"`
		PayCurrency     string  `json:"pay_currency"`
		PayAmount       string  `json:"pay_amount"`
		SettleTransRate string  `json:"settle_trans_rate"`
		TransPayRate    string  `json:"trans_pay_rate"`
		TotalAmount     float64 `json:"total_amount"`
		TransCurrency   string  `json:"trans_currency"`
		SettleCurrency  string  `json:"settle_currency"`
		ReceiptAmount   string  `json:"receipt_amount"`
		BuyerPayAmount  string  `json:"buyer_pay_amount"`
		PointAmount     float64 `json:"point_amount"`
		InvoiceAmount   string  `json:"invoice_amount"`
		GmtPayment      string  `json:"gmt_payment"`
		FundBillList    []struct {
			FundChannel string  `json:"fund_channel"`
			BankCode    string  `json:"bank_code"`
			Amount      int     `json:"amount"`
			RealAmount  float64 `json:"real_amount"`
		} `json:"fund_bill_list"`
		CardBalance         float64 `json:"card_balance"`
		StoreName           string  `json:"store_name"`
		BuyerUserID         string  `json:"buyer_user_id"`
		DiscountGoodsDetail string  `json:"discount_goods_detail"`
		VoucherDetailList   []struct {
			ID                         string  `json:"id"`
			Name                       string  `json:"name"`
			Type                       string  `json:"type"`
			Amount                     string  `json:"amount"`
			MerchantContribute         int     `json:"merchant_contribute"`
			OtherContribute            int     `json:"other_contribute"`
			Memo                       string  `json:"memo"`
			TemplateID                 string  `json:"template_id"`
			PurchaseBuyerContribute    float64 `json:"purchase_buyer_contribute"`
			PurchaseMerchantContribute float64 `json:"purchase_merchant_contribute"`
			PurchaseAntContribute      float64 `json:"purchase_ant_contribute"`
		} `json:"voucher_detail_list"`
		AdvanceAmount    string `json:"advance_amount"`
		AuthTradePayMode string `json:"auth_trade_pay_mode"`
		ChargeAmount     string `json:"charge_amount"`
		ChargeFlags      string `json:"charge_flags"`
		SettlementID     string `json:"settlement_id"`
		BusinessParams   string `json:"business_params"`
		BuyerUserType    string `json:"buyer_user_type"`
		MdiscountAmount  string `json:"mdiscount_amount"`
		DiscountAmount   string `json:"discount_amount"`
		BuyerUserName    string `json:"buyer_user_name"`
	} `json:"alipay_trade_pay_response"`
	Sign string `json:"sign"`
}

// AliTradeCancelRequest 支付宝撤单请求
type AliTradeCancelRequest struct {
	OutTradeNo string `json:"out_trade_no"`
}

// AliTradeCancelResponse 支付宝撤单返回
type AliTradeCancelResponse struct {
	AlipayTradeCancelResponse struct {
		Code               string `json:"code"`
		Msg                string `json:"msg"`
		TradeNo            string `json:"trade_no"`
		OutTradeNo         string `json:"out_trade_no"`
		RetryFlag          string `json:"retry_flag"`
		Action             string `json:"action"`
		GmtRefundPay       string `json:"gmt_refund_pay"`
		RefundSettlementID string `json:"refund_settlement_id"`
	} `json:"alipay_trade_cancel_response"`
	Sign string `json:"sign"`
}

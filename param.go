package pay

const (
	apiRegister = "/api/service/create" // 注册api
	apiDoPay    = "/api/pay/doPay"      // 支付api
	apiDoOutPay = "/api/pay/doOutPay"   // 出账api
)

// 内部支付域名
const (
	urlPay = "http://jxpay.com"
)

// CommonResponse beego框架统一返回结构
type CommonResponse struct {
	State   string
	Message string
	Data    interface{}
}

// DoPayRequest 支付请求
type DoPayRequest struct {
	ServiceID     int     `json:"serviceID"`     // 服务ID
	Name          string  `json:"name"`          // 服务名称
	TenantID      int     `json:"tenantID"`      // 商户ID
	OrderID       int     `json:"orderID"`       // 订单号 或者 batchID
	Money         float64 `json:"money"`         // 金额
	CashReasonID  int     `json:"cashReasonID"`  // 支付原因ID
	CashChannelID int     `json:"cashChannelID"` // 支付方式ID
	CallBackURL   string  `json:"callBackURL"`   // 支付回调地址(内部回调地址)
	// 非必须参数
	UserID        int     `json:"userID"`        // 系统内的用户ID(可空)
	ThirdUserID   string  `json:"thirdUserID"`   // 三方用户ID(可空)
	TradeNumber   string  `json:"tradeNumber"`   // 商户单号(只用于存两份cashflow的服务,目前只有停车服务需要)
	DiscountMoney float64 `json:"discountMoney"` // 优惠金额
	PresentMoney  float64 `json:"presentMoney"`  // 赠送金额
	ExtendParams  string  `json:"extendParams"`  // 额外参数(目前停车用)
}

// 服务注册参数
type RegisterRequest struct {
	Name string
	// 以上为通用参数 如果有商户独立结算 用下面参数传递
	PayItems []PayItem
}

// 注册返回
type RegisterResponse struct {
	State   string
	Message string
	Data    int
}

// 支付参数
type PayItem struct {
	TenantId         int
	WxAppId          string
	WxMchId          string
	WxSubMchId       string
	WxSecretKey      string
	WxPayKey         string
	WxCertPEM        string
	WxKeyPEM         string
	AliPayPublicKey  string
	AliPayPrivateKey string
	AliPayAppId      string
	AliPayPartnerId  string
	WxPublicAppId    string
	WxSubAppId       string
}

// SendCallBackNotify 回调通知
type SendCallBackNotify struct {
	TradeNumber   string  `json:"tradeNumber"`   // 商户单号
	Money         float64 `json:"money"`         // 金额
	OrderID       int     `json:"orderID"`       // 订单ID
	CashChannelID int     `json:"cashChannelID"` // 支付方式
	ConfirmAt     string  `json:"confirmAt"`     // 支付完成时间
	CashReasonID  int     `json:"cashReasonID"`  // 支付原因
}

// QRCodePayResponse 二维码支付返回
type QRCodePayResponse struct {
	CashReasonID uint   `json:"cashReasonID"` // 支付意向ID
	OrderID      uint   `json:"orderID"`      //订单好
	QrCodeURL    string `json:"qrCodeURL"`    //二维码链接
}

// DoOutPayRequest 发起出账请求
type DoOutPayRequest struct {
	ServiceID     int          `json:"serviceID"`     // 服务ID
	Name          string       `json:"name"`          // 服务名称
	TenantID      int          `json:"tenantID"`      // 商户ID
	UserID        int          `json:"userID"`        // 系统内的用户ID(可空)
	ThirdUserID   string       `json:"thirdUserID"`   // 三方用户ID(可空)
	OrderID       int          `json:"orderID"`       // 订单号 或者 batchID
	CashReasonID  int          `json:"cashReasonID"`  // 支付原因ID
	CashChannelID int          `json:"cashChannelID"` // 支付方式ID
	OutPayLists   []OutPayList `json:"outPayLists"`   // 出账列表
	CallBackURL   string       `json:"callBackURL"`   // 支付回调地址(内部回调地址)
}

// OutPayList 出账列表
type OutPayList struct {
	PayMoney    float64 `json:"payMoney"`    // 原订单金额
	TradeNumber string  `json:"tradeNumber"` // 商户单号
	RefundMoney float64 `json:"refundMoney"` // 需要退款金额
}

// WxRefundRequest 微信退款请求
type WxRefundRequest struct {
	PayMoney    float64
	TradeNumber string
	RefundMoney float64
}

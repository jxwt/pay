package pay

const (
	ALI_WEB             = iota + 1 // 支付宝网页
	ALI_APP                        // 支付宝App
	WECHAT_WEB                     // 微信网页
	WECHAT_APP                     // 微信App
	WECHAT_MINI_PROGRAM            // 微信小程序
	AlI_SCANBAR                    //支付宝扫码
	WECHAT_SCANBAR                 // 微信扫码
	WECHAT_H5                      //微信H5支付
	ALI_H5                         //支付宝h5
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
	SceneType   string  `json:"omitempty"` //h5支付使用
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

package pay

// 支付方式枚举
const (
	// 微信类
	CashChannelWxAppPay    = 11 // 微信支付
	CashChannelWxH5Pay     = 12 // 微信H5
	CashChannelWxMiniPay   = 13 // 微信小程序
	CashChannelWxPublicPay = 14 // 微信公众号
	CashChannelWxCodePay   = 15 // 微信扫码支付
	// 阿里类
	CashChannelAliAppPay  = 30 // 支付宝
	CashChannelAliH5Pay   = 31 // 支付宝H5
	CashChannelAliCodePay = 32 // 支付宝扫码支付
	CashChannelAliMiniPay = 33 // 支付宝小程序

	// 其他
	CashChannelDepositPay = 50 // 钱包支付
	CashChannelCash       = 51 // 现金支付
	CashChannelPackHour   = 52 // 包时长
)

// CashChannelRemarks 支付途径描述
var CashChannelRemarks = map[uint]string{
	CashChannelWxAppPay:    "微信支付",
	CashChannelAliAppPay:   "支付宝",
	CashChannelDepositPay:  "钱包支付",
	CashChannelWxH5Pay:     "微信H5",
	CashChannelAliH5Pay:    "支付宝H5",
	CashChannelWxMiniPay:   "微信小程序",
	CashChannelWxPublicPay: "微信公众号",
	CashChannelAliCodePay:  "支付宝扫码支付",
	CashChannelWxCodePay:   "微信扫码支付",
	CashChannelAliMiniPay:  "支付宝小程序",
	CashChannelCash:        "现金支付",
	CashChannelPackHour:    "包时长",
}

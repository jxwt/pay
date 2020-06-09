package pay

// 支付原因枚举
const (
	// 通用入账原因
	CashReasonOrder    = 0 // 支付订单
	CashReasonRecharge = 1 // 用户充值

	// parking_cloud 停车项目支付原因
	CashReasonReservation          = 101 // 预约支付
	CashReasonOrders               = 102 // 批量缴纳
	CashReasonBillPackPlateFund    = 103 // 包车牌支付
	CashReasonInvoiceFreight       = 104 // 发票运费支付
	CashReasonBillPackSpotFund     = 105 // 包车位支付
	CashReasonBillPackHourFund     = 106 // 包时订单支付
	CashReasonAppointShareOrderPay = 107 // 预约共享车位订单支付
	CashReasonMarginFund           = 108 // 保证金订单

	// 通用出账原因
	CashReasonOrderRefund = 10001 // 退款
	CashReasonCashOut     = 10002 // 提现,转账
)

// CashReasonRemarks 支付原因描述
var CashReasonRemarks = map[uint]string{
	CashReasonOrder:                "支付订单",
	CashReasonRecharge:             "用户充值",
	CashReasonReservation:          "预约支付",
	CashReasonOrders:               "批量支付",
	CashReasonBillPackPlateFund:    "包期支付",
	CashReasonInvoiceFreight:       "发票运费支付",
	CashReasonBillPackSpotFund:     "包车位支付",
	CashReasonBillPackHourFund:     "包时订单支付",
	CashReasonAppointShareOrderPay: "预约共享车位订单支付",
	CashReasonMarginFund:           "保证金支付",
	CashReasonOrderRefund:          "退款",
	CashReasonCashOut:              "提现转账",
}

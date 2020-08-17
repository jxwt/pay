package wxpay

import "testing"

func TestPayRefund(t *testing.T) {
	client := &WxClient{
		AppID:   "",
		MchID:   "",
		Key:     "",
	}
	req := &PayRefundRequest{
		OutRefundNo: "dfsefwdcweiojc3oe233",
		RefundDesc:  "退款",
		TotalFee:    0.01,
		RefundFee:   0.01,
		OutTradeNo:  "1320200729160058Re58BM7fNj",
		OpenId:      "oQ1W_4tgAd51uSJFhXQDGyfmRiSM",
	}
	res, err := client.PayRefund(req)
	if err != nil {
		t.Logf("%v", err)
	}
	t.Logf("%v", res)

}

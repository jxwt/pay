package wxpay

import "testing"

func TestH5Pay(t *testing.T){
	// 停车
	client:=&WxClient{
	AppID: "",
	MchID: "",
	Key: "",
	CallbackURL: "www.baidu.com",
	//SubMchId: "1601177571",
	}
	charge:=&Charge{
		TradeNum:    "sdfsdfec2e",
		MoneyFee:    0.02,
		CallbackURL: client.CallbackURL,
		Describe:    "test",
	}
	 m,err:=client.H5Pay(charge)
	 t.Logf("%v",m)
	 t.Logf("%v",err)
}
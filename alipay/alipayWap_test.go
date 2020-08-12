package alipay

import (
	"encoding/json"
	"testing"
)

func TestToH5Pay(t *testing.T){
	aliPay := Init("", "", "", "","www.baidu.com")
	industryRefluxInfo := &IndustryRefluxInfo{
		SceneCode: "parking_fee_order",
		Channel:   "common_park_provider",
		SceneData: SceneData{
			LicensePlate:   "浙A111111",
			StartTime:      "2020-08-02 15:04:05",
			ParkingLotName: "测试车场",
			CityCode:       "330100",
			ParkingLotId:  "1",
		},
	}
	d, _ := json.Marshal(industryRefluxInfo)
	extern := &ExtendParam{
		SysServiceProviderId: "2088521066336121",
		IndustryRefluxInfo:   string(d),
	}
	externParams, _ := json.Marshal(extern)
	charge:=&Charge{
		TradeNum:    "sdsfsdfe34343cdd2121e4",
		MoneyFee:    0.02,
		CallbackURL: aliPay.NotifyURL,
		Describe:    "test",
		AuthToken: "202008BBadd7073546f54a2bbfd99b7fc89e5C11",
		ExtendParam: string(externParams),
		//BuyerId: "2088702297070826",
	}
	//res,err:=aliPay.Client.ToH5Pay(charge)
	//t.Logf("%v",err)
	//t.Logf("%v",res)
	aliWapCient:=&AliWapClient{
		SellerID: aliPay.Client.SellerID,
		AppID: aliPay.Client.AppID,
		PrivateKey: aliPay.Client.PrivateKey,
		PublicKey: aliPay.Client.PublicKey,
	}
	res,err:=aliWapCient.ToH5Pay(charge)
	t.Logf("%v",err)
	t.Logf("%v",res)
}

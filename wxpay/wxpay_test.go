package wxpay

import "testing"

func TestH5Pay(t *testing.T) {
	// 停车
	client := &WxClient{
		AppID:       "",
		MchID:       "",
		Key:         "",
		CallbackURL: "www.baidu.com",
		//SubMchId: "1601177571",
	}
	charge := &Charge{
		TradeNum:    "sdfsdfec2e",
		MoneyFee:    0.02,
		CallbackURL: client.CallbackURL,
		Describe:    "test",
	}
	m, err := client.H5Pay(charge)
	t.Logf("%v", m)
	t.Logf("%v", err)
}

func TestWxClient_MiniLogin(t *testing.T) {
	c := &WxClient{
		AppID: "wx0a8581e498061282",
		Key:   "69715695cd4c32a14d0328db206882f3",
	}

	c.GetLoginInfo("9QSjCYMXJUvvkD5tC5LN/g==", "EEaeNQWaj4wAHzTeFl5dn8ZN60KLJhlKoHGeKjm4XZoJn0Xdv+OK+5k09dbapONPmQn+7lAwkpMSItsTkzwMHwq8Yuhg5eaUUzz+dBC9Y7ulyi5vwi1ev/c1UlAakJzIcd6BuAOJH0Op7EdMEtpqmCPSJcReUSjpg/dgXm7aBTkeQZyUM+Qbtk6D/Y88XBixkB+HkKLnNRbpUHs3Bwj8vtP5qJtFpk0qGViKN04oCaj5uW2UeYXqSa14nvpkeq+zljguap7F50x98v3rg+sfP4RTUPO41EuBCEWFqBnELmNx6nQIQs3Zb9qefkAM4LgVZjHWB3ABUeUM6mXT81ZHp0ThB30XWX3kUPGzaea4sUoRBigXbgXsP/li67fLdC9emxgsLa0LDZd++ILKPQJETVLdZsAhbq3W4a/v9svymZhitorJfYZJ4h1h7DZI9YGhTPstl/NxJPg63lhYJAT5qqhVA3eZ9Ye0F9GUH1gddpG0upk1nNGuh00Yut+MNIFBu8ErzcVro6r3Fn6EyHGGkQ==", "011bot100Y1lwK1cUk3000mAs30bot1u")
}

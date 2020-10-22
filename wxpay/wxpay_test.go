package wxpay

import (
	"fmt"
	"testing"
)

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
		Key:   "c45646ffea88cf98e4285dbf75893ff5",
	}
	//c.AppID  = "wx3474908ee1c58dea"
	//c.Key = "49bec205b4aeb983c1193096f80cd9f6"

	session := c.GetOpenSession("081CXKkl2deaQ54SZell2niID84CXKkM")
	fmt.Println(session.SessionKey)

	key := "pl3rOB+eRTIKmYYGg3at2Q=="
	iv := "KzevCXXfxsgygyh7EHwVuQ=="
	data := "y4guz/vPWpCu9tOEBGRMvmS6U4LKnVJ7zd12P7kJKydBocSgl28GkqltcBENJLbbcq7CuM0zXqn7vIIK2ZhU0NvhIcJ/BSg3Ry0M0PzlWyqbhUGXndKEQgU8ERRLN7oF3lZfw/WnPqBY4HFNrvDGMtCf0kQNbkpLbuxwn88Cpb5BAA2Rhbq9rnhEB8c0Hd/aIw5SxfTpE+kCp4IY5HsCQ9AufqB55U+nl9IK7MJtv0h3pxZvUQiw4wIgoQI9VNP71RSZLMUwQrVecGnLUXKf3tqT2VrgLapZdG0/BTn07oLeQ8ra1STdrnQQcbmyV4g+ALCTiY9Ezfqaea5swrgCQNHeGnzgOjxQp5BFsfynsGQT/x9XI0Pm/5iFenvJTxAjSWlh1mTc+zI3grCMe5osL+w12xzpzwZ3do/OeKHhRcHQmTWNP+x0aOa9dzLlD3N7qkb1bdkc4Ua464ptkgEFY2hfjlnn3ryT1/QQdHl9QKo45pxGydwqhpFd0iVzyq29khiBredWK/HhYnqTIdxWTA=="
	c.DecryptWXOpenData(key, data, iv)

	c.AppID  = "wx0a8581e498061282"
	c.Key = "c45646ffea88cf98e4285dbf75893ff5"

	c.DecryptWXOpenData(key, data, iv)
}

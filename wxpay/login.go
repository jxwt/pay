package wxpay

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jxwt/tools"
	"io/ioutil"
	"net/http"
)

type Errs struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type UserOpenInfo struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
	Errs
}

type RespWXSmall struct {
	Openid     string `json:"openid"`      //用户唯一标识
	Sessionkey string `json:"session_key"` //会话密钥
	Unionid    string `json:"unionid"`     //用户在开放平台的唯一标识符，在满足 UnionID 下发条件的情况下会返回，详见 UnionID 机制说明。
	Errcode    int    `json:"errcode"`     //错误码
	ErrMsg     string `json:"errMsg"`      //错误信息
}

type WxLoginInfoResult struct {
	OpenID     string `json:"openid"`
	NickName   string `json:"nickName"`
	Gender     uint8  `json:"gender"`
	Sex        int    `json:"sex"`
	Language   string `json:"language"`
	City       string `json:"city"`
	Province   string `json:"province"`
	Country    string `json:"country"`
	AvatarURL  string `json:"headimgurl"`
	UnionID    string `json:"unionId"`
	HeadImgUrl string `json:"headimgurl"`
	Watermark  struct {
		Timestamp int    `json:"timestamp"`
		Appid     string `json:"appid"`
	} `json:"watermark"`
	Errs
}

type WxSession struct {
	SessionKey string `json:"session_key"`
	Openid     string `json:"openid"`
	Unionid    string `json:"unionid"`
}

type WxAppLoginAccessResult struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
}

type WxLoginGetPhone struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	Watermark       struct {
		Timestamp int    `json:"timestamp"`
		Appid     string `json:"appid"`
	} `json:"watermark"`
}

/**
微信小程序登陆获取
*/
func (i *WxClient) MiniLogin(code string) (wxInfo RespWXSmall, err error) {
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	resp, err := http.Get(fmt.Sprintf(url, i.AppID, i.Key, code))
	if err != nil {
		return wxInfo, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &wxInfo)
	if err != nil {
		return wxInfo, err
	}
	if wxInfo.Errcode != 0 {
		return wxInfo, errors.New(fmt.Sprintf("code: %d, errmsg: %s", wxInfo.Errcode, wxInfo.ErrMsg))
	}
	return wxInfo, nil
}

func (i *WxClient) GetOpenSession(jsCode string) *WxSession {
	body, _ := tools.HttpBeegoPost("https://api.weixin.qq.com/sns/jscode2session", map[string]string{
		"appid":      i.AppID,
		"secret":     string(i.PrivateKey),
		"js_code":    jsCode,
		"grant_type": "authorization_code",
	}, nil)
	var wxSession WxSession
	json.Unmarshal(body, &wxSession)
	return &wxSession
}

func (i *WxClient) GetLoginInfo(iv string, encryptData string, code string) (WxLoginInfoResult, error) {
	session := i.GetOpenSession(code)
	if session.SessionKey == "" {
		var wxLoginInfoResult WxLoginInfoResult
		return wxLoginInfoResult, errors.New("session获取不到")
	}
	return i.DecryptWXOpenData(session.SessionKey, encryptData, iv)
}

func (i *WxClient) DecryptWXOpenData(sessionKey, encryptData, iv string) (WxLoginInfoResult, error) {
	var wxLoginInfoResult WxLoginInfoResult
	decodeBytes, err := base64.StdEncoding.DecodeString(encryptData)
	if err != nil {
		return wxLoginInfoResult, err
	}
	sessionKeyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return wxLoginInfoResult, err
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return wxLoginInfoResult, err
	}
	dataBytes, err := i.AesDecrypt(decodeBytes, sessionKeyBytes, ivBytes)

	d := tools.UnicodeEmojiCode(string(dataBytes))
	err = json.Unmarshal([]byte(d), &wxLoginInfoResult)
	logs.Warning(wxLoginInfoResult, err)
	if wxLoginInfoResult.Watermark.Appid != i.AppID {
		return wxLoginInfoResult, fmt.Errorf("invalid appid, get !%s!", wxLoginInfoResult.Watermark.Appid)
	}
	if err != nil {
		return wxLoginInfoResult, err
	}
	return wxLoginInfoResult, nil

}

func (i *WxClient) AesDecrypt(crypted, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	return i.PKCS7UnPadding(origData), nil
}

func (i *WxClient) PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	if length > 0 {
		unPadding := int(plantText[length-1])
		return plantText[:(length - unPadding)]

	}
	return plantText
}

func (i *WxClient) GetPhoneNumber(iv string, encryptData string, code string) (WxLoginGetPhone, error) {
	session := i.GetOpenSession(code)
	if session == nil {
		var wxLoginInfoResult WxLoginGetPhone
		return wxLoginInfoResult, errors.New("session获取不到")
	}
	return i.DecryptWXPhone(session.SessionKey, encryptData, iv)
}

func (i *WxClient) DecryptWXPhone(sessionKey, encryptData, iv string) (WxLoginGetPhone, error) {
	var wxLoginInfoResult WxLoginGetPhone
	decodeBytes, err := base64.StdEncoding.DecodeString(encryptData)
	if err != nil {
		return wxLoginInfoResult, err
	}
	sessionKeyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return wxLoginInfoResult, err
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return wxLoginInfoResult, err
	}
	dataBytes, err := i.AesDecrypt(decodeBytes, sessionKeyBytes, ivBytes)

	err = json.Unmarshal(dataBytes, &wxLoginInfoResult)
	logs.Warning(wxLoginInfoResult, err)
	if wxLoginInfoResult.Watermark.Appid != i.AppID {
		return wxLoginInfoResult, fmt.Errorf("invalid appid, get %s", wxLoginInfoResult.Watermark.Appid)
	}
	if err != nil {
		return wxLoginInfoResult, err
	}
	return wxLoginInfoResult, nil
}

//微信app登录
func (i *WxClient) AppLogin(code string) (*WxLoginInfoResult, error) {
	res, _ := tools.HttpBeegoPost("https://api.weixin.qq.com/sns/oauth2/access_token", map[string]string{
		"appid":      i.AppID,
		"secret":     i.Key,
		"code":       code,
		"grant_type": "authorization_code",
	}, nil)
	wxAppLoginAccessResult := new(WxAppLoginAccessResult)
	err := json.Unmarshal(res, wxAppLoginAccessResult)
	if err != nil {
		return nil, err
	}
	res, _ = tools.HttpBeegoPost("https://api.weixin.qq.com/sns/userinfo", map[string]string{
		"access_token": wxAppLoginAccessResult.AccessToken,
		"openid":       wxAppLoginAccessResult.Openid,
	}, nil)
	wxLoginInfoResult := new(WxLoginInfoResult)
	err = json.Unmarshal(res, wxLoginInfoResult)
	if err != nil {
		return nil, err
	}
	return wxLoginInfoResult, nil
}

// 微信公众号:code换取token和openId
func (i *WxClient) GetUserOpenId(code string) (*UserOpenInfo, error) {
	url := "https://api.weixin.qq.com/sns/oauth2/access_token?"
	url += "appid=" + i.AppID
	url += "&secret=" + i.Key
	url += "&code=" + code
	url += "&grant_type=authorization_code"

	body, err := tools.HttpBeegoGet(url, nil)
	if err != nil {
		return nil, err
	}
	user := new(UserOpenInfo)

	err = json.Unmarshal(body, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 微信公众号:获取用户详细信息
func (i *WxClient) GetUserInfoByOpenId(openId string, token string) (*WxLoginInfoResult, error) {
	//获取用户信息
	userUrl := "https://api.weixin.qq.com/sns/userinfo?"
	userUrl += "access_token=" + token
	userUrl += "&openid=" + openId
	userUrl += "&lang=zh_CN"

	userBody, _ := tools.HttpBeegoGet(userUrl, nil)
	userInfo := new(WxLoginInfoResult)

	err := json.Unmarshal(userBody, userInfo)
	if err != nil {
		return nil, err
	}
	if userInfo.ErrMsg != "" {
		return nil, errors.New(userInfo.ErrMsg)
	}
	return userInfo, nil
}

// 获取token，只能在通过model/wxpublic里的方法调用
func (i *WxClient) GetAccessToken() (string, error) {
	var resp struct {
		AccessToken string `json:"access_token"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential" + "&appid=" + i.AppID + "&secret=" + i.Key
	body, err := tools.HttpBeegoGet(url, nil)
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", err
	}

	if resp.ErrMsg != "" {
		return "", errors.New("wxpublicAccessToken获取失败" + resp.ErrMsg)
	}

	return resp.AccessToken, nil
}

func (i *WxClient) GetTicket(token string) (string, error) {
	var ret struct {
		Ticket    string `json:"ticket"`
		Errorcode int    `json:"errorcode"`
		Errmsg    string `json:"errmsg"`
	}
	url := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=" + token + "&type=jsapi"
	body, err := tools.HttpBeegoGet(url, nil)
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return "", err
	}

	if ret.Errmsg != "ok" || ret.Ticket == "" {
		return "", errors.New("获取ticket失败" + ret.Errmsg)
	}

	return ret.Ticket, nil
}

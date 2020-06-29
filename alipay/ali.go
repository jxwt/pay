package alipay

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jxwt/pay"
	"github.com/jxwt/tools"
	"reflect"
	"strconv"
	"strings"
)

type AliClient struct {
	AppId      string
	PublicKey  string
	PrivateKey string
	PartnerId  string
	NotifyURL  string
	Client     *AliAppClient
}

// 用户信息返回
type UserInfo struct {
	AliPaySystemOauthTokenResponse struct {
		AccessToken  string `json:"access_token"`
		AliPayUserID string `json:"AliPay_user_id"`
		ExpiresIn    int    `json:"expires_in"`
		ReExpiresIn  int    `json:"re_expires_in"`
		RefreshToken string `json:"refresh_token"`
		UserID       string `json:"user_id"`
	} `json:"AliPay_system_oauth_token_response"`
	Sign string `json:"sign"`
}

// 用户信息详细内容
type UserInfoDetail struct {
	AliPayUserInfoShareResponse struct {
		Code               string `json:"code"`
		Msg                string `json:"msg"`
		SubCode            string `json:"sub_code"`
		SubMsg             string `json:"sub_msg"`
		UserID             string `json:"user_id"`
		Avatar             string `json:"avatar"`
		Province           string `json:"province"`
		City               string `json:"city"`
		NickName           string `json:"nick_name"`
		IsStudentCertified string `json:"is_student_certified"`
		UserType           string `json:"user_type"`
		UserStatus         string `json:"user_status"`
		IsCertified        string `json:"is_certified"`
		Gender             string `json:"gender"`
		AuthToken          string `json:"authToken"`
	} `json:"AliPay_user_info_share_response"`
	Sign string `json:"sign"`
}

// 创建订单返回
type TradeCreateResult struct {
	AliPayTradeCreateResponse struct {
		Code       string `json:"code"`
		Msg        string `json:"msg"`
		OutTradeNo string `json:"out_trade_no"`
		TradeNo    string `json:"trade_no"`
	} `json:"alipay_trade_create_response"`
	Sign string `json:"sign"`
}

func Init(appId string, loginPid string, privateKey string, publicKey string, callBack string) *AliClient {
	aliPay := new(AliClient)
	aliPay.Client = &AliAppClient{
		SellerID:   loginPid,
		AppID:      appId,
		PrivateKey: GetRsaPrivateKey(privateKey),
		PublicKey:  GetRsaPublicKey(publicKey),
	}
	aliPay.NotifyURL = callBack
	return aliPay
}

func GetRsaPublicKey(key string) *rsa.PublicKey {
	key = "-----BEGIN PUBLIC KEY-----\n" + key + "\n-----END PUBLIC KEY-----"
	b, _ := pem.Decode([]byte(key))
	if b == nil {
		fmt.Println("rsaSign public_key error")
		return nil
	}
	rsaKey, _ := x509.ParsePKIXPublicKey(b.Bytes)
	return rsaKey.(*rsa.PublicKey)
}

func GetRsaPrivateKey(key string) *rsa.PrivateKey {
	key = "-----BEGIN PRIVATE KEY-----\n" + key + "\n-----END PRIVATE KEY-----"
	b, _ := pem.Decode([]byte(key))
	if b == nil {
		fmt.Println("rsaSign private_key error")
		return nil
	}
	rsaKey, err := x509.ParsePKCS8PrivateKey(b.Bytes)
	if err != nil {
		logs.Warning("ParsePKCS8PrivateKey ERR:%v\n", err)
	}
	return rsaKey.(*rsa.PrivateKey)
}

func (i *AliClient) CreatePayOrder(charge *Charge) (*TradeCreateResult, error) {
	if i.Client.PrivateKey == nil {
		return nil, errors.New("privateKey is nil")
	}
	if i.Client.PublicKey == nil {
		return nil, errors.New("publicKey is nil")
	}
	charge.PayMethod = pay.ALI_APP //支付方式
	charge.CallbackURL = i.NotifyURL    //回调地址必须跟下面一样

	res, err := i.Client.CreateOrder(charge)
	if err != nil {
		return nil, err
	}
	aliPayTradeCreateResult := new(TradeCreateResult)
	err = json.Unmarshal([]byte(res), aliPayTradeCreateResult)
	if err != nil {
		return nil, err
	}
	return aliPayTradeCreateResult, nil
}

func (i *AliClient) GetAppPayString(charge *Charge) (string, error) {
	if i.Client.PrivateKey == nil {
		return "", errors.New("privateKey is nil")
	}
	if i.Client.PublicKey == nil {
		return "", errors.New("publicKey is nil")
	}
	charge.PayMethod = pay.ALI_APP //支付方式
	charge.CallbackURL = i.NotifyURL    //回调地址必须跟下面一样
	return i.Client.AppPay(charge)
}

func (i *AliClient) GetAppWapString(charge *Charge) (string, error) {
	if i.Client.PrivateKey == nil {
		return "", errors.New("privateKey is nil")
	}
	if i.Client.PublicKey == nil {
		return "", errors.New("publicKey is nil")
	}
	charge.PayMethod = pay.ALI_H5 //支付方式
	charge.CallbackURL = i.NotifyURL   //回调地址必须跟下面一样
	return i.Client.ToH5Pay(charge)
}

func (i *AliClient) AppLoginUserInfo(code string) *UserInfoDetail {
	AliPayUserInfoDetail := new(UserInfoDetail)
	body, err := i.Client.Login(code)
	if err != nil {
		logs.Warning(err)
		return AliPayUserInfoDetail
	}
	AliPayUserInfo := new(UserInfo)
	logs.Error(body)
	if err := json.Unmarshal([]byte(body), AliPayUserInfo); err != nil {
		logs.Warning("AppLoginUserInfo body:", body)
	}

	AliPayUserInfoDetail.AliPayUserInfoShareResponse.UserID = AliPayUserInfo.AliPaySystemOauthTokenResponse.UserID
	AliPayUserInfoDetail.AliPayUserInfoShareResponse.AuthToken = AliPayUserInfo.AliPaySystemOauthTokenResponse.AccessToken
	return AliPayUserInfoDetail
}

func (i *AliClient) GetAppLoginParams() string {
	return i.Client.GetAppLoginParams(tools.GetUuidRandomString(32))
}

//支付宝收单线下交易
func (i *AliClient) PreCreate(preCreate *Charge) (string, error) {
	preCreate.CallbackURL = i.NotifyURL
	result, err := i.Client.AliPreCreate(*preCreate)
	if err != nil {
		return "", err
	} else {
		return result.QrCode, nil
	}
}


//支付宝退款
func (i *AliClient) Refund(tradeNo string, money float64, tenantId uint, orderId uint, outRequestNo string) (*AliRefundResponse, error) {
	request := new(AliRefundRequest)
	// 支持支付宝交易号退款
	if strings.Contains(tradeNo, "AliPay") || strings.Contains(tradeNo, "AliH5Pay") {
		request.OutTradeNo = tradeNo
	} else {
		request.TradeNo = tradeNo
	}
	// request.OutTradeNo = tradeNo
	request.OutRequestNo = outRequestNo
	request.RefundAmount = fmt.Sprintf("%.2f", money)
	request.OperatorId = strconv.Itoa(int(orderId))
	request.StoreId = strconv.Itoa(int(tenantId))
	return i.Client.Refund(request)
}

//支付宝退款查询
func (i *AliClient) QueryRefund(tradeNo string) (*AliRefundResponse, error) {
	subs := strings.Split(tradeNo, "Refund")
	if len(subs) == 1 {
		return nil, errors.New("无效订单号")
	}
	return i.Client.QueryRefund(subs[1])
}

//支付宝单笔转账
func (i *AliClient) AliSingleRefund(outBizNo, payeeType, payeeAccount, amount, payeeRealName, remark, payerShowName string) (*ToaccountTransferResponse, error) {
	payRefundRequest := ToaccountTransferRequest{
		OutBizNo:      outBizNo,
		PayeeType:     payeeType,
		PayeeAccount:  payeeAccount,
		Amount:        amount, //默认先退0.1元
		PayeeRealName: payeeRealName,
		Remark:        remark,
		PayerShowName: payerShowName,
	}
	return i.Client.ToaccountTransfer(&payRefundRequest)
}

// DecryptOpenDataToStruct 解密支付宝开放数据到 结构体
// encryptedData:包括敏感数据在内的完整用户信息的加密数据
// secretKey:AES密钥，支付宝管理平台配置
// beanPtr:需要解析到的结构体指针
func DecryptOpenDataToStruct(encryptedData, secretKey string, beanPtr interface{}) (err error) {
	beanValue := reflect.ValueOf(beanPtr)
	if beanValue.Kind() != reflect.Ptr {
		return errors.New("传入参数类型必须是以指针形式")
	}
	if beanValue.Elem().Kind() != reflect.Struct {
		return errors.New("传入interface{}必须是结构体")
	}
	var (
		block      cipher.Block
		blockMode  cipher.BlockMode
		originData []byte
	)
	aesKey, _ := base64.StdEncoding.DecodeString(secretKey)
	ivKey := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	secretData, _ := base64.StdEncoding.DecodeString(encryptedData)
	if block, err = aes.NewCipher(aesKey); err != nil {
		return fmt.Errorf("aes.NewCipher：%w", err)
	}
	if len(secretData)%len(aesKey) != 0 {
		return errors.New("encryptedData is error")
	}
	blockMode = cipher.NewCBCDecrypter(block, ivKey)
	originData = make([]byte, len(secretData))
	blockMode.CryptBlocks(originData, secretData)
	if len(originData) > 0 {
		originData = PKCS5UnPadding(originData)
	}
	if err = json.Unmarshal(originData, beanPtr); err != nil {
		return fmt.Errorf("json.Unmarshal(%s)：%w", string(originData), err)
	}
	return nil
}

// 解密填充模式（去除补全码） PKCS5UnPadding
// 解密时，需要在最后面去掉加密时添加的填充byte
func PKCS5UnPadding(origData []byte) (bs []byte) {
	length := len(origData)
	unPaddingNumber := int(origData[length-1]) // 找到Byte数组最后的填充byte
	if unPaddingNumber <= 16 {
		bs = origData[:(length - unPaddingNumber)] // 只截取返回有效数字内的byte数组
	} else {
		bs = origData
	}
	return
}

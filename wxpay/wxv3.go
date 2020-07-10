package wxpay

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jxwt/tools"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// PublicPemNo 聚鑫公钥证书编号
const PublicPemNo = "339B73BB805706D6FB26DBAE1041C923FC135792"

const (
	// WxMediaUploadURL 微信图片上传url
	WxMediaUploadURL = "https://api.mch.weixin.qq.com/v3/merchant/media/upload"
	// WxApplymentURL 提交申请单API
	WxApplymentURL = "https://api.mch.weixin.qq.com/v3/applyment4sub/applyment/"
)

// WxMediaUpLoadHeaderAuthorization 图片上传需要的header
const WxMediaUpLoadHeaderAuthorization = `WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",timestamp="%d",serial_no="%s",signature="%s"`

// WxMediaUpLoadBody 图片上传需要的body
const WxMediaUpLoadBody = "--boundary\r\nContent-Disposition:form-data;name=\"meta\";\r\nContent-Type:application/json\r\n\r\n{\"filename\":\"#file\",\"sha256\":\"#sha256\"}\r\n--boundary\r\nContent-Disposition:form-data;name=\"file\";filename=\"#file\";\r\nContent-Type:image/jpg\r\n\r\n#body\r\n--boundary--\r\n"

// Applyment4subRequest 提交申请单请求
// https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/tool/applyment4sub/chapter3_1.shtml
type Applyment4subRequest struct {
	BusinessCode    string                 `json:"business_code"`     // 业务申请编号
	ContactInfo     *ContactInfoStruct     `json:"contact_info"`      // 超级管理员信息
	SubjectInfo     SubjectInfoStruct      `json:"subject_info"`      // 主体资料
	BusinessInfo    BusinessInfoStruct     `json:"business_info"`     // 经营资料
	SettlementInfo  SettlementInfoStruct   `json:"settlement_info"`   // 结算规则
	BankAccountInfo *BankAccountInfoStruct `json:"bank_account_info"` // 结算银行账户
	AdditionInfo    AdditionInfoStruct     `json:"addition_info"`     // 补充材料
}

// Applyment4sub 申请成为特约商户
func (i *WxClient) Applyment4sub(req *Applyment4subRequest) error {
	// 对结构体内的敏感信息进行加密
	if req.ContactInfo != nil {
		req.ContactInfo = SerialStruct(req.ContactInfo, i.CertPEM).(*ContactInfoStruct)
	}
	if req.BankAccountInfo != nil {
		req.BankAccountInfo = SerialStruct(req.BankAccountInfo, i.CertPEM).(*BankAccountInfoStruct)
	}
	//
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	nonceStr := tools.GetRandomString(32)
	now := time.Now().Unix()
	sign := WxV3Sign("POST", "/v3/applyment4sub/applyment/", nonceStr, string(body), now, i.KeyPEM)
	headerAuthorization := fmt.Sprintf(WxMediaUpLoadHeaderAuthorization, i.AppID, nonceStr, now, i.KeyPemNo, sign)

	client := &http.Client{}
	request, err := http.NewRequest("POST", WxApplymentURL, bytes.NewBuffer(body))
	request.Header.Add("Wechatpay-Serial", i.KeyPemNo)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", headerAuthorization)

	resp, err := client.Do(request)
	if err != nil {
		logs.Warning("http Do err", err)
		return err
	}
	defer resp.Body.Close()
	resultBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return err
	}
	logs.Warning(string(resultBody))

	return nil
}

// RSAPublicSign RSA敏感信息加密
// rsaPublicKey=cert Pem
func RSAPublicSign(rsaPublicKey string, message string) string {
	secretMessage := []byte(message)
	rng := rand.Reader

	cipherdata, err := rsa.EncryptOAEP(sha1.New(), rng, WxPubStrToRSAPublic(rsaPublicKey), secretMessage, nil)
	if err != nil {
		fmt.Printf("Error from encryption: %s\n", err)
		return ""
	}

	ciphertext := base64.StdEncoding.EncodeToString(cipherdata)
	//fmt.Printf("Ciphertext: %x\n", ciphertext)
	return ciphertext
}

// WxPubStrToRSAPublic string转rsa.public格式
func WxPubStrToRSAPublic(publicStr string) *rsa.PublicKey {
	publicByte := FormatCeritficate(publicStr)
	block, _ := pem.Decode(publicByte)
	if block == nil {
		panic("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse DER encoded public key:" + err.Error())
	}
	rsaPublickey, _ := pub.PublicKey.(*rsa.PublicKey)
	return rsaPublickey
}

// WxPriStrToRSAPrivateKey string转rsa.private格式
func WxPriStrToRSAPrivateKey(key string) *rsa.PrivateKey {
	keyByte := FormatPrivateKey(key)
	b, _ := pem.Decode(keyByte)
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

type WxMediaUpLoadRequest struct {
	FileName string `json:"filename"`
	Sha256   string `json:"sha256"`
}

// WxMediaUpLoad 微信图片上传
func (i *WxClient) WxMediaUpLoad(file string, fileName string) (string,error) {
	// 对图片文件进行sha256计算
	h := sha256.New()
	h.Write([]byte(file))
	picSha256 := h.Sum(nil)

	timestamp := time.Now().Unix()
	nonceStr := tools.GetRandomString(32)
	//nonceStr := "L4s1UE5KRC2p5Kh30Kh8GAfTxpSGXXMd"
	// 请求构建 请求加签
	req := &WxMediaUpLoadRequest{
		FileName: fileName,
	}
	req.Sha256 = fmt.Sprintf("%x", string(picSha256))
	body, err := json.Marshal(req)
	if err != nil {
		logs.Warning("WxMediaUpLoad Marshal err", err)
		return "",err
	}
	sign := WxV3Sign("POST", `/v3/merchant/media/upload`, nonceStr, string(body), timestamp, i.KeyPEM)
	// 请求构建与发送
	headerAuthorization := fmt.Sprintf(WxMediaUpLoadHeaderAuthorization, i.AppID, nonceStr, timestamp, i.KeyPemNo, sign)

	reqBody := strings.ReplaceAll(WxMediaUpLoadBody, "#file", fileName)
	reqBody = strings.ReplaceAll(reqBody, "#sha256", req.Sha256)
	reqBody = strings.ReplaceAll(reqBody, "#body", file)

	client := &http.Client{}
	request, err := http.NewRequest("POST", WxMediaUploadURL, bytes.NewBuffer([]byte(reqBody)))
	request.Header.Add("Authorization", headerAuthorization)
	request.Header.Add("Content-Type", "multipart/form-data;boundary=boundary")
	request.Header.Set("Accept", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		logs.Warning("http Do err", err)
		return "",err
	}
	defer resp.Body.Close()
	resultBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return "",err
	}
	return string(resultBody),nil
}

// WxV3Sign 微信v3构建签名串
func WxV3Sign(method string, uri string, nonceStr string, body string, timestemp int64, privateKey string) string {
	// 构建签名meta
	pre := "%s\n%s\n%d\n%s\n%s\n" // method uri timestemp randomstr body
	pre = fmt.Sprintf(pre, method, uri, timestemp, nonceStr, body)
	blocks, _ := pem.Decode(FormatPrivateKey(privateKey))
	key, _ := x509.ParsePKCS8PrivateKey(blocks.Bytes)
	//st := key.(*rsa.PrivateKey)
	h := sha256.New()
	h.Write([]byte(pre))
	digest := h.Sum(nil)
	s, _ := rsa.SignPKCS1v15(rand.Reader, key.(*rsa.PrivateKey), crypto.SHA256, digest)
	return base64.StdEncoding.EncodeToString(s)

}

func (i *WxClient) WxV3GetCertificates() {
	nonceStr := tools.GetRandomString(32)
	now := time.Now().Unix()
	sign := WxV3Sign("GET", "/v3/certificates", nonceStr, "", now, i.KeyPEM)
	fmt.Println(sign)
	headerAuthorization := fmt.Sprintf(WxMediaUpLoadHeaderAuthorization, i.AppID, nonceStr, now, i.KeyPemNo, sign)
	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://api.mch.weixin.qq.com/v3/certificates", nil)
	request.Header.Add("Authorization", headerAuthorization)
	request.Header.Add("User-Agent", "https://zh.wikipedia.org/wiki/User_agent")
	request.Header.Set("Accept", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		logs.Warning("http Do err", err)
		return
	}
	defer resp.Body.Close()
	resultBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Info(string(resultBody))
}

// SerialStruct 对结构体内敏感信息进行加密
func SerialStruct(obj interface{}, rasPublic string) interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		// 传入的inStructPtr是指针，需要.Elem()取得指针指向的value
		t = t.Elem()
		v = v.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("serial")
		if tag == "1" {
			msg := v.Field(i).Interface().(string)
			msg = RSAPublicSign(rasPublic, msg)
			v.Field(i).Set(reflect.ValueOf(msg))
		}
	}
	return obj
}

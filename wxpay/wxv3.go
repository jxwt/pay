package wxpay

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jxwt/tools"
	"io/ioutil"
	"net/http"
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
	// WxApplymentCheckURL 查询申请状态
	WxApplymentCheckURL = "https://api.mch.weixin.qq.com/v3/applyment4sub/applyment/business_code/"
	// GetCertificatesURL 获取证书接口
	GetCertificatesURL = "https://api.mch.weixin.qq.com/v3/certificates"
)

// WxMediaUpLoadHeaderAuthorization 图片上传需要的header
const WxMediaUpLoadHeaderAuthorization = `WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",timestamp="%d",serial_no="%s",signature="%s"`

// WxMediaUpLoadBody 图片上传需要的body
const WxMediaUpLoadBody = "--boundary\r\nContent-Disposition:form-data;name=\"meta\";\r\nContent-Type:application/json\r\n\r\n{\"filename\":\"#file\",\"sha256\":\"#sha256\"}\r\n--boundary\r\nContent-Disposition:form-data;name=\"file\";filename=\"#file\";\r\nContent-Type:image/jpg\r\n\r\n#body\r\n--boundary--\r\n"

// Applyment4sub 申请成为特约商户
func (i *WxClient) Applyment4sub(req *Applyment4subRequest) error {
	// 获取证书,证书解析
	res, _ := i.GetCertificates()
	if len(res.Data) == 0 {
		return errors.New("证书获取失败")
	}
	i.SerialNo = res.Data[0].SerialNo
	i.Ciphertext, _ = CertificateDecryption(res, i.Key)
	// 对结构体内的敏感信息进行加密
	if req.ContactInfo != nil {
		req.ContactInfo = SerialStruct(req.ContactInfo, i.Ciphertext).(*ContactInfoStruct)
	}
	if req.BankAccountInfo != nil {
		req.BankAccountInfo = SerialStruct(req.BankAccountInfo, i.Ciphertext).(*BankAccountInfoStruct)
	}
	if req.SubjectInfo != nil {
		if req.SubjectInfo.IdentityInfo.IDCardInfo != nil {
			req.SubjectInfo.IdentityInfo.IDCardInfo = SerialStruct(req.SubjectInfo.IdentityInfo.IDCardInfo, i.Ciphertext).(*IDCardInfoStruct)
		}
	}
	//
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	nonceStr := tools.GetRandomString(32)
	now := time.Now().Unix()
	sign := WxV3Sign("POST", "/v3/applyment4sub/applyment/", nonceStr, string(body), now, i.KeyPEM)
	headerAuthorization := fmt.Sprintf(WxMediaUpLoadHeaderAuthorization, i.MchID, nonceStr, now, i.KeyPemNo, sign)

	client := &http.Client{}
	request, err := http.NewRequest("POST", WxApplymentURL, bytes.NewBuffer(body))
	request.Header.Add("Wechatpay-Serial", i.SerialNo)
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

// WxMediaUpLoad 微信图片上传
func (i *WxClient) WxMediaUpLoad(file string, fileName string) (string, error) {
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
		return "", err
	}
	sign := WxV3Sign("POST", `/v3/merchant/media/upload`, nonceStr, string(body), timestamp, i.KeyPEM)
	// 请求构建与发送
	headerAuthorization := fmt.Sprintf(WxMediaUpLoadHeaderAuthorization, i.MchID, nonceStr, timestamp, i.KeyPemNo, sign)

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
		return "", err
	}
	defer resp.Body.Close()
	resultBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return "", err
	}
	return string(resultBody), nil
}

// WxApplymentCheck 查询
func (i *WxClient) WxApplymentCheck(businessCode string) {

}

// GetCertificates 获取证书
func (i *WxClient) GetCertificates() (*GetCertificatesResponse, error) {
	nonceStr := tools.GetRandomString(32)
	timestamp := time.Now().Unix()
	sign := WxV3Sign("GET", `/v3/certificates`, nonceStr, "", timestamp, i.KeyPEM)
	headerAuthorization := fmt.Sprintf(WxMediaUpLoadHeaderAuthorization, i.MchID, nonceStr, timestamp, i.KeyPemNo, sign)
	client := &http.Client{}
	request, err := http.NewRequest("GET", GetCertificatesURL, bytes.NewBuffer([]byte("")))
	request.Header.Add("Authorization", headerAuthorization)
	request.Header.Add("User-Agent", "https://zh.wikipedia.org/wiki/User_agent")
	request.Header.Set("Accept", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		logs.Warning("http Do err", err)
		return nil, err
	}
	defer resp.Body.Close()
	resultBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	res := &GetCertificatesResponse{}
	json.Unmarshal(resultBody, res)
	return res, nil
}

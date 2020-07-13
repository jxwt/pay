package wxpay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jxwt/tools"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

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

// ContactInfoStruct 超级管理员信息
type ContactInfoStruct struct {
	ContactName     string `json:"contact_name" serial:"1"`      // 超级管理员姓名
	ContactIDNumber string `json:"contact_id_number" serial:"1"` // 超级管理员身份证件号码 (1/2)
	Openid          string `json:"openid"`                       // 超级管理员微信openid (1/2)
	MobilePhone     string `json:"mobile_phone" serial:"1"`      // 联系手机
	ContactEmail    string `json:"contact_email" serial:"1"`     // 联系邮箱
}

// BankAccountInfoStruct 结算银行账户
type BankAccountInfoStruct struct {
	BankAccountType string `json:"bank_account_type"`         // 账户类型
	AccountName     string `json:"account_name" serial:"1"`   // 开户名称
	AccountBank     string `json:"account_bank"`              // 开户银行
	BankAddressCode string `json:"bank_address_code"`         // 开户银行省市编码
	BankBranchID    string `json:"bank_branch_id"`            // 开户银行联行号
	BankName        string `json:"bank_name"`                 // 开户银行全称（含支行)
	AccountNumber   string `json:"account_number" serial:"1"` // 银行账号
}

// BusinessLicenseInfoStruct 营业执照 主体为个体户/企业 必填
type BusinessLicenseInfoStruct struct {
	LicenseCopy   string `json:"license_copy"`   // 营业执照照片
	LicenseNumber string `json:"license_number"` // 注册号/统一社会信用代码
	MerchantName  string `json:"merchant_name"`  // 商户名称
	LegalPerson   string `json:"legal_person"`   // 个体户经营者/法人姓名
}

// CertificateInfoStruct 登记证书 主体为党政、机关及事业单位/其他组织，必填。
type CertificateInfoStruct struct {
	CertCopy       string `json:"cert_copy"`       // 登记证书照片
	CertType       string `json:"cert_type"`       // 登记证书类型
	CertNumber     string `json:"cert_number"`     // 证书号
	MerchantName   string `json:"merchant_name"`   // 商户名称
	CompanyAddress string `json:"company_address"` // 注册地址
	LegalPerson    string `json:"legal_person"`    // 法人姓名
	PeriodBegin    string `json:"period_begin"`    // 有效期限开始日期
	PeriodEnd      string `json:"period_end"`      // 有效期限结束日期
}

// OrganizationInfoStruct 组织机构代码证	主体为企业/党政、机关及事业单位/其他组织，且证件号码不是18位时必填。
type OrganizationInfoStruct struct {
	OrganizationCopy string `json:"organization_copy"` // 组织机构代码证照片
	OrganizationCode string `json:"organization_code"` // 组织机构代码
	OrgPeriodBegin   string `json:"org_period_begin"`  // 组织机构代码证有效期开始日期
	OrgPeriodEnd     string `json:"org_period_end"`    // 组织机构代码证有效期结束日期
}

// IDCardInfoStruct 经营者/法人身份证件
type IDCardInfoStruct struct {
	IDCardCopy      string `json:"id_card_copy"`      // 身份证人像面照片(图片上传接口)
	IDCardNational  string `json:"id_card_national"`  // 身份证国徽面照片
	IDCardName      string `json:"id_card_name"`      // 身份证姓名
	IDCardNumber    string `json:"id_card_number"`    // 身份证号码
	CardPeriodBegin string `json:"card_period_begin"` // 身份证有效期开始时间
	CardPeriodEnd   string `json:"card_period_end"`   // 身份证有效期结束时间
}

// IDDocInfoStruct 其他类型证件信息
type IDDocInfoStruct struct {
	IDDocCopy      string `json:"id_doc_copy"`      // 证件照片
	IDDocName      string `json:"id_doc_name"`      // 证件姓名
	IDDocNumber    string `json:"id_doc_number"`    // 证件号码
	DocPeriodBegin string `json:"doc_period_begin"` // 证件有效期开始时间
	DocPeriodEnd   string `json:"doc_period_end"`   // 证件有效期结束时间
}

// UboInfoStruct 最终受益人信息(UBO)
type UboInfoStruct struct {
	IDType         string `json:"id_type"`          // 证件类型
	IDCardCopy     string `json:"id_card_copy"`     // 身份证人像面照片
	IDCardNational string `json:"id_card_national"` // 身份证国徽面照片
	IDDocCopy      string `json:"id_doc_copy"`      // 证件照片
	Name           string `json:"name"`             // 受益人姓名
	IDNumber       string `json:"id_number"`        // 证件号码
	IDPeriodBegin  string `json:"id_period_begin"`  // 证件有效期开始时间
	IDPeriodEnd    string `json:"id_period_end"`    // 证件有效期结束时间
}

// BizStoreInfoStruct 线下门店场景
type BizStoreInfoStruct struct {
	BizStoreName     string   `json:"biz_store_name"`
	BizAddressCode   string   `json:"biz_address_code"`
	BizStoreAddress  string   `json:"biz_store_address"`
	StoreEntrancePic []string `json:"store_entrance_pic"`
	IndoorPic        []string `json:"indoor_pic"`
	BizSubAppid      string   `json:"biz_sub_appid"`
}

// MpInfo 公众号场景
type MpInfoStruct struct {
	MpAppid    string   `json:"mp_appid"`
	MpSubAppid string   `json:"mp_sub_appid"`
	MpPics     []string `json:"mp_pics"`
}

// MiniProgramInfoStruct 小程序场景
type MiniProgramInfoStruct struct {
	MiniProgramAppid    string   `json:"mini_program_appid"`
	MiniProgramSubAppid string   `json:"mini_program_sub_appid"`
	MiniProgramPics     []string `json:"mini_program_pics"`
}

// AppInfoStruct APP场景
type AppInfoStruct struct {
	AppAppid    string   `json:"app_appid"`
	AppSubAppid string   `json:"app_sub_appid"`
	AppPics     []string `json:"app_pics"`
}

// WebInfoStruct 互联网网站场景
type WebInfoStruct struct {
	Domain           string `json:"domain"`
	WebAuthorisation string `json:"web_authorisation"`
	WebAppid         string `json:"web_appid"`
}

// WeworkInfoStruct 企业微信场景
type WeworkInfoStruct struct {
	CorpID     string   `json:"corp_id"`
	SubCorpID  string   `json:"sub_corp_id"`
	WeworkPics []string `json:"wework_pics"`
}

// SettlementInfoStruct 结算规则
type SettlementInfoStruct struct {
	SettlementID        string   `json:"settlement_id"`        // 入驻结算规则ID
	QualificationType   string   `json:"qualification_type"`   // 所属行业
	Qualifications      []string `json:"qualifications"`       // 特殊资质图片
	ActivitiesID        string   `json:"activities_id"`        // 优惠费率活动ID
	ActivitiesRate      string   `json:"activities_rate"`      // 优惠费率活动值
	ActivitiesAdditions []string `json:"activities_additions"` // 优惠费率活动补充材料
}

// AdditionInfoStruct 补充材料
type AdditionInfoStruct struct {
	LegalPersonCommitment string   `json:"legal_person_commitment"` // 法人开户承诺函
	LegalPersonVideo      string   `json:"legal_person_video"`      // 法人开户意愿视频
	BusinessAdditionPics  []string `json:"business_addition_pics"`  // 补充材料
	BusinessAdditionMsg   string   `json:"business_addition_msg"`   // 补充说明
}

//  SalesInfoStruct 经营场景
type SalesInfoStruct struct {
	SalesScenesType []string              `json:"sales_scenes_type"` // 经营场景类型 小程序：SALES_SCENES_MINI_PROGRAM 互联网：SALES_SCENES_WEB 公众号：SALES_SCENES_MP APP：SALES_SCENES_APP
	BizStoreInfo    BizStoreInfoStruct    `json:"biz_store_info"`    // 线下门店场景
	MpInfo          MpInfoStruct          `json:"mp_info"`           // 公众号场景
	MiniProgramInfo MiniProgramInfoStruct `json:"mini_program_info"` // 小程序场景
	AppInfo         AppInfoStruct         `json:"app_info"`          // APP场景
	WebInfo         WebInfoStruct         `json:"web_info"`          // 互联网网站场景
	WeworkInfo      WeworkInfoStruct      `json:"wework_info"`       // 企业微信场景
}

// BusinessInfoStruct 经营资料
type BusinessInfoStruct struct {
	MerchantShortname string          `json:"merchant_shortname"` // 商户简称
	ServicePhone      string          `json:"service_phone"`      // 客服电话
	SalesInfo         SalesInfoStruct `json:"sales_info"`         // 经营场景
}

// IdentityInfoStruct 经营者/法人身份证件
type IdentityInfoStruct struct {
	IDDocType  string           `json:"id_doc_type"`  // 证件类型
	IDCardInfo IDCardInfoStruct `json:"id_card_info"` // 经营者/法人身份证件
	IDDocInfo  IDDocInfoStruct  `json:"id_doc_info"`  // 其他类型证件信息
	Owner      bool             `json:"owner"`        // 经营者/法人是否为受益人 true false
}

// SubjectInfoStruct 主体资料
type SubjectInfoStruct struct {
	SubjectType           string                    `json:"subject_type"`            // SUBJECT_TYPE_INDIVIDUAL（个体户）SUBJECT_TYPE_ENTERPRISE（企业）SUBJECT_TYPE_INSTITUTIONS（党政、机关及事业单位）SUBJECT_TYPE_OTHERS（其他组织）
	BusinessLicenseInfo   BusinessLicenseInfoStruct `json:"business_license_info"`   // 营业执照 主体为个体户/企业 必填
	CertificateInfo       CertificateInfoStruct     `json:"certificate_info"`        // 登记证书 主体为党政、机关及事业单位/其他组织，必填。
	OrganizationInfo      OrganizationInfoStruct    `json:"organization_info"`       // 组织机构代码证	主体为企业/党政、机关及事业单位/其他组织，且证件号码不是18位时必填。
	CertificateLetterCopy string                    `json:"certificate_letter_copy"` // 单位证明函照片
	IdentityInfo          IdentityInfoStruct        `json:"identity_info"`           // 经营者/法人身份证件
	UboInfo               UboInfoStruct             `json:"ubo_info"`                // 最终受益人信息(UBO)
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

type WxMediaUpLoadRequest struct {
	FileName string `json:"filename"`
	Sha256   string `json:"sha256"`
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

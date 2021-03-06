package wxpay

import (
	"io/ioutil"
	"testing"
)

var publicKey = "MIID8DCCAtigAwIBAgIUM5tzu4BXBtb7JtuuEEHJI/wTV5IwDQYJKoZIhvcNAQEL\nBQAwXjELMAkGA1UEBhMCQ04xEzARBgNVBAoTClRlbnBheS5jb20xHTAbBgNVBAsT\nFFRlbnBheS5jb20gQ0EgQ2VudGVyMRswGQYDVQQDExJUZW5wYXkuY29tIFJvb3Qg\nQ0EwHhcNMjAwNjEwMDcyNjE2WhcNMjUwNjA5MDcyNjE2WjCBgTETMBEGA1UEAwwK\nMTU5NzM5ODMxMTEbMBkGA1UECgwS5b6u5L+h5ZWG5oi357O757ufMS0wKwYDVQQL\nDCTmna3lt57ogZrpkavlkL7lkIznp5HmioDmnInpmZDlhazlj7gxCzAJBgNVBAYM\nAkNOMREwDwYDVQQHDAhTaGVuWmhlbjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC\nAQoCggEBAOOCiw4Du4bB8fGdI3oPccv6vxq5u45+Y/2/cUVA59diir2if6XzcKI1\nIJGgb9l2DpMfUxml2kO5pJy8gSRrhxXVHQxkzqw5IvUMyDw68E2vg+jz/LnwXwBM\nuu5fKefR8b43E6EMWbVIUGyX9ZCW//sD5M2zq+iwDafjGgOzhNWugJA8isH+mF5b\n4KPELX3kNOP1jPHNLXsONiI+zxli7hkxAlzC2u+lgY2bDG9nfQMBIhcmJtGAEMOZ\nAw0RK2r2KUnk7NhuSbQtmD4UhX3J7H/F/G0yLCyeerckeIw1b6LWyuK/0BoXBXQ8\nluNKh9R8RVTxJRpMZqxdnxd4AWKekkECAwEAAaOBgTB/MAkGA1UdEwQCMAAwCwYD\nVR0PBAQDAgTwMGUGA1UdHwReMFwwWqBYoFaGVGh0dHA6Ly9ldmNhLml0cnVzLmNv\nbS5jbi9wdWJsaWMvaXRydXNjcmw/Q0E9MUJENDIyMEU1MERCQzA0QjA2QUQzOTc1\nNDk4NDZDMDFDM0U4RUJEMjANBgkqhkiG9w0BAQsFAAOCAQEAuQxg2etUynRX5fMb\n8C3gHRoarOtJ+W6L7xDLz2+r1VboOqXcFbmFA2ufjX/NCHDzzDCuz5ZjdoNtNK9E\nWTZ5TZqQ2+jORK6sgf0LL3pRKZNtj0skXx1C+HFtdBq3Ncx35nisR940V5sSoxB8\nH36Jc7DD5wuz7SHYOn+Do5xkiVhU9Q6sOoBXqmvwFgerp1KfIi9XMmFTq+pBq3ab\nrdmo8TG5By8wsqnr1uWPvIr3ixvgCkyFHWYXHQxqJ7QLPdo1g9Ig21hfID7PIayr\nYw5h2rNuYzlQtBp7DOn0HYX0rFE3olK/p32BLCeZIu9XiwR/vg0pNrhakXTaeafK\nhjYAZA=="
var keyPem = "MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDjgosOA7uGwfHx\nnSN6D3HL+r8aubuOfmP9v3FFQOfXYoq9on+l83CiNSCRoG/Zdg6TH1MZpdpDuaSc\nvIEka4cV1R0MZM6sOSL1DMg8OvBNr4Po8/y58F8ATLruXynn0fG+NxOhDFm1SFBs\nl/WQlv/7A+TNs6vosA2n4xoDs4TVroCQPIrB/pheW+CjxC195DTj9YzxzS17DjYi\nPs8ZYu4ZMQJcwtrvpYGNmwxvZ30DASIXJibRgBDDmQMNEStq9ilJ5OzYbkm0LZg+\nFIV9yex/xfxtMiwsnnq3JHiMNW+i1sriv9AaFwV0PJbjSofUfEVU8SUaTGasXZ8X\neAFinpJBAgMBAAECggEAW0vYd1BCIWqUp2tygBnQhZViuNvNivnnMD1xu+O25XSy\nzjR2WubczQravfWOzMoWQS2x0DoA42qMxyTSAgZwV++ET6PoV645+/IcLCdOpS4I\nliPKx+bQiLNB1EQ18cQK6VT6uIbXPOr+8wTr0xD1Ogqu92jhVGfJoxR8LP4OERzK\nzZUsHOyhgIJbFcPUEIKYP7lYL/MPPN4JvpKr9YXbzKcFIDfcbFV2R/Hk1STevtxv\nIdw24tCeCN44GUoNvDmBQAnxPakTB9BmBzF5C1IbEQe3BCDopQPW+9qg8jHW0CmS\n5/pb30BLxOKcLNf2smuqDAZsnD8LQM1cEEsFF58Y3QKBgQD4pl1Qy9KCGRFj11yL\nKtQ6gUeywJmr0xEhpN2SVbqVd3vn+rVW89zqAwh5U/8YSwYP5P96AKSYrNRn/sAD\nOWPtxwQUyALWPQHBOI+t57JuYCp8IHmfjHr6OtNlBkqaUCrvWDau2LRTgxb/Dfj6\nGJrKT1pNzHN/YTyKgxFdsKrO6wKBgQDqPDREdsD/LgSSE4YYFezXBJoXvR5LNqaG\n98YUHz+h1nL7ENSm2/um0e6ukSBhfVIV0ur8wzv3PgEOMmty6SV3rs6R9lQ37qKR\n0zUn09CftIYEbvJ12SWxLdpemFYnI5LNEhx1aVVFe00BnslCfBbH0WtDt5plMAD1\nc/w2KNgQgwKBgG5Oz8MSSRcyK8bROdr7ax9xTu98BjB1+HmmfC15HsdENJHbZSto\nEC84nT/GBbsvPUc73iKvulWJBsoD+Ab2JODNk3/so2WLtwWTJBqQWVYiD3b1qT8g\nwUXVZwbAXcRLoGCCD/BNbuJFm6QW/MdmtvTdc0BkXTC7YHJKZx/bSkt1AoGAf567\nr42wS5hL/zbJ+beAagpk1ohAyCQHiUPYVUBNUCTiUq5x3lO/Ab4huFTz+onoPmHD\njGHm+yd6Nbz81Af5VQMWI2q9qhfH1YHo1UFPyqP13NaCHflo0uczshR35C06n6a9\ngK8aOZgbdcWIzEOFuer88VFIutbzvsgp42xPhHcCgYEAhkDfc7bQs2HHUtKnsL2h\nOSOn6s9UiEXsU7AcJySHwrWlzgBnE50C5WUV4W9s/EBBKohuP+tvWOj2Xh7JMShp\ndDEbT0ebNm9CvV/mdeE7Rwitg+WWvTloHyYW81doBVo+MG3qSQIrzUzzBhWq5SiV\no4lqLy4CJKDro6Yl7Mlp0xM="

func TestWxClient_WxMediaUpLoad(t *testing.T) {
	c := new(WxClient)
	c.KeyPEM = "MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDjgosOA7uGwfHx\nnSN6D3HL+r8aubuOfmP9v3FFQOfXYoq9on+l83CiNSCRoG/Zdg6TH1MZpdpDuaSc\nvIEka4cV1R0MZM6sOSL1DMg8OvBNr4Po8/y58F8ATLruXynn0fG+NxOhDFm1SFBs\nl/WQlv/7A+TNs6vosA2n4xoDs4TVroCQPIrB/pheW+CjxC195DTj9YzxzS17DjYi\nPs8ZYu4ZMQJcwtrvpYGNmwxvZ30DASIXJibRgBDDmQMNEStq9ilJ5OzYbkm0LZg+\nFIV9yex/xfxtMiwsnnq3JHiMNW+i1sriv9AaFwV0PJbjSofUfEVU8SUaTGasXZ8X\neAFinpJBAgMBAAECggEAW0vYd1BCIWqUp2tygBnQhZViuNvNivnnMD1xu+O25XSy\nzjR2WubczQravfWOzMoWQS2x0DoA42qMxyTSAgZwV++ET6PoV645+/IcLCdOpS4I\nliPKx+bQiLNB1EQ18cQK6VT6uIbXPOr+8wTr0xD1Ogqu92jhVGfJoxR8LP4OERzK\nzZUsHOyhgIJbFcPUEIKYP7lYL/MPPN4JvpKr9YXbzKcFIDfcbFV2R/Hk1STevtxv\nIdw24tCeCN44GUoNvDmBQAnxPakTB9BmBzF5C1IbEQe3BCDopQPW+9qg8jHW0CmS\n5/pb30BLxOKcLNf2smuqDAZsnD8LQM1cEEsFF58Y3QKBgQD4pl1Qy9KCGRFj11yL\nKtQ6gUeywJmr0xEhpN2SVbqVd3vn+rVW89zqAwh5U/8YSwYP5P96AKSYrNRn/sAD\nOWPtxwQUyALWPQHBOI+t57JuYCp8IHmfjHr6OtNlBkqaUCrvWDau2LRTgxb/Dfj6\nGJrKT1pNzHN/YTyKgxFdsKrO6wKBgQDqPDREdsD/LgSSE4YYFezXBJoXvR5LNqaG\n98YUHz+h1nL7ENSm2/um0e6ukSBhfVIV0ur8wzv3PgEOMmty6SV3rs6R9lQ37qKR\n0zUn09CftIYEbvJ12SWxLdpemFYnI5LNEhx1aVVFe00BnslCfBbH0WtDt5plMAD1\nc/w2KNgQgwKBgG5Oz8MSSRcyK8bROdr7ax9xTu98BjB1+HmmfC15HsdENJHbZSto\nEC84nT/GBbsvPUc73iKvulWJBsoD+Ab2JODNk3/so2WLtwWTJBqQWVYiD3b1qT8g\nwUXVZwbAXcRLoGCCD/BNbuJFm6QW/MdmtvTdc0BkXTC7YHJKZx/bSkt1AoGAf567\nr42wS5hL/zbJ+beAagpk1ohAyCQHiUPYVUBNUCTiUq5x3lO/Ab4huFTz+onoPmHD\njGHm+yd6Nbz81Af5VQMWI2q9qhfH1YHo1UFPyqP13NaCHflo0uczshR35C06n6a9\ngK8aOZgbdcWIzEOFuer88VFIutbzvsgp42xPhHcCgYEAhkDfc7bQs2HHUtKnsL2h\nOSOn6s9UiEXsU7AcJySHwrWlzgBnE50C5WUV4W9s/EBBKohuP+tvWOj2Xh7JMShp\ndDEbT0ebNm9CvV/mdeE7Rwitg+WWvTloHyYW81doBVo+MG3qSQIrzUzzBhWq5SiV\no4lqLy4CJKDro6Yl7Mlp0xM="
	c.MchID = "1597398311"
	buf, err := ioutil.ReadFile("/Users/hcf/Pictures/商城图片/10.png")
	if err != nil {
		panic(err)
	}
	c.KeyPemNo = "339B73BB805706D6FB26DBAE1041C923FC135792"
	str, _ := c.WxMediaUpLoad(string(buf), "10.png")
	t.Logf("%v", str)
}

func TestSerialStruct(t *testing.T) {
	obj := &ContactInfoStruct{
		ContactName:     "黄晨帆",                //
		ContactIDNumber: "362204199502019999", //
		Openid:          "erfn2jif2e2sf2",
		MobilePhone:     "13979459983",     //
		ContactEmail:    "eeeeeee@163.com", //
	}
	publicKey := "MIID8DCCAtigAwIBAgIUM5tzu4BXBtb7JtuuEEHJI/wTV5IwDQYJKoZIhvcNAQEL\nBQAwXjELMAkGA1UEBhMCQ04xEzARBgNVBAoTClRlbnBheS5jb20xHTAbBgNVBAsT\nFFRlbnBheS5jb20gQ0EgQ2VudGVyMRswGQYDVQQDExJUZW5wYXkuY29tIFJvb3Qg\nQ0EwHhcNMjAwNjEwMDcyNjE2WhcNMjUwNjA5MDcyNjE2WjCBgTETMBEGA1UEAwwK\nMTU5NzM5ODMxMTEbMBkGA1UECgwS5b6u5L+h5ZWG5oi357O757ufMS0wKwYDVQQL\nDCTmna3lt57ogZrpkavlkL7lkIznp5HmioDmnInpmZDlhazlj7gxCzAJBgNVBAYM\nAkNOMREwDwYDVQQHDAhTaGVuWmhlbjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC\nAQoCggEBAOOCiw4Du4bB8fGdI3oPccv6vxq5u45+Y/2/cUVA59diir2if6XzcKI1\nIJGgb9l2DpMfUxml2kO5pJy8gSRrhxXVHQxkzqw5IvUMyDw68E2vg+jz/LnwXwBM\nuu5fKefR8b43E6EMWbVIUGyX9ZCW//sD5M2zq+iwDafjGgOzhNWugJA8isH+mF5b\n4KPELX3kNOP1jPHNLXsONiI+zxli7hkxAlzC2u+lgY2bDG9nfQMBIhcmJtGAEMOZ\nAw0RK2r2KUnk7NhuSbQtmD4UhX3J7H/F/G0yLCyeerckeIw1b6LWyuK/0BoXBXQ8\nluNKh9R8RVTxJRpMZqxdnxd4AWKekkECAwEAAaOBgTB/MAkGA1UdEwQCMAAwCwYD\nVR0PBAQDAgTwMGUGA1UdHwReMFwwWqBYoFaGVGh0dHA6Ly9ldmNhLml0cnVzLmNv\nbS5jbi9wdWJsaWMvaXRydXNjcmw/Q0E9MUJENDIyMEU1MERCQzA0QjA2QUQzOTc1\nNDk4NDZDMDFDM0U4RUJEMjANBgkqhkiG9w0BAQsFAAOCAQEAuQxg2etUynRX5fMb\n8C3gHRoarOtJ+W6L7xDLz2+r1VboOqXcFbmFA2ufjX/NCHDzzDCuz5ZjdoNtNK9E\nWTZ5TZqQ2+jORK6sgf0LL3pRKZNtj0skXx1C+HFtdBq3Ncx35nisR940V5sSoxB8\nH36Jc7DD5wuz7SHYOn+Do5xkiVhU9Q6sOoBXqmvwFgerp1KfIi9XMmFTq+pBq3ab\nrdmo8TG5By8wsqnr1uWPvIr3ixvgCkyFHWYXHQxqJ7QLPdo1g9Ig21hfID7PIayr\nYw5h2rNuYzlQtBp7DOn0HYX0rFE3olK/p32BLCeZIu9XiwR/vg0pNrhakXTaeafK\nhjYAZA=="
	a := SerialStruct(obj, publicKey)
	t.Logf("%v", a)
}

func TestApplyment4sub(t *testing.T) {
	wxClient := &WxClient{
		MchID:    "1597398311",
		CertPEM:  publicKey,
		KeyPEM:   keyPem,
		Key:      "Hxf8DV9q21Zi4YYNBpBwpg4Ne1qQqRWN",
		KeyPemNo: "339B73BB805706D6FB26DBAE1041C923FC135792",
	}
	contactInfo := &ContactInfoStruct{
		ContactName:     "黄晨帆",                //
		ContactIDNumber: "362204199502019999", //
		Openid:          "erfn2jif2e2sf2",
		MobilePhone:     "13979459983",     //
		ContactEmail:    "eeeeeee@163.com", //
	}
	req := &Applyment4subRequest{
		BusinessCode: "ssss11",
		ContactInfo:  contactInfo,
	}
	res, err := wxClient.Applyment4sub(req)
	if err != nil {
		t.Logf("%v", err)
	}
	t.Logf("%+v", res)
}

func TestGetCertificates(t *testing.T) {
	wxClient := &WxClient{
		MchID:    "1597398311",
		CertPEM:  publicKey,
		KeyPEM:   keyPem,
		KeyPemNo: "339B73BB805706D6FB26DBAE1041C923FC135792",
	}
	res, _ := wxClient.GetCertificates()
	t.Logf("%+v", res)
}

func TestWxApplymentCheck(t *testing.T) {
	wxClient := &WxClient{
		MchID:    "1597398311",
		CertPEM:  publicKey,
		KeyPEM:   keyPem,
		KeyPemNo: "339B73BB805706D6FB26DBAE1041C923FC135792",
	}
	res, err := wxClient.WxApplymentCheck("1594621668EAjzDs8D")
	if err != nil {
		t.Logf("%v", err)
	}
	t.Logf("%+v", res)
}

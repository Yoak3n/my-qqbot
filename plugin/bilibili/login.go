package bilibili

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/skip2/go-qrcode"
	"github.com/tidwall/gjson"
	"io"
	"my-qqbot/config"
	"net/http"
	"net/url"
	"os"
	re "regexp"
	"strings"
	"time"
)

const (
	userAgent    = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36 Edg/97.0.1072.69`
	publicKeyPEM = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDLgd2OAkcGVtoE3ThUREbio0Eg
Uc/prcajMKXvkCKFCWhJYJcLkcM2DKKcSeFpD/j6Boy538YXnR6VhcuUJOhH2x71
nzPjfdTcqMz7djHum0qSZA0AyCBDABUqCrfNgCiJ00Ra7GmRj+YCK1NJEuewlb40
JNrRuoEUXpabUzGB8QIDAQAB
-----END PUBLIC KEY-----
`
)

var (
	CK           string
	refreshToken string
	Scan         chan string
)

func init() {
	CK = ""
}

func GetCookie() string {
	Scan = make(chan string, 1)
	cookie, _ := login()
	config.UpdateBilibiliCookie(cookie, refreshToken)
	return cookie
}

func checkCookieNeedRefresh() (bool, int64, error) {
	uri := "https://passport.bilibili.com/x/passport-login/web/cookie/info?csrf=" + getCsrf()
	client := http.Client{}
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Set("User-Agent", userAgent)
	if CK == "" {
		CK = config.Conf.Bilibili.Cookie
	}
	req.Header.Set("Cookie", CK)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("checkCookieNeedRefresh error:", err)
		return false, 0, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	data := gjson.ParseBytes(body)
	if data.Get("code").Int() != 0 {
		return false, 0, err
	}
	if data.Get("data.refresh").Bool() {
		return true, data.Get("data.timestamp").Int(), nil
	}
	return false, 0, nil
}

func getLoginKeyAndLoginUrl() (loginKey string, loginUrl string) {
	uri := "https://passport.bilibili.com/x/passport-login/web/qrcode/generate"
	client := http.Client{}
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	data := gjson.ParseBytes(body)
	loginKey = data.Get("data.qrcode_key").String()
	loginUrl = data.Get("data.url").String()
	return
}

func verifyLogin(loginKey string) {
	for {
		uri := "https://passport.bilibili.com/x/passport-login/web/qrcode/poll"
		client := http.Client{}
		uri += "?" + "qrcode_key=" + loginKey
		req, _ := http.NewRequest("GET", uri, nil)
		req.Header.Set("User-Agent", userAgent)
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		data := gjson.ParseBytes(body)
		if data.Get("data.url").String() != "" {
			var cookieContent []byte
			cookie := make(map[string]string)
			for _, v := range resp.Header["Set-Cookie"] {
				kv := strings.Split(v, ";")[0]
				kvArr := strings.Split(kv, "=")
				cookie[kvArr[0]] = kvArr[1]
			}
			cookieContent = []byte(`DedeUserID=` + cookie["DedeUserID"] + `;DedeUserID__ckMd5=` + cookie["DedeUserID__ckMd5"] + `;Expires=` + cookie["Expires"] + `;SESSDATA=` + cookie["SESSDATA"] + `;bili_jct=` + cookie["bili_jct"] + `;`)
			CK = string(cookieContent)
			refreshToken = data.Get("data.refresh_token").String()
			break
		}
		time.Sleep(time.Second * 3)
	}
}

func RefreshCookie() error {
	// 获取 refresh_csrf
	refresh, _, err := checkCookieNeedRefresh()
	if err != nil {
		return err
	}
	if !refresh {
		return errors.New("cookie有效")
	}
	refreshCsrf, err := getRefreshCsrf()
	if err != nil {
		return err
	}
	// 获取新cookie
	newRefreshToken, err := refreshCookie(refreshCsrf)
	if err != nil {
		return err
	}
	// 确认更新
	err = commitCookie()
	if err != nil {
		return err
	}
	refreshToken = newRefreshToken
	config.UpdateBilibiliCookie(CK, refreshToken)
	return nil
}

func getRefreshCsrf() (string, error) {
	uri := "https://www.bilibili.com/correspond/1/"
	correspondPath, err := getCorrespondPath(time.Now().UnixMilli())
	if err != nil {
		return "", err
	}
	uri += correspondPath
	client := http.Client{}
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Cookie", CK)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	dom, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	refreshCsrf := dom.Find("#1-name").Text()
	return refreshCsrf, nil
}

func refreshCookie(refreshCsrf string) (string, error) {
	uri := "https://passport.bilibili.com/x/passport-login/web/cookie/refresh"
	postData := url.Values{}
	csrf := getCsrf()
	postData.Add("refresh_token", refreshToken)
	postData.Add("source", "main_page")
	postData.Add("refresh_csrf", refreshCsrf)
	postData.Add("csrf", csrf)
	fmt.Println("正在刷新cookie...")
	client := http.DefaultClient
	req, _ := http.NewRequest("POST", uri, strings.NewReader(postData.Encode()))
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Cookie", CK)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	data := gjson.ParseBytes(body)
	if data.Get("code").Int() != 0 {
		return "", errors.New(data.Get("message").String())
	}
	var cookieContent []byte
	cookie := make(map[string]string)
	for _, v := range res.Header["Set-Cookie"] {
		kv := strings.Split(v, ";")[0]
		kvArr := strings.Split(kv, "=")
		cookie[kvArr[0]] = kvArr[1]
	}
	cookieContent = []byte(`DedeUserID=` + cookie["DedeUserID"] + `;DedeUserID__ckMd5=` + cookie["DedeUserID__ckMd5"] + `;Expires=` + cookie["Expires"] + `;SESSDATA=` + cookie["SESSDATA"] + `;bili_jct=` + cookie["bili_jct"] + `;`)
	CK = string(cookieContent)
	newRefreshToken := data.Get("data.refresh_token").String()
	return newRefreshToken, nil
}

func commitCookie() error {
	client := http.DefaultClient
	uri := "https://passport.bilibili.com/x/passport-login/web/confirm/refresh"
	postData := url.Values{}
	postData.Add("csrf", getCsrf())
	postData.Add("refresh_token", refreshToken)
	req, _ := http.NewRequest("POST", uri, strings.NewReader(postData.Encode()))
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Cookie", CK)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	data := gjson.ParseBytes(body)
	if data.Get("code").Int() != 0 {
		return errors.New(data.Get("message").String())
	}
	return nil
}

func getCorrespondPath(ts int64) (string, error) {
	pubKeyBlock, _ := pem.Decode([]byte(publicKeyPEM))
	hash := sha256.New()
	random := rand.Reader
	msg := []byte(fmt.Sprintf("refresh_%d", ts))
	var pub *rsa.PublicKey
	pubInterface, parseErr := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if parseErr != nil {
		return "", parseErr
	}
	pub = pubInterface.(*rsa.PublicKey)
	encryptedData, encryptErr := rsa.EncryptOAEP(hash, random, pub, msg, nil)
	if encryptErr != nil {
		return "", encryptErr
	}
	return hex.EncodeToString(encryptedData), nil
}

func isLogin() (bool, gjson.Result, string, string) {
	uri := "https://api.bilibili.com/x/web-interface/nav"
	csrf := getCsrf()
	cookieStr := CK
	client := http.Client{}
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Cookie", cookieStr)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	data := gjson.ParseBytes(body)
	return data.Get("code").Int() == 0, data, cookieStr, csrf
}

func login() (string, string) {
	for {
		loginKey, loginUrl := getLoginKeyAndLoginUrl()
		fmt.Println(loginKey, loginUrl)
		fp, err := os.OpenFile("qrcode.png", os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		var png []byte
		png, err = qrcode.Encode(loginUrl, qrcode.Medium, 256)
		if err != nil {
			panic(err)
		}
		_, err = fp.Write(png)
		if err != nil {
			panic(err)
		}
		fp.Close()
		Scan <- "scan"
		verifyLogin(loginKey)
		logged, data, cookieStr, csrf := isLogin()
		if logged {
			Scan <- "done"
			_ = os.Remove("qrcode.png")
			uname := data.Get("data.uname").String()
			fmt.Println(uname + "已登录")
			return cookieStr, csrf
		}

	}
}

func getCsrf() string {
	reg := re.MustCompile(`bili_jct=([0-9a-zA-Z]+);`)
	if CK == "" {
		CK = config.Conf.Bilibili.Cookie
	}
	csrf := reg.FindStringSubmatch(CK)[1]
	return csrf
}

package miniapi

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"yonghemolimis/src/settings"
)

type miniSession struct {
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type wechatSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type wechatPhone struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
	} `json:"phone_info"`
}

var miniTokens = struct {
	sync.RWMutex
	store map[string]miniSession
}{store: map[string]miniSession{}}

func createMiniToken(userID string) string {
	token := randomHex(32)
	miniTokens.Lock()
	defer miniTokens.Unlock()
	miniTokens.store[token] = miniSession{
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	return token
}

func miniUserIDFromToken(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return ""
	}
	miniTokens.RLock()
	sess, ok := miniTokens.store[token]
	miniTokens.RUnlock()
	if !ok || time.Now().After(sess.ExpiresAt) {
		if ok {
			miniTokens.Lock()
			delete(miniTokens.store, token)
			miniTokens.Unlock()
		}
		return ""
	}
	return sess.UserID
}

func randomHex(size int) string {
	buf := make([]byte, size)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}

func wxCodeToSession(code string) (*wechatSession, error) {
	conf := settings.Conf.MiniWechat
	if conf == nil || conf.AppID == "" || conf.AppSecret == "" {
		return &wechatSession{OpenID: "dev_openid_" + stableCodePart(code)}, nil
	}
	endpoint := "https://api.weixin.qq.com/sns/jscode2session"
	q := url.Values{}
	q.Set("appid", conf.AppID)
	q.Set("secret", conf.AppSecret)
	q.Set("js_code", code)
	q.Set("grant_type", "authorization_code")
	var out wechatSession
	if err := getJSON(endpoint+"?"+q.Encode(), &out); err != nil {
		return nil, err
	}
	if out.ErrCode != 0 {
		return nil, fmt.Errorf("wechat jscode2session failed: %s", out.ErrMsg)
	}
	if out.OpenID == "" {
		return nil, fmt.Errorf("wechat openid empty")
	}
	return &out, nil
}

func wxPhoneNumber(phoneCode string) (string, error) {
	conf := settings.Conf.MiniWechat
	if conf == nil || conf.AppID == "" || conf.AppSecret == "" {
		return "", fmt.Errorf("mini wechat appid/secret not configured")
	}
	token, err := wxAccessToken()
	if err != nil {
		return "", err
	}
	endpoint := "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=" + url.QueryEscape(token)
	reqBody := strings.NewReader(`{"code":"` + phoneCode + `"}`)
	resp, err := http.Post(endpoint, "application/json", reqBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var out wechatPhone
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.ErrCode != 0 {
		return "", fmt.Errorf("wechat phone failed: %s", out.ErrMsg)
	}
	if out.PhoneInfo.PurePhoneNumber != "" {
		return out.PhoneInfo.PurePhoneNumber, nil
	}
	return out.PhoneInfo.PhoneNumber, nil
}

func wxAccessToken() (string, error) {
	conf := settings.Conf.MiniWechat
	endpoint := "https://api.weixin.qq.com/cgi-bin/token"
	q := url.Values{}
	q.Set("grant_type", "client_credential")
	q.Set("appid", conf.AppID)
	q.Set("secret", conf.AppSecret)
	var out struct {
		AccessToken string `json:"access_token"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	if err := getJSON(endpoint+"?"+q.Encode(), &out); err != nil {
		return "", err
	}
	if out.ErrCode != 0 {
		return "", fmt.Errorf("wechat access_token failed: %s", out.ErrMsg)
	}
	if out.AccessToken == "" {
		return "", fmt.Errorf("wechat access_token empty")
	}
	return out.AccessToken, nil
}

func getJSON(endpoint string, out any) error {
	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("http status %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func stableCodePart(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return randomHex(8)
	}
	if len(code) > 24 {
		return code[:24]
	}
	return code
}

package miniapi

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"yonghemolimis/src/dao/db"
	"yonghemolimis/src/settings"

	"github.com/gin-gonic/gin"
)

type smsCodeState struct {
	Code        string
	ExpiresAt   time.Time
	SentAt      time.Time
	HourSends   []time.Time
	DaySends    []time.Time
	Failures    int
	LockedUntil time.Time
}

var smsCodes = struct {
	sync.Mutex
	values map[string]*smsCodeState
}{values: map[string]*smsCodeState{}}
var numericCodePattern = regexp.MustCompile(`^\d{4,8}$`)

func phoneCode(c *gin.Context) {
	var req struct {
		Phone string `json:"phone"`
		Scene string `json:"scene"`
	}
	if c.ShouldBindJSON(&req) != nil {
		businessFail(c, 400, 40000, "请求参数错误")
		return
	}
	req.Phone = strings.TrimSpace(req.Phone)
	if !mainlandPhonePattern.MatchString(req.Phone) {
		businessFail(c, 400, 40001, "手机号格式不正确")
		return
	}
	if req.Scene == "" {
		req.Scene = "LOGIN"
	}
	if req.Scene != "LOGIN" {
		businessFail(c, 400, 40000, "scene 仅支持 LOGIN")
		return
	}
	now := time.Now()
	smsCodes.Lock()
	state := smsCodes.values[req.Phone]
	if state == nil {
		state = &smsCodeState{}
	}
	state.HourSends = since(state.HourSends, now.Add(-time.Hour))
	state.DaySends = since(state.DaySends, now.Add(-24*time.Hour))
	if !state.SentAt.IsZero() && now.Sub(state.SentAt) < time.Minute {
		retry := int(time.Minute.Seconds() - now.Sub(state.SentAt).Seconds())
		smsCodes.Unlock()
		businessFail(c, 429, 42900, "请求过于频繁，请在 "+strconvItoa(retry)+" 秒后重试")
		return
	}
	if len(state.HourSends) >= 5 || len(state.DaySends) >= 10 {
		smsCodes.Unlock()
		businessFail(c, 429, 42900, "请求过于频繁")
		return
	}
	code := configuredSMSCode()
	if code == "" {
		code = randomNumericCode()
	}
	if code == "" {
		smsCodes.Unlock()
		businessFail(c, 500, 50000, "短信服务暂不可用")
		return
	}
	// 先占用一分钟发送窗口，避免并发请求绕过频率限制；外部网关调用不持锁。
	state.SentAt = now
	smsCodes.values[req.Phone] = state
	smsCodes.Unlock()
	if err := sendSMS(req.Phone, code); err != nil {
		smsCodes.Lock()
		if state.SentAt.Equal(now) {
			state.SentAt = time.Time{}
		}
		smsCodes.Unlock()
		businessFail(c, 500, 50000, "短信服务暂不可用")
		return
	}
	smsCodes.Lock()
	state.Code = code
	state.ExpiresAt = now.Add(5 * time.Minute)
	state.HourSends = append(state.HourSends, now)
	state.DaySends = append(state.DaySends, now)
	state.Failures = 0
	smsCodes.values[req.Phone] = state
	smsCodes.Unlock()
	businessOK(c, 200, gin.H{"retryAfterSeconds": 60, "expiresInSeconds": 300})
}

func phoneLogin(c *gin.Context) {
	var req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	if c.ShouldBindJSON(&req) != nil {
		businessFail(c, 400, 40000, "请求参数错误")
		return
	}
	req.Phone = strings.TrimSpace(req.Phone)
	req.Code = strings.TrimSpace(req.Code)
	if !mainlandPhonePattern.MatchString(req.Phone) {
		businessFail(c, 400, 40001, "手机号格式不正确")
		return
	}
	if !numericCodePattern.MatchString(req.Code) {
		businessFail(c, 400, 40002, "验证码错误或已失效")
		return
	}
	now := time.Now()
	smsCodes.Lock()
	state := smsCodes.values[req.Phone]
	if state == nil || now.After(state.ExpiresAt) || state.Code == "" {
		smsCodes.Unlock()
		businessFail(c, 400, 40002, "验证码错误或已失效")
		return
	}
	if now.Before(state.LockedUntil) {
		smsCodes.Unlock()
		businessFail(c, 429, 42900, "验证码错误次数过多，请稍后重试")
		return
	}
	if state.Code != req.Code {
		state.Failures++
		if state.Failures >= 5 {
			state.LockedUntil = now.Add(10 * time.Minute)
		}
		smsCodes.Unlock()
		if state.Failures >= 5 {
			businessFail(c, 429, 42900, "验证码错误次数过多，请稍后重试")
		} else {
			businessFail(c, 400, 40002, "验证码错误或已失效")
		}
		return
	}
	state.Code = ""
	state.ExpiresAt = time.Time{}
	state.Failures = 0
	smsCodes.Unlock()
	row, err := db.GetMiniUserByPhone(req.Phone)
	if err != nil {
		row = &db.CustomerDO{ID: newMiniUserID(), Phone: req.Phone, Nickname: maskPhone(req.Phone), Status: db.CustomerStatusActive, CreatedAt: now}
	}
	if row.Status != db.CustomerStatusActive {
		businessFail(c, 403, 40300, "账号已被禁用")
		return
	}
	row.Phone = req.Phone
	row.LastLoginAt = now.Format("2006-01-02 15:04:05")
	row.UpdatedAt = now
	if err := db.UpsertMiniUserProfile(row); err != nil {
		businessFail(c, 500, 50000, "登录失败")
		return
	}
	user := userFromDO(*row)
	user.Phone = maskPhone(row.Phone)
	businessOK(c, 200, gin.H{"token": createMiniToken(row.ID), "expiresIn": 30 * 24 * 60 * 60, "user": user})
}

func configuredSMSCode() string {
	if settings.Conf.Mode != "release" && settings.Conf.MiniSMS != nil {
		return strings.TrimSpace(settings.Conf.MiniSMS.TestCode)
	}
	return ""
}
func sendSMS(phone, code string) error {
	conf := settings.Conf.MiniSMS
	if conf == nil || strings.TrimSpace(conf.Endpoint) == "" {
		if settings.Conf.Mode != "release" && configuredSMSCode() != "" {
			return nil
		}
		return errors.New("sms gateway not configured")
	}
	body, _ := json.Marshal(gin.H{"phone": phone, "code": code, "scene": "LOGIN"})
	req, err := http.NewRequest(http.MethodPost, conf.Endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if conf.Token != "" {
		req.Header.Set("Authorization", "Bearer "+conf.Token)
	}
	resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("sms gateway rejected request")
	}
	return nil
}
func randomNumericCode() string {
	var raw [8]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return ""
	}
	return strconvItoa(int(binary.BigEndian.Uint64(raw[:])%900000) + 100000)
}
func since(values []time.Time, cutoff time.Time) []time.Time {
	out := values[:0]
	for _, value := range values {
		if value.After(cutoff) {
			out = append(out, value)
		}
	}
	return out
}
func strconvItoa(value int) string {
	const digits = "0123456789"
	if value <= 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for value > 0 {
		i--
		buf[i] = digits[value%10]
		value /= 10
	}
	return string(buf[i:])
}

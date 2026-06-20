package session

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// Session 会话数据
type Session struct {
	AdminID      uint   // 本地管理员 ID
	SSOUserID    uint   // SSO 用户 ID
	Username     string // 显示名
	Email        string
	Avatar       string
	IsSuperAdmin bool
	RoleID       *int64
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

var (
	store = make(map[string]*Session)
	mu    sync.RWMutex
	ttl   = 24 * time.Hour
)

// Create 创建新会话，返回 sessionID
func Create(adminID, ssoUserID uint, username, email, avatar string, isSuperAdmin bool, roleID *int64) string {
	sid := generateID()
	mu.Lock()
	defer mu.Unlock()
	store[sid] = &Session{
		AdminID:      adminID,
		SSOUserID:    ssoUserID,
		Username:     username,
		Email:        email,
		Avatar:       avatar,
		IsSuperAdmin: isSuperAdmin,
		RoleID:       roleID,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(ttl),
	}
	return sid
}

// Get 获取会话，同时检查过期
func Get(sid string) *Session {
	mu.RLock()
	defer mu.RUnlock()
	s, ok := store[sid]
	if !ok {
		return nil
	}
	if time.Now().After(s.ExpiresAt) {
		go Destroy(sid)
		return nil
	}
	return s
}

// Destroy 销毁会话
func Destroy(sid string) {
	mu.Lock()
	defer mu.Unlock()
	delete(store, sid)
}

// generateID 生成随机 session ID
func generateID() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

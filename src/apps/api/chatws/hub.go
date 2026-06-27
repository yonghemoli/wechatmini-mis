package chatws

import (
	"errors"
	"net/http"
	"strings"
	"sync"

	"yonghemolimis/src/dao/db"
	"yonghemolimis/src/logger"
	"yonghemolimis/src/middlewares"
	"yonghemolimis/src/pkgs/session"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

type ClientRole string

const (
	RoleAdmin ClientRole = "admin"
	RoleMini  ClientRole = "mini"
)

type Event struct {
	Type      string            `json:"type"`
	SessionID string            `json:"sessionId,omitempty"`
	Message   *db.ChatMessageDO `json:"message,omitempty"`
	Session   *db.ChatSessionDO `json:"session,omitempty"`
	Error     string            `json:"error,omitempty"`
}

type inboundMessage struct {
	Type      string `json:"type"`
	SessionID string `json:"sessionId"`
	Content   string `json:"content"`
	MsgType   string `json:"msgType"`
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
}

type client struct {
	role      ClientRole
	sessionID string
	send      chan Event
}

type hub struct {
	mu      sync.RWMutex
	clients map[*client]struct{}
}

var defaultHub = &hub{clients: make(map[*client]struct{})}

func RegisterRoutes(r gin.IRouter, role ClientRole) {
	r.GET("/ws/chat", func(c *gin.Context) {
		if role == RoleAdmin && !isAdminAuthed(c) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录或会话已过期"})
			return
		}
		handler := websocket.Server{
			Handler: func(ws *websocket.Conn) {
				serveConn(ws, role, c)
			},
			Handshake: func(config *websocket.Config, req *http.Request) error {
				config.Origin, _ = websocket.Origin(config, req)
				return nil
			},
		}
		handler.ServeHTTP(c.Writer, c.Request)
	})
}

func BroadcastMessage(message *db.ChatMessageDO) {
	if message == nil {
		return
	}
	broadcast(Event{
		Type:      "message",
		SessionID: message.SessionID,
		Message:   message,
	})
}

func BroadcastSession(session *db.ChatSessionDO) {
	if session == nil {
		return
	}
	broadcast(Event{
		Type:      "session",
		SessionID: session.ID,
		Session:   session,
	})
}

func serveConn(ws *websocket.Conn, role ClientRole, c *gin.Context) {
	sessionID := strings.TrimSpace(c.Query("sessionId"))
	if role == RoleMini {
		sessionID = ensureMiniSessionID(c, sessionID)
	}
	conn := &client{
		role:      role,
		sessionID: sessionID,
		send:      make(chan Event, 16),
	}
	defaultHub.add(conn)
	defer defaultHub.remove(conn)

	if sessionID != "" {
		if row, err := db.GetChatSession(sessionID); err == nil {
			conn.send <- Event{Type: "session", SessionID: row.ID, Session: row}
		}
	}

	go writeLoop(ws, conn)
	readLoop(ws, conn, c)
}

func readLoop(ws *websocket.Conn, conn *client, c *gin.Context) {
	for {
		var msg inboundMessage
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			return
		}
		if msg.Type == "" {
			msg.Type = "message"
		}
		switch msg.Type {
		case "ping":
			conn.send <- Event{Type: "pong"}
		case "message":
			if err := handleMessage(conn, c, msg); err != nil {
				conn.send <- Event{Type: "error", Error: err.Error()}
			}
		}
	}
}

func writeLoop(ws *websocket.Conn, conn *client) {
	for msg := range conn.send {
		_ = websocket.JSON.Send(ws, msg)
	}
}

func handleMessage(conn *client, c *gin.Context, msg inboundMessage) error {
	content := strings.TrimSpace(msg.Content)
	if content == "" {
		return errors.New("消息内容不能为空")
	}
	sessionID := conn.sessionID
	if conn.role == RoleAdmin && msg.SessionID != "" {
		sessionID = msg.SessionID
	}
	if sessionID == "" {
		return errors.New("会话ID不能为空")
	}
	if conn.role == RoleMini {
		userID, userName := miniUser(c, msg)
		if err := db.EnsureChatSession(sessionID, userID, userName); err != nil {
			return err
		}
	}
	row, err := db.CreateChatMessage(sessionID, string(sender(conn.role)), msg.MsgType, content)
	if err != nil {
		return err
	}
	BroadcastMessage(row)
	if sessionRow, err := db.GetChatSession(sessionID); err == nil {
		BroadcastSession(sessionRow)
	}
	return nil
}

func broadcast(event Event) {
	defaultHub.mu.RLock()
	defer defaultHub.mu.RUnlock()
	for conn := range defaultHub.clients {
		if conn.sessionID == "" || conn.sessionID == event.SessionID {
			select {
			case conn.send <- event:
			default:
				logger.Warnf("客服WS客户端消息队列已满: role=%s session=%s", conn.role, conn.sessionID)
			}
		}
	}
}

func (h *hub) add(conn *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = struct{}{}
}

func (h *hub) remove(conn *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, conn)
	close(conn.send)
}

func isAdminAuthed(c *gin.Context) bool {
	sid, err := c.Cookie(middlewares.SessionCookieName)
	if err != nil || sid == "" {
		return false
	}
	return session.Get(sid) != nil
}

func ensureMiniSessionID(c *gin.Context, sessionID string) string {
	if sessionID != "" {
		return sessionID
	}
	userID := strings.TrimSpace(c.Query("userId"))
	if userID == "" {
		userID = "u_1"
	}
	return "chat_" + userID
}

func miniUser(c *gin.Context, msg inboundMessage) (string, string) {
	userID := firstNonEmpty(msg.UserID, c.Query("userId"), "u_1")
	userName := firstNonEmpty(msg.UserName, c.Query("userName"), "微信用户")
	return userID, userName
}

func sender(role ClientRole) ClientRole {
	if role == RoleAdmin {
		return RoleAdmin
	}
	return ClientRole("user")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func CreateMiniMessage(sessionID, userID, userName, msgType, content string) (*db.ChatMessageDO, error) {
	if sessionID == "" {
		sessionID = "chat_" + firstNonEmpty(userID, "u_1")
	}
	if err := db.EnsureChatSession(sessionID, firstNonEmpty(userID, "u_1"), firstNonEmpty(userName, "微信用户")); err != nil {
		return nil, err
	}
	row, err := db.CreateChatMessage(sessionID, "user", msgType, content)
	if err != nil {
		return nil, err
	}
	BroadcastMessage(row)
	if sessionRow, err := db.GetChatSession(sessionID); err == nil {
		BroadcastSession(sessionRow)
	}
	return row, nil
}

func TouchSession(sessionID, userID, userName string) (*db.ChatSessionDO, error) {
	if sessionID == "" {
		sessionID = "chat_" + firstNonEmpty(userID, "u_1")
	}
	if err := db.EnsureChatSession(sessionID, firstNonEmpty(userID, "u_1"), firstNonEmpty(userName, "微信用户")); err != nil {
		return nil, err
	}
	return db.GetChatSession(sessionID)
}

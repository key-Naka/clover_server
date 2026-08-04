package chat_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/ctype/chat_msg"
	"clover_server/models/enum"
	"clover_server/service/chat_service"
	"clover_server/utils/jwts"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

type ChatResponse struct {
	ChatListResponse
}

var (
	onlineConnMap = map[uint]map[string]*websocket.Conn{}
	onlineConnMux sync.RWMutex
)

func (ChatApi) ChatView(c *gin.Context) {
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil || claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}
	var currentUser models.UserModel
	if err = global.DB.Take(&currentUser, claims.UserID).Error; err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}
	websocket.Handler(func(conn *websocket.Conn) {
		remoteAddr := "unknown"
		if conn.Request() != nil && conn.Request().RemoteAddr != "" {
			remoteAddr = conn.Request().RemoteAddr
		}
		putOnlineConn(claims.UserID, remoteAddr, conn)
		defer func() {
			removeOnlineConn(claims.UserID, remoteAddr)
			_ = conn.Close()
		}()
		for {
			var req ChatRequest
			if err = websocket.JSON.Receive(conn, &req); err != nil {
				return
			}
			if err = validateChatRequest(claims.UserID, req); err != nil {
				writeWSFail(conn, err.Error())
				continue
			}
			chat, saveErr := saveWSChat(req, claims.UserID)
			if saveErr != nil {
				writeWSFail(conn, saveErr.Error())
				continue
			}
			response := ChatResponse{ChatListResponse: ChatListResponse{
				ChatModel:        *chat,
				SendUserNickname: currentUser.Nickname,
				SendUserAvatar:   currentUser.Avatar,
				RevUserNickname:  "",
				RevUserAvatar:    "",
				IsMe:             true,
				IsRead:           false,
			}}
			var revUser models.UserModel
			if err = global.DB.Take(&revUser, req.RevUserID).Error; err == nil {
				response.RevUserNickname = revUser.Nickname
				response.RevUserAvatar = revUser.Avatar
				revResp := response
				revResp.IsMe = false
				sendToUser(req.RevUserID, revResp)
			}
			_ = websocket.JSON.Send(conn, map[string]any{"code": 0, "data": response, "msg": "成功"})
		}
	}).ServeHTTP(c.Writer, c.Request)
}

func validateChatRequest(userID uint, req ChatRequest) error {
	if req.RevUserID == 0 || req.RevUserID == userID {
		return errChat("接收用户错误")
	}
	var revUser models.UserModel
	if err := global.DB.Take(&revUser, req.RevUserID).Error; err != nil {
		return errChat("接收用户不存在")
	}
	switch req.MsgType {
	case enum.ChatTextMsgType:
		if req.Msg.TextMsg == nil || strings.TrimSpace(req.Msg.TextMsg.Content) == "" {
			return errChat("文本消息不能为空")
		}
	case enum.ChatImageMsgType:
		if req.Msg.ImageMsg == nil || strings.TrimSpace(req.Msg.ImageMsg.Src) == "" {
			return errChat("图片消息不能为空")
		}
	case enum.ChatMarkdownMsgType:
		if req.Msg.MarkdownMsg == nil || strings.TrimSpace(req.Msg.MarkdownMsg.Content) == "" {
			return errChat("markdown消息不能为空")
		}
	case enum.ChatReadMsgType:
		if req.Msg.MsgReadMsg == nil || req.Msg.MsgReadMsg.ReadChatID == 0 {
			return errChat("已读消息错误")
		}
	default:
		return errChat("不支持的消息类型")
	}
	return nil
}

func saveWSChat(req ChatRequest, userID uint) (*models.ChatModel, error) {
	switch req.MsgType {
	case enum.ChatTextMsgType:
		if err := chat_service.ToTextChat(userID, req.RevUserID, strings.TrimSpace(req.Msg.TextMsg.Content)); err != nil {
			return nil, err
		}
	case enum.ChatImageMsgType:
		if err := chat_service.ToImageChat(userID, req.RevUserID, strings.TrimSpace(req.Msg.ImageMsg.Src)); err != nil {
			return nil, err
		}
	case enum.ChatMarkdownMsgType:
		if err := chat_service.ToMarkdownChat(userID, req.RevUserID, strings.TrimSpace(req.Msg.MarkdownMsg.Content)); err != nil {
			return nil, err
		}
	case enum.ChatReadMsgType:
		if err := chat_service.ToChat(userID, req.RevUserID, enum.ChatReadMsgType, chat_msg.ChatMsg{MsgReadMsg: &chat_msg.MsgReadMsg{ReadChatID: req.Msg.MsgReadMsg.ReadChatID}}); err != nil {
			return nil, err
		}
	}
	var chat models.ChatModel
	if err := global.DB.Last(&chat).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}

func putOnlineConn(userID uint, remoteAddr string, conn *websocket.Conn) {
	onlineConnMux.Lock()
	defer onlineConnMux.Unlock()
	if _, ok := onlineConnMap[userID]; !ok {
		onlineConnMap[userID] = map[string]*websocket.Conn{}
	}
	onlineConnMap[userID][remoteAddr] = conn
}

func removeOnlineConn(userID uint, remoteAddr string) {
	onlineConnMux.Lock()
	defer onlineConnMux.Unlock()
	if _, ok := onlineConnMap[userID]; !ok {
		return
	}
	delete(onlineConnMap[userID], remoteAddr)
	if len(onlineConnMap[userID]) == 0 {
		delete(onlineConnMap, userID)
	}
}

func sendToUser(userID uint, data any) {
	onlineConnMux.RLock()
	defer onlineConnMux.RUnlock()
	for _, conn := range onlineConnMap[userID] {
		if err := websocket.JSON.Send(conn, map[string]any{"code": 0, "data": data, "msg": "成功"}); err != nil {
			slog.Warn("ws推送失败", "userID", userID, "error", err.Error())
		}
	}
}

func writeWSFail(conn *websocket.Conn, msg string) {
	_ = websocket.JSON.Send(conn, map[string]any{"code": 1001, "data": map[string]any{}, "msg": msg, "time": time.Now().Unix()})
}

type chatError struct{ msg string }

func (e chatError) Error() string { return e.msg }

func errChat(msg string) error { return chatError{msg: msg} }

var _ http.Handler
var _ = json.Marshal

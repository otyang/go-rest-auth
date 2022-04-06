package controller

import (
	"auth-project/src/usecase/interactor"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/spf13/viper"
	"sync"
	"time"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Message struct {
	Message string `json:"message"`
}

type Channel struct {
	conn     *websocket.Conn
	stop     chan struct{}
	complete chan string
	token    string
}

func NewChannel(token string, conn *websocket.Conn) *Channel {
	c := &Channel{
		conn:     conn,
		stop:     make(chan struct{}, 1),
		complete: make(chan string, 1),
		token:    token,
	}

	return c
}

func CloseChannel(ch *Channel) {
	close(ch.stop)
	close(ch.complete)
}

type qrCodeAuthController struct {
	qrCodeAuthInteractor interactor.QrCodeAuthInteractor
	connArr              map[string]*Channel
	mx                   *sync.Mutex
}

type QrCodeAuthController interface {
	CreateAuthTokenByAuthQrCode(ctx *fiber.Ctx) error
	QrCodeAuthWebsocket(c *websocket.Conn)
}

func NewQrCodeAuthController(qi interactor.QrCodeAuthInteractor) QrCodeAuthController {

	return &qrCodeAuthController{
		qi,
		make(map[string]*Channel),
		&sync.Mutex{},
	}
}

func (qc *qrCodeAuthController) CreateAuthTokenByAuthQrCode(ctx *fiber.Ctx) error {
	userId, ok := ctx.Context().Value("token_user_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "context value type invalid")
	}

	// Send user id to ws goroutine
	ok = qc.Load(ctx.Params("qrCodeToken"), userId)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "your token does not registered")
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}

func (qc *qrCodeAuthController) QrCodeAuthWebsocket(c *websocket.Conn) {

	qrCode, qrCodeToken, err := qc.qrCodeAuthInteractor.GenerateQrCode(context.Background())
	if err != nil {
		_ = c.Conn.WriteJSON(Message{err.Error()})
		_ = c.Conn.WriteMessage(websocket.CloseInternalServerErr, nil)
		return
	}

	err = c.Conn.WriteJSON(map[string]interface{}{
		"qr_code":       qrCode,
		"qr_code_token": qrCodeToken,
	})
	if err != nil {
		_ = c.Conn.WriteJSON(Message{"err sending qr token"})
		_ = c.Conn.WriteMessage(websocket.CloseInternalServerErr, nil)
		return
	}

	// Waiting for complete or close
	userId, ok := <-qc.Store(qrCodeToken, c)
	if !ok {
		_ = c.Conn.WriteJSON(Message{"timeout"})
		_ = c.Conn.WriteMessage(websocket.CloseMessage, nil)
		return
	}

	details, err := qc.qrCodeAuthInteractor.GenerateTokenPairByUserID(c, userId)
	if err != nil {
		_ = c.Conn.WriteJSON(Message{err.Error()})
		_ = c.Conn.WriteMessage(websocket.CloseMessage, nil)
		return
	}

	_ = c.Conn.WriteJSON(map[string]string{
		"access_token":  details.AccessToken,
		"refresh_token": details.RefreshToken,
	})

	_ = c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(writeWait))
	_ = c.Conn.Close()
}

func (qc *qrCodeAuthController) Store(token string, c *websocket.Conn) <-chan string {
	qc.mx.Lock()
	defer qc.mx.Unlock()

	ch := NewChannel(token, c)
	qc.connArr[token] = ch
	go qc.Timeout(ch)
	go qc.Reader(ch)

	return ch.complete
}

func (qc *qrCodeAuthController) Load(key, userId string) bool {
	qc.mx.Lock()
	defer qc.mx.Unlock()

	ch, ok := qc.connArr[key]
	if !ok {
		return false
	}

	ch.complete <- userId
	go qc.Unregister(ch.token)

	return ok
}

func (qc *qrCodeAuthController) Unregister(token string) bool {
	qc.mx.Lock()
	defer qc.mx.Unlock()

	ch, ok := qc.connArr[token]
	if !ok {
		return false
	}

	CloseChannel(ch)
	delete(qc.connArr, ch.token)

	return true
}

func (qc *qrCodeAuthController) Timeout(ch *Channel) {
	var err error
	ticker := time.NewTicker(pingPeriod)
	endTimer := time.NewTimer(viper.GetDuration("ws.timeout_duration"))

	defer func() {
		endTimer.Stop()
		ticker.Stop()
		qc.Unregister(ch.token)
	}()

	for {
		select {
		case <-ch.stop: // Write to conn
			return
		case <-endTimer.C: // Close conn after deadline
			return
		case <-ticker.C: // Ping ticker
			if err = ch.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}

func (qc *qrCodeAuthController) Reader(ch *Channel) {
	defer qc.Unregister(ch.token)

	for {
		_, _, err := ch.conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

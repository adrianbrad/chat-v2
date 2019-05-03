package integration_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adrianbrad/chat-v2/internal/client"
	"github.com/adrianbrad/chat-v2/internal/message"
	"github.com/adrianbrad/chat-v2/internal/repository/roomrepository"

	"github.com/adrianbrad/chat-v2/pkg/websocketshandler"
	"github.com/gorilla/websocket"

	"github.com/stretchr/testify/assert"

	"github.com/adrianbrad/chat-v2/internal/chatservice"
	"github.com/adrianbrad/chat-v2/internal/repository/userrepository"
	"github.com/adrianbrad/chat-v2/pkg/hashauth"
	"github.com/adrianbrad/chat-v2/pkg/otpauth/httpotpauth"
	"github.com/adrianbrad/chat-v2/test"
	testutils "github.com/adrianbrad/chat-v2/test/utils"
)

func TestAuthenticateTokenSuccess(t *testing.T) {
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer db.Close()

	ur := userrepository.NewUserRepositoryDB(db)
	rRepo := roomrepository.NewRoomRepositoryDB(db)
	mp := message.NewMessageProcessor(nil)

	cs := chatservice.NewChatService(ur, rRepo, client.NewFactory(mp))

	upgrader := &websocket.Upgrader{}
	websocketsHandler := websocketshandler.NewWebsocketsHandler(
		cs,
		upgrader,
		func(r *http.Request) (data map[string]interface{}, err error) {
			data = make(map[string]interface{})
			data["userID"] = r.Header.Get("X-OTPAuth-ID")
			data["roomID"] = r.FormValue("roomID")
			return
		},
	)

	authFunc := func(ID string) bool {
		_, err := ur.GetOne(ID)
		if err != nil {
			return false
		}
		return true
	}

	tokenAuth := httpotpauth.NewHTTPOTPAuthenticator(10*time.Second, authFunc, websocketsHandler)

	secret := "chat"
	hAuth := hashauth.NewHTTPHashAuthenticator(
		secret,

		tokenAuth,

		func(r *http.Request) (hash, data string, err error, skipAuth bool) {
			bodyBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return
			}
			data = string(bodyBytes)
			hash = r.Header.Get("Authenticate")

			r.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
			return hash, data, nil, false
		})

	body := []byte("user_a")
	r := testutils.NewTestRequest(t, http.MethodPost, "", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h := hmac.New(sha256.New, []byte(secret))
	hash := hashauth.GenerateHash(h, body)

	r.Header.Set("Authenticate", hash)

	hAuth.ServeHTTP(rr, r)

	authToken := rr.Header().Get("Authorization")

	http.Handle("/ws-chat", tokenAuth)
	go func() {
		t.Fatal(http.ListenAndServe(":8080", nil))
	}()

	time.Sleep(time.Second / 2)

	_, resp, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws-chat?key="+authToken+"&roomID=room_a", nil)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
}

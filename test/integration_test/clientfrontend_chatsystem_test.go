package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/adrianbrad/chat-v2/internal/message"
	"github.com/adrianbrad/chat-v2/internal/user"

	"github.com/adrianbrad/chat-v2/internal/server"
	"github.com/adrianbrad/chat-v2/pkg/hashauth"
	testutils "github.com/adrianbrad/chat-v2/test/utils"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func requestToken(t *testing.T, userID string) (token string) {
	body := []byte(userID)
	r := testutils.NewTestRequest(t, http.MethodPost, baseAddress+"/auth", bytes.NewReader(body))

	hash := hashauth.GenerateHash(hasher, body)

	r.Header.Set("Authenticate", hash)

	resp, err := httpClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
	return resp.Header.Get("Authorization")
}

func createUser(t *testing.T, user *user.User) {
	userJsonBytes, _ := json.Marshal(user)
	r := testutils.NewTestRequest(t, http.MethodPost, baseAddress+"/user", bytes.NewReader(userJsonBytes))

	hash := hashauth.GenerateHash(hasher, userJsonBytes)
	r.Header.Set("Authenticate", hash)

	httpClient.Do(r)
}

func TestConnectionSuccess(t *testing.T) {
	db := initDB(t)
	defer db.Close()

	initDependencies(db)

	stopServer := startServer(
		server.PathHandler{
			Path:    "/auth",
			Handler: httpOTPAuthenticator.Auth(nil)},

		server.PathHandler{
			Path:    "/chat",
			Handler: httpOTPAuthenticator.Auth(websocketHandler)},
	)
	defer stopServer()

	authToken := requestToken(t, "user_a")

	_, resp, err := websocket.DefaultDialer.Dial("ws://localhost:8080/chat?key="+authToken+"&roomID=room_a", nil)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
}

func TestChatInteractions(t *testing.T) {
	db := initDB(t)
	defer db.Close()

	initDependencies(db)

	stopServer := startServer(
		server.PathHandler{
			Path:    "/auth",
			Handler: httpOTPAuthenticator.Auth(nil)},

		server.PathHandler{
			Path:    "/chat",
			Handler: httpOTPAuthenticator.Auth(websocketHandler)},

		server.PathHandler{
			Path:    "/user",
			Handler: hashAuthenticator.Auth(userService),
		},
	)
	defer stopServer()

	user1 := &user.User{
		ID:       "test_user1",
		Nickname: "TestUser169",
		Permissions: map[string]struct{}{
			user.SendMessage.String(): struct{}{},
		},
	}
	user2 := &user.User{
		ID:       "test_user2",
		Nickname: "TestUser269",
		Permissions: map[string]struct{}{
			user.SendMessage.String(): struct{}{},
		},
	}
	user3 := &user.User{
		ID:       "test_user3",
		Nickname: "TestUser369",
	}

	createUser(t, user1)
	createUser(t, user2)
	createUser(t, user3)

	authToken1 := requestToken(t, user1.ID)
	authToken2 := requestToken(t, user2.ID)
	authToken3 := requestToken(t, user3.ID)

	conn1, _, _ := websocket.DefaultDialer.Dial("ws://localhost:8080/chat?key="+authToken1+"&roomID=room_a", nil)
	conn2, _, _ := websocket.DefaultDialer.Dial("ws://localhost:8080/chat?key="+authToken2+"&roomID=room_a", nil)
	conn3, _, _ := websocket.DefaultDialer.Dial("ws://localhost:8080/chat?key="+authToken3+"&roomID=room_a", nil)

	user1Message := message.BareMessage{
		Action: message.Text.String(),
		Body:   message.MessageBody{"hello"},
		RoomID: "room_a",
	}

	conn1.WriteJSON(user1Message)

	go func() {
		var receivedMes map[string]interface{}
		err := conn1.ReadJSON(&receivedMes)
		fmt.Println("conn1")
		fmt.Println(receivedMes, err)
	}()

	go func() {
		var receivedMes map[string]interface{}
		conn2.ReadJSON(&receivedMes)
		fmt.Println("conn2")

		fmt.Println(receivedMes)
	}()

	go func() {
		var receivedMes map[string]interface{}
		conn3.ReadJSON(&receivedMes)
		fmt.Println("conn3")
		fmt.Println(receivedMes)
	}()

	select {}
}

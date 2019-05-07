package integration

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/adrianbrad/chat-v2/internal/chatservice"
	"github.com/adrianbrad/chat-v2/internal/client"
	"github.com/adrianbrad/chat-v2/internal/message"
	"github.com/adrianbrad/chat-v2/internal/repository/messagerepository"
	"github.com/adrianbrad/chat-v2/internal/repository/roomrepository"
	"github.com/adrianbrad/chat-v2/internal/repository/userrepository"
	"github.com/adrianbrad/chat-v2/internal/server"
	"github.com/adrianbrad/chat-v2/internal/user"
	"github.com/adrianbrad/chat-v2/pkg/hashauth"
	"github.com/adrianbrad/chat-v2/pkg/otpauth/httpotpauth"
	"github.com/adrianbrad/chat-v2/pkg/websocketshandler"
	"github.com/adrianbrad/chat-v2/test"
	"github.com/gorilla/websocket"
)

var (
	userRepository       *userrepository.UserRepositoryDB
	roomRepository       *roomrepository.RoomRepositoryDB
	messageProcessor     *message.MessageProcessor
	messageRepository    *messagerepository.MessageRepositoryDB
	chatService          *chatservice.ChatService
	clientFactory        client.Factory
	websocketUpgrader    *websocket.Upgrader
	websocketHandler     *websocketshandler.WebsocketsHandler
	httpOTPAuthenticator *httpotpauth.HTTPOTPAuthenticator
	userService          *user.UserService

	secret            = "chat"
	hashAuthenticator *hashauth.HTTPHashAuthenticator

	getDataFromWebsocketRequestFunc func(r *http.Request) (data map[string]interface{}, err error)
	otpAuthFunc                     func(dataToAuthenticate string) (valid bool)
	retrieveDataForHashAuthFunc     func(r *http.Request) (hash, data string, err error, skipAuth bool)
	bareMessageFactoryFunc          = message.NewBareMessage

	serverPort  = ":8080"
	baseAddress = "http://localhost" + serverPort

	httpClient = &http.Client{}
	hasher     = hmac.New(sha256.New, []byte(secret))
)

func initDependencies(db *sql.DB) {
	userRepository = userrepository.NewUserRepositoryDB(db)
	userService = user.NewUserService(userRepository)
	roomRepository = roomrepository.NewRoomRepositoryDB(db)

	messageRepository = messagerepository.NewMessageRepositoryDB(db)
	messageProcessor = message.NewMessageProcessor(messageRepository)
	clientFactory = client.NewFactory(messageProcessor, bareMessageFactoryFunc)

	chatService = chatservice.NewChatService(
		userRepository,
		roomRepository,
		clientFactory,
	)

	getDataFromWebsocketRequestFunc = func(r *http.Request) (data map[string]interface{}, err error) {
		data = make(map[string]interface{})
		data["userID"] = r.Header.Get("X-OTPAuth-ID")
		data["roomID"] = r.FormValue("roomID")
		return
	}

	websocketUpgrader = &websocket.Upgrader{}

	websocketHandler = websocketshandler.NewWebsocketsHandler(
		chatService,
		websocketUpgrader,
		getDataFromWebsocketRequestFunc,
	)

	otpAuthFunc = func(ID string) bool {
		_, err := userRepository.GetOne(ID)
		if err != nil {
			return false
		}
		return true
	}

	httpOTPAuthenticator = httpotpauth.NewHTTPOTPAuthenticator(
		10*time.Second,
		otpAuthFunc,
	)

	retrieveDataForHashAuthFunc = func(r *http.Request) (hash, data string, err error, skipAuth bool) {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		data = string(bodyBytes)
		hash = r.Header.Get("Authenticate")

		r.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		return hash, data, nil, false
	}

	hashAuthenticator = hashauth.NewHTTPHashAuthenticator(
		secret,
		retrieveDataForHashAuthFunc,
	)
}

func initDB(t *testing.T) (db *sql.DB) {
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatal(err.Error())
	}
	return db
}

func startServer(pathHandlers ...server.PathHandler) (stopServerFunc func()) {
	mux := server.NewMux(pathHandlers...)

	stopServer := make(chan os.Signal)
	waitServerStop := make(chan struct{})
	go func() {
		server.Run(serverPort, mux, stopServer, 0)
		waitServerStop <- struct{}{}
	}()

	time.Sleep(time.Second / 2)
	return func() {
		stopServer <- nil
		<-waitServerStop
		fmt.Println("Server stopped")
	}
}

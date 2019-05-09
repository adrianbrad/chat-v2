package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/cors"

	"github.com/adrianbrad/chat-v2/configs"
	d "github.com/adrianbrad/chat-v2/db"
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
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewChatDatabaseCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "chatdatabase",
		Short: "Chat using database persistency",
		Long:  "",
		Run:   run,
	}
	cmd.Flags().StringP("dbconfig", "d", "", "Set database config file")
	cmd.Flags().StringP("appconfig", "a", "", "Set application config file")
	cmd.Flags().StringP("basedir", "b", "", "Set application base directory file")
	return cmd
}

func run(cmd *cobra.Command, args []string) {
	dbConfig, appConfig := loadConfigs(cmd)
	db, err := d.ConnectDB(dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Successfully connected to db")

	err = d.SetMigrationVersion(filepath.Join(appConfig.Basedir, "db", "migrations"), db, dbConfig.MigrationVersion)
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

	log.Infof("Using migration version %d", dbConfig.MigrationVersion)

	userRepository := userrepository.NewUserRepositoryDB(db)
	userService := user.NewUserService(userRepository)
	roomRepository := roomrepository.NewRoomRepositoryDB(db)

	messageRepository := messagerepository.NewMessageRepositoryDB(db)
	messageProcessor := message.NewMessageProcessor(messageRepository)

	bareMessageFactoryFunc := message.NewBareMessage
	clientFactory := client.NewFactory(messageProcessor, bareMessageFactoryFunc)

	chatService := chatservice.NewChatService(
		userRepository,
		roomRepository,
		clientFactory,
	)

	getDataFromWebsocketRequestFunc := func(r *http.Request) (data map[string]interface{}, err error) {
		data = make(map[string]interface{})
		data["userID"] = r.Header.Get("X-OTPAuth-ID")
		data["roomID"] = r.FormValue("roomID")
		return
	}

	websocketUpgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	websocketHandler := websocketshandler.NewWebsocketsHandler(
		chatService,
		websocketUpgrader,
		getDataFromWebsocketRequestFunc,
	)

	otpAuthFunc := func(ID string) bool {
		_, err := userRepository.GetOne(ID)
		if err != nil {
			return false
		}
		return true
	}

	httpOTPAuthenticator := httpotpauth.NewHTTPOTPAuthenticator(
		10*time.Second,
		otpAuthFunc,
	)

	retrieveDataForHashAuthFunc := func(r *http.Request) (hash, data string, err error, skipAuth bool) {
		epochTimeString := r.Header.Get("Date")
		if epochTimeString == "" {
			err = fmt.Errorf("No date passed as header")
			return
		}

		i, err := strconv.ParseInt(epochTimeString, 10, 64)
		if err != nil {
			return
		}
		tm := time.Unix(i, 0)

		fmt.Println(time.Now().UTC())
		if time.Since(tm).Seconds() > 5 {
			err = fmt.Errorf("Request is too old")
			return
		}
		hash = r.Header.Get("Authorization")

		data = epochTimeString
		return hash, data, nil, false
	}

	httpHashAuthenticator := hashauth.NewHTTPHashAuthenticator(
		appConfig.Secret,
		retrieveDataForHashAuthFunc,
	)
	mux := server.NewMux(server.PathHandler{
		Path:    "/auth",
		Handler: httpHashAuthenticator.Auth(httpOTPAuthenticator.Auth(nil))},

		server.PathHandler{
			Path:    "/chat",
			Handler: httpOTPAuthenticator.Auth(websocketHandler)},

		server.PathHandler{
			Path:    "/users",
			Handler: httpHashAuthenticator.Auth(userService),
		},
		server.PathHandler{
			Path: "/client/main.wasm",
			Handler: cors.Default().Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, filepath.Join(appConfig.Basedir, "client", "wasm", "main.wasm"))
			})),
		},
	)

	stopServer := make(chan os.Signal)
	server.Run(appConfig.Port, mux, stopServer, 3)
	db.Close()
	os.Exit(0)
}

func loadConfigs(cmd *cobra.Command) (dbConfig configs.DBconfig, appConfig configs.ApplicationConfig) {
	basedir, _ := cmd.Flags().GetString("basedir")

	dbConfigFile, _ := cmd.Flags().GetString("dbconfig")
	dbConfig, err := configs.LoadDBconfig(filepath.Join(basedir, "configs", dbConfigFile))
	if err != nil {
		log.Fatal(err)
	}

	appConfigfile, _ := cmd.Flags().GetString("appconfig")
	appConfig, err = configs.LoadAppconfig(filepath.Join(basedir, "configs", appConfigfile))
	if err != nil {
		log.Fatal(err)
	}

	appConfig.Basedir = basedir

	log.Infof("Using db config file: %s\nUsing app config file: %s", dbConfigFile, appConfigfile)
	return
}

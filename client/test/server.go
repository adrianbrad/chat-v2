package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"flag"
	"net/http"

	"github.com/adrianbrad/chat-v2/pkg/hashauth"

	"github.com/adrianbrad/chat-v2/pkg/chatdatabase/cmd"
)

func main() {
	port := flag.String("p", "3000", "port to serve on")
	flag.Parse()
	stop := make(chan struct{}, 1)
	go func() {
		command := cmd.NewChatDatabaseCommand()
		command.SetArgs([]string{
			"-d=../configs/test-db-config.yaml",
			"-a=../configs/application-config.yaml",
		})
		command.Execute()
		stop <- struct{}{}
	}()

	http.Handle("/", http.FileServer(http.Dir("./test/test")))

	http.HandleFunc("/getToken", func(w http.ResponseWriter, r *http.Request) {
		user := r.FormValue("user")

		req, _ := http.NewRequest("POST", "http://localhost:8080/auth", bytes.NewReader([]byte(user)))

		h := hmac.New(sha256.New, []byte("chat"))
		hash := hashauth.GenerateHash(h, []byte(user))

		r.Header.Set("Authenticate", hash)

		httpClient := &http.Client{}
		resp, _ := httpClient.Do(req)

		w.Write([]byte(resp.Header.Get("Authorization")))
	})

	go func() {
		err := http.ListenAndServe(":"+*port, nil)
		if err != nil {
			panic(err)
		}
	}()
	<-stop
}

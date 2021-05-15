package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func httpStatus(code int, w http.ResponseWriter) {
	buff, _ := marshalStruct(struct {
		Result    string `json:"result"`
		Timestamp string `json:"timestamp"`
	}{
		Result:    http.StatusText(code),
		Timestamp: time.Now().Format(time.RFC3339),
	})

	w.WriteHeader(code)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	_, _ = w.Write(buff)
}

func handlePlayground(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpStatus(http.StatusMethodNotAllowed, w)
		return
	}

	_ = playgroundTemplate.Execute(w, struct {
		Endpoint string
	}{
		Endpoint: "/graphql",
	})
}

func handleGraphQL(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "31536000")
		return
	}

	if r.Method != http.MethodPost {
		httpStatus(http.StatusMethodNotAllowed, w)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		httpStatus(http.StatusBadRequest, w)
		return
	}

	body := strings.TrimSpace(string(bodyBytes))
	ret, statusCode := executeGraphQL(r.Context(), body)

	if statusCode != http.StatusOK {
		httpStatus(statusCode, w)
		return
	}

	buff, err := marshalStruct(ret)

	if err != nil {
		httpStatus(http.StatusInternalServerError, w)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	_, _ = w.Write(buff)
}

func Http() {
	http.HandleFunc("/graphql", handleGraphQL)
	http.HandleFunc("/", handlePlayground)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}

package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Rest struct {
	Store *Store
}

type Input struct {
	Msg  string
	Time int64
}

func sendErr(w http.ResponseWriter, err error) {
	log.Println("[ERROR]", err)

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

func (rest *Rest) Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/printMeAt", rest.printMsg).Methods("GET")

	return router
}

func (rest *Rest) printMsg(w http.ResponseWriter, r *http.Request) {
	printTimeStr := r.URL.Query().Get("time")
	if printTimeStr == "" {
		sendErr(w, errors.New("time parameter required"))
		return
	}
	printTime, err := strconv.ParseInt(printTimeStr, 10, 64)
	if err != nil {
		sendErr(w, errors.Wrap(err, "error parsing time parameter"))
		return
	}

	message := r.URL.Query().Get("message")
	if message == "" {
		sendErr(w, errors.New("message parameter required"))
		return
	}

	input := &Input{message, printTime}

	id, err := genID()
	if err != nil {
		sendErr(w, errors.Wrap(err, "error generating id"))
		return
	}

	err = rest.Store.Write(id, input)
	if err != nil {
		sendErr(w, errors.Wrap(err, "error writing input to db"))
		return
	}

	go rest.Store.Print(id, input)

	w.Write([]byte("done"))
}

func genID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

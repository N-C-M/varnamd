package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/gorilla/mux"
)

func startServer() {
	if learnOnly {
		startLearnOnlyDaemon()
	} else {
		startDaemon()
	}
}

func startLearnOnlyDaemon() {
	initLearnChannels()
	varnamRPC := new(VarnamRPC)
	rpc.Register(varnamRPC)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", learnPort))
	if e != nil {
		log.Fatalln("Learn server error :", e)
	}
	log.Printf("Starting learn-only server at %d", learnPort)
	if err := http.Serve(l, nil); err != nil {
		log.Fatalln("Unable to start learn only server ", err)
	}
}

func startDaemon() {
	initLanguageChannels()
	r := mux.NewRouter()
	r.HandleFunc("/tl/{langCode}/{word}", transliterationHandler).Methods("GET")
	r.HandleFunc("/rtl/{langCode}/{word}", reverseTransliterationHandler).Methods("GET")
	r.HandleFunc("/learn", learnHandler()).Methods("POST")

	address := fmt.Sprintf(":%d", port)
	log.Printf("Starting server at %s", address)
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatalln(err)
	}
}

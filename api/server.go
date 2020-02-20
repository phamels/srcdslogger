package api

import (
	"../events"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func WebServer(ip string, port int) {
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	names := r.URL.Query()["name"]
	//	var name string
	//	if len(names) == 1 {
	//		name = names[0]
	//	}
	//	w.Write([]byte("Hello " + name))
	//})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var data, _ = json.Marshal(events.Rounds)
		w.Write([]byte(data))
	})


	go func() {
		err := http.ListenAndServe(ip+":"+strconv.Itoa(port), nil)

		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Printf("API Server listening on %s\n", ip+":"+strconv.Itoa(port))
}
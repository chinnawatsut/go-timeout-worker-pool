package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		secs, _ := r.URL.Query()["key"]
		s, _ := strconv.Atoi(secs[0])

		sid, _ := r.URL.Query()["id"]
		var wg sync.WaitGroup
		wg.Add(1)
		go func(second int) {
			fmt.Println(s, ":📲")
			time.Sleep(time.Duration(second) * time.Second)
			fmt.Println(s, ":🏁")
			wg.Done()

		}(s)
		wg.Wait()
		fmt.Println(s, ":🚀")
		fmt.Fprintf(w, "time:"+secs[0]+" id:"+sid[0])
	})
	log.Fatal(http.ListenAndServe(":8084", nil))
}

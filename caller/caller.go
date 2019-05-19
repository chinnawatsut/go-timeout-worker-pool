package main

import (
	"fmt"
	pool "go-timeout-worker-pool/workerpool"
	"log"
	"net/http"
	"strconv"
	"time"
)

func caller(w http.ResponseWriter, r *http.Request) {
	secs, _ := r.URL.Query()["worker"]
	tms, _ := r.URL.Query()["timeout"]
	numbPool, _ := strconv.Atoi(secs[0])
	timeOut, _ := strconv.Atoi(tms[0])
	listURL := []interface{}{
		[]string{"1", "3"},
		[]string{"2", "3"},
		[]string{"3", "3"},
		[]string{"4", "3"},
		[]string{"5", "3"},
		[]string{"6", "3"},
	}

	fmt.Println("worker:", numbPool, " timeout:", timeOut, "data:", len(listURL))
	type response struct {
		code int
		err  error
	}
	doneJobs := []response{}

	procc := func(resource interface{}) {
		timeout := time.Duration(7 * time.Second)
		url := resource.([]string)

		client := http.Client{
			Timeout: timeout,
		}
		a, err := client.Get("http://localhost:8084?key=" + url[1] + "&id=" + url[0])
		if err != nil {
			doneJobs = append(doneJobs, response{
				err: err,
			})
			return
		}
		defer a.Body.Close()
		doneJobs = append(doneJobs, response{
			code: a.StatusCode,
		})
		return
	}

	pw := pool.NewPool(numbPool)
	pw.Start(listURL, int64(timeOut), procc)

	fmt.Println("done with: ", len(doneJobs))
}

func handleRequests() {
	http.HandleFunc("/", caller)
	log.Fatal(http.ListenAndServe(":8083", nil))
}

func main() {
	handleRequests()
}

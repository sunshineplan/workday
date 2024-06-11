package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sunshineplan/utils/cache"
	"github.com/sunshineplan/utils/httpsvr"
	"github.com/sunshineplan/workday"
	"github.com/sunshineplan/workday/apihubs"
	"github.com/sunshineplan/workday/timor"
)

var (
	server = httpsvr.New()
	c      = cache.New[string, workday.Response](true)
)

func run() error {
	router := httprouter.New()
	server.Handler = router
	router.POST("/workday", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var data struct{ Date string }
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		t := time.Now()
		if data.Date != "" {
			var err error
			t, err = time.Parse("20060102", data.Date)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}
		resp, err := isWorkday(t)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		b, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})
	return server.Run()
}

func isWorkday(t time.Time) (resp workday.Response, err error) {
	date := t.Format("20060102")
	if v, ok := c.Get(date); ok {
		return v, nil
	}
	resp.Workday, err = workday.IsWorkday(t, timor.API, apihubs.API)
	if err != nil {
		return
	}
	resp.Date = date
	c.Set(date, resp, 7*24*time.Hour, nil)
	return
}

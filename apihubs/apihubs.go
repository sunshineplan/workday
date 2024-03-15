package apihubs

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/sunshineplan/workday"
)

var API APIHubs

var _ workday.WorkdayAPI = API

// doc: http://doc.apihubs.cn/api-78392937
const endpoint = "https://api.apihubs.cn/holiday/get?field=date,workday&date="

type response struct {
	Code    int
	Message string `json:"msg"`
	Data    struct {
		List []struct {
			Date    int
			Workday int
		}
	}
}

type APIHubs struct{}

func (APIHubs) IsWorkday(t time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	date := t.Format("20060102")
	req, _ := http.NewRequestWithContext(ctx, "GET", endpoint+date, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var res response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, err
	}
	if res.Code != 0 {
		return false, errors.New(res.Message)
	}
	for _, i := range res.Data.List {
		if strconv.Itoa(i.Date) == date {
			return i.Workday == 1, nil
		}
	}
	return false, errors.New("not found in list")
}

func (APIHubs) IsTodayWorkday() (bool, error) {
	return API.IsWorkday(time.Now())
}

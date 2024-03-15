package timor

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/sunshineplan/workday"
)

var API Timor

var _ workday.WorkdayAPI = API

// doc: https://timor.tech/api/holiday/
const endpoint = "https://timor.tech/api/holiday/info/"

type response struct {
	Code    int
	Message string
	Type    struct {
		Type int
		Week int
	}
}

type Timor struct{}

func (Timor) IsWorkday(t time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", endpoint+t.Format("2006-01-02"), nil)
	req.Header.Set("User-Agent", "timor")
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
		if err := res.Message; err != "" {
			return false, errors.New(err)
		}
		return false, errors.New("internal server error")
	}
	return res.Type.Type == 0 || res.Type.Type == 3, nil
}

func (Timor) IsTodayWorkday() (bool, error) {
	return API.IsWorkday(time.Now())
}

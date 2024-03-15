package workday

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"time"
)

type WorkdayAPI interface {
	IsWorkday(time.Time) (bool, error)
	IsTodayWorkday() (bool, error)
}

func IsWorkday(t time.Time, apis ...WorkdayAPI) (bool, error) {
	if len(apis) == 0 {
		return false, errors.New("no api provided")
	}
	var res []bool
	var errs []error
	for _, api := range apis {
		b, err := api.IsWorkday(t)
		if err != nil {
			errs = append(errs, err)
		} else {
			res = append(res, b)
		}
	}
	if len(res) == 0 {
		return false, errors.Join(append([]error{errors.New("all api failed")}, errs...)...)
	} else if len(errs) != 0 {
		for _, err := range errs {
			log.Print(err)
		}
	}
	slices.SortFunc(res, func(a, b bool) int {
		if a == b {
			return 0
		} else if a {
			return 1
		} else {
			return -1
		}
	})
	res = slices.Compact(res)
	switch len(res) {
	case 1:
		return res[0], nil
	default:
		return false, errors.New("workday results are not same")
	}
}

func IsTodayWorkday(apis ...WorkdayAPI) (bool, error) {
	return IsWorkday(time.Now(), apis...)
}

var _ WorkdayAPI = api{}

type api struct {
	endpoint string
}

func NewWorkdayAPI(endpoint string) WorkdayAPI {
	return api{endpoint}
}

type Response struct {
	Date    string `json:"date"`
	Workday bool   `json:"workday"`
}

func (api api) IsWorkday(t time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	b, _ := json.Marshal(map[string]any{"date": t.Format("20060102")})
	req, _ := http.NewRequestWithContext(ctx, "POST", api.endpoint, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, err
	}
	return res.Workday, nil
}

func (api api) IsTodayWorkday() (bool, error) {
	return api.IsWorkday(time.Now())
}

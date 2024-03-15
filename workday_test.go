package workday_test

import (
	"os"
	"testing"

	"github.com/sunshineplan/workday"
	"github.com/sunshineplan/workday/apihubs"
	"github.com/sunshineplan/workday/timor"
)

func TestWorkdayAPI(t *testing.T) {
	apis := []workday.WorkdayAPI{timor.API, apihubs.API}
	if _, err := workday.IsTodayWorkday(apis...); err != nil {
		t.Error(err)
	}
}

func TestNewAPI(t *testing.T) {
	api := os.Getenv("WORKDAY_API")
	if api == "" {
		return
	}
	if _, err := workday.NewWorkdayAPI(api).IsTodayWorkday(); err != nil {
		t.Error(err)
	}
}

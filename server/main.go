package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/sunshineplan/service"
	"github.com/sunshineplan/utils/flags"
)

var svc = service.New()

func init() {
	svc.Name = "Workday"
	svc.Desc = "workday api"
	svc.Exec = run
	svc.Options = service.Options{
		Dependencies: []string{"Wants=network-online.target", "After=network.target"},
	}
}

func main() {
	self, err := os.Executable()
	if err != nil {
		svc.Fatalln("Failed to get self path:", err)
	}
	flag.StringVar(&server.Unix, "unix", "/var/run/workday.sock", "UNIX-domain Socket")
	flag.StringVar(&server.Host, "host", "0.0.0.0", "Server Host")
	flag.StringVar(&server.Port, "port", "12345", "Server Port")
	flag.StringVar(&svc.Options.UpdateURL, "update", "", "Update URL")
	flag.StringVar(&svc.Options.PIDFile, "pid", "/var/run/workday.pid", "PID file path")
	flags.SetConfigFile(filepath.Join(filepath.Dir(self), "config.ini"))
	flags.Parse()

	if err := svc.ParseAndRun(flag.Args()); err != nil {
		svc.Fatal(err)
	}
}

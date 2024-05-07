module workday

go 1.22

require (
	github.com/julienschmidt/httprouter v1.3.0
	github.com/sunshineplan/service v1.0.19
	github.com/sunshineplan/utils v0.1.65
	github.com/sunshineplan/workday v0.0.0-00010101000000-000000000000
)

require golang.org/x/sys v0.18.0 // indirect

replace github.com/sunshineplan/workday => ../

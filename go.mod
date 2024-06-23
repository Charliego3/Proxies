module github.com/charliego3/proxies

go 1.22

require github.com/progrium/macdriver v0.5.0-preview.0.20240307055056-32e7360ca836

require (
	github.com/elazarl/goproxy v0.0.0-20231117061959-7cc037d33fb5
	github.com/google/uuid v1.6.0
)

require golang.org/x/text v0.16.0

replace github.com/progrium/macdriver => ../macdriver

module github.com/karellincoln/7-day-golang/gee-cache

go 1.17

require github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da

require geecache v0.0.0

require (
	github.com/golang/protobuf v1.5.2 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace geecache => ./geecache

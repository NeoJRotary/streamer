package streamer

import (
	"../streamUDP"
)

// Streamer ...
type Streamer struct {
	UDP         *stream.UDP
	Debug       bool
	Services    map[string]string
	Routers     map[string]Router
	ServerError int
	NoRouter    int
	InvalidBody int
}

// Router ...
type Router struct {
	Body []string
	Func func(map[string]interface{}) (int, map[string]interface{})
}

// Default Services ...
var defaultSrvs = map[string]string{
	"rest": "localhost:9899",
	"host": "localhost:9900",
	"db":   "localhost:9901",
	"auth": "localhost:9902",
}

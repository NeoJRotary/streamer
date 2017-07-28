package streamer

import (
	"log"
	"runtime/debug"

	"../streamUDP"
)

// GetStreamer ...
func GetStreamer() *Streamer {
	return &Streamer{
		Debug:       true,
		Services:    defaultSrvs,
		Routers:     make(map[string]Router),
		ServerError: 99,
		NoRouter:    98,
		InvalidBody: 95,
	}
}

// Serve ...
func (s *Streamer) Serve(srv string, middleware func(*stream.Request, *stream.Src, []byte) bool) error {
	udp, err := stream.GetStream(s.Services[srv], s.Services)
	if err != nil {
		return err
	}
	s.UDP = udp

	go func() {
		for {
			req, src, data, err := udp.Listen()
			if req == nil && err == nil {
				break
			}
			if err != nil {
				log.Println("Streamer listener Error :" + err.Error())
				continue
			}
			go s.route(req, src, data, middleware)
		}
	}()

	return nil
}

// SetRouter ...
func (s *Streamer) SetRouter(key string, body []string, call func(map[string]interface{}) (int, map[string]interface{})) {
	s.Routers[key] = Router{
		Body: body,
		Func: call,
	}
}

// PING ...
func (s *Streamer) PING() bool {
	err := s.UDP.Write([]byte{}, "listen")
	if err != nil {
		return false
	}
	return true
}

// IsErr ...
func (s *Streamer) IsErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Close ...
func (s *Streamer) Close() {
	s.UDP.Close()
}

func sameBody(obj map[string]interface{}, keys []string) bool {
	if len(obj) != len(keys) {
		return false
	}
	for _, v := range keys {
		if _, has := obj[v]; !has {
			return false
		}
	}
	return true
}

func (s *Streamer) sendResult(src *stream.Src, result int, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	err := s.UDP.Response(&stream.Response{Result: result, Data: data}, src)
	if err != nil {
		panic(err)
	}
}

func (s *Streamer) handler(src *stream.Src) {
	if err := recover(); err != nil {
		s.sendResult(src, s.ServerError, nil)
		log.Println(err)
		if s.Debug {
			debug.PrintStack()
		}
	}
}

func (s *Streamer) route(req *stream.Request, src *stream.Src, data []byte, middleware func(*stream.Request, *stream.Src, []byte) bool) {
	if middleware != nil {
		done := middleware(req, src, data)
		if done {
			return
		}
	}

	defer s.handler(src)
	if _, has := s.Routers[req.Type]; !has {
		s.sendResult(src, s.NoRouter, nil)
		return
	}

	if same := sameBody(req.Data, s.Routers[req.Type].Body); !same {
		s.sendResult(src, s.InvalidBody, nil)
		return
	}

	i, obj := s.Routers[req.Type].Func(req.Data)
	s.sendResult(src, i, obj)
}

# Streamer 0.1.0
> Light UDP Framework base on StreamUDP

Support Routing and Req/Res Handler for StreamUDP.

## Methods
**func GetStreamer() (\*Streamer, error)**  
Return `*Streamer` struct. 
```
type Streamer struct {
  UDP         *stream.UDP
  Debug       bool
  Services    map[string]string
  Routers     map[string]Router
  ServerError int
  NoRouter    int
  InvalidBody int
}
```
- UDP : StreamUDP struct. See [StreamUDP](https://github.com/NeoJRotary/streamUDP) for detail.
- Debug : if true, it will do `debug.PrintStack()` when recovering panic. Default is true.
- Services : service list for StreamUDP. It will load `defaultSrvs` in `struct.go` as default value. It is useful to set services here without do settings in each service script.
- Routers : Router list. You can use `SetRouter()` to setup. See below for more info.
- ServerError : Server Error Code. Default is 99.
- NoRouter : Router Not Found Code. Default is 98.
- InvalidBody : Invalid Data Body Code. Default is 95.
   
**func (\*Streamer) SetRouter(key string, body []string, call func(map[string]interface{}) (int, map[string]interface{}))**   
Set Router.
- key : name of the router.
- body : data body key list for validation.
- call : function to excute after body validation is passed.
   
**func (\*Streamer) Serve(srv string, middleware func(\*stream.Request, \*stream.Src, []byte) bool) error**   
Start listening host.
- srv : service name. Which you set in `Streamer.Services`.
- middleware : function to excute before routing. Set `nil` if you dont need it. Middleware should return boolean for skipping route or not.
   
**func (\*Streamer) PING() bool**   
Do PING to test listener is working or not. If it get error from listener it will return `false`.
   
**func (\*Streamer) IsErr(err error)**   
Check error is `nil` or not. If is, it will do `panic(err)`. Router will recover the panic and response error code to client.
   
**func (\*Streamer) Close()**   
Close Listner.
    
    
## Example Usage:
```
package main

import (
  "log"
  "time"

  "../../packages/streamUDP"
  "../../packages/streamer"
)

var S *streamer.Streamer

func init() {
  S := streamer.GetStreamer()
  S.Debug = true
  S.SetRouter("getData", []string{"id"}, getData)
  S.SetRouter("setData", []string{"data", "id"}, setData)
  initSystem()
}

func main() {
  log.Println("System Start")
  for {
    time.Sleep(30 * time.Second)
    if ping := S.PING(); !ping {
      restartSystem()
    }
  }
}

func restartSystem() {
  log.Println("System no response, restarting...")
  S.Close()
  time.Sleep(1 * time.Second)
  initSystem()
  log.Println("restarted!")
}

func initSystem() {
  err = S.Serve("rest", nil)
  if err != nil {
    log.Fatalln("stream init fail: " + err.Error())
  }
}

func getData(body map[string]interface{}) (int, map[string]interface{}) {
  //just dummy. Do what you want here

  data, err := getDB(body["id"].(string))
  S.IsErr(err)

  return 1, data
}

func setData(body map[string]interface{}) (int, map[string]interface{}) {
  //just dummy. Do what you want here
  err := setDB(body["id"].(string), body["data"].(map[string]string))
  S.IsErr(err)

  return 1, nil
}

```
// main.go
package main


import (
    "github.com/beego/beego/v2/server/web"
    _ "testBeego/routers"
)

func main() {
    web.BConfig.Log.AccessLogs = true
    web.BConfig.RunMode = "dev"
    web.Run() 
}
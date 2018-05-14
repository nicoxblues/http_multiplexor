package main

import (
	"fmt"
	"github.com/nicoxblues/http_multiplexor"
)

type TestSend struct {
	Param1 string `json:"param1"`
	Param2 string `json:"param2"`
}

func (t *TestSend) WriteEntity() {

}

func main() {

	httpMux := http_multiplexor.NewMux()
	stest := &TestSend{}

	httpMux.AddMethodRestFul("GET", "/", func(context *http_multiplexor.ClientCustomContext) {
		fmt.Println(stest.Param1)

	}, stest)

	httpMux.AddMethodRestFul("POST", "/TEST", func(context *http_multiplexor.ClientCustomContext) {
		fmt.Println(stest.Param1)

	}, stest)

	httpMux.RunServer()

}

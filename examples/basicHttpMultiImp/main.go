package main

import (
	"fmt"
	"http_multiplexor"
)



type TestSend struct {
	Param1 string `json:"param1"`
	Param2 string `json:"param2"`
}

func (t *TestSend) WriteEntity(ctx *http_multiplexor.ClientCustomContext) {

}

func (t *TestSend) WriteListEntity(ctx *http_multiplexor.ClientCustomContext) []http_multiplexor.Entity{


	return nil

}


var getTEST http_multiplexor.FuncMethod = func(context *http_multiplexor.ClientCustomContext) {

	fmt.Println("hola")

}


var gettestIi http_multiplexor.FuncMethod = func(context *http_multiplexor.ClientCustomContext) {

	val := context.CliRequest.EntityObject.(*TestSend).Param1

	fmt.Println("hola"  + val)

}

func main() {

	httpMux := http_multiplexor.NewMux()
	stest := &TestSend{}

	getTEST.AddSupport(http_multiplexor.SupportUploadFile)
	getTEST.AddSupport(http_multiplexor.SupportList)

	gettestIi.AddSupport(http_multiplexor.SupportList)
	httpMux.AddMethodRestFul("GET", "/", &getTEST, nil)
	httpMux.AddMethodRestFul("POST", "/te", &gettestIi, stest)















	httpMux.RunServer()

}

package main

import (
	"fmt"
	"http_multiplexor"
	"os"
)



type TestSend struct {
	Param1 string `json:"param1" form:"param1"`
	Param2 string `json:"param2" form:"param2"`
	file1 os.File `form:"file"`


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
	form ,_:= context.Ctx.MultipartForm()
	fmt.Println(form.Value["asd"])
	val := context.CliRequest.EntityObject.(*TestSend).file1
	context.Ctx.SaveUploadedFile()
	fmt.Println("hola "  + val.Name())

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

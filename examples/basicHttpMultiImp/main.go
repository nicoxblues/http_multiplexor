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

func main() {

	httpMux := http_multiplexor.NewMux()
	stest := &TestSend{}

	httpM := httpMux.AddMethodRestFul("GET", "/", func(context *http_multiplexor.ClientCustomContext) {

		fmt.Println(stest.Param1)

	}, stest)



	httpM2 := httpM.AddMethodRestFul("GET", "/TEST", func(context *http_multiplexor.ClientCustomContext) {
		fmt.Println(stest.Param1)

	}, stest)


	httpM2.UploadSupport()


	httpM2.AddMethodRestFul("GET", "/TEST", func(context *http_multiplexor.ClientCustomContext) {
		fmt.Println(stest.Param1)

	}, stest)




	httpMux.RunServer()

}

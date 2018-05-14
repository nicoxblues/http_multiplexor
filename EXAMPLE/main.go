package http_multiplexor

import (
	"fmt"

	"http_multiplexor"
)

type TestSend struct {
	Param1 string `json:"param1"`
	Param2 string `json:"param2"`
}

func (t *TestSend) WriteEntity(clientContext *http_multiplexor.ClientCustomContext) {
	param := clientContext.CliRequest.UrlParameters
	fmt.Println(param)

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


	httpMux.AddMethodRestFul("GET", "/TEST_GET", func(context *http_multiplexor.ClientCustomContext) {
		fmt.Println(stest.Param1)

	}, stest)



	httpMux.RunServer()

}

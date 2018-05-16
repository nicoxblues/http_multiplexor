package http_multiplexor

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

var multiPlex *multiplexor




type Entity interface {
	WriteEntity(*ClientCustomContext)
}

type funcMethod func(*ClientCustomContext)

type parsedClientRequest struct {
	Entity        *Entity
	RawUrl        string
	UrlParameters map[string][]string
	clientSession *AppSession

	ClientIP string
}

type ClientCustomContext struct {
	Ctx        *gin.Context
	CliRequest *parsedClientRequest
	OriginalClientRequest *http.Request
}

type HandlerMethod func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes

// TODO: quiero usarlo para que controle la sessiones,  de alguna manera, no se, es una idea que aun no termino de armar

type handlerComunitaction struct {
	hadlerMethodRef HandlerMethod
	obj             Entity
	getData         func(*handlerComunitaction) interface{}
	OriginalCtx     *gin.Context
	clientCtx       *ClientCustomContext
}

func (hc *handlerComunitaction) getDataContext(ctx *ClientCustomContext) interface{} {
	hc.clientCtx = ctx
	return hc.getData(hc)
}

type InterpreterExecuteFunction func(*ClientCustomContext)


func (hc *handlerComunitaction) getSessionFromRequest(ctx *ClientCustomContext) *AppSession{

	session , _ := Store.Get(ctx.OriginalClientRequest,"sessionStore")

	return  &AppSession{session,"sessionID"}


}

func (hc *handlerComunitaction) executeInterpreter(relativePath string, funcExec InterpreterExecuteFunction) {

	hc.hadlerMethodRef(relativePath, func(context *gin.Context) {



		hc.OriginalCtx = context
		customContext := &ClientCustomContext{Ctx: context}
		customContext.OriginalClientRequest = context.Request


		appSession := hc.getSessionFromRequest(customContext)

		ip, _ := getClientIPByRequest(context.Request)
		clientWrapperReq :=
			&parsedClientRequest{clientSession:appSession,
												UrlParameters:context.Request.URL.Query(),
												ClientIP:ip,
												RawUrl:context.Request.URL.Path}


		customContext.CliRequest = clientWrapperReq



		var valStr string = "[ '{{repeat(5, 7)}}', { _id: '{{objectId()}}', index: '{{index()}}', guid: '{{guid()}}', isActive: '{{bool()}}', balance: '{{floating(1000, 4000, 2, ', picture: 'http://placehold.it/32x32', age: '{{integer(20, 40)}}', eyeColor: '{{random('blue', 'brown', 'green')}}', name: '{{firstName()}} {{surname()}}', gender: '{{gender()}}', company: '{{company().toUpperCase()}}', email: '{{email()}}', phone: '+1 {{phone()}}', address: '{{integer(100, 999)}} {{street()}}, {{city()}}, {{state()}}, {{integer(100, 10000)}}', about: '{{lorem(1, 'paragraphs')}}', registered: '{{date(new Date(2014, 0, 1), new Date(), 'YYYY-MM-ddThh:mm:ss Z')}}', latitude: '{{floating(-90.000001, 90)}}', longitude: '{{floating(-180.000001, 180)}}', tags: [ '{{repeat(7)}}', '{{lorem(1, 'words')}}' ], friends: [ '{{repeat(3)}}', { id: '{{index()}}', name: '{{firstName()}} {{surname()}}' } ], greeting: function (tags) { return 'Hello, ' + this.name + '! You have ' + tags.integer(1, 10) + ' unread messages.'; }, favoriteFruit: function (tags) { var fruits = ['apple', 'banana', 'strawberry']; return fruits[tags.integer(0, fruits.length - 1)]; } } ]"

		appSession.Values["codigoRojo"] = []byte(valStr)

		hc.getDataContext(customContext)


		funcExec(customContext)

		appSession.Save(customContext)

	})

}

func (hc *handlerComunitaction) getMethodHandler() *HandlerMethod {

	return &hc.hadlerMethodRef

}

type GinWrapperHandler func() handlerComunitaction

type multiplexor struct {
	routerEngine *gin.Engine
	methodMap    map[string]*GinWrapperHandler
}

func getClientIPByRequest(req *http.Request) (ip string, err error) {

	// Try via request
	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		log.Printf("debug: Getting req.RemoteAddr %v", err)
		return "", err
	} else {
		log.Printf("debug: With req.RemoteAddr found IP:%v; Port: %v", ip, port)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		message := fmt.Sprintf("debug: Parsing IP from Request.RemoteAddr got nothing.")
		log.Printf(message)
		return "", fmt.Errorf(message)

	}
	log.Printf("debug: Found IP: %v", userIP)
	return userIP.String(), nil

}
var Store sessions.Store
func NewMux() *multiplexor {


	 Store = NewStoreForSessionType("asd",[]byte("super-secret_ohsi_key_megaHardcodeada"))
	//Store  = NewCookieStore([]byte("super-secret_ohsi_key_megaHardcodeada"))




	multiPlex = new(multiplexor)
	r := gin.Default()

	multiPlex.methodMap = make(map[string]*GinWrapperHandler)

	var getFunction GinWrapperHandler = func() handlerComunitaction {

		handler := handlerComunitaction{hadlerMethodRef: r.GET}

		handler.getData = func(hc *handlerComunitaction) interface{} {

			hc.obj.WriteEntity(hc.clientCtx)
			return hc.obj
		}

		return handler
	}

	var postFunction GinWrapperHandler = func() handlerComunitaction {

		handler := handlerComunitaction{hadlerMethodRef: r.POST}

		handler.getData = func(hc *handlerComunitaction) interface{} {

			if hc.obj != nil { //GET, POST
				hc.OriginalCtx.Bind(hc.obj)
				fmt.Println(hc.obj)

			}
			return &(hc.obj)

		}

		return handler

	}

	multiPlex.routerEngine = r

	multiPlex.methodMap["GET"] = &getFunction

	multiPlex.methodMap["POST"] = &postFunction

	return multiPlex

}
func (multi *multiplexor) RunServer() {
	multi.routerEngine.Run()
}

func (multi *multiplexor) AddMethodRestFul(methodName string, relativePath string, fMethod funcMethod, obj Entity) {

	method := strings.ToUpper(methodName)

	if methodFunc, ok := multi.methodMap[method]; ok {

		//methodFunc.(func(string, gin.HandlerFunc))(relativePath, func(context *gin.Context) {
		wrap := (*methodFunc)()
		wrap.obj = obj
		//wrap.getMethodHandler()
		wrap.executeInterpreter(relativePath, func(context *ClientCustomContext) {
			log.Println("Interpreter ejecutado con exito ! ")


			fMethod(context)

		})

	}

}

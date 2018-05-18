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
	Entity              *Entity
	RawUrl              string
	UrlParameters       map[string][]string
	ClientCookieSession *AppSession

	ClientIP string
}

type ClientCustomContext struct {
	Ctx        *gin.Context
	CliRequest *parsedClientRequest
	OriginalClientRequest *http.Request

}

type HandlerMethod func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes

// TODO: quiero usarlo para que controle la sessiones,  de alguna manera, no se, es una idea que aun no termino de armar

type handlerCommunication struct {
	handlerMethodRef HandlerMethod
	obj              Entity
	getData          func(*handlerCommunication) interface{}
	OriginalCtx      *gin.Context
	clientCtx        *ClientCustomContext
}

func (hc *handlerCommunication) getDataFromContext(ctx *ClientCustomContext) interface{} {
	hc.clientCtx = ctx
	return hc.getData(hc)
}

type InterpreterExecuteFunction func(*ClientCustomContext)


func (hc *handlerCommunication) getSessionCookieFromRequest(ctx *ClientCustomContext) *AppSession{

	//sessionDB, err := mgo.Dial("localhost")

	//fmt.Println(err)
	//c := sessionDB.DB("baseName?")
//	fmt.Println(c)

	session , _ := Store.Get(ctx.OriginalClientRequest,"sessionStore")

	return  &AppSession{session,"sessionID"}


}

func (hc *handlerCommunication) executeInterpreter(relativePath string, funcExec InterpreterExecuteFunction) {

	hc.handlerMethodRef(relativePath, func(context *gin.Context) {



		hc.OriginalCtx = context
		customContext := &ClientCustomContext{Ctx: context}
		customContext.OriginalClientRequest = context.Request


		appSession := hc.getSessionCookieFromRequest(customContext)

		//ip, _ := getClientIPByRequest(context.Request)
		ip := context.ClientIP()
		//fmt.Println(ip2)
		clientWrapperReq :=
			&parsedClientRequest{ClientCookieSession:appSession,
												UrlParameters:context.Request.URL.Query(),
												ClientIP:ip,
												RawUrl:context.Request.URL.Path}


		customContext.CliRequest = clientWrapperReq



		var valStr = "asdasd"

		appSession.Values["codigoRojo"] = []byte(valStr)

		hc.getDataFromContext(customContext)


		funcExec(customContext)

		appSession.save(customContext)

	})

}

func (hc *handlerCommunication) getMethodHandler() *HandlerMethod {

	return &hc.handlerMethodRef

}

type GinWrapperHandler func() *handlerCommunication

type multiplexor struct {
	routerEngine *gin.Engine
	methodMap    map[string]*GinWrapperHandler
	perentMultiplex *multiplexor
	basePath string
	uploadtestSup string
}

func (multi *multiplexor) UploadSupport () *multiplexor{
	multi.uploadtestSup = "gruoup with support"

	return multi

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
	multiPlex.perentMultiplex = nil
	r := gin.Default()

	multiPlex.methodMap = make(map[string]*GinWrapperHandler)

	var getFunction GinWrapperHandler = func() *handlerCommunication {

		handler := handlerCommunication{handlerMethodRef: r.GET}

		handler.getData = func(hc *handlerCommunication) interface{} {


			hc.obj.WriteEntity(hc.clientCtx)
			return hc.obj
		}

		return &handler
	}

	var postFunction GinWrapperHandler = func() *handlerCommunication {

		handler := handlerCommunication{handlerMethodRef: r.POST}

		handler.getData = func(hc *handlerCommunication) interface{} {

			if hc.obj != nil { //GET, POST
				hc.OriginalCtx.Bind(hc.obj)
				fmt.Println(hc.obj)

			}
			return &(hc.obj)

		}

		return &handler

	}

	multiPlex.routerEngine = r

	multiPlex.methodMap["GET"] = &getFunction

	multiPlex.methodMap["POST"] = &postFunction

	return multiPlex

}
func (multi *multiplexor) RunServer(port ...string) {
	// me cubro por las dudas, no  se puede

	multi.routerEngine.Run(port ...)


}



func (multi *multiplexor) AddMethodRestFul(methodName string, relativePath string, fMethod funcMethod, obj Entity) *multiplexor {

	method := strings.ToUpper(methodName)

	if methodFunc, ok := multi.methodMap[method]; ok {

		//methodFunc.(func(string, gin.HandlerFunc))(relativePath, func(context *gin.Context) {
		wrap := (*methodFunc)()
		wrap.obj = obj
		path := multi.basePath + relativePath

		wrap.executeInterpreter(path, func(context *ClientCustomContext) {
			fmt.Println(multi.uploadtestSup)
			log.Println("Interpreter ejecutado con exito ! ")


			fMethod(context)

		})

	}
	// TODO: ver de hacer que se devuelve a si mismo, con un objeto interno "childMultiplexor"
	return &multiplexor{routerEngine:multi.routerEngine,perentMultiplex:multi, basePath:multi.basePath + relativePath,methodMap:multi.methodMap}


}

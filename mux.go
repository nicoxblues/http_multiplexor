package http_multiplexor

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

//var multiPlex *multiplexor




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


	session , _ := Store.Get(ctx.OriginalClientRequest,"sessionStore")
	//session.ID ="testhnolamundo"

	session.Values["sessionType"] = 1


	return  &AppSession{session,
				"sessionID", 0}


}

func (hc *handlerCommunication) executeInterpreter(relativePath string, funcExec InterpreterExecuteFunction) {

	hc.handlerMethodRef(relativePath, func(context *gin.Context) {



		hc.OriginalCtx = context
		customContext := &ClientCustomContext{	Ctx: context,
											  	OriginalClientRequest:context.Request}



		appSession := hc.getSessionCookieFromRequest(customContext)


		if !(SessionManager.sessions[appSession.ID] ){
			SessionManager.register <- appSession

		}




		clientWrapperReq := &parsedClientRequest{
									ClientCookieSession:appSession,
									UrlParameters: context.Request.URL.Query(),
									ClientIP: context.ClientIP(),
									RawUrl: context.Request.URL.Path,
								}


		customContext.CliRequest = clientWrapperReq


		hc.getDataFromContext(customContext)


		funcExec(customContext)

		appSession.Save(customContext)

	})

}

func (hc *handlerCommunication) getMethodHandler() *HandlerMethod {

	return &hc.handlerMethodRef

}

type GinWrapperHandler func() handlerCommunication

type multiplexor struct {
	routerEngine *gin.Engine
	sessionHandler *SessionHandler
	methodMap    map[string]*GinWrapperHandler
}

var Store sessions.Store
var SessionManager * SessionHandler
func NewMux() *multiplexor {


	Store = NewStoreForSessionType("asd1",[]byte("super-secret_ohsi_key_megaHardcodeada"))



	multiPlex := new(multiplexor)
	r := gin.Default()

	multiPlex.startHandlerSessionConn()

	SessionManager = multiPlex.sessionHandler

	multiPlex.methodMap = make(map[string]*GinWrapperHandler)

	var getFunction GinWrapperHandler = func() handlerCommunication {

		handler := handlerCommunication{handlerMethodRef: r.GET}

		handler.getData = func(hc *handlerCommunication) interface{} {


			hc.obj.WriteEntity(hc.clientCtx)
			return hc.obj
		}

		return handler
	}

	var postFunction GinWrapperHandler = func() handlerCommunication {

		handler := handlerCommunication{handlerMethodRef: r.POST}

		handler.getData = func(hc *handlerCommunication) interface{} {

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

func (multi *multiplexor) startHandlerSessionConn(){

	multi.sessionHandler = newSessionHandler()
	go multi.sessionHandler.start()





}

func (multi *multiplexor) AddMethodRestFul(methodName string, relativePath string, fMethod funcMethod, obj Entity) {

	method := strings.ToUpper(methodName)

	if methodFunc, ok := multi.methodMap[method]; ok {

		wrap := (*methodFunc)()
		wrap.obj = obj

		wrap.executeInterpreter(relativePath, func(context *ClientCustomContext) {
			log.Println("Interpreter ejecutado con exito ! ")


			fMethod(context)

		})

	}

}

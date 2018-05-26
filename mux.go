package http_multiplexor

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)




type Entity interface {
	WriteEntity(*ClientCustomContext)
	WriteListEntity(*ClientCustomContext) [] Entity

}


// proporciona soporte para determinadas procesos o tareas, upadte de archivos, return type como listas de un determinado objeto etc...
type handlerSupport func(communication *handlerCommunication)

type SupportType int



const (
	SupportUploadFile SupportType = 0
	SupportList SupportType = 1


)





type FuncMethod func(*ClientCustomContext)

var supportMap = make(map[*FuncMethod]handlerSupport)


func (fm *FuncMethod)  AddSupport(supportType SupportType){

	switch supportType  {


	case SupportUploadFile:

			supportMap[fm] = func(hc *handlerCommunication) {

				log.Println("********* Soporte para upload no implementado *************")
			}

	case SupportList:

			supportMap[fm] = func(communication *handlerCommunication) {
				if communication.obj != nil {
					communication.listEntity = communication.obj.WriteListEntity(communication.clientCtx)
				}
				//log.Println("******** Soporte para listas no implementado ***********")

			}


	}





}

type parsedClientRequest struct {
	EntityObject        Entity
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

	// la info puede ser solicitada en lista o individual
	obj              Entity
	listEntity		 []Entity


	getData          func(*handlerCommunication) interface{}
	OriginalCtx      *gin.Context
	handlerSupport 	 handlerSupport
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

	if session.Values["sessionID"] == nil{
		session.Values["sessionID"] = "test"


	}

	return  &AppSession{session,
				"sessionID", 0}


}
func (hc *handlerCommunication) executeHandlersSupport() {

	hc.handlerSupport(hc)


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


		hc.executeHandlersSupport()
		hc.getDataFromContext(customContext)
		customContext.CliRequest.EntityObject = hc.obj

		// ejecucion  de la funcion del cliente
		funcExec(customContext)

		appSession.save(customContext)

	})

}

func (hc *handlerCommunication) getMethodHandler() *HandlerMethod {

	return &hc.handlerMethodRef

}

type GinWrapperHandler func(*multiplexor) *handlerCommunication

type multiplexor struct {

	routerEngine    *gin.Engine
	sessionHandler  *SessionHandler
	methodMap       map[string]*GinWrapperHandler
	parentMultiplex *multiplexor
	basePath        string
	uploadTestSup   string
	ChildMultiplex  *multiplexor

	shouldSendList bool

}

var Store sessions.Store
var SessionManager * SessionHandler
func NewMux() *multiplexor {

// TODO reveer esto
	Store = NewStoreForSessionType("1",[]byte("super-secret_ohsi_key_megaHardcodeada"))



	multiPlex := new(multiplexor)
	r := gin.Default()

	multiPlex.startHandlerSessionConn()

	SessionManager = multiPlex.sessionHandler

	multiPlex.methodMap = make(map[string]*GinWrapperHandler)


	// TODO: el unico sentido de tenerlo asi, es que tengo acceso al multiplexador root
	var getFunction GinWrapperHandler = func(multi *multiplexor) *handlerCommunication {

		handler := handlerCommunication{handlerMethodRef: r.GET}

		handler.getData = func(hc *handlerCommunication) interface{} {


			if hc.obj != nil {
				//hc.handlerSupport(hc)

				// Esto quedo medio deprecado TODO, ver de intergrar a handlerSupport
				hc.obj.WriteEntity(hc.clientCtx)


			}

			return hc.obj
		}



		return &handler
	}

	var postFunction GinWrapperHandler = func(multi *multiplexor) *handlerCommunication {

		handler := handlerCommunication{handlerMethodRef: r.POST}

		handler.getData = func(hc *handlerCommunication) interface{} {

			if hc.obj != nil { // POST
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
//
//func (multi *multiplexor) ListSupport () *multiplexor{
//
//	multi.shouldSendList = true
//
//	return  multi
//
//
//}
//
//func (multi *multiplexor) UploadSupport () *multiplexor{
//	multi.uploadTestSup = "gruoup with support"
//
//	return multi
//
//}
func (multi *multiplexor) RunServer(port ...string) {
	// me cubro por las dudas, no  se puede

	multi.routerEngine.Run(port ...)


}

func (multi *multiplexor) startHandlerSessionConn(){

	multi.sessionHandler = newSessionHandler()
	go multi.sessionHandler.start()



}

func (multi *multiplexor) AddMethodRestFul(methodName string, relativePath string, fMethod *FuncMethod, obj Entity) *multiplexor {

	method := strings.ToUpper(methodName)

	if methodFunc, ok := multi.methodMap[method]; ok {

		wrap := (*methodFunc)(multi)
		wrap.obj = obj
		wrap.handlerSupport = supportMap[fMethod]

		wrap.executeInterpreter(relativePath, func(context *ClientCustomContext) {
			log.Println("Interpreter ejecutado con exito ! ")





			(*fMethod)(context)

		})

	}
	multiChild :=  &multiplexor{routerEngine:multi.routerEngine, parentMultiplex:multi, basePath:multi.basePath + relativePath,methodMap:multi.methodMap}
	multi.ChildMultiplex = multiChild

	return multi


}

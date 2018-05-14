package http_multiplexor

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"strings"
)

var multiPlex *multiplexor

type Entity interface {
	WriteEntity()
}

type funcMethod func(*ClientCustomContext)

type parsedClientRequest struct {
	entity        *Entity
	urlRaw        string
	urlParameters map[string]string

	clientIP string
}

type ClientCustomContext struct {
	ctx        *gin.Context
	cliRequest parsedClientRequest
}

type HandlerMethod func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes

// TODO: quiero usarlo para que controle la sessiones,  de alguna manera, no se, es una idea que aun no termino de armar

type handlerComunitaction struct {
	hadlerMethodRef HandlerMethod
	obj             Entity
	getData         func(*handlerComunitaction) interface{}
	OriginalCtx     *gin.Context
}

func (hc *handlerComunitaction) getDataContext(ctx *ClientCustomContext) interface{} {
	return hc.getData(hc)
}

type InterpreterExecuteFunction func(*ClientCustomContext)

func (hc *handlerComunitaction) executeInterpreter(relativePath string, funcExec InterpreterExecuteFunction) {

	hc.hadlerMethodRef(relativePath, func(context *gin.Context) {

		hc.OriginalCtx = context
		customContext := &ClientCustomContext{ctx: context}

		clientWrapperReq := parsedClientRequest{}
		ip, _ := getClientIPByRequest(context.Request)

		clientWrapperReq.clientIP = ip
		hc.getDataContext(customContext)

		customContext.cliRequest = clientWrapperReq
		funcExec(customContext)

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

func NewMux() *multiplexor {

	multiPlex = new(multiplexor)
	r := gin.Default()

	multiPlex.methodMap = make(map[string]*GinWrapperHandler)

	var getFunction GinWrapperHandler = func() handlerComunitaction {

		handler := handlerComunitaction{hadlerMethodRef: r.GET}

		handler.getData = func(hc *handlerComunitaction) interface{} {

			hc.obj.WriteEntity()
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

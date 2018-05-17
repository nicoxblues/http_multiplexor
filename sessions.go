package http_multiplexor

import (
	"github.com/gorilla/sessions"

)

type AppSession struct {
	*sessions.Session
	sessionID string

}

func (ass *AppSession) Save(clientContext *ClientCustomContext){
	ass.Session.Save(clientContext.OriginalClientRequest , clientContext.Ctx.Writer)


}


func (ass *AppSession) BindJson(obj interface{}){


}




func NewStoreForSessionType (typeSession string,keyPairs ...[]byte) sessions.Store{

	if "1" == typeSession{
		return sessions.NewCookieStore(keyPairs ...)
	}else{
		return sessions.NewFilesystemStore("",[]byte("fileKeyStore"))
	}


	return nil
}

func NewCookieStore(keyPairs ...[]byte) *sessions.CookieStore{
	return sessions.NewCookieStore(keyPairs ...)


}



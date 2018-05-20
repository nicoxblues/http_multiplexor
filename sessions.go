package http_multiplexor

import (
	"github.com/gorilla/sessions"

)



const (
	LoguinUser = 0
	AdminUser  = 2
	GuestUser  = 3


)

type AppSession struct {
	*sessions.Session
	sessionID   string
	sessionType int

}

func (ass *AppSession) save(clientContext *ClientCustomContext){
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


}

func NewCookieStore(keyPairs ...[]byte) *sessions.CookieStore{
	return sessions.NewCookieStore(keyPairs ...)


}



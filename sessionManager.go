package http_multiplexor

type SessionHandler struct {
	sessions      map[string]bool
	broadcast    chan []byte
	register     chan *AppSession
	unregister   chan *AppSession

}


func newSessionHandler () *SessionHandler {

	return &SessionHandler{
				sessions:   make(map[string]bool),
				broadcast:  make(chan []byte),
				register:   make(chan *AppSession),
				unregister: make(chan *AppSession),
			}



}


func (manager *SessionHandler) start() {

	for {
		select {
			case sessionReg := <-manager.register:
				manager.sessions[sessionReg.ID] = true


			case sessionUnReg := <-manager.unregister:
				if _, ok := manager.sessions[sessionUnReg.ID]; ok {
					delete(manager.sessions, sessionUnReg.ID)

				}

		}
	}
}
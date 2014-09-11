//Session management implementation for drops
package session

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/satori/go.uuid"

	"github.com/mkasner/drops/element"
)

type session struct {
	ActiveDOM *element.DOM
	Data      interface{}
	Messages  element.Messages
}

var mutex *sync.RWMutex

var sessionConfig struct {
	CleanInterval int
}

var sessionStore map[string]*session

//Creates new session store
func NewSessionStore() {
	sessionStore = make(map[string]*session)
	mutex = &sync.RWMutex{}
}

func SessionStore() map[string]*session {
	return sessionStore
}

//Creates session and return session id
func CreateSession(id string) string {
	if id == "" {
		id = uuid.NewV4().String()
	}
	session := &session{}
	mutex.Lock()
	sessionStore[id] = session
	fmt.Printf("Session size: %+v\n", len(sessionStore))
	mutex.Unlock()
	return id
}

//sets active dom for session
func SetSessionActiveDOM(id string, dom *element.DOM) {
	mutex.Lock()
	// fmt.Printf("Setting activeDOM: %+v", dom)
	if _, ok := sessionStore[id]; ok {
		sessionStore[id].ActiveDOM = dom
	}
	mutex.Unlock()
}

func GetSessionActiveDOM(id string) *element.DOM {
	mutex.RLock()
	defer mutex.RUnlock()
	// fmt.Printf("Getting activeDOM for session: %+v", id)
	if session, ok := sessionStore[id]; ok {
		// fmt.Printf("Getting activeDOM: %+v", session.ActiveDOM)
		return session.ActiveDOM
	}
	return nil
}

func GetSessionData(id string) (interface{}, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	// fmt.Printf("Getting activeDOM for session: %+v", id)
	if session, ok := sessionStore[id]; ok {
		// fmt.Printf("Getting activeDOM: %+v", session.ActiveDOM)
		return session.Data, nil
	}
	return nil, errors.New("Session does not exist")
}

func SetSessionData(id string, data interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()
	// fmt.Printf("Setting activeDOM: %+v", dom)
	if _, ok := sessionStore[id]; ok {
		sessionStore[id].Data = data
		return nil
	}

	return errors.New("Session does not exist")
}

func GetSessionMessages(id string) element.Messages {
	mutex.RLock()
	defer mutex.RUnlock()
	// fmt.Printf("Getting activeDOM for session: %+v", id)
	if session, ok := sessionStore[id]; ok {
		// fmt.Printf("Getting activeDOM: %+v", session.ActiveDOM)
		return session.Messages
	}
	return nil
}

func SetSessionMessages(id string, messages element.Messages) {
	mutex.Lock()
	// fmt.Printf("Setting activeDOM: %+v", dom)
	if _, ok := sessionStore[id]; ok {
		sessionStore[id].Messages = messages
	}
	mutex.Unlock()
}

//Checks if session exist
func SessionExist(id string) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	if _, ok := sessionStore[id]; ok {
		return true
	}
	return false
}

//Cleans session store based on expired flag
func CleanSessionStore() {

}

//Deletes session
func DeleteSession(sessionId string) {
	mutex.Lock()
	delete(sessionStore, sessionId)
	mutex.Unlock()
}

//extracts session id from provided sessionId or http.Request
//sessionId in parameters can be provided from websocket params
func GetSessionId(r *http.Request, sessionId string) string {
	fmt.Printf("Session id on handler param %v\n", sessionId)
	sessionFound := ""
	if sessionId != "" && SessionExist(sessionId) {
		sessionFound = sessionId
	}
	if sessionFound == "" && r != nil {
		sessionCookie, err := r.Cookie("session")
		if err == nil {
			sessionFound = sessionCookie.Value
		}
	}
	return sessionFound
}

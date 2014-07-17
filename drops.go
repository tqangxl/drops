//Currently implementing http serving with goji
// I'll drop that in the future
package drops

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mkasner/drops/element"
	"github.com/mkasner/drops/router"
	"github.com/mkasner/drops/session"
)

var rtr *router.Router

func NewDrops() {
	rtr = router.New()
	session.NewSessionStore()
}

// GET is a shortcut for r.Handle("GET", path, handle)
func GET(path string, handle router.Handle) {
	rtr.Handle("GET", path, handle)
}

// POST is a shortcut for r.Handle("POST", path, handle)
func POST(path string, handle router.Handle) {
	rtr.Handle("POST", path, handle)
}

// PUT is a shortcut for r.Handle("PUT", path, handle)
func PUT(path string, handle router.Handle) {
	rtr.Handle("PUT", path, handle)
}

// PATCH is a shortcut for r.Handle("PATCH", path, handle)
func PATCH(path string, handle router.Handle) {
	rtr.Handle("PATCH", path, handle)
}

// DELETE is a shortcut for r.Handle("DELETE", path, handle)
func DELETE(path string, handle router.Handle) {
	rtr.Handle("DELETE", path, handle)
}

func EVENT(path string, handle router.Handle) {
	rtr.Handle("EVENT", path, handle)
}

func HandleFunc(path string, handle http.HandlerFunc) {
	rtr.HandlerFunc("GET", path, handle)
}

func Serve(port string) {
	rtr.ServeFiles("/assets/*filepath", http.Dir("assets"))
	rtr.HandlerFunc("GET", "/ws", ServeWs)
	// rtr.HandlerFunc("GET", "/*", ResourceHandler)
	log.Fatal(http.ListenAndServe(":"+port, rtr))
}

func ResourceHandler(w http.ResponseWriter, r *http.Request, dom *element.DOM) *element.DOM {
	// log.Printf("Resource handler: %s - %s\n", r.Method, r.URL.Path)
	var response string

	// handle, paramsFromRequest, _ := rtr.Lookup(r.Method, r.URL.Path)
	// var dom *element.DOM
	// if handle != nil {
	// 	log.Printf("Routing success: %v\n", paramsFromRequest)
	// 	dom = handle(nil, nil, paramsFromRequest)
	// } else {
	// 	log.Println("Routing failure, no handler")
	// }
	// if drops.ActiveDOM == nil {
	if w != nil && dom != nil {
		var sessionId string
		sessionCookie, err := r.Cookie("session")
		// fmt.Printf("sessionCookie: %+v\n", pretty.Formatter(sessionCookie))
		if err != nil {
			// fmt.Printf("Error fetching cookie %v\n", err)
			//No cookie found
			sessionId = session.CreateSession("")
			sessionCookie = &http.Cookie{Name: "session", Value: sessionId, Path: "/"}
			// fmt.Printf("Created Cookie: %s: %+v\n", sessionCookie.Name, sessionCookie.Value)
		} else {
			// log.Printf("Cookie: %s: %+v\n", sessionCookie.Name, sessionCookie.Value)
			sessionId = sessionCookie.Value
			if !session.SessionExist(sessionId) {
				sessionId = session.CreateSession(sessionId)
			}
			sessionCookie = &http.Cookie{Name: "session", Value: sessionId, Path: "/"}
			// sessionCookie.Value = sessionId
		}
		// fmt.Printf("sessionId: %s\n", sessionId)
		session.SetSessionActiveDOM(sessionId, dom)
		// fmt.Printf("\ndom set as active: %+v\n", pretty.Formatter(dom.Id))

		// ActiveDOM = dom
		buffer := Render(&dom.View)
		response = buffer.String()
		// cookie := &http.Cookie{Name: "session", Value: sessionId}
		// http.Set
		http.SetCookie(w, sessionCookie)

		fmt.Fprint(w, response)

		return nil
	}
	return dom
}

//Used for static file serving
func Assets(w http.ResponseWriter, r *http.Request) {
	// logger.Errorf("Assets handler: %s", r.URL.Path[1:])
	http.ServeFile(w, r, r.URL.Path[1:])

}

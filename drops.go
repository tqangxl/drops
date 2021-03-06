//Currently implementing http serving with goji
// I'll drop that in the future
package drops

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mkasner/drops/element"
	"github.com/mkasner/drops/event"
	"github.com/mkasner/drops/protocol"
	"github.com/mkasner/drops/router"
	"github.com/mkasner/drops/session"
)

var (
	rtr *router.Router
)

func NewDrops() {
	router.InitRouter()
	rtr = router.GetRouter()
	session.NewSessionStore()
	event.NewEventStore()
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

//register handlers from events
func RegisterHandlers(events event.Events) error {
	for _, e := range events {

		rtr.Handle(e.Type, e.Route, e.Handler)

	}
	return nil
}

func Serve(port, assetsPath, assetsUrl string) {
	if assetsPath != "" {
		fmt.Printf("Assets path: %s\n", assetsPath)
		rtr.ServeFiles("/"+assetsUrl+"/*filepath", http.Dir(assetsPath))
	} else {
		rtr.ServeFiles("/assets/*filepath", http.Dir("assets"))

	}
	rtr.HandlerFunc("GET", "/ws", ServeWs)

	// rtr.HandlerFunc("GET", "/*", ResourceHandler)
	log.Fatal(http.ListenAndServe(":"+port, rtr))
}

func ResourceHandler(w http.ResponseWriter, r *http.Request, dropsResponse *protocol.DropsResponse) *protocol.DropsResponse {
	// log.Printf("Resource handler: %s - %s\n", r.Method, r.URL.Path)
	var response string

	if w != nil && dropsResponse.Dom != nil {
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
		session.SetSessionActiveDOM(sessionId, dropsResponse.Dom)
		// fmt.Printf("\ndom set as active: %+v\n", pretty.Formatter(dom.Id))

		// ActiveDOM = dom
		buffer := element.Render(&dropsResponse.Dom.View)
		response = buffer.String()
		// cookie := &http.Cookie{Name: "session", Value: sessionId}
		// http.Set
		http.SetCookie(w, sessionCookie)

		fmt.Fprint(w, response)

		return nil
	}
	return dropsResponse
}

//Used for static file serving
func Assets(w http.ResponseWriter, r *http.Request) {
	// logger.Errorf("Assets handler: %s", r.URL.Path[1:])
	http.ServeFile(w, r, r.URL.Path[1:])

}

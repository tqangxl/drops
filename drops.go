//Currently implementing http serving with goji
// I'll drop that in the future
package drops

import (
	"fmt"
	"github.com/zenazn/goji"
	"log"
	"net/http"
)

var router *Router

func NewDrops() {
	router = NewRouter()
}

// GET is a shortcut for router.Handle("GET", path, handle)
func GET(path string, handle Handle) {
	router.Handle("GET", path, handle)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func POST(path string, handle Handle) {
	router.Handle("POST", path, handle)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func PUT(path string, handle Handle) {
	router.Handle("PUT", path, handle)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func PATCH(path string, handle Handle) {
	router.Handle("PATCH", path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func DELETE(path string, handle Handle) {
	router.Handle("DELETE", path, handle)
}

// Api path for returning  regular htp handler
func HandleFunc(path string, handle http.HandlerFunc) {
	goji.Get(path, handle)
}

func Serve() {
	goji.Get("/assets/*", Assets)
	goji.Handle("/ws", ServeWs)
	goji.Handle("/*", ResourceHandler)
	goji.Serve()
}

func ResourceHandler(w http.ResponseWriter, r *http.Request) {
	// io.WriteString(w, "Deployments")
	log.Printf("Resource handler: %s - %s\n", r.Method, r.URL.Path)
	var response string

	handle, paramsFromRequest, _ := router.Lookup(r.Method, r.URL.Path)
	var dom *DOM
	if handle != nil {
		log.Printf("Routing success: %v\n", paramsFromRequest)
		dom = handle(paramsFromRequest)
	} else {
		log.Println("Routing failure, no handler")
	}
	// if drops.ActiveDOM == nil {
	if dom != nil {
		ActiveDOM = dom
		buffer := dom.Render()
		response = buffer.String()
	}
	// } else {    //Only used for testing purposes
	// buffer := drops.ActiveDOM.Render()
	// response = buffer.String()
	// }
	fmt.Fprint(w, response)

	// fmt.Fprintf(w, "%v\n", result)
}

//Used for static file serving
func Assets(w http.ResponseWriter, r *http.Request) {
	// logger.Errorf("Assets handler: %s", r.URL.Path[1:])
	http.ServeFile(w, r, r.URL.Path[1:])

}

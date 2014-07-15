//Event dispatcher
package drops

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mkasner/drops/element"
	"github.com/mkasner/drops/router"
)

type EventHandler interface {
	Handle(data map[string]interface{}) []byte
}
type EventDispatcher map[string]EventHandler

var eventDispatcher EventDispatcher

func NewEventDispatcher() {
	eventDispatcher = EventDispatcher{}
}

func AddEventHandler(event string, handler EventHandler) {
	eventDispatcher[event] = handler
}

//Handles event and creates DOM patch
func HandleExecute(handle router.Handle, paramsFromRequest router.Params) ([]byte, error) {
	var dom *element.DOM
	if handle != nil {
		fmt.Printf("Routing success: %v\n", paramsFromRequest)

		dom = handle(paramsFromRequest)
		// PrintDOM(ActiveDOM, "1")

		fmt.Printf("ActiveDOM: %+v\n", ActiveDOM)
		fmt.Printf("New DOM: %+v\n", dom)
		// fmt.Printf("Active dom is the same to new DOM: %v\n", *ActiveDOM == *dom)
		// PrintDOM(dom, "2")
	} else {
		fmt.Println("Routing failure, no handler")

		dom = ActiveDOM
	}

	patches := Diff(&ActiveDOM.View, &dom.View)
	fmt.Printf("Patches: %+v\n", patches)
	ActiveDOM = dom
	message, err := json.Marshal(patches)
	if err != nil {
		log.Println("Error marshaling patch")
	}
	fmt.Printf("Patches generated: %+v\n", string(message))
	return message, nil
}

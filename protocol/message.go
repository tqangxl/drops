//Contains messages and rules for communicating with front end
package protocol

import (
	"encoding/json"
	"log"

	"github.com/mkasner/drops/element"
)

//Response message from server to client
type DropsResponse struct {
	Dom       *element.DOM
	Route     string
	ActiveDom *element.DOM
}

type ResponseMessage struct {
	Patches []Patch `json:"patches"`
	Route   string  `json:"route"`
}

//Patch message
type Patch struct {
	Element string `json:"element"`
	Payload string `json:"payload"`
}

//Generates final message that's sent to client
func GenerateMessage(response *DropsResponse) []byte {

	//generate patches
	pathces := generatePatches(response)
	responseMessage := &ResponseMessage{Patches: pathces, Route: response.Route}
	message, err := json.Marshal(responseMessage)
	if err != nil {
		log.Println("Error marshaling patch")
	}
	return message
}

//compares old and new dom and generates pathces
func generatePatches(response *DropsResponse) []Patch {
	var patches []Patch

	if response.ActiveDom != nil {
		patches = Diff(&response.ActiveDom.View, &response.Dom.View)
		// fmt.Printf("Patches: %+v\n", patches)

	} else {
		patches = make([]Patch, 0)
	}
	return patches
}

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package drops

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mkasner/drops/protocol"
	"github.com/mkasner/drops/router"
	"github.com/mkasner/drops/session"
	"github.com/mkasner/drops/store"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	sessionId string
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
		// session.DeleteSession(c.sessionId)
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		var m Message
		if err := json.Unmarshal(message, &m); err != nil {
			panic(err)
		}
		fmt.Printf("Message %+v\n", m)
		fmt.Printf("Session on connection %+v\n", c.sessionId)
		// fmt.Printf("Message received: %s\n", string(message))
		//Check to see if its url change
		switch m.Type {
		case "GET":

			var params router.Params
			if _, ok := m.Data["params"]; ok { //if params exist
				paramsMap := m.Data["params"].(map[string]interface{})
				if paramsMap != nil && len(paramsMap) > 0 {
					params = make(router.Params, len(paramsMap))
					i := 0
					for k, v := range paramsMap {
						param := &router.Param{k, v}
						params[i] = *param
						i++
					}
				}
			}
			// fmt.Printf("Params: %+v", m.Data["params"])
			// dom := Route(m.Data["route"].(string)[1:], params)

			handle, paramsFromRequest, _ := rtr.Lookup(m.Type, m.Data["route"].(string))
			message, err := c.handleExecute(handle, paramsFromRequest)
			if err != nil {
				log.Printf("Error handling: %v\n", err)
			}
			// fmt.Printf("Patches generated: %+v\n", string(message))
			c.send <- message
		case "EVENT":
			fmt.Printf("Event received: %+v\n", m.Data)

			handle, paramsFromRequest, _ := rtr.Lookup(m.Type, m.Data["route"].(string))
			if handle != nil {
				newParam := router.Param{Key: "data", Value: m.Data}

				paramsFromRequest = append(paramsFromRequest, newParam)
				message, err := c.handleExecute(handle, paramsFromRequest)
				if err != nil {
					log.Printf("Error handling: %v\n", err)
				}
				// fmt.Printf("Patches generated: %+v\n", string(message))
				c.send <- message
			} else {
				if val, ok := m.Data["action"]; ok {
					// if handler, ok := eventDispatcher[val.(string)]; ok {
					// 	message := handler.Handle(m.Data)
					// 	c.send <- message
					// 	continue
					// }
					switch val {
					case "new":

						model := m.Data["model"].(map[string]interface{})
						modelname := m.Data["model-name"].(string)
						store.AddModel(modelname, model)
					case "edit":
						model := m.Data["model"].(map[string]interface{})
						modelname := m.Data["model-name"].(string)
						store.SaveModel(modelname, model)

					case "delete":
						model := m.Data["model"].(map[string]interface{})
						modelname := m.Data["model-name"].(string)
						store.DeleteModel(modelname, model)
						// //Start deployment
						// patch := &Patch{Element: "#alert", Payload: "Do you want  to undo delete of" + model["id"].(string) + "? Not yet..."}
						// message, err := json.Marshal(patch)
						// if err != nil {
						// 	log.Println("Error marshaling patch")
						// }
						c.send <- message
					}
				}
			}

		default:
			h.broadcast <- message
		}
	}
	// var message Message
	// for {
	// 	err := c.ws.ReadJSON(message)
	// 	if err != nil {
	// 		break
	// 	}
	// 	fmt.Printf("Message received: %v\n", message)
	// 	//Check to see if its url change

	// 	// if err := c.write(websocket.TextMessage, []byte(messageMap["data"])); err != nil {
	// 	// 	break
	// 	// }

	// 	// h.broadcast <- []byte(message)
	// }
}

//Handles event and creates DOM patch
func (c *connection) handleExecute(handle router.Handle, paramsFromRequest router.Params) ([]byte, error) {
	dropsResponse := &protocol.DropsResponse{}
	activeDOM := session.GetSessionActiveDOM(c.sessionId)
	// fmt.Printf("\nGetting activeDOM on websocket connection: %+v\n", pretty.Formatter(activeDOM))
	if handle != nil {
		sessionParam := router.Param{Key: "session", Value: c.sessionId}

		paramsFromRequest = append(paramsFromRequest, sessionParam)
		// fmt.Printf("Routing success: %v\n", paramsFromRequest)

		dropsResponse = handle(nil, nil, paramsFromRequest)
		dropsResponse.ActiveDom = activeDOM
		// PrintDOM(ActiveDOM, "1")

		// fmt.Printf("ActiveDOM: %+v\n", activeDOM)
		// fmt.Printf("New DOM: %+v\n", dom)
		// fmt.Printf("Active dom is the same to new DOM: %v\n", *ActiveDOM == *dom)
		// PrintDOM(dom, "2")
	} else {
		log.Println("Routing failure, no handler")

		dropsResponse.Dom = activeDOM
	}
	var message []byte
	// var err error
	session.SetSessionActiveDOM(c.sessionId, dropsResponse.Dom)

	message = protocol.GenerateMessage(dropsResponse)

	fmt.Printf("Drops response: %+v\n", string(message))
	return message, nil

}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// ServerWs handles webocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request) {
	go h.run() //Starting hub, this is different than in example
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	cookie, err := r.Cookie("session")
	if err != nil {
		// fmt.Printf("Cookie doesn't exist: %v\n", err)
	} else {
		// fmt.Printf("\nQuerying session for id: %s\n", cookie.Value)
		if session.SessionExist(cookie.Value) {
			// fmt.Printf("Session set on websocket connection: %s\n", cookie.Value)
			c.sessionId = cookie.Value
		} else {
			// fmt.Printf("Session doesn't exist: %s\n", cookie.Value)
			// fmt.Printf("Session store: %+v\n", session.SessionStore())
		}
	}
	h.register <- c

	// log.Printf("Client connected: %+v\n", c)
	go c.writePump()
	c.readPump()
	// fmt.Fprintf(w, "connection success")
}

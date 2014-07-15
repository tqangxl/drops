// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package drops

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mkasner/drops/router"
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
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
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
			message, err := HandleExecute(handle, paramsFromRequest)
			if err != nil {
				log.Printf("Error handling: %v\n", err)
			}
			fmt.Printf("Patches generated: %+v\n", string(message))
			c.send <- message
		case "EVENT":
			fmt.Printf("Event received: %+v\n", m.Data)
			// example event. Using this before I Implement eventDispatcher

			if val, ok := m.Data["className"]; ok {
				if strings.Contains(val.(string), "deploy-btn") {
					//Start deployment
					patch := &Patch{Element: "#alert", Payload: "Started deployment " + m.Data["id"].(string)}
					message, err := json.Marshal(patch)
					if err != nil {
						log.Println("Error marshaling patch")
					}
					h.broadcast <- message

				}
			}

			handle, paramsFromRequest, _ := rtr.Lookup(m.Type, m.Data["route"].(string))
			if handle != nil {
				newParam := router.Param{Key: "data", Value: m.Data}

				paramsFromRequest = append(paramsFromRequest, newParam)
				message, err := HandleExecute(handle, paramsFromRequest)
				if err != nil {
					log.Printf("Error handling: %v\n", err)
				}
				fmt.Printf("Patches generated: %+v\n", string(message))
				c.send <- message
			} else {
				if val, ok := m.Data["action"]; ok {
					if handler, ok := eventDispatcher[val.(string)]; ok {
						message := handler.Handle(m.Data)
						c.send <- message
						continue
					}
					switch val {
					case "new":

						model := m.Data["model"].(map[string]interface{})
						modelname := m.Data["model-name"].(string)
						AddModel(modelname, model)
					case "edit":
						model := m.Data["model"].(map[string]interface{})
						modelname := m.Data["model-name"].(string)
						SaveModel(modelname, model)

					case "delete":
						model := m.Data["model"].(map[string]interface{})
						modelname := m.Data["model-name"].(string)
						DeleteModel(modelname, model)
						//Start deployment
						patch := &Patch{Element: "#alert", Payload: "Do you want  to undo delete of" + model["id"].(string) + "? Not yet..."}
						message, err := json.Marshal(patch)
						if err != nil {
							log.Println("Error marshaling patch")
						}
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
	h.register <- c
	log.Printf("Client connected: %+v\n", c)
	go c.writePump()
	c.readPump()
	// fmt.Fprintf(w, "connection success")
}

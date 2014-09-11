//Package that handles front end messages in drops application
package message

import (
	"github.com/mkasner/drops/element"
	"github.com/mkasner/drops/session"
)

//remove messages with expired counter
func cleanupMessages(messages element.Messages) element.Messages {
	cleanMessages := make(element.Messages, 0)
	for _, m := range messages {
		if m.Expires > 0 {
			cleanMessages = append(cleanMessages, m)
		}
	}
	return cleanMessages
}

func AddMessages(sessionId string, messages element.Messages) {
	activeMessages := getMessageFromSession(sessionId)
	activeMessages = append(activeMessages, messages...)

	//set back session messages
	session.SetSessionMessages(sessionId, activeMessages)
}

func NewMessage(sessionId string, text string) {
	msg := &element.Message{
		Text:     text,
		Expires:  1,
		Priority: 1,
	}
	AddMessages(sessionId, element.Messages{msg})
}

//Reduces expire after showing up
func reduceExpire(messages element.Messages) {
	for _, m := range messages {
		m.Expires = m.Expires - 1
	}
}

func getMessageFromSession(sessionId string) element.Messages {
	return session.GetSessionMessages(sessionId)
}

//Creates message views from valid messages
func ProcessMessages(sessionId string, dom *element.DOM) *element.DOM {
	messages := getMessageFromSession(sessionId)
	messages = cleanupMessages(messages)

	for _, m := range messages {
		view := newMessageView(m)
		dom = element.Add(dom, view)
	}
	reduceExpire(messages)
	session.SetSessionMessages(sessionId, messages)

	return dom
}

//Constructs a new message and place it on screen
func newMessageView(message *element.Message) *element.View {
	view := &element.View{Template: "message.tpl", InjectInto: "#message", Model: &element.Model{MAP: make(map[string]interface{}), TYPE: message}}
	return view
}

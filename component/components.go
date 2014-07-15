package component

import "github.com/mkasner/drops/element"

//Creates input text component
func NewInputText(model map[string]interface{}) *element.View {
	view := &element.View{Template: "inputText.tpl", Model: &element.Model{MAP: make(map[string]interface{})}}
	view.Model.MAP = model
	return view
}

//Create select component
func NewSelect(model map[string]interface{}) *element.View {
	view := &element.View{Template: "select.tpl", Model: &element.Model{MAP: make(map[string]interface{})}}
	view.Model.MAP = model

	return view
}

func NewFieldset() *element.View {
	view := &element.View{Template: "fieldset.tpl", Return: "Fieldset", Model: &element.Model{MAP: make(map[string]interface{})}}
	return view
}

type Message struct {
	*element.View
	Text string
	Type string
}

//Constructs a new message and place it on screen
func NewMessage(messageText string, messageType string) *Message {
	message := &Message{Text: messageText, Type: messageType}
	message.View = &element.View{Template: "message.tpl", InjectInto: "#message", Model: &element.Model{MAP: make(map[string]interface{}), TYPE: message}}
	return message
}

func NewMessagePanel(injectInto string) *element.View {
	view := &element.View{Template: "messagePanel.tpl", InjectInto: injectInto, Return: "MessagePanel", Provides: "#message", Model: &element.Model{MAP: make(map[string]interface{})}}
	return view
}

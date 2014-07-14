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

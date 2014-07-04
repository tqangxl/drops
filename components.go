package drops

//Creates input text component
func NewInputText(model map[string]interface{}) *View {
  view := &View{Template: "inputText.tpl", Model: model}
  return view
}

//Create select component
func NewSelect(model map[string]interface{}) *View {
  view := &View{Template: "select.tpl", Model: model}
  return view
}

func NewFieldset() *View {
  view := &View{Template: "fieldset.tpl", Return: "Fieldset"}
  return view
}

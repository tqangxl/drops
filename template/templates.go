//Templates for drops actions stored as string const
//It would be better to store them as html files but I have to find
//a way to include those html files into project that imports drops library
//For now I will be using this method: const string as default
//Later, I'll implement some kind of overriding
package template

import (
	tpl "text/template"
)

const newTplString = `<h1>{{.MAP.Title}} NEW</h1>
<form id="form-new" data-model="{{.MAP.Title}}" data-route="/model/new/">
	<button type="submit" class="small">Save</button>
	{{.MAP.Fieldset}}
	<button type="submit"  class="small">Save</button>
</form>`

const editTplString = `<h1>{{.MAP.Title}} EDIT</h1>
<dl class="sub-nav">
  <dd><a href="#" data-action="delete" data-id="{{.MAP.Id}}" data-model="{{.MAP.Title}}"  class="delete-action">Delete</a></dd>
</dl>
<form id="form-edit" data-model="{{.MAP.Title}}" data-id="{{.MAP.Id}}" data-route="/model/save/">
	<button type="submit" class="small">Save</button>
{{.MAP.Fieldset}}
	<button type="submit"  class="small">Save</button>
</form>
{{.MAP.Content}}`

const viewTplString = `<h1>{{.MAP.Title}} {{.MAP.Id}}</h1>
<dl class="sub-nav">
  <dd><a href="#" data-action="edit" data-id="{{.MAP.Id}}" class="action">Edit</a></dd>
  <dd><a href="#" data-route="/{{.MAP.Title}}/{{.MAP.Id}}/" class="delete-action">Delete</a></dd>
</dl>
	{{range .MAP.Fields}}
	<div class="row">
		<div class="large-2 column">
			<span>{{.Label}}</span>
		</div>
		<div class="large-4 column end">
			<span>{{.Value}}</span>
		</div>
	</div>
	{{end}}`

//Select tag
const selectTplString = `
<label>{{.MAP.Label}}
<select name="{{.MAP.Name}}">
	{{range .MAP.Options}}
	<option value="{{.Id}}" {{if .Selected}}{{.Selected}}{{end}}">{{.Label}}</option>
	{{end}}
</select>
</label>
`

//Input tag
const inputTextTplString = `
<label>{{.MAP.Label}}<input type="text" name="{{.MAP.Name}}" value="{{if .MAP.Value}}{{.MAP.Value}}{{end}}"/></label>
`

var newTpl, editTpl, viewTpl, selectTpl *tpl.Template

//Template registry
var Templates map[string]*tpl.Template

//Load templates
func init() {
	// Templates = &tst.Trie{}
	var templ *tpl.Template
	Templates = make(map[string]*tpl.Template)
	newTpl = tpl.Must(tpl.New("new.tpl").Parse(newTplString))
	Templates[newTpl.Name()] = newTpl
	editTpl = tpl.Must(tpl.New("edit.tpl").Parse(editTplString))
	Templates[editTpl.Name()] = editTpl
	viewTpl = tpl.Must(tpl.New("view.tpl").Parse(viewTplString))
	Templates[viewTpl.Name()] = viewTpl
	templ = tpl.Must(tpl.New("select.tpl").Parse(selectTplString))
	Templates[templ.Name()] = templ
	templ = tpl.Must(tpl.New("inputText.tpl").Parse(inputTextTplString))
	Templates[templ.Name()] = templ
	templ = tpl.Must(tpl.New("fieldset.tpl").Parse("<fieldset>{{.MAP.Content}}</fieldset>"))
	Templates[templ.Name()] = templ
	templ = tpl.Must(tpl.New("message.tpl").Parse(`<li class="{{.TYPE.Type}}">{{.TYPE.Text}}</li>`))
	Templates[templ.Name()] = templ
	templ = tpl.Must(tpl.New("messagePanel.tpl").Parse(`<ul id="message" class="panel no-bullet">{{.MAP.Content}}</ul>`))
	Templates[templ.Name()] = templ

}

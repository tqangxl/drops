//Templates for drops actions stored as string const
//It would be better to store them as html files but I have to find
//a way to include those html files into project that imports drops library
//For now I will be using this method: const string as default
//Later, I'll implement some kind of overriding
package drops

import (
	"text/template"
)

const newTplString = `<h1>{{.Title}} NEW</h1>
<form id="form-new" data-model="{{.Title}}">
	<button type="submit" class="small">Save</button>
	{{range .Fields}}
		<label>{{.Label}}<input type="text" name="{{.Name}}"/></label>
	{{end}}
	<button type="submit"  class="small">Save</button>
</form>`

const editTplString = `<h1>{{.Title}} EDIT</h1>
<dl class="sub-nav">
  <dd><a href="#" data-action="delete" data-id="{{.Id}}" class="action">Delete</a></dd>
</dl>
<form id="form-new" data-model="{{.Title}}">
	<button type="submit" class="small">Save</button>
	{{range .Fields}}
		<label>{{.Label}}<input type="text" name="{{.Name}}" value="{{.Value}}"/></label>
	{{end}}
	<button type="submit"  class="small">Save</button>
</form>`

const viewTplString = `<h1>{{.Title}} {{.Id}}</h1>
<dl class="sub-nav">
  <dd><a href="#" data-action="edit" data-id="{{.Id}}" class="action">Edit</a></dd>
  <dd><a href="#" data-action="delete" data-id="{{.Id}}" class="action">Delete</a></dd>
</dl>
	{{range .Fields}}
	<div class="row">
		<div class="large-2 column">
			<span>{{.Label}}</span> 
		</div>
		<div class="large-4 column end">
			<span>{{.Value}}</span>
		</div>
	</div>
	{{end}}
</form>`

var newTpl, editTpl, viewTpl *template.Template

//Template registry
var Templates map[string]*template.Template

//Load templates
func init() {
	// Templates = &tst.Trie{}
	Templates = make(map[string]*template.Template)
	newTpl = template.Must(template.New("new.tpl").Parse(newTplString))
	Templates[newTpl.Name()] = newTpl
	editTpl = template.Must(template.New("edit.tpl").Parse(editTplString))
	Templates[editTpl.Name()] = editTpl
	viewTpl = template.Must(template.New("view.tpl").Parse(viewTplString))
	Templates[viewTpl.Name()] = viewTpl

}

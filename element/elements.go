package element

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/kr/pretty"

	"github.com/mkasner/drops/template"
)

type DOM struct {
	View
	IdMap map[string]*View
	Id    string
}
type Head struct {
	View
}
type Body struct {
	View
}

type Model struct {
	MAP  map[string]interface{}
	TYPE interface{}
}
type View struct {
	Children   []*View
	Template   string
	Content    string
	InjectInto string //ID to inject this view into
	Provides   string //ID that this view provides, into which we can inject other views
	*Model
	Return string //key which view returns and will be used in model for templating
	Parent *View
}

type Messages []*Message

type Message struct {
	Text     string
	Type     string //class
	Expires  int    //counter for how long the messages should be shown
	Priority int    //Priority of messages
}

//Renders view recurively
func Render(v *View) bytes.Buffer {
	var buffer bytes.Buffer
	// Iterate through list and and print its contents.
	// l := list.New()
	// fmt.Printf("Children: %+v\n", v.Children)
	if v.Children != nil {
		// fmt.Println("rendering children...")
		for _, view := range v.Children {
			// fmt.Println("rendering child...")

			// fmt.Printf("View: %+v\n", view)
			buff := Render(view)
			// fmt.Println(rendered)
			// v.Content = buff.String()
			if v.Model == nil {
				v.Model = &Model{
					MAP: make(map[string]interface{}),
				}
			}
			if v.Model.MAP == nil {
				v.Model.MAP = make(map[string]interface{})
			}
			if view.Return == "" {
				buffer.WriteString(buff.String())
			} else {
				v.MAP[view.Return] = buff.String()
				// fmt.Printf("Return : %+v\n", buff.String())
			}
			// fmt.Printf("Buffer: %+v\n", buffer)
		}
	} else {
		// fmt.Printf("No children: %+v\n", v.Template)
	}
	// fmt.Printf("Template: %s\n", v.Template)
	if strings.HasSuffix(v.Template, ".html") || strings.HasSuffix(v.Template, ".tpl") {
		if v.Model == nil {
			v.Model = &Model{MAP: make(map[string]interface{})}
		}
		if v.Model.MAP == nil {
			v.Model.MAP = make(map[string]interface{})
		}
		v.Content = buffer.String()

		v.Model.MAP["Content"] = v.Content

		buffer = RenderTemplate(v.Template, v.Model)
		// fmt.Printf("Template: %s, Rendered: %s\n", v.Template, buffer.String())
		// }
	} else {
		// fmt.Printf("Rendering just string: %+v\n", v.Template)
		buffer.WriteString(v.Template)
	}
	return buffer
}

func RenderTemplate(templateName string, data *Model) bytes.Buffer {
	var buffer bytes.Buffer
	// t, _ := template.ParseFiles(templateName)
	// fmt.Printf("Rendering template: %+v\n", templateName)
	if t, ok := template.Templates[templateName]; ok {
		t.Execute(&buffer, data)
	}
	// fmt.Printf("\n%+v\n", buffer.String())
	return buffer
}

//Adds element to dom structure to parent specified in InjectTo field
func Add(dom *DOM, view *View) *DOM {
	parent := &dom.View
	// fmt.Printf("Provides: %s\n", view.Provides)
	// fmt.Printf("Returns: %s\n", view.Return)
	// fmt.Printf("InjectInto: %s\n", view.InjectInto)

	if view.Provides != "" {
		dom.IdMap[view.Provides] = view
	}
	//if view.InjectInto != "" && view.Return == "" {
	if view.InjectInto != "" {
		if node, ok := dom.IdMap[view.InjectInto]; ok {
			// fmt.Printf("Injecting into Node: %+v\n", node)
			// if parentView.Children == nil {
			// 	parentView.Children = list.New()
			// }
			parent = node
			// fmt.Printf("ParentView children before: %+v\n", parent)

		}

	}
	view.Parent = parent
	parent.Children = append(parent.Children, view)
	// fmt.Printf("ParentView children after: %+v\n", parent)
	// fmt.Println("Added... ")
	// fmt.Printf("View: %+v\n", view)
	// fmt.Printf("Dom: %+v\n", dom)
	return dom
}

//Replaces element at selected injectointo point, old element gets removed
func Replace(dom *DOM, view *View) *DOM {
	parent := &dom.View
	// fmt.Printf("Provides: %s\n", view.Provides)
	// fmt.Printf("Returns: %s\n", view.Return)
	// fmt.Printf("InjectInto: %s\n", view.InjectInto)

	if view.Provides != "" {
		dom.IdMap[view.Provides] = view
	}
	//if view.InjectInto != "" && view.Return == "" {
	if view.InjectInto != "" {
		if node, ok := dom.IdMap[view.InjectInto]; ok {
			// fmt.Printf("Injecting into Node: %+v\n", node)
			// if parentView.Children == nil {
			// 	parentView.Children = list.New()
			// }
			parent = node
			// fmt.Printf("ParentView children before: %+v\n", parent)

		}

	}
	view.Parent = parent
	parent.Children = []*View{view}
	// fmt.Printf("ParentView children after: %+v\n", parent)
	// fmt.Println("Added... ")
	// fmt.Printf("View: %+v\n", view)
	// fmt.Printf("Dom: %+v\n", dom)
	return dom
}

//Adds view to existing view
//It-s used when we don't have DOM existent
func AddToView(view *View, newView *View) *View {
	// fmt.Printf("\n\n\nBefore AddToView length: %+v\n", len(view.Children))
	if newView != nil {

		parent := view.Children

		parent = append(parent, newView)
		view.Children = parent
	}
	// fmt.Printf("After AddToView: %+v\n", len(view.Children))
	return view
}

func PrintDOM(dom *DOM, tag string) {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "DOM: %+v\n", dom.Id)
	fmt.Fprintln(&buffer, "View:")
	PrintView(&buffer, &dom.View)

	fmt.Fprintln(&buffer, "\n")
	fmt.Fprintln(&buffer, "IdTree:")
	PrintTree(&buffer, dom.IdMap)
	filename := "/tmp/dom" + tag + ".txt"
	err := ioutil.WriteFile(filename, buffer.Bytes(), 0644)
	if err != nil {
		log.Println(err)
	}
}

func PrintView(buffer *bytes.Buffer, view *View) {
	fmt.Fprintf(buffer, "View: %+v\n", pretty.Formatter(view))
	fmt.Fprintf(buffer, "Model: %+v\n", *view.Model)
	if view.Children != nil {
		// fmt.Println("rendering children...")
		for _, view := range view.Children {
			// fmt.Println("rendering child...")

			// fmt.Printf("View: %+v\n", view)
			PrintView(buffer, view)

		}
	}
}

func PrintTree(buffer *bytes.Buffer, idMap map[string]*View) {
	fmt.Fprintf(buffer, "%+v\n", idMap)
}

//Makes a snapshot of provided DOM, and enables us to make new views on it, and compare it to old one
func CopyDom(dom DOM) *DOM {
	newDOM := dom
	// newView := copyView(*dom.View)
	// newTrie := copyIdTrie(*dom.IdTree)
	// newDOM.View = dom.View
	// newDOM.IdTree = copyIdTrie(dom.IdTree)
	newDOM.Id = "2"
	// copyChildren]
	idMap := make(map[string]*View)
	newView := copyView(dom.View, idMap)
	newDOM.View = newView
	newDOM.IdMap = idMap

	// fmt.Printf("New DOM copied: %+v\n", newDOM)
	// fmt.Printf("ActiveDOM old: %+v\n", dom)
	return &newDOM
}

func copyView(v View, idMap map[string]*View) View {
	if v.Children != nil {
		newChildren := make([]*View, len(v.Children))
		// copy(newChildren, v.Children)
		// fmt.Printf("New children length: %v  Old children length: %v\n", len(newChildren), len(v.Children))

		// fmt.Println("rendering children...")
		for i, view := range v.Children {
			// fmt.Println("rendering child...")

			// fmt.Printf("View: %+v\n", view)
			childView := copyView(*view, idMap)
			childView.Model = copyModel(*view.Model)
			newChildren[i] = &childView
			if childView.Provides != "" {
				// fmt.Printf("Added to new idMap: %+v\n", v.Provides)
				idMap[childView.Provides] = &childView
			}
			v.Children = newChildren

		}
	} else {
		v.Children = nil
		if v.Provides != "" {
			// fmt.Printf("Added to new idMap: %+v\n", v.Provides)
			idMap[v.Provides] = &v
		}
	}

	return v
}

func copyModel(m Model) *Model {
	return &m
}

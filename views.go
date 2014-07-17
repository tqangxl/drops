package drops

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/mkasner/drops/component"
	"github.com/mkasner/drops/element"
)

//Struct which holds facts for fields that are going to be rendered
type ComponentRule struct {
	Ignore bool
	View   *element.View
	Label  string
}

//Renders view recurively
func Render(v *element.View) bytes.Buffer {
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
				v.Model = &element.Model{
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
			v.Model = &element.Model{MAP: make(map[string]interface{})}
		}
		if v.Model.MAP == nil {
			v.Model.MAP = make(map[string]interface{})
		}
		v.Content = buffer.String()

		v.Model.MAP["Content"] = v.Content
		// fmt.Printf("Model for template: %s \n%+v\n", v.Template, v.Model)

		// fmt.Printf("Model struct: %v\n", v.ModelStruct)
		// if v.ModelStruct != nil {
		// 	buffer = RenderTemplate(v.Template, v.ModelStruct)
		// 	//Find which ids this view provides and generate another template onto which content should be injected
		// 	if v.Provides != "" {
		// 		templateBeforeProcess := buffer.String()
		// 		idIndex := strings.Index(templateBeforeProcess, v.Provides[1:])
		// 		if idIndex > -1 {
		// 			//if id exist find out for index of character '>'
		// 			tagEndIndex := strings.IndexRune(templateBeforeProcess[idIndex:], '>')
		// 			if tagEndIndex > -1 {
		// 				//If i find tag end index
		// 				fmt.Printf("idIndex: %v  tagEndIndex: %v \n", idIndex, tagEndIndex)
		// 				var newBuffer bytes.Buffer
		// 				newBuffer.WriteString(templateBeforeProcess[:idIndex+tagEndIndex+1])
		// 				newBuffer.WriteString(v.Content)
		// 				newBuffer.WriteString(templateBeforeProcess[idIndex+tagEndIndex+2:])
		// 				buffer = newBuffer
		// 			}
		// 		}
		// 	}
		// } else {

		buffer = RenderTemplate(v.Template, v.Model)
		// fmt.Printf("Template: %s, Rendered: %s\n", v.Template, buffer.String())
		// }
	} else {
		// fmt.Printf("Rendering just string: %+v\n", v.Template)
		buffer.WriteString(v.Template)
	}
	return buffer
}

func RenderTemplate(templateName string, data *element.Model) bytes.Buffer {
	var buffer bytes.Buffer
	// t, _ := template.ParseFiles(templateName)
	// fmt.Printf("Rendering template: %+v\n", templateName)
	if t, ok := Templates[templateName]; ok {
		t.Execute(&buffer, data)
	}
	// fmt.Printf("\n%+v\n", buffer.String())
	return buffer
}

//Adds element to dom structure to parent specified in InjectTo field
func Add(dom *element.DOM, view *element.View) *element.DOM {
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
func Replace(dom *element.DOM, view *element.View) *element.DOM {
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
	parent.Children = []*element.View{view}
	// fmt.Printf("ParentView children after: %+v\n", parent)
	// fmt.Println("Added... ")
	// fmt.Printf("View: %+v\n", view)
	// fmt.Printf("Dom: %+v\n", dom)
	return dom
}

//Adds view to existing view
//It-s used when we don't have DOM existent
func AddToView(view *element.View, newView *element.View) *element.View {
	// fmt.Printf("\n\n\nBefore AddToView length: %+v\n", len(view.Children))
	if newView != nil {

		parent := view.Children

		parent = append(parent, newView)
		view.Children = parent
	}
	// fmt.Printf("After AddToView: %+v\n", len(view.Children))
	return view
}

//Component for creating New objects of a model, based on a type
//Rules: map with rules for certain field
//It is used for not polluting struct definition, you can set it after struct initialization and pass it to New handler
//Use it like this:
//Key: field name (retreived by reflection)
//Value: Struct tag like: name:"value" name:"value"
//Example: form ignore
func New(typ reflect.Type, rules map[string]string, injectInto string, title string) *element.View {
	view := &element.View{Template: "new.tpl", InjectInto: injectInto, Model: &element.Model{MAP: make(map[string]interface{})}}
	view.Model.MAP["Title"] = title
	var fieldset *element.View
	fieldset = component.NewFieldset()

	fields := make([]map[string]interface{}, 0)
	fmt.Printf("Number of fields %+v\n", typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		fieldName := typ.Field(i).Name
		fmt.Printf("Field name %v\n", fieldName)

		fieldMap := make(map[string]interface{}) //map passed to field generator
		fieldMap["Label"] = fieldName
		fieldMap["Name"] = fieldName
		// fmt.Printf("Field map %v\n", fieldMap)
		fields = append(fields, fieldMap)
		fieldRule := &ComponentRule{}
		var ruleTag string
		if rules != nil {
			//Process rules

			if tag, ok := rules[fieldName]; ok {
				fmt.Printf("Rule tag %v\n", ruleTag)
				ruleTag = tag

			}
		}
		//Execute rules
		fieldRule = execIgnoreRule(ruleTag, fieldRule)
		// fmt.Printf("Rule at the moment %v\n", fieldRule)
		fieldRule = execForeignRule(ruleTag, fieldMap, fieldRule)
		// fmt.Printf("Rule at the moment %v\n", fieldRule)
		fieldset = AddToView(fieldset, fieldRule.View)
	}

	view.Model.MAP["Fields"] = fields
	view = AddToView(view, fieldset)
	fmt.Printf("Fields generated %v\n", fields)

	return view

}

//Executes ignore rule and returns if filed should be ignored
func execIgnoreRule(ruleTag string, fieldRule *ComponentRule) *ComponentRule {
	ignoreRule := GetRule(ruleTag, "ignore")
	// fmt.Printf("Ignore rule %v\n", ignoreRule)
	var ignore bool
	var err error
	if ignoreRule != "" {
		ignore, err = strconv.ParseBool(ignoreRule)
		if err != nil {
			log.Println("Ignore rule not properly set")
		}
		// fmt.Printf("Ignore parsed %v\n", ignore)
		fieldRule.Ignore = ignore
	}

	return fieldRule
}

//executes foreign rule, to create a special field that creates
//reference to foreign key table and populates
// select element with options
//model - model name for which  I must retreive object store
func execForeignRule(ruleTag string, fieldMap map[string]interface{}, fieldRule *ComponentRule) *ComponentRule {
	model := GetRule(ruleTag, "foreign")
	// var field *element.View
	// fmt.Printf("Ignoring on creating view %v\n", fieldRule.Ignore)
	if fieldRule.Ignore {
		return fieldRule
	}
	if model != "" {
		//Changing label to be the name of model if not explicitely set
		if fieldRule.Label == "" {
			fieldMap["Label"] = model
		}
		// fmt.Printf("Selecting for model %v\n", model)
		allModel := GetAll(model)
		for _, v := range allModel {
			// fmt.Printf("Model value %v\n", v)
			if _, ok := v["Label"]; !ok {
				v["Label"] = v["Name"] //example usage, if label not set
			}
			//Add selected attribute if it matches value of the field

			if foreignValue, ok := fieldMap["Value"]; ok {
				if foreignValue == v["Id"] {
					v["Selected"] = "selected"
				} else {
					v["Selected"] = ""
				}
			}
		}
		fieldMap["Options"] = allModel

		fieldRule.View = component.NewSelect(fieldMap)
	} else {
		fieldRule.View = component.NewInputText(fieldMap)
	}
	return fieldRule
}

func ViewModel(value interface{}, rules map[string]string, injectInto string, title string, template string) *element.View {
	chosenTemplate := "view.tpl"
	if template != "" {
		chosenTemplate = template //override chosen template
	}
	view := &element.View{Template: chosenTemplate, InjectInto: injectInto, Model: &element.Model{MAP: make(map[string]interface{})}}
	view.Model.MAP["Title"] = title
	typ := reflect.TypeOf(value).Elem()
	val := reflect.Indirect(reflect.ValueOf(value))

	fields := make([]*map[string]interface{}, 0)
	for i := 0; i < typ.NumField(); i++ {
		fieldName := typ.Field(i).Name
		// fmt.Printf("Field name %v\n", fieldName)

		fieldMap := make(map[string]interface{}) //map passed to field generator
		fieldMap["Label"] = fieldName
		fieldMap["Value"] = val.Field(i).Interface()
		// fmt.Printf("Field map %v\n", fieldMap)
		fields = append(fields, &fieldMap)
		if rules != nil {
			if ruleTag, ok := rules[fieldName]; ok {
				// fmt.Printf("Rule tag %v\n", ruleTag)

				ignoreRule := GetRule(ruleTag, "ignore")
				// fmt.Printf("Ignore rule %v\n", ignoreRule)
				if ignoreRule != "" {
					ignore, err := strconv.ParseBool(ignoreRule)
					if err != nil {
						log.Println("Ignore rule not properly set")
					}
					//If not ignored input will be generated
					if ignore {
						fields = fields[:len(fields)-1]

					}
				}
			}
		}
	}
	view.MAP["Fields"] = fields
	view.MAP["Id"] = val.FieldByName("Id").Interface()
	// fmt.Printf("Fields generated %v\n", fields)

	return view

}

//Magical method that creates form for editing based on provided value
func Edit(value interface{}, rules map[string]string, injectInto string, title string) *element.View {
	view := &element.View{Template: "edit.tpl", InjectInto: injectInto, Model: &element.Model{MAP: make(map[string]interface{})}}
	view.MAP["Title"] = title
	fieldset := component.NewFieldset()

	typ := reflect.TypeOf(value).Elem()
	val := reflect.Indirect(reflect.ValueOf(value))

	// fmt.Printf("Converted typ %v\n", typ)

	fields := make([]*map[string]interface{}, 0)
	for i := 0; i < typ.NumField(); i++ {
		fieldName := typ.Field(i).Name
		// fmt.Printf("Field name %v\n", fieldName)

		fieldMap := make(map[string]interface{}) //map passed to field generator
		fieldMap["Label"] = fieldName
		fieldMap["Name"] = fieldName
		fieldMap["Value"] = val.Field(i).Interface()
		// fmt.Printf("Field map %v\n", fieldMap)
		fields = append(fields, &fieldMap)
		fieldRule := &ComponentRule{}
		var ruleTag string
		if rules != nil {
			//Process rules

			if tag, ok := rules[fieldName]; ok {
				fmt.Printf("Rule tag %v\n", ruleTag)
				ruleTag = tag

			}
		}
		//Execute rules
		fieldRule = execIgnoreRule(ruleTag, fieldRule)
		// fmt.Printf("Rule at the moment %v\n", fieldRule)
		fieldRule = execForeignRule(ruleTag, fieldMap, fieldRule)
		// fmt.Printf("Rule at the moment %v\n", fieldRule)
		fieldset = AddToView(fieldset, fieldRule.View)
	}
	view.MAP["Fields"] = fields
	view = AddToView(view, fieldset)
	view.MAP["Id"] = val.FieldByName("Id").Interface()
	// fmt.Printf("Fields generated %v\n", fields)

	return view

}

func PrintDOM(dom *element.DOM, tag string) {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "DOM: %+v\n", dom.Id)
	fmt.Println(&buffer, "View: %v\n")
	PrintView(&buffer, &dom.View)

	fmt.Fprintln(&buffer, "\n")
	fmt.Fprintln(&buffer, "IdTree: %v\n")
	PrintTree(&buffer, dom.IdMap)
	filename := "/tmp/dom" + tag + ".txt"
	err := ioutil.WriteFile(filename, buffer.Bytes(), 0644)
	if err != nil {
		log.Println(err)
	}
}

func PrintView(buffer *bytes.Buffer, view *element.View) {
	fmt.Fprintf(buffer, "View: %+v\n", view)
	if view.Children != nil {
		// fmt.Println("rendering children...")
		for _, view := range view.Children {
			// fmt.Println("rendering child...")

			// fmt.Printf("View: %+v\n", view)
			PrintView(buffer, view)

		}
	}
}

func PrintTree(buffer *bytes.Buffer, idMap map[string]*element.View) {
	fmt.Fprintf(buffer, "%+v\n", idMap)
}

//Makes a snapshot of provided DOM, and enables us to make new views on it, and compare it to old one
func CopyDom(dom element.DOM) *element.DOM {
	newDOM := dom
	// newView := copyView(*dom.View)
	// newTrie := copyIdTrie(*dom.IdTree)
	// newDOM.View = dom.View
	// newDOM.IdTree = copyIdTrie(dom.IdTree)
	newDOM.Id = "2"
	// copyChildren]
	idMap := make(map[string]*element.View)
	newView := copyView(dom.View, idMap)
	newDOM.View = newView
	newDOM.IdMap = idMap

	// fmt.Printf("New DOM copied: %+v\n", newDOM)
	// fmt.Printf("ActiveDOM old: %+v\n", dom)
	return &newDOM
}

func copyView(v element.View, idMap map[string]*element.View) element.View {

	var newChildren []*element.View
	copy(newChildren, v.Children)
	if newChildren != nil {
		// fmt.Println("rendering children...")
		for i, view := range v.Children {
			// fmt.Println("rendering child...")

			// fmt.Printf("View: %+v\n", view)
			childView := copyView(*view, idMap)
			newChildren[i] = &childView

		}
	} else {
		v.Children = newChildren
	}
	if v.Provides != "" {
		idMap[v.Provides] = &v
	}
	return v
}

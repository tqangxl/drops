package drops

import (
	"bytes"
	"container/list"
	"fmt"
	"github.com/cyfdecyf/tst"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type View struct {
	Children   *list.List
	Template   string
	Content    string
	InjectInto string //ID to inject this view into
	Provides   string //ID that this view provides, into which we can inject other views
	Model      map[string]interface{}
	Return     string //key which view returns and will be used in model for templating

}

//Renders view recurively
func (v *View) Render() bytes.Buffer {
	var buffer bytes.Buffer
	// Iterate through list and and print its contents.
	// l := list.New()
	if v.Children != nil {
		// fmt.Println("rendering children...")
		for e := v.Children.Front(); e != nil; e = e.Next() {
			// fmt.Println("rendering child...")

			view := e.Value.(*View)
			// fmt.Printf("View: %+v\n", view)
			buff := view.Render()
			// fmt.Println(rendered)
			// v.Content = buff.String()
			if v.Model == nil {
				v.Model = make(map[string]interface{})
			}
			if view.Return == "" {
				buffer.WriteString(buff.String())
			} else {
				v.Model[view.Return] = buff.String()
			}
			// fmt.Printf("Buffer: %+v\n", buffer)
		}
	}
	// fmt.Printf("Template: %s\n", v.Template)
	if strings.HasSuffix(v.Template, ".html") || strings.HasSuffix(v.Template, ".tpl") {
		if v.Model == nil {
			v.Model = make(map[string]interface{})
		}
		v.Content = buffer.String()
		v.Model["Content"] = v.Content

		// fmt.Printf("Content: %v\n", v.Model)
		buffer = RenderTemplate(v.Template, v.Model)
	} else {
		buffer.WriteString(v.Template)
	}
	return buffer
}

func RenderTemplate(templateName string, data interface{}) bytes.Buffer {
	var buffer bytes.Buffer
	// t, _ := template.ParseFiles(templateName)
	// fmt.Printf("Templates: %+v", Templates)
	if t, ok := Templates[templateName]; ok {
		t.Execute(&buffer, data)
	}
	return buffer
}

//Adds element to dom structure to parent specified in InjectTo field
func Add(dom *DOM, view *View) *DOM {
	parent := dom.Children
	// fmt.Printf("Provides: %s\n", view.Provides)
	// fmt.Printf("InjectInto: %s\n", view.InjectInto)

	if view.Provides != "" {
		dom.IdTree.Put(view.Provides, view)
	}
	if view.InjectInto != "" && view.Return == "" {
		if node := dom.IdTree.Get(view.InjectInto); node != nil {
			// fmt.Printf("Node: %+v\n", node)
			parentView := node.(*View)
			if parentView.Children == nil {
				parentView.Children = list.New()
			}
			parent = parentView.Children
		}

	}

	parent.PushBack(view)
	return dom
}

type DOM struct {
	*View
	IdTree *tst.Trie
}
type Head struct {
	*View
}
type Body struct {
	*View
}

//Dom currently active
var ActiveDOM *DOM

//Component for creating New objects of a model, based on a type
//Rules: map with rules for certain field
//It is used for not polluting struct definition, you can set it after struct initialization and pass it to New handler
//Use it like this:
//Key: field name (retreived by reflection)
//Value: Struct tag like: name:"value" name:"value"
//Example: form ignore
func New(typ reflect.Type, rules map[string]string, injectInto string, title string) *View {
	view := &View{Template: "new.tpl", InjectInto: injectInto, Model: map[string]interface{}{
		"Title": title,
	}}

	fields := make([]*map[string]string, 0)
	for i := 0; i < typ.NumField(); i++ {
		fieldName := typ.Field(i).Name
		fmt.Printf("Field name %v\n", fieldName)

		fieldMap := make(map[string]string) //map passed to field generator
		fieldMap["Label"] = fieldName
		fieldMap["Name"] = fieldName
		// fmt.Printf("Field map %v\n", fieldMap)
		fields = append(fields, &fieldMap)
		if rules != nil {
			if ruleTag, ok := rules[fieldName]; ok {
				fmt.Printf("Rule tag %v\n", ruleTag)

				ignoreRule := GetRule(ruleTag, "ignore")
				fmt.Printf("Ignore rule %v\n", ignoreRule)
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
	view.Model["Fields"] = fields
	fmt.Printf("Fields generated %v\n", fields)

	return view

}

func ViewModel(value interface{}, rules map[string]string, injectInto string, title string, template string) *View {
	chosenTemplate := "view.tpl"
	if template != "" {
		chosenTemplate = template //override chosen template
	}
	view := &View{Template: chosenTemplate, InjectInto: injectInto, Model: map[string]interface{}{
		"Title": title,
	}}

	typ := reflect.TypeOf(value).Elem()
	val := reflect.Indirect(reflect.ValueOf(value))

	fields := make([]*map[string]interface{}, 0)
	for i := 0; i < typ.NumField(); i++ {
		fieldName := typ.Field(i).Name
		fmt.Printf("Field name %v\n", fieldName)

		fieldMap := make(map[string]interface{}) //map passed to field generator
		fieldMap["Label"] = fieldName
		fieldMap["Value"] = val.Field(i).Interface()
		// fmt.Printf("Field map %v\n", fieldMap)
		fields = append(fields, &fieldMap)
		if rules != nil {
			if ruleTag, ok := rules[fieldName]; ok {
				fmt.Printf("Rule tag %v\n", ruleTag)

				ignoreRule := GetRule(ruleTag, "ignore")
				fmt.Printf("Ignore rule %v\n", ignoreRule)
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
	view.Model["Fields"] = fields
	view.Model["Id"] = val.FieldByName("Id").Interface()
	fmt.Printf("Fields generated %v\n", fields)

	return view

}

func Edit(value interface{}, rules map[string]string, injectInto string, title string) *View {
	view := &View{Template: "edit.tpl", InjectInto: injectInto, Model: map[string]interface{}{
		"Title": title,
	}}

	typ := reflect.TypeOf(value).Elem()
	val := reflect.Indirect(reflect.ValueOf(value))

	fmt.Printf("Converted typ %v\n", typ)

	fields := make([]*map[string]interface{}, 0)
	for i := 0; i < typ.NumField(); i++ {
		fieldName := typ.Field(i).Name
		fmt.Printf("Field name %v\n", fieldName)

		fieldMap := make(map[string]interface{}) //map passed to field generator
		fieldMap["Label"] = fieldName
		fieldMap["Name"] = fieldName
		fieldMap["Value"] = val.Field(i).Interface()
		// fmt.Printf("Field map %v\n", fieldMap)
		fields = append(fields, &fieldMap)
		if rules != nil {
			if ruleTag, ok := rules[fieldName]; ok {
				fmt.Printf("Rule tag %v\n", ruleTag)

				ignoreRule := GetRule(ruleTag, "ignore")
				fmt.Printf("Ignore rule %v\n", ignoreRule)
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
	view.Model["Fields"] = fields
	view.Model["Id"] = val.FieldByName("Id").Interface()
	fmt.Printf("Fields generated %v\n", fields)

	return view

}

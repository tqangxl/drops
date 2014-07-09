package drops

import (
	"bytes"
	"container/list"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/cyfdecyf/tst"
)

type View struct {
	Children    *list.List
	Template    string
	Content     string
	InjectInto  string //ID to inject this view into
	Provides    string //ID that this view provides, into which we can inject other views
	ModelStruct interface{}
	Model       map[string]interface{}
	Return      string //key which view returns and will be used in model for templating

}

//Struct which holds facts for fields that are going to be rendered
type ComponentRule struct {
	Ignore bool
	View   *View
	Label  string
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
		fmt.Printf("Model for template: %s \n%v\n", v.Template, v.Model)
		v.Model["Content"] = v.Content

		fmt.Printf("Model struct: %v\n", v.ModelStruct)
		if v.ModelStruct != nil {
			buffer = RenderTemplate(v.Template, v.ModelStruct)
			//Find which ids this view provides and generate another template onto which content should be injected
			if v.Provides != "" {
				templateBeforeProcess := buffer.String()
				idIndex := strings.Index(templateBeforeProcess, v.Provides[1:])
				if idIndex > -1 {
					//if id exist find out for index of character '>'
					tagEndIndex := strings.IndexRune(templateBeforeProcess[idIndex:], '>')
					if tagEndIndex > -1 {
						//If i find tag end index
						fmt.Printf("idIndex: %v  tagEndIndex: %v \n", idIndex, tagEndIndex)
						var newBuffer bytes.Buffer
						newBuffer.WriteString(templateBeforeProcess[:idIndex+tagEndIndex+1])
						newBuffer.WriteString(v.Content)
						newBuffer.WriteString(templateBeforeProcess[idIndex+tagEndIndex+2:])
						buffer = newBuffer
					}
				}
			}
		} else {
			buffer = RenderTemplate(v.Template, v.Model)

		}
	} else {
		buffer.WriteString(v.Template)
	}
	return buffer
}

func RenderTemplate(templateName string, data interface{}) bytes.Buffer {
	var buffer bytes.Buffer
	// t, _ := template.ParseFiles(templateName)
	fmt.Printf("Rendering template: %+v\n", templateName)
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
	//if view.InjectInto != "" && view.Return == "" {
	if view.InjectInto != "" {
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

//Adds view to existing view
//It-s used when we don't have DOM existent
func AddToView(view *View, newView *View) *View {
	if newView != nil {
		parent := view.Children

		if parent == nil {
			parent = list.New()
		}
		parent.PushBack(newView)
		view.Children = parent
	}
	return view
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
	var fieldset *View
	fieldset = NewFieldset()

	fields := make([]*map[string]interface{}, 0)
	fmt.Printf("NNumber of fields %+v\n", typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		fieldName := typ.Field(i).Name
		fmt.Printf("Field name %v\n", fieldName)

		fieldMap := make(map[string]interface{}) //map passed to field generator
		fieldMap["Label"] = fieldName
		fieldMap["Name"] = fieldName
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

	view.Model["Fields"] = fields
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
		fmt.Printf("Ignore parsed %v\n", ignore)
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
	// var field *View
	// fmt.Printf("Ignoring on creating view %v\n", fieldRule.Ignore)
	if fieldRule.Ignore {
		return fieldRule
	}
	if model != "" {
		//Changing label to be the name of model if not explicitely set
		if fieldRule.Label == "" {
			fieldMap["Label"] = model
		}
		fmt.Printf("Selecting for model %v\n", model)
		allModel := GetAll(model)
		for _, v := range allModel {
			fmt.Printf("Model value %v\n", v)
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

		fieldRule.View = NewSelect(fieldMap)
	} else {
		fieldRule.View = NewInputText(fieldMap)
	}
	return fieldRule
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

//Magical method that creates form for editing based on provided value
func Edit(value interface{}, rules map[string]string, injectInto string, title string) *View {
	view := &View{Template: "edit.tpl", InjectInto: injectInto, Model: map[string]interface{}{
		"Title": title,
	}}

	fieldset := NewFieldset()

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
	view.Model["Fields"] = fields
	view = AddToView(view, fieldset)
	view.Model["Id"] = val.FieldByName("Id").Interface()
	fmt.Printf("Fields generated %v\n", fields)

	return view

}

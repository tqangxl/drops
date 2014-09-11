package component

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/kr/pretty"
	"github.com/mkasner/drops"
	"github.com/mkasner/drops/element"
	"github.com/mkasner/drops/event"
	"github.com/mkasner/drops/message"
	"github.com/mkasner/drops/protocol"
	"github.com/mkasner/drops/router"
	"github.com/mkasner/drops/session"
	"github.com/mkasner/drops/store"
)

func InitGenerator() {
	NewEvents()
	EditEvents()
}

//Struct which holds facts for fields that are going to be rendered
type ComponentRule struct {
	Ignore bool
	View   *element.View
	Label  string
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
	fieldset = NewFieldset()

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
		fieldset = element.AddToView(fieldset, fieldRule.View)
	}

	view.Model.MAP["Fields"] = fields
	view = element.AddToView(view, fieldset)
	fmt.Printf("Fields generated %v\n", fields)

	return view
}

//Injects edit handler event
func NewEvents() {
	editEvents := event.Events{
		{
			JsEvent: "submit",
			Origin:  "#form-new",
			Type:    "EVENT",
			Route:   "/model/new/",
			Handler: AddHandler,
		},
	}
	handle, _, _ := router.GetRouter().Lookup("EVENT", "/model/new/")
	if handle == nil {
		event.AddEvents(editEvents)
		drops.RegisterHandlers(editEvents)

	}
}

func AddHandler(w http.ResponseWriter, r *http.Request, params router.Params) *protocol.DropsResponse {
	data := params.ByName("data").(map[string]interface{})
	modelname := data["model"].(string)

	fmt.Printf("%# v", pretty.Formatter(data))

	store.AddModel(modelname, data)
	sessionId := session.GetSessionId(r, params.ByName("session").(string))

	message.NewMessage(sessionId, modelname+" saved")

	dom := session.GetSessionActiveDOM(sessionId)
	message.ProcessMessages(sessionId, dom)

	response := &protocol.DropsResponse{}
	response.Dom = dom
	response.Route = "/" + strings.ToLower(modelname) + "/"
	return response
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
		allModel := store.GetAll(model)
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

		fieldRule.View = NewSelect(fieldMap)
	} else {
		fieldRule.View = NewInputText(fieldMap)
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
	view.Model.TYPE = value
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
func Edit(value interface{}, rules map[string]string, injectInto string, provides string, title string) *element.View {
	view := &element.View{Template: "edit.tpl", InjectInto: injectInto, Provides: provides, Model: &element.Model{MAP: make(map[string]interface{})}}
	view.MAP["Title"] = title
	fieldset := NewFieldset()

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
		fieldset = element.AddToView(fieldset, fieldRule.View)
	}
	view.MAP["Fields"] = fields
	view = element.AddToView(view, fieldset)
	view.MAP["Id"] = val.FieldByName("Id").Interface()
	// fmt.Printf("Fields generated %v\n", fields)

	return view

}

//Injects edit handler event
func EditEvents() {
	editEvents := event.Events{
		{
			JsEvent: "submit",
			Origin:  "#form-edit",
			Type:    "EVENT",
			Route:   "/model/save/",
			Handler: SaveHandler,
		},
	}
	handle, _, _ := router.GetRouter().Lookup("EVENT", "/model/save/")
	if handle == nil {
		event.AddEvents(editEvents)
		drops.RegisterHandlers(editEvents)

	}
}

func SaveHandler(w http.ResponseWriter, r *http.Request, params router.Params) *protocol.DropsResponse {
	data := params.ByName("data").(map[string]interface{})
	modelname := data["model"].(string)
	id, err := strconv.ParseInt(data["id"].(string), 0, 64)
	if err != nil {
		log.Printf("Error parsing id %v", err)
	}
	data["id"] = id
	fmt.Printf("%# v", pretty.Formatter(data))

	store.SaveModel(modelname, data)

	sessionId := session.GetSessionId(r, params.ByName("session").(string))
	message.NewMessage(sessionId, modelname+" saved")

	dom := session.GetSessionActiveDOM(sessionId)
	message.ProcessMessages(sessionId, dom)

	response := &protocol.DropsResponse{}
	response.Dom = dom
	return response
}

// Get returns the value associated with key in the tag string.
// If there is no such key in the tag, Get returns the empty string.
// If the tag does not have the conventional format, the value
// returned by Get is unspecified.
func GetRule(tag string, key string) string {
	for tag != "" {
		// skip leading space
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// scan to colon.
		// a space or a quote is a syntax error
		i = 0
		for i < len(tag) && tag[i] != ' ' && tag[i] != ':' && tag[i] != '"' {
			i++
		}
		if i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// scan quoted string to find value
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			value, _ := strconv.Unquote(qvalue)
			return value
		}
	}
	return ""
}

package event

import (
	"bytes"
	"errors"
	"net/http"
	"text/template"

	"strings"
	"sync"

	"github.com/mkasner/drops/protocol"
	"github.com/mkasner/drops/router"
)

//Event struct that connects drops route with JS event
type Event struct {
	JsEvent string                                                                          //js event
	Origin  string                                                                          //js element
	Type    string                                                                          //Message type get, delete
	Route   string                                                                          //drops route
	Handler func(http.ResponseWriter, *http.Request, router.Params) *protocol.DropsResponse //drops handler
}

var mutex *sync.RWMutex

var eventStore Events
var eventStoreDirty bool //if eventstore is changed, mark it as dirty

//Creates new event store
func NewEventStore() {
	eventStore = make(Events, 0)
	mutex = &sync.RWMutex{}
	parseEventListenerTemplates()
}

//Template registry
var eventListenersStore map[string]*template.Template

func parseEventListenerTemplates() {
	var templ *template.Template
	eventListenersStore = make(map[string]*template.Template)
	templ = template.Must(template.New("clickEventListener").Parse(clickEventListenerTpl))
	eventListenersStore[templ.Name()] = templ
	templ = template.Must(template.New("submitEventListener").Parse(submitEventListenerTpl))
	eventListenersStore[templ.Name()] = templ
}

type Events []Event

func AddEvents(events Events) {
	newEvents := make(Events, 0)
	for _, e := range events {
		// fmt.Printf("Added event %# v\n", pretty.Formatter(e))
		if !EventExist(e) {
			// mutex.Lock()
			newEvents = append(newEvents, e)
			// mutex.Unlock()
		}
	}
	mutex.Lock()
	defer mutex.Unlock()

	eventStore = append(eventStore, newEvents...)
}

//Checks to see if event exists in the store
func EventExist(e Event) bool {
	var result bool
	for _, eStore := range eventStore {
		// fmt.Println("Comparing events ")
		// fmt.Printf("%# v", pretty.Formatter(e))
		// fmt.Printf("%# v", pretty.Formatter(eStore))
		if eStore.JsEvent == e.JsEvent && eStore.Origin == e.Origin {
			result = true
			// fmt.Printf("\nSame %v\n", result)

		}
	}
	return result
}

const clickEventListenerTpl = `$('body').on('{{.JsEvent}}','{{.Origin}}',function(e) {
		e.preventDefault();
		var data = e.currentTarget.dataset;
		var message = JSON.stringify({"type":"{{.Type}}", 
			"route": data.route,
			"data": data});
		Drops.sendMessage(message);
	})`
const submitEventListenerTpl = `$('body').on('{{.JsEvent}}','{{.Origin}}',function(e) {
		e.preventDefault();
		var data = $(this).serializeObject();
		 data = _.extend(data, e.currentTarget.dataset);
		var message = JSON.stringify({"type":"{{.Type}}", 
			"route": data.route,
			"data": data});
		Drops.sendMessage(message);
	})`

//Creates JS event listener from events array
func CreateEventListeners() (string, error) {
	var buffer bytes.Buffer
	// mutex.RLock()
	// defer mutex.Unlock()
	// fmt.Printf("Event store %# v\n", pretty.Formatter(eventStore))

	for _, e := range eventStore {
		if e.Origin != "" {
			eventListenerName := e.JsEvent + "EventListener"
			listener := renderTemplate(eventListenerName, &e)
			buffer.Write(listener.Bytes())
			buffer.WriteString("\n")
		}
	}
	// mutex.Lock()
	eventStoreDirty = false
	// mutex.Unlock()
	return buffer.String(), nil
}

//If event listeners dirty Create new eventListeners
func EventListenersModified() bool {

	return eventStoreDirty
}

//Renders eventlistener template
func renderTemplate(templateName string, data *Event) bytes.Buffer {
	var buffer bytes.Buffer
	// t, _ := template.ParseFiles(templateName)
	// fmt.Printf("Rendering template: %+v\n", templateName)
	if t, ok := eventListenersStore[templateName]; ok {
		t.Execute(&buffer, data)
	}
	// fmt.Printf("\n%+v\n", buffer.String())
	return buffer
}

//Splits event name and validates if it's correct
func ProcessEventRoute(event Event) ([]string, error) {
	eventSplit := strings.Split(event.Route, " ") //Split to type and name
	if len(eventSplit) != 2 {
		return nil, errors.New("Wrong event format in handler. [type] [name]")
	}
	return eventSplit, nil
}

package event

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestRegisterHandlers(t *testing.T) {
	fmt.Println("test")
	events := Events{
		{
			Origin:  ".menu-item",
			JsEvent: "click",
			Route:   "GET /drops/test",
			Handler: nil,
		},
	}
	RegisterHandlers(events)
}

func TestProcessEventRoute(t *testing.T) {
	fmt.Println("TestProcessEventName")
	events := Events{
		{
			Origin:  ".menu-item",
			JsEvent: "click",
			Route:   "GET /drops/test",
			Handler: nil,
		},
		{
			Origin:  ".delete-action",
			JsEvent: "click",
			Route:   "DELETE /drops/test/:id/",
			Handler: nil,
		},
	}
	expectedResults := [][]string{
		[]string{"GET", "/drops/test"},
		[]string{"DELETE", "/drops/test/:id/"},
	}
	for i, event := range events {
		result, err := ProcessEventRoute(event)
		if err != nil {
			t.FailNow()
		}
		equals(t, expectedResults[i], result)
	}
}

func TestCreateEventListeners(t *testing.T) {
	fmt.Println("TestCreateEventListeners")
	NewEventStore()
	events := Events{
		{
			Origin:  ".menu-item",
			JsEvent: "click",
			Route:   "GET /drops/test",
			Handler: nil,
		},
		{
			Origin:  ".delete-action",
			JsEvent: "click",
			Route:   "DELETE /drops/test/:id/",
			Handler: nil,
		},
	}
	AddEvents(events)
	eventListeners, err := CreateEventListeners()
	if err != nil {
		t.FailNow()
	}
	for _, event := range events {
		assert(t, strings.Contains(eventListeners, event.Origin), "Event listener not set: %v", event.Origin)
	}
	fmt.Println(eventListeners)

}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

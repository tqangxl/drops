//Event dispatcher
package drops

type EventHandler interface {
	Handle(data map[string]interface{}) []byte
}
type EventDispatcher map[string]EventHandler

var eventDispatcher EventDispatcher

func NewEventDispatcher() {
	eventDispatcher = EventDispatcher{}
}

func AddEventHandler(event string, handler EventHandler) {
	eventDispatcher[event] = handler
}

package events

import (
	"fmt"
	"strings"
	"sync"
)

const eventTypeOn string = "ON"
const eventTypeOnce string = "ONCE"

var availableEmitTypes = []string{eventTypeOn, eventTypeOnce}

func includesString(list []string, target string) bool {
	for _, item := range list {
		if target == item {
			return true
		}
	}
	return false
}

func validateHandler(handler interface{}) {
	switch handler.(type) {
	case func(), func(interface{}), func(...interface{}):
		// DO NOTHING
	default:
		panic(fmt.Sprintf("Param handler should be a function, got %T %v", handler, handler))
	}
}

func executeHandler(wg *sync.WaitGroup, handler interface{}, params ...interface{}) {
	defer wg.Done()
	switch h := handler.(type) {
	case func():
		h()
	case func(interface{}):
		h(params[0])
	case func(...interface{}):
		h(params...)
	default:
		panic(fmt.Sprintf("Param handler should be a function, got %T %v", handler, handler))
	}
}

// EventEmitter is the controller to add and remove listeners to event names.
type EventEmitter struct {
	handlerList map[string][]interface{}
}

// NewEventEmitter creates an pointer to a new instance of EventEmitter.
func NewEventEmitter() *EventEmitter {
	eventEmitter := EventEmitter{}
	eventEmitter.handlerList = make(map[string][]interface{})

	return &eventEmitter
}

// On add a handler on an event, and it will be called every time the event got triggered.
func (emitter *EventEmitter) On(eventName string, handler interface{}) {
	emitter.AddListener(eventTypeOn, eventName, handler)
}

// Once add a handler on an event, and it will be called only the first time the event got triggerd.
func (emitter *EventEmitter) Once(eventName string, handler interface{}) {
	emitter.AddListener(eventTypeOnce, eventName, handler)
}

// AddListener provides a way to add listener on an event, and by different emit type the handler will be called only once or every time.
func (emitter *EventEmitter) AddListener(emitType string, eventName string, handler interface{}) {
	emitType = strings.ToUpper(emitType)
	if !includesString(availableEmitTypes, emitType) {
		panic(fmt.Sprintf("Invalid emit type: %v", emitType))
	}

	if handler == nil {
		return
	}

	emitEventName := emitType + ":" + eventName

	emitter.pushHandlerByEmitEventName(emitEventName, handler)
}

func (emitter *EventEmitter) pushHandlerByEmitEventName(emitEventName string, handler interface{}) {
	if _, ok := emitter.handlerList[emitEventName]; !ok {
		emitter.handlerList[emitEventName] = []interface{}{}
	}

	validateHandler(handler)
	emitter.handlerList[emitEventName] = append(emitter.handlerList[emitEventName], handler)
}

// Emit an event and  listener handlers will be executed.
func (emitter *EventEmitter) Emit(eventName string, params ...interface{}) *sync.WaitGroup {
	var wg sync.WaitGroup
	for _, emitType := range availableEmitTypes {
		emitEventName := emitType + ":" + eventName
		emitter.emit(&wg, emitEventName, params...)
		if emitType == eventTypeOnce {
			delete(emitter.handlerList, emitEventName)
		}
	}

	return &wg
}

func (emitter *EventEmitter) emit(wg *sync.WaitGroup, emitEventName string, params ...interface{}) {
	if _, ok := emitter.handlerList[emitEventName]; !ok {
		return
	}
	for _, handler := range emitter.handlerList[emitEventName] {
		wg.Add(1)
		go executeHandler(wg, handler, params...)
	}
}

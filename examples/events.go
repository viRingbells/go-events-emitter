package examples

import (
	"errors"

	"github.com/viRingbells/go-libs/events"
)

func foo(input interface{}) {
	message := input.(string)
	println("Call func foo -> " + message)
}

func bar(input ...interface{}) {
	message := input[0].(string)
	e := input[1].(error)
	println("Call func bar -> " + message + " Error:" + e.Error())
}

func haz(input interface{}) {
	message := input.(string)
	println("Call func haz -> " + message)
}

// Events shows the example for events
func Events() {
	emitter := events.NewEventEmitter()

	println("on('Hello', anonymous func)")
	emitter.On("Hello", func() {
		println("Call anonymous func")
	})

	println("on('Hello', foo)")
	emitter.On("Hello", foo)

	println("once('Hello', bar)")
	emitter.Once("Hello", bar)

	println("on('World', haz)")
	emitter.On("World", haz)

	println("emit('Hello')")
	wg := emitter.Emit("Hello", "This is Message", errors.New("This is error"))
	wg.Wait()

	println("emit('Hello') again")
	wg = emitter.Emit("Hello", "This is Message", errors.New("This is error"))
	wg.Wait()

	println("emit('World')")
	wg = emitter.Emit("World", "This is Message", errors.New("This is error"))
	wg.Wait()

	println("emit('WhatEver')")
	wg = emitter.Emit("WhatEver", "This is Message", errors.New("This is error"))
	wg.Wait()

	println("Events Done")
}

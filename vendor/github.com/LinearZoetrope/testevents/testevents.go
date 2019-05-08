// Package testevents is an example for this question on StackOverflow:
// http://stackoverflow.com/questions/23699237/go-lang-test-events-listener/23706120#23706120
//
// It provides a simple wrapper of *testing.T and a basic event manager API to dispatch events on tests.
// An example of its use would be, at the top of your testing file:
//
//    import "github.com/Jragonmiris/testevents"
//    import "testing"
//
//    func init() {
//        testevents.Register(MyListener)
//    }
//
//    func TestXxx(tst *testing.T) {
//        t := testevents.Start(tst, "TestXxx", false)
//        defer t.Done()
//
//        err := DoSomething()
//        if err != nil {
//            t.Error(err)
//        }
//    }
//
// This package is engineered to be as close to the original testing framework as possible.
// Note that the Parallel function is intentionally left unimplemented, as the package is not
// concurrency safe.
package testevents

// EventListeners are functions that are called when the appropriate event
// is dispatched
//
// Funcs can only be compared to nil, so a listener
// is expected to have a unique numeric ID attached to it.
type EventListener struct {
	Listener func(Event)
	ID       int
}

var listeners map[EventType][]EventListener

func init() {
	listeners = make(map[EventType][]EventListener)
}

// The EventType is exactly what occurred during testing, such as the test passing
// or failing
type EventType uint

const (
	TestStarted EventType = iota
	TestPassed
	TestFinished
	TestFailed
	TestSkipped
)

// An Event is a simple pair of the name of the testing function being executed,
// and the test status (e.g. TestFailed, TestFinished, etc)
type Event struct {
	Name string
	Typ  EventType
}

// Registers a listener for a certain type of test event
//
// If the listener is already registered, this is a no op
func Register(typ EventType, listener EventListener) {
	if IsRegistered(typ, listener) {
		return
	}

	listeners[typ] = append(listeners[typ], listener)
}

// Unregisters a listener for a certain type of event
//
// If the listener is not registered, this is a no op
func Unregister(typ EventType, listener EventListener) {
	index := -1
	for i, l := range listeners[typ] {
		if l.ID == listener.ID {
			index = i
			break
		}
	}

	listeners[typ] = append(listeners[typ][:index], listeners[typ][index+1:]...)
}

// Returns whether or not a listener is registered already
func IsRegistered(typ EventType, listener EventListener) bool {
	for _, l := range listeners[typ] {
		if l.ID == listener.ID {
			return true
		}
	}

	return false
}

// Dispatches an event to all registered listeners
func Dispatch(e Event) {
	for _, l := range listeners[e.Typ] {
		l.Listener(e)
	}
}

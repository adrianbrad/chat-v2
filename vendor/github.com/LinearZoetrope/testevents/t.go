package testevents

import (
	"testing"
)

// A T is a wrapper over a *testing.T. Is dispatches events registered with the package
// depending on criteria. It can be configured for multiple reporting or no multiple reporting.
//
// If multiple reporting is enabled, a TestFailed event will be sent for each Fail(), Error() etc. If multiple
// reporting isn't enabled a TestFailed event will be send for the first failure, this status can be reset by calling Reset()
//
// A T is created with the Start function, which also dispatches a TestStarted event, and all tests should be finished
// with a call to Done() (preferable with "defer t.Done()") which dispatches either a TestFinished or TestPassed event
// depending on the failure status.
//
// Parallel() is left not implemented because the package is not concurrency safe
type T struct {
	t              *testing.T
	testName       string
	hasFailed      bool
	reportMultiple bool
}

// Start initialized a T and returns it. The first argument is the testing.T to wrap,
// the second is the name of the test being run (e.g. TestXxx) and reportMultiple
// sets whether or not events should be dispatched for multiple failures.
func Start(t *testing.T, testName string, reportMultiple bool) *T {
	out := &T{t: t, testName: testName, reportMultiple: reportMultiple, hasFailed: false}
	Dispatch(Event{Typ: TestStarted, Name: testName})

	return out
}

func (t *T) sendFail() {
	if !t.hasFailed || t.reportMultiple {
		Dispatch(Event{Typ: TestFailed, Name: t.testName})
	}

	t.hasFailed = true
}

// Error dispatches an event if the test hasn't already failed (or multiple reporting is enabled)
// and calls t.Error
func (t *T) Error(args ...interface{}) {
	t.sendFail()
	t.t.Error(args...)
}

// Error dispatches an event if the test hasn't already failed (or multiple reporting is enabled)
// and calls t.Errorf
func (t *T) Errorf(format string, args ...interface{}) {
	t.sendFail()
	t.t.Errorf(format, args...)
}

// Error dispatches an event if the test hasn't already failed (or multiple reporting is enabled)
// and calls t.Fail
func (t *T) Fail() {
	t.sendFail()
	t.t.Fail()
}

// Error dispatches an event if the test hasn't already failed (or multiple reporting is enabled)
// and calls t.FailNow
func (t *T) FailNow() {
	t.sendFail()
	t.t.FailNow()
}

// Error dispatches an event if the test hasn't already failed (or multiple reporting is enabled)
// and calls t.Fatal
func (t *T) Fatal(args ...interface{}) {
	t.sendFail()
	t.t.Fatal(args...)
}

// Error dispatches an event if the test hasn't already failed (or multiple reporting is enabled)
// and calls t.Fatalf
func (t *T) Fatalf(format string, args ...interface{}) {
	t.sendFail()
	t.t.Fatalf(format, args...)
}

// Log is equivalent to t.Log, nothing extra is done
func (t *T) Log(args ...interface{}) {
	t.t.Log(args...)
}

// Logf is equivalent to t.Logf, nothing extra is done
func (t *T) Logf(format string, args ...interface{}) {
	t.t.Logf(format, args...)
}

// Skip dispatches an event that the test has been skipped and then calls
// t.Skip
func (t *T) Skip(args ...interface{}) {
	Dispatch(Event{Typ: TestSkipped, Name: t.testName})
	t.t.Skip(args...)
}

// SkipNow dispatches an event that the test has been skipped and then calls
// t.SkipNow
func (t *T) SkipNow() {
	Dispatch(Event{Typ: TestSkipped, Name: t.testName})
	t.t.SkipNow()
}

// Skipf dispatches an event that the test has been skipped and then calls
// t.Skipf
func (t *T) Skipf(format string, args ...interface{}) {
	Dispatch(Event{Typ: TestSkipped, Name: t.testName})
	t.t.Skipf(format, args...)
}

// This is equivalent to t.Skipped. Nothing extra is done.
func (t *T) Skipped() bool {
	return t.t.Skipped()
}

// Resets the "hasFailed" status of T, meaning that if multiple reporting isn't enabled
// it will be treated as if it hasn't failed yet and will send another event on failure
//
// If multiple reporting was enabled on Start, this is irrelevant
func (t *T) Recover() {
	t.hasFailed = false
}

// If t has failed (and Reset hasn't been called since the most recent failure), this
// dispatched a TestFinished event. Otherwise it dispatches a TestPassed event.
func (t *T) Done() {
	if t.hasFailed {
		Dispatch(Event{Typ: TestFinished, Name: t.testName})
	} else {
		Dispatch(Event{Typ: TestPassed, Name: t.testName})
	}
}

package gournal

import (
	"errors"
)

// Tracker is the main objet uppon which errors are tracked
type Tracker struct {
	Context map[string]string
	Count   int
	Errors  map[int][]error
	Reports map[string]ReportFunction
}

// ReportFunction describes functions that are able to provide reports.
type ReportFunction func(tracker Tracker) interface{}

// NewTracker is the constructor for tracker objects
func NewTracker(context map[string]string, reportFunctions map[string]ReportFunction) Tracker {
	tracker := Tracker{}
	tracker.Reports = reportFunctions
	tracker.Context = context
	tracker.Errors = make(map[int][]error)
	return tracker
}

// Next step, increases count.
func (tracker *Tracker) Next() {
	tracker.Count++
}

// Add stores non nil errors
func (tracker *Tracker) Add(err error) {
	if err != nil {
		cycleErrors, ok := tracker.Errors[tracker.Count]
		if !ok {
			cycleErrors = make([]error, 0)
		}
		cycleErrors = append(cycleErrors, err)
		tracker.Errors[tracker.Count] = cycleErrors
	}
}

// ErrorsInCycleN returns errors recorded during the specified cycle
func (tracker Tracker) ErrorsInCycleN(cycle int) ([]error, error) {
	if cycle < 0 || cycle > tracker.Count {
		return nil, errors.New("Invalid cycle")
	}
	return tracker.Errors[cycle], nil
}

// ErrorsInCurrentCycle returns errors recorded during the current cycle
func (tracker Tracker) ErrorsInCurrentCycle() []error {
	res, _ := tracker.ErrorsInCycleN(tracker.Count)
	return res
}

// HasErrorInCurrentCycle tells if there are errors in the current cycle
func (tracker Tracker) HasErrorInCurrentCycle() bool {
	return len(tracker.ErrorsInCurrentCycle()) > 0
}

// HasErrorInCycleN tells if there are errors in the current cycle
func (tracker Tracker) HasErrorInCycleN(cycle int) (bool, error) {
	res, err := tracker.ErrorsInCycleN(cycle)
	return len(res) > 0, err
}

// Report executes a reporting function and returns the report
func (tracker Tracker) Report(code string) interface{} {
	if _, ok := tracker.Reports[code]; ok {
		return tracker.Reports[code](tracker)
	}
	return nil
}

// CountErrors returns total count of errors accounted since the start of the tracker
func (tracker Tracker) CountErrors() int {
	l := 0
	for _, e := range tracker.Errors {
		l += len(e)
	}
	return l
}

// CountErrorCycles returns count of cycle with at least 1 error since the start of the tracker
func (tracker Tracker) CountErrorCycles() int {
	l := 0
	for _, e := range tracker.Errors {
		if len(e) > 0 {
			l++
		}
	}
	return l
}

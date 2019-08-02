package gournal

import "regexp"

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

// Error stores non nil errors
func (tracker *Tracker) Error(err error) {
	if err != nil {
		cycleErrors, ok := tracker.Errors[tracker.Count]
		if !ok {
			cycleErrors = make([]error, 0)
		}
		cycleErrors = append(cycleErrors, err)
		tracker.Errors[tracker.Count] = cycleErrors
	}
}

// CurrentErrors returns errors recorded during the cycle
func (tracker Tracker) CurrentErrors(regexp *regexp.Regexp) []error {
	if regexp == nil {
		return tracker.Errors[tracker.Count]
	}
	var ret []error
	for _, err := range tracker.Errors[tracker.Count] {
		if regexp.MatchString(err.Error()) {
			ret = append(ret, err)
		}
	}
	return ret
}

// ErrorInCycle tells if there are errors in the current cycle
func (tracker Tracker) ErrorInCycle(regexp *regexp.Regexp) bool {
	return len(tracker.CurrentErrors(regexp)) > 0
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

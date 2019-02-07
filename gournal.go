package gournal

// Tracker helper pour suivre les erreurs lors d'un traitement
type Tracker struct {
	Context map[string]string
	Count   int
	Errors  map[int][]string
	Reports map[string]ReportFunction
}

// ReportFunction describes functions that are able to provide reports.
type ReportFunction func(tracker Tracker) interface{}

// Next step, increases count.
func (tracker *Tracker) Next() {
	tracker.Count++
}

// Error prend note des erreurs non nil
func (tracker *Tracker) Error(err error) {
	if err != nil {
		cycleErrors, ok := tracker.Errors[tracker.Count]
		if tracker.Errors == nil {
			tracker.Errors = make(map[int][]string)
		}
		if !ok {
			cycleErrors = make([]string, 0)
		}

		cycleErrors = append(cycleErrors, err.Error())

		tracker.Errors[tracker.Count] = cycleErrors
	}
}

// ErrorInCycle permet de savoir si des erreurs ont été constatées pendant le cycle
func (tracker Tracker) ErrorInCycle() bool {
	return len(tracker.Errors[tracker.Count]) > 0
}

// Report executes a reporting function and returns the report
func (tracker Tracker) Report(code string) interface{} {
	if _, ok := tracker.Reports[code]; ok {
		return tracker.Reports[code](tracker)
	} else {
		return nil
	}
}

// CountErrors returns total count of errors accounted since the start of the tracker
func (tracker Tracker) CountErrors() int {
	l := 0
	for _, e := range tracker.Errors {
		l += len(e)
	}
	return l
}

// CurrentErrors returns errors recorded during the cycle
func (tracker Tracker) CurrentErrors() []string {
	return tracker.Errors[tracker.Count]
}

// NewTracker provides a brand new tracker instance
func NewTracker(context map[string]string, reportFunctions map[string]ReportFunction) Tracker {
	tracker := Tracker{}
	tracker.Reports = reportFunctions
	tracker.Context = context
	tracker.Errors = make(map[int][]string)
	return tracker
}

package gournal

// Tracker helper pour suivre les erreurs lors d'un traitement
type Tracker struct {
	Cycles int
	Errors map[int][]string
}

// Ok incrémente le nombre de cycles ok
func (tracker *Tracker) Ok() {
	tracker.Cycles++
}

// Error prend note des erreurs non nil
func (tracker *Tracker) Error(err error) {
	if err != nil {
		cycleErrors, ok := tracker.Errors[tracker.Cycles]
		if tracker.Errors == nil {
			tracker.Errors = make(map[int][]string)
		}
		if !ok {
			cycleErrors = make([]string, 0)
		}

		cycleErrors = append(cycleErrors, err.Error())

		tracker.Errors[tracker.Cycles] = cycleErrors
	}
}

// ErrorInCycle permet de savoir si des erreurs ont été constatées pendant le cycle
func (tracker Tracker) ErrorInCycle() bool {
	return len(tracker.Errors[tracker.Cycles]) > 0
}

// CountErrors retourne le nombre d'erreurs du tracker
func (tracker Tracker) CountErrors() int {
	l := 0
	for _, e := range tracker.Errors {
		l += len(e)
	}
	return l
}

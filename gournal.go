package gournal

import (
	"errors"
	"time"
)

// Loglevel identifies priority of events
type Loglevel string

// Code identifies a tracker
type Code string

// Event test
type Event struct {
	Date     time.Time `json:"date" bson:"date"`
	Message  string    `json:"message" bson:"message"`
	Loglevel string    `json:"priority" bson:"priority"`
	Code     string    `json:"code" bson:"code"`
}

// Trackers gets all trackers together
type Trackers struct {
	Tracks map[string]tracker
	Output func(event Event) error
}

// NewTracker creates a new tracker
func (trackers *Trackers) NewTracker(code string) error {
	if _, ok := trackers.Tracks[code]; !ok {
		trackers.Tracks[code] = tracker{
			code:   code,
			count:  0,
			Events: map[int][]Event{},
		}
		return nil
	}
	return errors.New("Tracker " + code + " already exists")
}

type tracker struct {
	code   string
	count  int
	Events map[int][]Event
}

// Ok incrémente le nombre de cycles ok
func (tracker *tracker) Next() {
	tracker.count++
}

// Next incrémente le nombre de cycles de tous les trackers
func (trackers *Trackers) Next() {
	for _, tracker := range trackers.Tracks {
		tracker.count++
	}
}

// Error prend note des erreurs non nil
func (tracker *tracker) Error(err error) {
	if err != nil {
		cycleErrors, ok := tracker.Events[tracker.count]

		if tracker.Events == nil {
			tracker.Events = make(map[int][]Event)
		}

		if !ok {
			cycleErrors = make([]Event, 0)
		}

		event := Event{
			Date: time.Now(),
			Code: tracker.code,
		}

		cycleErrors = append(cycleErrors, event)
		tracker.Events[tracker.count] = cycleErrors
	}
}

// ErrorInCycle permet de savoir si des erreurs ont été constatées pendant le cycle
func (tracker tracker) LoglevelInCycle(level Loglevel) bool {
	return len(tracker.Errors[tracker.Cycles]) > 0
}

func (tracker tracker) CodeInCycle(Code Loglevel) bool {
	return len(tracker.Errors[tracker.Cycles]) > 0
}

// CountErrors retourne le nombre d'erreurs du tracker
func (tracker tracker) CountErrors() int {
	l := 0
	for _, e := range tracker.Errors {
		l += len(e)
	}
	return l
}

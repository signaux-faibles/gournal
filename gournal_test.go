package gournal

import (
	"errors"
	"fmt"
	"testing"
)

func testReportFunction(tracker Tracker) interface{} {
	return tracker.Context["test"]
}

func Test_NewTracker(t *testing.T) {
	context := map[string]string{
		"test": "test",
	}
	reportFunctions := map[string]ReportFunction{
		"test": testReportFunction,
	}
	tracker := NewTracker(context, reportFunctions)

	if tracker.Report("test") == "test" && tracker.Count == 0 {
		t.Log("tracker.NewTracker constructs a tracker object: Ok")
	} else {
		t.Fatal("tracker.NewTracker constructs a tracker object: Fail")
	}
}

func Test_Report(t *testing.T) {
	context := map[string]string{
		"test": "test",
	}
	reportFunctions := map[string]ReportFunction{
		"test": testReportFunction,
	}
	tracker := NewTracker(context, reportFunctions)

	if tracker.Report("toto") == nil {
		t.Log("tracker.Reports() returns nil when reportFunction is unknown: Ok")
	} else {
		t.Fatal("tracker.Reports() returns nil when reportFunction is unknown: Fail")
	}

	if tracker.Report("test") == "test" {
		t.Log("tracker.Reports() returns test when for testReportFunction: Ok")
	} else {
		t.Fatal("tracker.Reports() returns test when for testReportFunction: Fail")
	}
}

func Test_CountError(t *testing.T) {
	tracker := NewTracker(nil, nil)

	if tracker.CountErrors() == 0 {
		t.Log("CountError == 0: Ok")
	} else {
		t.Fatal("CountError == 0: Fail")
	}

	tracker.Error(errors.New("error 1"))
	tracker.Error(errors.New("error 2"))
	tracker.Error(errors.New("error 3"))
	tracker.Next()
	tracker.Error(errors.New("error 1"))
	tracker.Error(errors.New("error 2"))
	tracker.Error(errors.New("error 3"))
	tracker.Next()
	tracker.Next()

	if tracker.CountErrors() == 6 {
		t.Log("CountError == 6: Ok")
	} else {
		t.Fatal("CountError == 6: Fail")
	}
	if tracker.CountErrorCycles() == 2 {
		t.Log("CountErrorCycles == 2: Ok")
	} else {
		t.Fatal("CountErrorCycles == 2: Fail")
	}
}

func Test_HasErrorInCycleN(t *testing.T) {
	tracker := NewTracker(nil, nil)

	if b, err := tracker.HasErrorInCycleN(0); !b && err == nil {
		t.Log("New tracker has no error: Ok")
	} else {
		t.Fatal("New tracker has no error: Fail")
	}

	tracker.Error(errors.New("error 1"))
	tracker.Error(errors.New("error 2"))
	tracker.Error(errors.New("error 3"))

	if b, err := tracker.HasErrorInCycleN(0); b && err == nil {
		t.Log("tracker has errors: Ok ")
	} else {
		t.Fatal("Tracker has errors: Fail")
	}
	if _, err := tracker.HasErrorInCycleN(-1); err != nil {
		t.Log("Incorrect cycle fails: Ok ")
	} else {
		t.Fatal("Incorrect cycle fails: Fail")
	}
}

func Test_HasErrorInCurrentCycle(t *testing.T) {
	tracker := NewTracker(nil, nil)

	if !tracker.HasErrorInCurrentCycle() {
		t.Log("New tracker has no error: Ok")
	} else {
		t.Fatal("New tracker has no error: Fail")
	}

	tracker.Error(errors.New("error 1"))
	tracker.Error(errors.New("error 2"))
	tracker.Error(errors.New("error 3"))

	if tracker.HasErrorInCurrentCycle() {
		t.Log("Tracker with 3 errors in the current cycle has errors: Ok")
	} else {
		t.Fatal("Tracker with 3 errors in the current cycle has errors: Fail")
	}

	tracker.Next()
	if !tracker.HasErrorInCurrentCycle() {
		t.Log("Tracker with new cycle has no errors: Ok")
	} else {
		t.Fatal("Tracker with new cycle has no errors: Fail")
	}

	tracker.Error(nil)
	if !tracker.HasErrorInCurrentCycle() {
		t.Log("Tracker cycle with nil Error has no errors: Ok")
	} else {
		t.Fatal("Tracker with nil Error has no errors: Fail")
	}

}
func Test_ErrorsInCycleN(t *testing.T) {
	tracker := NewTracker(nil, nil)

	if errs, err := tracker.ErrorsInCycleN(0); errs == nil && err == nil {
		t.Log("New tracker has no error: Ok")
	} else {
		t.Fatal("New tracker has no error: Fail")
	}
	tracker.Error(errors.New("test error 1"))
	tracker.Error(errors.New("test error 2"))
	tracker.Next()

	if errs, err := tracker.ErrorsInCycleN(0); err == nil &&
		errs[0].Error() == "test error 1" &&
		errs[1].Error() == "test error 2" {
		t.Log("Tracker with 2 errors in current cycle has 2 errors: Ok")
	} else {
		t.Fatal("Tracker with 2 errors in current cycle has 2 errors: Fail")
	}

	if errs, err := tracker.ErrorsInCycleN(2); errs == nil && err != nil {
		t.Log("Invalid cycle fails: OK")
	} else {
		t.Fatal("Invalid cycle fails: Fail")
	}
}

func Test_ErrorsInCurrentCycle(t *testing.T) {
	tracker := NewTracker(nil, nil)

	if tracker.ErrorsInCurrentCycle() == nil {
		t.Log("New tracker has no error: Ok")
	} else {
		t.Fatal("New tracker has no error: Fail")
	}

	tracker.Error(errors.New("test error 1"))
	tracker.Error(errors.New("test error 2"))

	if len(tracker.ErrorsInCurrentCycle()) == 2 {
		t.Log("Tracker with 2 errors in current cycle has 2 errors: Ok")
	} else {
		t.Fatal("Tracker with 2 errors in current cycle has 2 errors: Fail")
	}

	if tracker.ErrorsInCurrentCycle()[0].Error() == "test error 1" && tracker.ErrorsInCurrentCycle()[1].Error() == "test error 2" {
		t.Log("Tracker keeps the good errors in the good order: Ok")
	} else {
		t.Fatal("Tracker keeps the good errors in the good order: Fail")
	}
	tracker.Next()
	fmt.Println(tracker.ErrorsInCurrentCycle() == nil)
	if tracker.ErrorsInCurrentCycle() == nil {
		t.Log("Tracker with error and a brand new cycle has no error: Ok")
	} else {
		t.Fatal("Tracker with error and a brand new cycle has no error: Fail")
	}
}

package gournal

import (
	"errors"
	"fmt"
	"regexp"
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

func Test_ErrorInCycle(t *testing.T) {
	tracker := NewTracker(nil, nil)

	if !tracker.ErrorInCycle(nil) {
		t.Log("New tracker has no error: Ok")
	} else {
		t.Fatal("New tracker has no error: Fail")
	}

	tracker.Error(errors.New("error 1"))
	tracker.Error(errors.New("error 2"))
	tracker.Error(errors.New("error 3"))

	if tracker.ErrorInCycle(nil) {
		t.Log("Tracker with 3 errors in the current cycle has errors: Ok")
	} else {
		t.Fatal("Tracker with 3 errors in the current cycle has errors: Fail")
	}

	regexp1, err1 := regexp.Compile("^[^error]")
	regexp2, err2 := regexp.Compile("[^1]$")
	if err1 != nil || err2 != nil {
		panic("Wrong regexp !")
	}
	if !tracker.ErrorInCycle(regexp1) {
		t.Log("Tracker with filtered errors returns no error: OK")
	} else {
		t.Fatal("Tracker with filtered errors returns no error: FAIL")
	}
	if tracker.ErrorInCycle(regexp2) {
		t.Log("Tracker with remaining errors returns an error: OK")
	} else {
		t.Fatal("Tracker with remaining errors returns an error: FAIL")
	}

	tracker.Next()
	if !tracker.ErrorInCycle(nil) {
		t.Log("Tracker with new cycle has no errors: Ok")
	} else {
		t.Fatal("Tracker with new cycle has no errors: Fail")
	}

	tracker.Error(nil)
	if !tracker.ErrorInCycle(nil) {
		t.Log("Tracker cycle with nil Error has no errors: Ok")
	} else {
		t.Fatal("Tracker with nil Error has no errors: Fail")
	}

}

func Test_CurrentErrors(t *testing.T) {
	tracker := NewTracker(nil, nil)

	if tracker.CurrentErrors(nil) == nil {
		t.Log("New tracker has no error: Ok")
	} else {
		t.Fatal("New tracker has no error: Fail")
	}

	tracker.Error(errors.New("test error 1"))
	tracker.Error(errors.New("test error 2"))

	if len(tracker.CurrentErrors(nil)) == 2 {
		t.Log("Tracker with 2 errors in current cycle has 2 errors: Ok")
	} else {
		t.Fatal("Tracker with 2 errors in current cycle has 2 errors: Fail")
	}

	if tracker.CurrentErrors(nil)[0].Error() == "test error 1" && tracker.CurrentErrors(nil)[1].Error() == "test error 2" {
		t.Log("Tracker keeps the good errors in the good order: Ok")
	} else {
		t.Fatal("Tracker keeps the good errors in the good order: Fail")
	}
	regexp, err := regexp.Compile("[^1]$")
	if err != nil {
		panic("Wrong regexp !")
	}
	if tracker.CurrentErrors(regexp)[0].Error() == "test error 2" {
		t.Log("Tracker keeps the good errors in the good order after filter: Ok")
	} else {
		t.Fatal("Tracker keeps the good errors in the good order after filter: Fail")
	}
	tracker.Next()
	fmt.Println(tracker.CurrentErrors(nil) == nil)
	if tracker.CurrentErrors(nil) == nil {
		t.Log("Tracker with error and a brand new cycle has no error: Ok")
	} else {
		t.Fatal("Tracker with error and a brand new cycle has no error: Fail")
	}
}

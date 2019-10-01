package test

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"
)

// ErrorTypeComparer defines an error type comparator for use with go-cmp
// Example:
//  opts := cmp.Options{
//  	 test.ErrorTypeComparer,
//  }
//
//  if !cmp.Equal(err, MyCustomError{}, opts) {
//  	t.Fatal(e.String())
//  }
//
var ErrorTypeComparer = cmp.Comparer(func(x error, y error) bool {
	return reflect.TypeOf(x) == reflect.TypeOf(y)
})

// ErrorReporter defines a reporter for use with go-cmp.
// This enables producing elegant diff output when the ErrorTypeComparator
// reports unequal.
type ErrorReporter struct {
	path  cmp.Path
	diffs []string
}

// String implements the Stringer interface.
func (e ErrorReporter) String() string {
	return strings.Join(e.diffs, "\n")
}

// PushStep is always called before report.
// This method is defined in the reporterIface interface for go-cmp.
func (e *ErrorReporter) PushStep(ps cmp.PathStep) {
	e.path = append(e.path, ps)
}

// Report is always called after PushStep and before PopStep.
// It is called wether the comparison reports equal, unequal or is ignored.
// This method is defined in the reporterIface interface for go-cmp.
func (e *ErrorReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		vx, vy := e.path.Last().Values()
		e.diffs = append(e.diffs, fmt.Sprintf(
			"\nerror type mismatch:\n\texpected: %T\n\tgot: %T\n",
			vx.Interface(),
			vy.Interface(),
		))
	}
}

// PopStep is always called after Report.
// This method is defined in the reporterIface interface for go-cmp.
func (e *ErrorReporter) PopStep() {
	e.path = e.path[:len(e.path)-1]
}

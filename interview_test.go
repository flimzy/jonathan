package jonathan

import (
	"os"
	"testing"
)

const CSVFile = "customers.csv"

func TestSort(t *testing.T) {
	_, err := os.Open(CSVFile)
	if err != nil {
		t.Fatalf("Unable to open %s: %s", CSVFile, err)
	}
}

type ColumnTest struct {
	Name     string
	Headings []string
	Expected int
	Error    string
}

var ColumnTests = []ColumnTest{
	ColumnTest{
		Name:     "Exact match",
		Headings: []string{"one", "email", "three"},
		Expected: 1,
	},
	ColumnTest{
		Name:     "Close match",
		Headings: []string{"one", "the_email_address", "three"},
		Expected: 1,
	},
	ColumnTest{
		Name:     "Exact match before Close match",
		Headings: []string{"one", "email", "the_email", "four"},
		Expected: 1,
	},
	ColumnTest{
		Name:     "Close match before Exact match",
		Headings: []string{"one", "some_email_address", "email", "four"},
		Expected: 2,
	},
	ColumnTest{
		Name:     "No match",
		Headings: []string{"one", "two", "three", "four"},
		Expected: 0,
		Error:    "No email column found",
	},
}

func TestFindEmailColumn(t *testing.T) {
	for _, test := range ColumnTests {
		if i, err := findEmailColumn(test.Headings); err != nil {
			if test.Error == "" {
				t.Fatalf("%s failed: %s", test.Name, err)
			}
			if test.Error != err.Error() {
				t.Fatalf("%s produced unexpected error '%s', expected '%s'", test.Name, err, test.Error)
			}
		} else if i != test.Expected {
			t.Fatalf("%s returned %d, not %d", test.Name, i, test.Expected)
		}
	}
}

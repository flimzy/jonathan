package jonathan

import (
	"os"
	"testing"
	// "github.com/davecgh/go-spew/spew"
)

const CSVFile = "customers.csv"

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

const ExpectedDomainStats = 500
const ExpectedAddrCount = 3000

var TopStats = []DomainStats{
	DomainStats{
		DomainName: "123-reg.co.uk",
		Addresses:  8,
	},
	DomainStats{
		DomainName: "163.com",
		Addresses:  6,
	},
	DomainStats{
		DomainName: "1688.com",
		Addresses:  3,
	},
	DomainStats{
		DomainName: "1und1.de",
		Addresses:  5,
	},
	DomainStats{
		DomainName: "360.cn",
		Addresses:  6,
	},
}

var BottomStats = []DomainStats{
	DomainStats{
		DomainName: "youku.com",
		Addresses:  5,
	},
	DomainStats{
		DomainName: "youtu.be",
		Addresses:  6,
	},
	DomainStats{
		DomainName: "youtube.com",
		Addresses:  3,
	},
	DomainStats{
		DomainName: "zdnet.com",
		Addresses:  8,
	},
	DomainStats{
		DomainName: "zimbio.com",
		Addresses:  3,
	},
}

func TestTally(t *testing.T) {
	file, err := os.Open(CSVFile)
	if err != nil {
		t.Fatalf("Unable to open %s: %s", CSVFile, err)
	}
	ds, err := TallyDomainStats(file)
	if err != nil {
		t.Fatalf("Error tallying stats: %s", err)
	}
	if len(ds) != ExpectedDomainStats {
		t.Errorf("Expected %d stats, got %d", ExpectedDomainStats, len(ds))
	}
	var addrCount int
	for _, s := range ds {
		addrCount += s.Addresses
	}
	if addrCount != ExpectedAddrCount {
		t.Errorf("Expected %d addresses, got %d", ExpectedAddrCount, addrCount)
	}
	for i, s := range TopStats {
		if ds[i].DomainName != s.DomainName {
			t.Errorf("Expected domain `%s` in position %d, got `%s`", s.DomainName, i, ds[i].DomainName)
		}
		if ds[i].Addresses != s.Addresses {
			t.Errorf("Expected %d addresses in position %d, got %d", s.Addresses, i, ds[i].Addresses)
		}
	}
	for j, s := range BottomStats {
		i := len(ds) - len(BottomStats) + j
		if ds[i].DomainName != s.DomainName {
			t.Errorf("Expected domain `%s` in position %d, got `%s`", s.DomainName, i, ds[i].DomainName)
		}
		if ds[i].Addresses != s.Addresses {
			t.Errorf("Expected %d addresses in position %d, got %d", s.Addresses, i, ds[i].Addresses)
		}
	}
}

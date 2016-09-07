package jonathan

import (
	"bytes"
	"log"
	"os"
	"regexp"
	"testing"
)

var buf bytes.Buffer

func init() {
	log.SetOutput(&buf)
}

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
const ExpectedLog = `YYYY/MM/DD hh:mm:ss [line 1002] Error parsing email address: mail: missing phrase
YYYY/MM/DD hh:mm:ss [line 2003] Error parsing email address: mail: missing phrase
`

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

var logRE = regexp.MustCompile(`(?m)^\d\d\d\d/\d\d/\d\d \d\d:\d\d:\d\d`)

func TestTally(t *testing.T) {
	file, err := os.Open(CSVFile)
	if err != nil {
		t.Fatalf("Unable to open %s: %s", CSVFile, err)
	}
	ds, err := TallyDomainStats(file)
	if err != nil {
		t.Fatalf("Error tallying stats: %s", err)
	}

	if l := logRE.ReplaceAllString(buf.String(), "YYYY/MM/DD hh:mm:ss"); l != ExpectedLog {
		t.Errorf("Log output different than expected.\nGot:\n%s---\nExpected:\n%s\n---\n",
			l, ExpectedLog)
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

type AddrTest struct {
	Input  string
	Domain string
	Error  string
}

var AddrTests = []AddrTest{
	AddrTest{
		Input:  "foo@foo.com",
		Domain: "foo.com",
	},
	AddrTest{
		Input:  "<foo@foo.com>",
		Domain: "foo.com",
	},
	AddrTest{
		Input: "an invalid address",
		Error: "mail: missing phrase",
	},
	AddrTest{
		Input: "",
		Error: "mail: no address",
	},
	AddrTest{
		Input:  `"John Doe" <foo@foo.com>`,
		Domain: "foo.com",
	},
	AddrTest{
		Input:  `"John Doe@Work" <foo@foo.com>`,
		Domain: "foo.com",
	},
	AddrTest{
		Input: `@foo.com`,
		Error: "mail: missing word in phrase: mail: invalid string",
	},
	AddrTest{
		Input: `foo@@foo.com`,
		Error: "mail: no angle-addr",
	},
}

func TestExtractDomain(t *testing.T) {
	for _, test := range AddrTests {
		result, err := extractDomain(test.Input)
		if err != nil {
			if test.Error == "" {
				t.Errorf("Error extracting domain from `%s`: %s", test.Domain, err)
			}
			if test.Error != err.Error() {
				t.Errorf("Unexpected error extracting domain from `%s`. Expected `%s`, got `%s`",
					test.Domain, test.Error, err)
			}
		}
		if result != test.Domain {
			t.Errorf("extractDomain() returned `%s`, expected `%s`", result, test.Domain)
		}
	}
}

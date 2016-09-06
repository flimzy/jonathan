// Package jonathan reads from the given customers.csv file and returns a sorted
// (data structure of your choice) of email domains along with the number of
// customers with e-mail addresses for each domain.  Any errors should be logged
// (or handled). Performance matters (this is only ~3k lines, but *could* be 1m
// lines or run on a small machine).
package jonathan

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"
)

// CaseInsensitiveAddresses can be set to true to ignore case in the local
// part of the email address. See RFC 2821.
var CaseInsensitiveAddresses bool

// DomainStats represents email address statistics for a given domain name.
// The domain name is converted to lowercase for consistency (see RFC 4343).
// The ad
type DomainStats struct {
	// DomainName contains the lowercase domain name for which address stats
	// have been calculated.
	DomainName string
	// Addresses counts the absolute number of addresses found to match the
	// domain name.
	Addresses int
	// UniqueAddresses counts the number of unique addresses found to match the
	// domain name.
	UniqueAddresses int
}

type ignoreCase bool

// TallyDomainStatsIgnoreCase works exactly as TallyDomainStats, but ignores
// case for the local portion of the email address.
func TallyDomainStatsIgnoreCase(r io.Reader) ([]DomainStats, error) {
	var ic ignoreCase = true
	return ic.tallyStats(r)
}

// TallyDomainStats reads a CSV file from the passed io.Reader, and returns a
// slice of DomainStats sorted by domain name. Stats are calculated on the
// domain-sensitive email address, to be compliant with RFC 2821.
func TallyDomainStats(r io.Reader) ([]DomainStats, error) {
	var ic ignoreCase = false
	return ic.tallyStats(r)
}

func (ic ignoreCase) tallyStats(r io.Reader) ([]DomainStats, error) {
	csvReader := csv.NewReader(r)
	headings, err := csvReader.Read()
	if err != nil {
		return nil, err
	}
	_, err = findEmailColumn(headings)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func findEmailColumn(headings []string) (int, error) {
	var closeMatch *int
	// First look for an exact match
	for i, heading := range headings {
		h := strings.ToLower(heading)
		if h == "email" || h == "e-mail" {
			// Return an exact match immediately
			return i, nil
		}
		if closeMatch == nil {
			if strings.Contains(h, "email") || strings.Contains(h, "e-mail") {
				x := i
				closeMatch = &x
			}
		}
	}
	if closeMatch != nil {
		return *closeMatch, nil
	}
	return 0, errors.New("No email column found")
}

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
	"log"
	"net/mail"
	"sort"
	"strings"
)

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
}

// DomainStatsSlice attaches the methods of sort.Interface to []DomainStats,
// sorting in increasing order.
type DomainStatsSlice []*DomainStats

// Len returns the length of the underlying slice.
func (s DomainStatsSlice) Len() int { return len(s) }

// Less reports whether the element i should sort before element j in the list.
func (s DomainStatsSlice) Less(i, j int) bool { return s[i].DomainName < s[j].DomainName }

// Swap swaps the elements i and j.
func (s DomainStatsSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// TallyDomainStats reads a CSV file from the passed io.Reader, and returns a
// slice of DomainStats sorted by domain name.
func TallyDomainStats(r io.Reader) ([]*DomainStats, error) {
	stats := make(map[string]*DomainStats)
	csvReader := csv.NewReader(r)
	headings, err := csvReader.Read()
	if err != nil {
		return nil, err
	}
	emailColumn, err := findEmailColumn(headings)
	if err != nil {
		return nil, err
	}
	counter := 1
	for {
		counter++
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("[line %d] Error parsing CSV: %s", counter, err)
			continue
		}
		domain, err := extractDomain(row[emailColumn])
		if err != nil {
			log.Printf("[line %d] Error parsing email address: %s", counter, err)
			continue
		}
		dc, ok := stats[domain]
		if !ok {
			dc = &DomainStats{DomainName: domain}
			stats[domain] = dc
		}
		dc.Addresses++
	}
	ds := DomainStatsSlice(make([]*DomainStats, 0, len(stats)))
	for _, stat := range stats {
		ds = append(ds, stat)
	}
	sort.Sort(ds)
	return ds, nil
}

func extractDomain(full string) (string, error) {
	addr, err := mail.ParseAddress(full)
	if err != nil {
		return "", err
	}
	// Always look for the last @, because the local part might contain quoted
	// or escaped @ signs.
	atIdx := strings.LastIndex(addr.Address, "@")
	if atIdx == -1 {
		return "", errors.New("No @ in email address")
	}
	return addr.Address[atIdx+1:], nil
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

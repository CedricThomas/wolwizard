package domain

import "strings"

type BootEntry struct {
	Name string
}

func (e BootEntry) MatchesOS(osName OSName) bool {
	return strings.Contains(strings.ToLower(e.Name), strings.ToLower(string(osName)))
}

func MatchGrubEntryToOS(entries []BootEntry, osName OSName) (string, error) {
	for _, entry := range entries {
		if entry.MatchesOS(osName) {
			return entry.Name, nil
		}
	}
	return "", &ErrGRUBEntryNotFound{OSName: osName}
}

type ErrGRUBEntryNotFound struct {
	OSName OSName
}

func (e *ErrGRUBEntryNotFound) Error() string {
	return "GRUB entry not found for OS: " + string(e.OSName)
}

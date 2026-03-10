package hosts

import (
	"fmt"
	"os"
	"strings"
)

const HostsPath = "/etc/hosts"

// Entry represents a single line in /etc/hosts
type Entry struct {
	IP     string
	Hosts  []string
	Line   string
	Number int
}

// ReadHostsFile reads the entire /etc/hosts file
func ReadHostsFile() ([]string, error) {
	data, err := os.ReadFile(HostsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read hosts file: %w", err)
	}
	lines := strings.Split(string(data), "\n")
	return lines, nil
}

// ParseHostsEntries parses all non-comment, non-empty lines from hosts file
func ParseHostsEntries(lines []string) []Entry {
	var entries []Entry
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		entries = append(entries, Entry{
			IP:     fields[0],
			Hosts:  fields[1:],
			Line:   line,
			Number: i,
		})
	}
	return entries
}

// FindDomainEntry checks if a domain is already in the hosts file
func FindDomainEntry(domain string) (Entry, error) {
	lines, err := ReadHostsFile()
	if err != nil {
		return Entry{}, err
	}

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		for _, host := range fields[1:] {
			if host == domain {
				return Entry{
					IP:    fields[0],
					Hosts: fields[1:],
					Line:  line,
				}, nil
			}
		}
	}

	return Entry{}, nil
}

// DomainExists checks if a domain entry exists in hosts file
func DomainExists(domain string) bool {
	entry, err := FindDomainEntry(domain)
	if err != nil {
		return false
	}
	return entry.IP != ""
}

// AddDomainEntry adds a domain to 127.0.0.1 in the hosts file
// Returns true if entry was added, false if it already existed
func AddDomainEntry(domain string) (bool, error) {
	if DomainExists(domain) {
		return false, nil
	}

	lines, err := ReadHostsFile()
	if err != nil {
		return false, err
	}

	// Check if there's an existing 127.0.0.1 entry we can append to
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "127.0.0.1" {
			// Append domain to existing entry
			lines[i] = fmt.Sprintf("%s\t%s", line, domain)
			return true, writeHostsFile(lines)
		}
	}

	// No existing 127.0.0.1 entry, append to end
	lines = append(lines, fmt.Sprintf("127.0.0.1\t%s", domain))
	return true, writeHostsFile(lines)
}

// RemoveDomainEntry removes a domain from the hosts file
// Returns true if entry was removed, false if it didn't exist
func RemoveDomainEntry(domain string) (bool, error) {
	lines, err := ReadHostsFile()
	if err != nil {
		return false, err
	}

	removed := false
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Check if domain is in this line
		for j, host := range fields[1:] {
			if host == domain {
				// Found it, remove this host from the line
				if len(fields) > 2 {
					// Line has multiple hosts, remove just this one
					newFields := append(fields[:j+1], fields[j+2:]...)
					lines[i] = fmt.Sprintf("%s\t%s", newFields[0], strings.Join(newFields[1:], "\t"))
				} else {
					// Line has only one host, comment out the line
					lines[i] = "# " + line
				}
				removed = true
				break
			}
		}
	}

	if !removed {
		return false, nil
	}

	return true, writeHostsFile(lines)
}

// writeHostsFile writes lines back to /etc/hosts
func writeHostsFile(lines []string) error {
	content := strings.Join(lines, "\n")
	err := os.WriteFile(HostsPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write hosts file: %w", err)
	}
	return nil
}

// ListDomains returns all domains pointing to 127.0.0.1 in hosts file
func ListDomains() ([]string, error) {
	entries, err := ReadAndParse()
	if err != nil {
		return nil, err
	}

	var domains []string
	for _, entry := range entries {
		if entry.IP == "127.0.0.1" {
			domains = append(domains, entry.Hosts...)
		}
	}
	return domains, nil
}

// ReadAndParse is a convenience function that reads and parses hosts file
func ReadAndParse() ([]Entry, error) {
	lines, err := ReadHostsFile()
	if err != nil {
		return nil, err
	}
	return ParseHostsEntries(lines), nil
}

// IsWritable checks if we can write to the hosts file (without root check)
func IsWritable() bool {
	file, err := os.OpenFile(HostsPath, os.O_WRONLY, 0)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

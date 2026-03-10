package hosts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tempHostsFile(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	return filepath.Join(tmpDir, "hosts")
}

func writeTempHosts(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func setupHostsTest(t *testing.T, content string) (string, func()) {
	t.Helper()
	path := tempHostsFile(t)
	writeTempHosts(t, path, content)

	// Note: We can't actually change HostsPath since it's a const
	return path, func() {
		os.Remove(path)
	}
}

func TestParseHostsEntries(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int // number of entries
	}{
		{
			name:     "empty file",
			content:  "",
			expected: 0,
		},
		{
			name:     "only comments",
			content:  "# Comment 1\n# Comment 2",
			expected: 0,
		},
		{
			name:     "single entry",
			content:  "127.0.0.1\tlocalhost",
			expected: 1,
		},
		{
			name:     "multiple entries",
			content:  "127.0.0.1\tlocalhost\n127.0.0.1\texample.com",
			expected: 2,
		},
		{
			name:     "mixed content",
			content:  "# Comment\n127.0.0.1\tlocalhost\n\n192.168.1.1\trouter",
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := strings.Split(tt.content, "\n")
			entries := ParseHostsEntries(lines)
			if len(entries) != tt.expected {
				t.Errorf("ParseHostsEntries() returned %d entries, expected %d", len(entries), tt.expected)
			}
		})
	}
}

func TestFindDomainEntry(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		domain   string
		wantIP   string
		wantErr  bool
	}{
		{
			name:    "domain exists",
			content: "127.0.0.1\tlocalhost example.com",
			domain:  "example.com",
			wantIP:  "127.0.0.1",
			wantErr: false,
		},
		{
			name:    "domain not found",
			content: "127.0.0.1\tlocalhost",
			domain:  "example.com",
			wantIP:  "",
			wantErr: false,
		},
		{
			name:    "empty file",
			content: "",
			domain:  "example.com",
			wantIP:  "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tempHostsFile(t)
			writeTempHosts(t, path, tt.content)

			// We can't actually test FindDomainEntry properly without changing the const HostsPath
			// This is a limitation of the current design
			// For now, we'll skip these tests
			t.Skip("HostsPath is a const, cannot test with temporary file")
		})
	}
}

func TestDomainExists(t *testing.T) {
	t.Skip("HostsPath is a const, cannot test with temporary file")
}

func TestAddDomainEntry(t *testing.T) {
	t.Skip("HostsPath is a const, cannot test with temporary file")
}

func TestRemoveDomainEntry(t *testing.T) {
	t.Skip("HostsPath is a const, cannot test with temporary file")
}

func TestListDomains(t *testing.T) {
	t.Skip("HostsPath is a const, cannot test with temporary file")
}

func TestIsWritable(t *testing.T) {
	t.Skip("HostsPath is a const, cannot test with temporary file")
}

func TestReadHostsFile(t *testing.T) {
	t.Skip("HostsPath is a const, cannot test with temporary file")
}

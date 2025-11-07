package scanner

import (
	"testing"
)

func TestParsePortRange(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
		wantErr  bool
	}{
		{"Single port", "80", []int{80}, false},
		{"Port range", "1-3", []int{1, 2, 3}, false},
		{"Multiple ports", "22,80,443", []int{22, 80, 443}, false},
		{"Mixed range and single", "1-3,80,443", []int{1, 2, 3, 80, 443}, false},
		{"Invalid port", "99999", nil, true},
		{"Invalid range", "3-1", nil, true},
		{"Empty input", "", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParsePortRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePortRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !equalSlices(result, tt.expected) {
				t.Errorf("ParsePortRange() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Benchmark tests
func BenchmarkScanUnlimited(b *testing.B) {
	config := &Config{
		Host:      "127.0.0.1",
		PortRange: generatePortRange(1, 100),
		Timeout:   100,
		Workers:   0, // Unlimited
		GetBanner: false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Scan(config)
	}
}

func BenchmarkScanWithWorkers10(b *testing.B) {
	config := &Config{
		Host:      "127.0.0.1",
		PortRange: generatePortRange(1, 100),
		Timeout:   100,
		Workers:   10,
		GetBanner: false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Scan(config)
	}
}

func BenchmarkScanWithWorkers50(b *testing.B) {
	config := &Config{
		Host:      "127.0.0.1",
		PortRange: generatePortRange(1, 100),
		Timeout:   100,
		Workers:   50,
		GetBanner: false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Scan(config)
	}
}

func BenchmarkScanWithWorkers100(b *testing.B) {
	config := &Config{
		Host:      "127.0.0.1",
		PortRange: generatePortRange(1, 100),
		Timeout:   100,
		Workers:   100,
		GetBanner: false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Scan(config)
	}
}

func generatePortRange(start, end int) []int {
	var ports []int
	for i := start; i <= end; i++ {
		ports = append(ports, i)
	}
	return ports
}

// Memory benchmark
func BenchmarkScanMemory(b *testing.B) {
	config := &Config{
		Host:      "127.0.0.1",
		PortRange: generatePortRange(1, 50),
		Timeout:   100,
		Workers:   10,
		GetBanner: false,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Scan(config)
	}
}

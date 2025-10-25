package scanner

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ScanResult struct {
	Port    int
	Open    bool
	Service string
	Banner  string
}

type Config struct {
	Host      string
	PortRange []int
	Timeout   int
	Workers   int
	GetBanner bool
}

func ParsePortRange(portsStr string) ([]int, error) {
	var ports []int

	// Handle comma-separated ports
	if strings.Contains(portsStr, ",") {
		parts := strings.Split(portsStr, ",")
		for _, part := range parts {
			if strings.Contains(part, "-") {
				// Handle range
				rangeParts := strings.Split(part, "-")
				if len(rangeParts) != 2 {
					return nil, fmt.Errorf("invalid port range: %s", part)
				}
				start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
				if err != nil {
					return nil, fmt.Errorf("invalid port number: %s", rangeParts[0])
				}
				end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
				if err != nil {
					return nil, fmt.Errorf("invalid port number: %s", rangeParts[1])
				}
				for i := start; i <= end; i++ {
					if i >= 1 && i <= 65535 {
						ports = append(ports, i)
					}
				}
			} else {
				// Handle single port
				port, err := strconv.Atoi(strings.TrimSpace(part))
				if err != nil {
					return nil, fmt.Errorf("invalid port number: %s", part)
				}
				if port >= 1 && port <= 65535 {
					ports = append(ports, port)
				}
			}
		}
	} else if strings.Contains(portsStr, "-") {
		// Handle simple range
		rangeParts := strings.Split(portsStr, "-")
		if len(rangeParts) != 2 {
			return nil, fmt.Errorf("invalid port range: %s", portsStr)
		}
		start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid port number: %s", rangeParts[0])
		}
		end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid port number: %s", rangeParts[1])
		}
		for i := start; i <= end; i++ {
			if i >= 1 && i <= 65535 {
				ports = append(ports, i)
			}
		}
	} else {
		// Handle single port
		port, err := strconv.Atoi(strings.TrimSpace(portsStr))
		if err != nil {
			return nil, fmt.Errorf("invalid port number: %s", portsStr)
		}
		if port >= 1 && port <= 65535 {
			ports = append(ports, port)
		}
	}

	if len(ports) == 0 {
		return nil, fmt.Errorf("no valid ports specified")
	}

	return ports, nil
}

func Scan(config *Config) []ScanResult {
	var results []ScanResult
	resultChan := make(chan ScanResult, len(config.PortRange))
	var wg sync.WaitGroup

	// Worker pool with semaphore
	if config.Workers > 0 {
		sem := make(chan struct{}, config.Workers)
		for _, port := range config.PortRange {
			wg.Add(1)
			go func(p int) {
				defer wg.Done()
				sem <- struct{}{}        // Acquire semaphore
				defer func() { <-sem }() // Release semaphore
				scanPort(config.Host, p, config.Timeout, config.GetBanner, resultChan)
			}(port)
		}
	} else {
		// Unlimited goroutines
		for _, port := range config.PortRange {
			wg.Add(1)
			go func(p int) {
				defer wg.Done()
				scanPort(config.Host, p, config.Timeout, config.GetBanner, resultChan)
			}(port)
		}
	}

	// Close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	for result := range resultChan {
		if result.Open {
			results = append(results, result)
		}
	}

	return results
}

func scanPort(host string, port int, timeoutMs int, getBanner bool, resultChan chan<- ScanResult) {
	address := fmt.Sprintf("%s:%d", host, port)
	timeout := time.Duration(timeoutMs) * time.Millisecond

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		resultChan <- ScanResult{Port: port, Open: false}
		return
	}
	defer conn.Close()

	result := ScanResult{
		Port: port,
		Open: true,
	}

	if getBanner {
		result.Service, result.Banner = grabBanner(conn, port, timeout)
	} else {
		result.Service = guessService(port)
	}

	resultChan <- result
}

func grabBanner(conn net.Conn, port int, timeout time.Duration) (string, string) {
	_ = conn.SetReadDeadline(time.Now().Add(timeout))

	// Try to read initial banner
	reader := bufio.NewReader(conn)
	banner, _ := reader.ReadString('\n')
	banner = strings.TrimSpace(banner)

	service := guessService(port)

	// For HTTP ports, send a request to get more info
	if port == 80 || port == 443 || port == 8080 || port == 8443 {
		httpResponse := grabHTTPBanner(conn, timeout)
		if httpResponse != "" {
			banner = httpResponse
			if port == 443 || port == 8443 {
				service = "HTTPS"
			} else {
				service = "HTTP"
			}
		}
	}

	// Identify service based on banner
	if service == "Unknown" && banner != "" {
		service = identifyService(banner)
	}

	return service, banner
}

func grabHTTPBanner(conn net.Conn, timeout time.Duration) string {
	_ = conn.SetWriteDeadline(time.Now().Add(timeout))
	_, err := conn.Write([]byte("GET / HTTP/1.0\r\nHost: localhost\r\n\r\n"))
	if err != nil {
		return ""
	}

	_ = conn.SetReadDeadline(time.Now().Add(timeout))
	scanner := bufio.NewScanner(conn)
	var response strings.Builder
	for i := 0; i < 10 && scanner.Scan(); i++ { // Read first 10 lines
		response.WriteString(scanner.Text())
		response.WriteString("\n")
	}

	return strings.TrimSpace(response.String())
}

func guessService(port int) string {
	switch port {
	case 21:
		return "FTP"
	case 22:
		return "SSH"
	case 23:
		return "Telnet"
	case 25:
		return "SMTP"
	case 53:
		return "DNS"
	case 80:
		return "HTTP"
	case 110:
		return "POP3"
	case 143:
		return "IMAP"
	case 443:
		return "HTTPS"
	case 993:
		return "IMAPS"
	case 995:
		return "POP3S"
	case 3306:
		return "MySQL"
	case 3389:
		return "RDP"
	case 5432:
		return "PostgreSQL"
	case 6379:
		return "Redis"
	case 8080:
		return "HTTP"
	case 8443:
		return "HTTPS"
	default:
		return "Unknown"
	}
}

func identifyService(banner string) string {
	banner = strings.ToUpper(banner)
	switch {
	case strings.Contains(banner, "SSH"):
		return "SSH"
	case strings.Contains(banner, "HTTP"):
		return "HTTP"
	case strings.Contains(banner, "FTP"):
		return "FTP"
	case strings.Contains(banner, "SMTP"):
		return "SMTP"
	case strings.Contains(banner, "POP3"):
		return "POP3"
	case strings.Contains(banner, "IMAP"):
		return "IMAP"
	case strings.Contains(banner, "MYSQL"):
		return "MySQL"
	case strings.Contains(banner, "POSTGRESQL"):
		return "PostgreSQL"
	case strings.Contains(banner, "REDIS"):
		return "Redis"
	default:
		return "Unknown"
	}
}

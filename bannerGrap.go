/*

__________                                      ________                                        
\______   \_____    ____   ____   ___________  /  _____/___________  ______        ____   ____  
 |    |  _/\__  \  /    \ /    \_/ __ \_  __ \/   \  __\_  __ \__  \ \____ \      / ___\ /  _ \ 
 |    |   \ / __ \|   |  \   |  \  ___/|  | \/\    \_\  \  | \// __ \|  |_> >    / /_/  >  <_> )
 |______  /(____  /___|  /___|  /\___  >__|    \______  /__|  (____  /   __/ /\  \___  / \____/ 
        \/      \/     \/     \/     \/               \/           \/|__|    \/ /_____/         
                                                                      Version 1.0

    Copyright 2025 MrEchoFi_Ebwer
	
	MIT License

Copyright (c) 2025 MrEchoFi_Md. Abu Naser Nayeem [Tanjib Isham]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/



package main

import (
	"bufio"
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// Protocols and their default payloads
var protocolPayloads = map[string]string{
	"http":  "GET / HTTP/1.1\r\nHost: %s\r\n\r\n",
	"https": "GET / HTTP/1.1\r\nHost: %s\r\n\r\n",
	"smtp":  "EHLO %s\r\n",
	"ftp":   "USER anonymous\r\n",
	"ssh":   "", // SSH typically sends banner first
	"telnet": "",
}

// Result struct for output
type BannerResult struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	Banner   string `json:"banner"`
	Error    string `json:"error,omitempty"`
}

// grabBanner connects to the host:port and retrieves the banner (optionally using TLS)
func grabBanner(host, port, protocol, payload string, timeout time.Duration, useTLS bool) BannerResult {
	address := net.JoinHostPort(host, port)
	var conn net.Conn
	var err error

	// Connect (plain or TLS)
	if useTLS {
		conn, err = tls.DialWithDialer(&net.Dialer{Timeout: timeout}, "tcp", address, &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		})
	} else {
		conn, err = net.DialTimeout("tcp", address, timeout)
	}
	if err != nil {
		return BannerResult{Host: host, Port: port, Protocol: protocol, Error: err.Error()}
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(timeout))

	// Send payload if provided
	if payload != "" {
		fmt.Fprintf(conn, payload, host)
	}

	// Read banner into a builder
	var banner strings.Builder
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			banner.Write(buf[:n])
			// break on HTTP double newline
			if strings.Contains(banner.String(), "\r\n\r\n") {
				break
			}
		}
		if err != nil {
			break
		}
	}
	return BannerResult{Host: host, Port: port, Protocol: protocol, Banner: banner.String()}
}

// parseTarget splits "host[:port]" into host and port (empty if none)
func parseTarget(target string) (host, port string) {
	parts := strings.Split(target, ":")
	switch len(parts) {
	case 2:
		return parts[0], parts[1]
	default:
		return parts[0], ""
	}
}

func main() {

	
    asciiArt := `
      
__________                                      ________                                        
\______   \_____    ____   ____   ___________  /  _____/___________  ______        ____   ____  
 |    |  _/\__  \  /    \ /    \_/ __ \_  __ \/   \  __\_  __ \__  \ \____ \      / ___\ /  _ \ 
 |    |   \ / __ \|   |  \   |  \  ___/|  | \/\    \_\  \  | \// __ \|  |_> >    / /_/  >  <_> )
 |______  /(____  /___|  /___|  /\___  >__|    \______  /__|  (____  /   __/ /\  \___  / \____/ 
        \/      \/     \/     \/     \/               \/           \/|__|    \/ /_____/         
																	
    `
    fmt.Println(asciiArt)

	// CLI flags
	targetsFile := flag.String("f", "", "File with list of targets (host:port per line)")
	protocol := flag.String("proto", "http", "Protocol: http, https, ftp, smtp, ssh, telnet, custom")
	portFlag := flag.String("port", "", "Port (overrides port in targets file/CLI)")
	payload := flag.String("payload", "", "Custom payload (default depends on protocol)")
	timeout := flag.Int("timeout", 5, "Timeout in seconds per connection")
	concurrency := flag.Int("threads", 10, "Number of concurrent scans")
	output := flag.String("o", "", "Output file (CSV or JSON, inferred by extension)")
	flag.Parse()

	// Build list of targets
	var targets []string
	if *targetsFile != "" {
		file, err := os.Open(*targetsFile)
		if err != nil {
			fmt.Println("[-] Error opening targets file:", err)
			os.Exit(1)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				targets = append(targets, line)
			}
		}
	} else {
		targets = flag.Args()
	}
	if len(targets) == 0 {

      

		fmt.Println("------------------ B A N N E R G R A P ------------------")
		fmt.Println("------------------ V E R S I O N   1.0 ------------------")
		fmt.Println("----------- Banner grabbing tool BY: MrEchoFi -----------")
		fmt.Println("---------------- Copyright 2025 MrEchoFi ----------------")
		fmt.Println("***********************************************************")
		fmt.Println("Usage: All Usages is on Guid_or_Usage.txt file.  |Flags:  bannerGrap [options] target1[:port] target2[:port]...")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Protocol & payload setup
	proto := strings.ToLower(*protocol)
	useTLS := proto == "https"
	// Determine payload
	protoPayload := protocolPayloads[proto]
	if *payload != "" {
		protoPayload = *payload
	}

	overridePort := *portFlag
	timeoutDur := time.Duration(*timeout) * time.Second

	// Prepare concurrency
	var wg sync.WaitGroup
	sem := make(chan struct{}, *concurrency)
	results := make([]BannerResult, len(targets))

	// Scan each target
	for i, tgt := range targets {
		wg.Add(1)
		go func(idx int, target string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			host, portStr := parseTarget(target)
			// Apply override or default
			if overridePort != "" {
				portStr = overridePort
			} else if portStr == "" {
				if useTLS {
					portStr = "443"
				} else {
					portStr = "80"
				}
			}

			results[idx] = grabBanner(host, portStr, proto, protoPayload, timeoutDur, useTLS)
		}(i, tgt)
	}
	wg.Wait()

	// Output results
	switch {
	case *output != "":
		if strings.HasSuffix(*output, ".json") {
			file, _ := os.Create(*output)
			defer file.Close()
			enc := json.NewEncoder(file)
			enc.SetIndent("", "  ")
			enc.Encode(results)
			fmt.Println("[+] Results written to", *output)
		} else if strings.HasSuffix(*output, ".csv") {
			file, _ := os.Create(*output)
			defer file.Close()
			w := csv.NewWriter(file)
			defer w.Flush()
			w.Write([]string{"host", "port", "protocol", "banner", "error"})
			for _, r := range results {
				w.Write([]string{r.Host, r.Port, r.Protocol, r.Banner, r.Error})
			}
			fmt.Println("[+] Results written to", *output)
		}
	default:
		for _, r := range results {
			fmt.Printf("Host: %s:%s | Protocol: %s\n", r.Host, r.Port, r.Protocol)
			if r.Error != "" {
				fmt.Printf("  [ERROR] %s\n", r.Error)
			} else {
				fmt.Printf("  Banner:\n%s\n", r.Banner)
			}
			fmt.Println(strings.Repeat("-", 60))
		}
	}
}

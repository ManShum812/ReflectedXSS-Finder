package main

import (
    "bufio"
    "context"
    "crypto/tls"
    "fmt"
    "io"
    "math/rand"
    "net/http"
    "net/url"
    "os"
    "strings"
    "sync"
    "time"

    "golang.org/x/net/http2"
)

const (
    maxWorkers   = 10 // Increase number of workers for higher concurrency
    batchSize    = 10 // Process URLs in batches
    timeout      = 10 * time.Second // Timeout for HTTP requests
    idleConnTimeout = 90 * time.Second
    maxIdleConns    = 100
    maxConnsPerHost = 100
)

// List of user agents to rotate through
var userAgentList = []string{
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/53.0",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/16.16299",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Safari/537.3",
}

func getRandomUserAgent() string {
    rand.Seed(time.Now().UnixNano())
    return userAgentList[rand.Intn(len(userAgentList))]
}

// Struct to hold URL scan result
type URLResult struct {
    URL        string
    StatusCode int
    Vulnerable bool
}

// Function to check for reflected XSS vulnerability
func checkXSS(targetURL string, payload string, mu *sync.Mutex, results *[]URLResult) {
    // Create an HTTP client with disabled SSL verification, custom redirect policy, and connection pooling
    transport := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        MaxIdleConns:    maxIdleConns,
        IdleConnTimeout: idleConnTimeout,
        MaxIdleConnsPerHost: maxConnsPerHost,
    }

    // Enable HTTP/2 support
    http2.ConfigureTransport(transport)

    client := &http.Client{
        Timeout:   timeout,
        Transport: transport,
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            if len(via) >= 10 {
                return fmt.Errorf("stopped after 10 redirects")
            }
            return nil
        },
    }

    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
    if err != nil {
        fmt.Printf("Failed to create request for %s: %s\n", targetURL, err)
        return
    }

    // Set a random user agent
    req.Header.Set("User-Agent", getRandomUserAgent())

    resp, err := client.Do(req)
    if err != nil {
        // Check for HTTP/2 specific error and retry with HTTP/1.1
        if strings.Contains(err.Error(), "HTTP/2 stream") || strings.Contains(err.Error(), "INTERNAL_ERROR") {
            fmt.Printf("HTTP/2 error for %s, retrying with HTTP/1.1...\n", targetURL)
            transport.ForceAttemptHTTP2 = false
            client.Transport = transport
            resp, err = client.Do(req)
        }

        if err != nil {
            fmt.Printf("Request to %s failed: %s\n", targetURL, err)
            return
        }
    }

    defer resp.Body.Close()

    // Ensure resp.Body is not nil before reading
    if resp.Body == nil {
        fmt.Printf("No response body for %s\n", targetURL)
        return
    }

    // Limit the size of the response body to 1MB to prevent excessive memory usage
    body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
    if err != nil {
        fmt.Printf("Failed to read response body from %s: %s\n", targetURL, err)
        return
    }

    decodedPayload, err := url.QueryUnescape(payload)
    if err != nil {
        fmt.Printf("Failed to decode payload: %s\n", err)
        return
    }

    vulnerable := strings.Contains(string(body), decodedPayload)
    if vulnerable {
        fmt.Printf("Vulnerable to XSS: %s\n", targetURL)
    } else {
        fmt.Printf("Not vulnerable: %s\n", targetURL)
    }

    result := URLResult{URL: targetURL, StatusCode: resp.StatusCode, Vulnerable: vulnerable}

    mu.Lock()
    *results = append(*results, result)
    mu.Unlock()

    // Write result to file
    writeResult(result)
}

func writeResult(result URLResult) {
    file, err := os.OpenFile("results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("Failed to open results file: %s\n", err)
        return
    }
    defer file.Close()

    _, err = file.WriteString(fmt.Sprintf("%s, Status Code: %d, Vulnerable: %t\n", result.URL, result.StatusCode, result.Vulnerable))
    if err != nil {
        fmt.Printf("Failed to write result to file: %s\n", err)
    }
}

func processBatch(urls []string, payload string, wg *sync.WaitGroup, mu *sync.Mutex, results *[]URLResult) {
    semaphore := make(chan struct{}, maxWorkers)

    for _, url := range urls {
        wg.Add(1)
        semaphore <- struct{}{}
        go func(url string) {
            defer wg.Done()
            defer func() { <-semaphore }()
            checkXSS(url, payload, mu, results)
        }(url)
    }

    wg.Wait()
}

func main() {
    inputFile := "input.txt"
    payload := "'\"><12345"

    file, err := os.Open(inputFile)
    if err != nil {
        fmt.Printf("Failed to open input file: %s\n", err)
        return
    }
    defer file.Close()

    var results []URLResult
    var wg sync.WaitGroup
    var mu sync.Mutex
    var urls []string

    scanner := bufio.NewScanner(file)
    count := 0

    // Create or truncate the results file
    err = os.WriteFile("results.txt", []byte(""), 0644)
    if err != nil {
        fmt.Printf("Failed to create results file: %s\n", err)
        return
    }

    for scanner.Scan() {
        urls = append(urls, scanner.Text())
        count++

        if count == batchSize {
            processBatch(urls, payload, &wg, &mu, &results)
            urls = nil // Clear the slice for the next batch
            count = 0
        }
    }

    // Process any remaining URLs
    if len(urls) > 0 {
        processBatch(urls, payload, &wg, &mu, &results)
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading input file: %s\n", err)
    }

    fmt.Println("Scanning complete. Results saved to results.txt.")
}

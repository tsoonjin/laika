package google

import (
    "log"
    "net/http"
    "io/ioutil"
    "strings"
    "fmt"
    "sync"
    "encoding/json"
    "github.com/m7shapan/njson"
)

var pageSpeedUrl string = "https://www.googleapis.com/pagespeedonline/v5/runPagespeed"

type PageSpeedInsightResponse struct {
    Overall float64 `njson:"lighthouseResult.categories.performance.score"`
    LCP     string  `njson:"lighthouseResult.audits.largest-contentful-paint.displayValue"`
    FCP     string  `njson:"lighthouseResult.audits.first-contentful-paint.displayValue"`
    TBT     string  `njson:"lighthouseResult.audits.total-blocking-time.displayValue"`
    SI      string  `njson:"lighthouseResult.audits.speed-index.displayValue"`
}

type Result struct {
    Overall float64 `json:"overall"`
    LCP     string  `json:"lcp"`
    FCP     string  `json:"fcp"`
    TBT     string  `json:"tbt"`
    SI      string  `json:"si"`
}

type ScanResult struct {
    Url string `json:"url"`
    Result Result `json:"result"`
}

type Config struct {
    ApiKey string `json:"api-key"`
}

type Client struct {
    Config Config
    url []string
}

func queryPageSpeed(url string, key string, resultChan chan ScanResult, wg *sync.WaitGroup) {
    defer wg.Done()
    log.Println("Scanning url: ", url)
    resp, err := http.Get(fmt.Sprintf("%s?url=%s&key=%s", pageSpeedUrl, url, key))
    if err != nil {
        log.Fatalln(err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalln(err)
    }

    var result PageSpeedInsightResponse
    njson.Unmarshal(body, &result)
    resultChan <- ScanResult{url, Result(result)}
}

func (c *Client) Start() []ScanResult {
    log.Println("Start running Page Speed Insight for following url(s):\n", strings.Join(c.url, "\n"))
    var scanResult []ScanResult
    var wg sync.WaitGroup
    resultChan := make(chan ScanResult,len(c.url))

    for _, url := range c.url {
        wg.Add(1)
        go queryPageSpeed(url, c.Config.ApiKey, resultChan, &wg)
    }
    wg.Wait()
    close(resultChan)
    for s := range resultChan {
        scanResult = append(scanResult, s)
    }
    return scanResult
}

func (c *Client) Handler(w http.ResponseWriter, r *http.Request) {
    var url string
    // Read query param
    queryParameters := r.URL.Query()
    if url = queryParameters.Get("url"); url == "" {
        log.Fatalln("Url(s) not provided")
    }
    c.url = strings.Split(url, ",")
    log.Println("Query Params", queryParameters)
    // Read request body
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatalln(err)
    }
    log.Println("Request Body", body)
    result := c.Start()
    // Set response
    respData, err := json.Marshal(result)
    if err != nil {
        log.Println("Unable to marshal result")
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Failed to parse result from PageSpeedInsight API!"))
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(respData)
}

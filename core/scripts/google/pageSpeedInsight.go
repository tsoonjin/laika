package google

import (
    "log"
    "net/http"
    "io/ioutil"
    "strings"
)

var pageSpeedUrl string = "https://www.googleapis.com/pagespeedonline/v5/runPagespeed"

type Config struct {
    ApiKey string `json:"api-key"`
}

type Client struct {
    Config Config
    url []string
}

func queryPageSpeed(url string) {
    log.Println("Scanning url: ", url)
}

func (c *Client) Start() {
    log.Println("Start running Page Speed Insight for following url(s):\n", strings.Join(c.url, "\n"))
    testUrl := c.url[0]
    queryPageSpeed(testUrl)
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
    c.Start()
}

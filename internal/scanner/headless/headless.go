// Package headless performs a bounded, same-origin browser crawl for SPA
// discovery. It records requests made by the page; it does not submit forms,
// click destructive controls, or execute model-generated instructions.
package headless

import (
	"context"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type Endpoint struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	Type   string `json:"resource_type,omitempty"`
	Status int64  `json:"status,omitempty"`
	Source string `json:"source"`
}

type Result struct {
	Target   string     `json:"target"`
	Started  time.Time  `json:"started_at"`
	Duration int64      `json:"duration_ms"`
	Browser  bool       `json:"browser_available"`
	URLs     []Endpoint `json:"endpoints"`
	Warning  string     `json:"warning,omitempty"`
}

func Crawl(target string, timeout time.Duration) Result {
	started := time.Now()
	result := Result{Target: target, Started: started}
	base, err := url.Parse(target)
	if err != nil || base.Host == "" || (base.Scheme != "http" && base.Scheme != "https") {
		result.Warning = "target must be an absolute HTTP URL"
		return result
	}
	if timeout <= 0 || timeout > 2*time.Minute {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	alloc, cancelAlloc := chromedp.NewExecAllocator(ctx, chromedp.NoFirstRun, chromedp.NoDefaultBrowserCheck, chromedp.Headless)
	defer cancelAlloc()
	browser, cancelBrowser := chromedp.NewContext(alloc)
	defer cancelBrowser()

	seen := map[string]bool{}
	var events []*network.EventRequestWillBeSent
	chromedp.ListenTarget(browser, func(ev interface{}) {
		if request, ok := ev.(*network.EventRequestWillBeSent); ok {
			events = append(events, request)
		}
	})
	result.Browser = true
	err = chromedp.Run(browser, network.Enable(), chromedp.Navigate(target), chromedp.Sleep(3*time.Second))
	if err != nil {
		result.Warning = "headless browser unavailable or navigation failed: " + err.Error()
	}
	for _, ev := range events {
		u, parseErr := url.Parse(ev.Request.URL)
		if parseErr != nil || u.Host != base.Host || (u.Scheme != base.Scheme && u.Scheme != "" && base.Scheme != "") {
			continue
		}
		key := ev.Request.Method + " " + ev.Request.URL
		if seen[key] || len(result.URLs) >= 500 {
			continue
		}
		seen[key] = true
		result.URLs = append(result.URLs, Endpoint{Method: ev.Request.Method, URL: ev.Request.URL, Type: string(ev.Type)})
	}
	sort.Slice(result.URLs, func(i, j int) bool { return result.URLs[i].URL < result.URLs[j].URL })
	result.Duration = time.Since(started).Milliseconds()
	return result
}

func IsAPIEndpoint(u string) bool {
	l := strings.ToLower(u)
	return strings.Contains(l, "/api/") || strings.Contains(l, "/rest/") || strings.Contains(l, "/graphql") || strings.Contains(l, "/v1/")
}

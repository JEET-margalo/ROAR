// Package spa detects single-page-app fallback behavior without crawling or
// mutating the target. It is used to explain path-based scan limitations.
package spa

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Profile struct {
	Target       string `json:"target"`
	Checked      int    `json:"checked"`
	Responses200 int    `json:"responses_200"`
	SPAFallback  bool   `json:"spa_fallback"`
	BaselineHash string `json:"baseline_hash,omitempty"`
	BaselineSize int    `json:"baseline_size,omitempty"`
	Warning      string `json:"warning,omitempty"`
}

func Detect(target string) Profile {
	p := Profile{Target: target}
	u, err := url.Parse(strings.TrimSpace(target))
	if err != nil || u.Scheme == "" || u.Host == "" {
		p.Warning = "target is not an absolute HTTP URL"
		return p
	}
	client := &http.Client{Timeout: 8 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}
	base := strings.TrimRight(u.Scheme+"://"+u.Host, "/")
	paths := []string{"/roar-nonexistent-a", "/roar-nonexistent-b", "/roar-nonexistent-c"}
	var baseline string
	for _, path := range paths {
		p.Checked++
		resp, reqErr := client.Get(base + path)
		if reqErr != nil {
			continue
		}
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			continue
		}
		p.Responses200++
		h := sha256.Sum256(body)
		hash := hex.EncodeToString(h[:])
		if baseline == "" {
			baseline = hash
			p.BaselineHash = hash
			p.BaselineSize = len(body)
		} else if hash != baseline {
			return p
		}
	}
	p.SPAFallback = p.Checked >= 2 && p.Responses200 == p.Checked && baseline != ""
	if p.SPAFallback {
		p.Warning = "target returns the same HTML for unknown paths; path-based results may be misleading"
	}
	return p
}

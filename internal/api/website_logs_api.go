package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/h0tak88r/ROAR/internal/utils"
)

var websiteLogMu sync.Mutex
var websiteSecretRE = regexp.MustCompile(`(?i)(authorization|api[_-]?key|token|password|secret|cookie)([=:]\s*|\s+)[^\s,;]+`)

// apiAppendWebsiteLog stores browser diagnostics as redacted JSONL for debugging.
func apiAppendWebsiteLog(c *gin.Context) {
	var entry map[string]interface{}
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid log entry"})
		return
	}
	entry["server_time"] = time.Now().UTC().Format(time.RFC3339Nano)
	entry["message"] = websiteSecretRE.ReplaceAllString(fmt.Sprint(entry["message"]), `$1$2[REDACTED]`)
	data, err := json.Marshal(entry)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "log entry could not be encoded"})
		return
	}
	if len(data) > 12000 {
		data = data[:12000]
	}
	path := websiteLogPath()
	websiteLogMu.Lock()
	defer websiteLogMu.Unlock()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err == nil {
		_, err = f.Write(append(data, '\n'))
		_ = f.Close()
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save website log"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"ok": true})
}

func websiteLogPath() string { return filepath.Join(utils.GetResultsDir(), "website.log") }

// apiExportWebsiteLog downloads the redacted browser log for AI-assisted debugging.
func apiListWebsiteLogs(c *gin.Context) {
	data, err := os.ReadFile(websiteLogPath())
	if os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{"entries": []interface{}{}, "total": 0})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	entries := make([]map[string]interface{}, 0, len(lines))
	for _, line := range lines {
		var entry map[string]interface{}
		if json.Unmarshal([]byte(line), &entry) == nil {
			entries = append(entries, entry)
		}
	}
	c.JSON(http.StatusOK, gin.H{"entries": entries, "total": len(entries)})
}

func apiExportWebsiteLog(c *gin.Context) {
	data, err := os.ReadFile(websiteLogPath())
	if os.IsNotExist(err) {
		data = []byte("No browser diagnostics have been recorded.\n")
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="roar-website.log"`)
	c.String(http.StatusOK, strings.TrimSpace(string(data))+"\n")
}

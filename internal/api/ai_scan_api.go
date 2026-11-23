package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/h0tak88r/ROAR/internal/brain"
	"github.com/h0tak88r/ROAR/internal/db"
	"github.com/h0tak88r/ROAR/internal/utils"
)

// apiAnalyzeScanWithAI performs optional, read-only AI triage over local scan
// artifacts. It never executes model output or submits reports.
func apiAnalyzeScanWithAI(c *gin.Context) {
	_ = db.Init()
	_ = db.EnsureSchema()
	scanID := strings.TrimSpace(c.Param("id"))
	if scanID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scan id required"})
		return
	}
	if _, err := db.GetScan(scanID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	entries, err := os.ReadDir(utils.GetScanResultsDir(scanID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "scan results are not available"})
		return
	}
	var b strings.Builder
	const maxContext = 120000
	for _, entry := range entries {
		if entry.IsDir() || strings.EqualFold(entry.Name(), "scan.log") {
			continue
		}
		data, readErr := os.ReadFile(filepath.Join(utils.GetScanResultsDir(scanID), entry.Name()))
		if readErr != nil || len(data) == 0 {
			continue
		}
		remaining := maxContext - b.Len()
		if remaining <= 0 {
			break
		}
		if len(data) > remaining {
			data = data[:remaining]
		}
		fmt.Fprintf(&b, "\n--- %s ---\n%s\n", entry.Name(), data)
	}
	if b.Len() == 0 {
		c.JSON(http.StatusOK, gin.H{"scan_id": scanID, "analysis": "No scan artifacts were available for AI triage."})
		return
	}

	system := `You are ROAR's defensive security triage assistant. Analyze only the supplied local scan artifacts. Do not invent vulnerabilities, do not recommend destructive actions, and do not execute commands. Separate confirmed evidence from hypotheses. For each candidate include severity, evidence, confidence, safe manual verification steps, remediation, and why it may be a false positive. State clearly when no issue is supported. All findings require human verification before reporting.`
	analysis, err := brain.ChatWithAI(nil, b.String(), system)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error(), "ai_available": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"scan_id": scanID, "analysis": analysis, "ai_available": true, "manual_review_required": true})
}

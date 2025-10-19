package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mhsanaei/3x-ui/v2/database"
	"github.com/mhsanaei/3x-ui/v2/database/model"
)

func TestGetXraySettingIncludesDatabaseManagedEntries(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	if err := database.InitDB(dbPath); err != nil {
		t.Fatalf("init db: %v", err)
	}
	defer database.CloseDB()

	db := database.GetDB()
	inbound := &model.Inbound{
		Tag:            "managed-in",
		Protocol:       model.VMESS,
		Port:           23456,
		Settings:       `{"clients":[]}`,
		StreamSettings: `{}`,
		Sniffing:       `{}`,
		Enable:         true,
	}
	if err := db.Create(inbound).Error; err != nil {
		t.Fatalf("seed inbound: %v", err)
	}

	outbound := &model.Outbound{
		Tag:      "managed-out",
		Protocol: "freedom",
		Settings: `{}`,
		Enable:   true,
	}
	if err := db.Create(outbound).Error; err != nil {
		t.Fatalf("seed outbound: %v", err)
	}

	router := gin.New()
	NewXraySettingController(router.Group("/panel"))

	req := httptest.NewRequest(http.MethodPost, "/panel/xray/", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("unexpected status %d", resp.Code)
	}

	var response struct {
		Success bool
		Obj     map[string]json.RawMessage
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if !response.Success {
		t.Fatalf("expected success response: %s", resp.Body.String())
	}

	var payload map[string]any
	if err := json.Unmarshal(response.Obj["xraySetting"], &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}

	inbounds, ok := payload["inbounds"].([]any)
	if !ok {
		t.Fatalf("inbounds type mismatch: %T", payload["inbounds"])
	}
	if !containsTag(inbounds, "managed-in") {
		t.Fatalf("expected managed inbound in payload: %v", payload["inbounds"])
	}

	outbounds, ok := payload["outbounds"].([]any)
	if !ok {
		t.Fatalf("outbounds type mismatch: %T", payload["outbounds"])
	}
	if !containsTag(outbounds, "managed-out") {
		t.Fatalf("expected managed outbound in payload: %v", payload["outbounds"])
	}
}

func containsTag(items []any, tag string) bool {
	for _, item := range items {
		if m, ok := item.(map[string]any); ok {
			if value, ok := m["tag"].(string); ok && value == tag {
				return true
			}
		}
	}
	return false
}

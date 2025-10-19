package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/mhsanaei/3x-ui/v2/database"
	"github.com/mhsanaei/3x-ui/v2/database/model"
	"github.com/mhsanaei/3x-ui/v2/logger"
	"github.com/mhsanaei/3x-ui/v2/web/session"
	logging "github.com/op/go-logging"
)

func TestGetXraySettingIncludesDatabaseManagedEntries(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	if err := database.InitDB(dbPath); err != nil {
		t.Fatalf("init db: %v", err)
	}
	defer database.CloseDB()

	setupTestLogger(t, tmpDir)

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
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session", store))
	router.Use(func(c *gin.Context) {
		session.SetLoginUser(c, &model.User{Id: 777})
		c.Next()
	})
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

func TestUpdateXraySettingPersistsManagedEntries(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	if err := database.InitDB(dbPath); err != nil {
		t.Fatalf("init db: %v", err)
	}
	defer database.CloseDB()

	setupTestLogger(t, tmpDir)

	router := gin.New()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session", store))
	router.Use(func(c *gin.Context) {
		session.SetLoginUser(c, &model.User{Id: 777})
		c.Next()
	})
	NewXraySettingController(router.Group("/panel"))

	initial := fetchXraySetting(t, router)

	var payload map[string]any
	if err := json.Unmarshal(initial, &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}

	payload["inbounds"] = append(copySlice(payload["inbounds"]), map[string]any{
		"tag":      "adv-added-in",
		"protocol": "vmess",
		"port":     31234,
		"settings": map[string]any{"clients": []any{}},
	})
	payload["outbounds"] = append(copySlice(payload["outbounds"]), map[string]any{
		"tag":      "adv-added-out",
		"protocol": "freedom",
		"settings": map[string]any{},
	})

	updated, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	form := strings.NewReader(url.Values{"xraySetting": {string(updated)}}.Encode())
	req := httptest.NewRequest(http.MethodPost, "/panel/xray/update", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("unexpected status %d", resp.Code)
	}

	var updateResp struct {
		Success bool
		Msg     string
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &updateResp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if !updateResp.Success {
		t.Fatalf("expected success, got %q", updateResp.Msg)
	}

	followup := fetchXraySetting(t, router)

	if !containsTagAnyRaw(followup, "inbounds", "adv-added-in") {
		t.Fatalf("expected inbound present after update")
	}
	if !containsTagAnyRaw(followup, "outbounds", "adv-added-out") {
		t.Fatalf("expected outbound present after update")
	}
}

func fetchXraySetting(t *testing.T, router *gin.Engine) json.RawMessage {
	t.Helper()

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
	return response.Obj["xraySetting"]
}

func containsTagAnyRaw(raw json.RawMessage, key, target string) bool {
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return false
	}
	return containsTagAny(payload[key], target)
}

func copySlice(value any) []any {
	if value == nil {
		return []any{}
	}
	if slice, ok := value.([]any); ok {
		return append([]any(nil), slice...)
	}
	return []any{}
}

func containsTagAny(section any, target string) bool {
	items, ok := section.([]any)
	if !ok {
		return false
	}
	for _, item := range items {
		if m, ok := item.(map[string]any); ok {
			if tag, ok := m["tag"].(string); ok && tag == target {
				return true
			}
		}
	}
	return false
}

func setupTestLogger(t *testing.T, logDir string) {
	t.Helper()
	os.Setenv("XUI_LOG_FOLDER", logDir)
	logger.InitLogger(logging.ERROR)
	t.Cleanup(func() {
		logger.CloseLogger()
	})
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

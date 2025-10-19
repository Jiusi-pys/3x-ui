package service

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"

	"github.com/mhsanaei/3x-ui/v2/database"
	"github.com/mhsanaei/3x-ui/v2/database/model"

	"gorm.io/gorm"
)

func TestApplyAdvancedSettingSyncsDatabaseResources(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	if err := database.InitDB(dbPath); err != nil {
		t.Fatalf("init db: %v", err)
	}
	defer database.CloseDB()

	settingSvc := SettingService{}
	xraySvc := XraySettingService{SettingService: settingSvc}
	if err := xraySvc.SaveXraySetting(`{"inbounds":[],"outbounds":[]}`); err != nil {
		t.Fatalf("seed template: %v", err)
	}

	inboundSvc := InboundService{}
	outboundSvc := OutboundService{}

	payload := `{"inbounds":[{"tag":"adv-in","protocol":"vmess","port":12345,"settings":{},"streamSettings":{},"sniffing":{}}],"outbounds":[{"tag":"adv-out","protocol":"freedom","settings":{}}]}`
	if err := xraySvc.ApplyAdvancedSetting(payload, &inboundSvc, &outboundSvc); err != nil {
		t.Fatalf("apply advanced setting: %v", err)
	}

	db := database.GetDB()

	var inbound model.Inbound
	if err := db.Where("tag = ?", "adv-in").First(&inbound).Error; err != nil {
		t.Fatalf("inbound not stored: %v", err)
	}
	if inbound.Port != 12345 {
		t.Fatalf("unexpected inbound port: %d", inbound.Port)
	}
	if !inbound.Enable {
		t.Fatalf("expected inbound enabled")
	}

	var outbound model.Outbound
	if err := db.Where("tag = ?", "adv-out").First(&outbound).Error; err != nil {
		t.Fatalf("outbound not stored: %v", err)
	}
	if outbound.Protocol != "freedom" {
		t.Fatalf("unexpected outbound protocol: %s", outbound.Protocol)
	}
	if !outbound.Enable {
		t.Fatalf("expected outbound enabled")
	}

	templateJSON, err := settingSvc.GetXrayConfigTemplate()
	if err != nil {
		t.Fatalf("get template: %v", err)
	}
	var template map[string]any
	if err := json.Unmarshal([]byte(templateJSON), &template); err != nil {
		t.Fatalf("unmarshal template: %v", err)
	}
	if containsTagAny(template["inbounds"], "adv-in") {
		t.Fatalf("managed inbound persisted in template: %v", template["inbounds"])
	}
	if containsTagAny(template["outbounds"], "adv-out") {
		t.Fatalf("managed outbound persisted in template: %v", template["outbounds"])
	}

	emptyPayload := `{"inbounds":[],"outbounds":[]}`
	if err := xraySvc.ApplyAdvancedSetting(emptyPayload, &inboundSvc, &outboundSvc); err != nil {
		t.Fatalf("remove managed resources: %v", err)
	}

	if err := db.Where("tag = ?", "adv-in").First(&model.Inbound{}).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected inbound deleted, got %v", err)
	}
	if err := db.Where("tag = ?", "adv-out").First(&model.Outbound{}).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected outbound deleted, got %v", err)
	}
}

func containsTagAny(section any, target string) bool {
	items := toSlice(section)
	for _, item := range items {
		if m, ok := toMap(item); ok {
			if getString(m["tag"]) == target {
				return true
			}
		}
	}
	return false
}

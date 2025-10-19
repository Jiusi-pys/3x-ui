package service

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mhsanaei/3x-ui/v2/database"
	"github.com/mhsanaei/3x-ui/v2/database/model"
	"github.com/mhsanaei/3x-ui/v2/util/common"
	"github.com/mhsanaei/3x-ui/v2/xray"

	"gorm.io/gorm"
)

// XraySettingService provides business logic for Xray configuration management.
// It handles validation and storage of Xray template configurations.
type XraySettingService struct {
	SettingService
}

func (s *XraySettingService) SaveXraySetting(newXraySettings string) error {
	if err := s.CheckXrayConfig(newXraySettings); err != nil {
		return err
	}
	return s.SettingService.saveSetting("xrayTemplateConfig", newXraySettings)
}

func (s *XraySettingService) CheckXrayConfig(XrayTemplateConfig string) error {
	xrayConfig := &xray.Config{}
	err := json.Unmarshal([]byte(XrayTemplateConfig), xrayConfig)
	if err != nil {
		return common.NewError("xray template config invalid:", err)
	}
	return nil
}

// ApplyAdvancedSetting synchronizes the advanced editor payload with template storage
// and database-backed inbound/outbound resources.
func (s *XraySettingService) ApplyAdvancedSetting(payload string, inboundService *InboundService, outboundService *OutboundService) error {
	if inboundService == nil || outboundService == nil {
		return errors.New("missing inbound or outbound service")
	}

	posted, err := decodeJSONPayload(payload)
	if err != nil {
		return err
	}

	templateJSON, err := s.SettingService.GetXrayConfigTemplate()
	if err != nil {
		return err
	}
	templateMap, err := decodeJSONPayload(templateJSON)
	if err != nil {
		return err
	}

	templateInboundTags := collectTemplateTags(templateMap["inbounds"])
	templateOutboundTags := collectTemplateTags(templateMap["outbounds"])

	postedInbounds := toSlice(posted["inbounds"])
	postedOutbounds := toSlice(posted["outbounds"])

	var managedInboundMaps []map[string]any
	var managedOutboundMaps []map[string]any

	templateInbounds := make([]any, 0, len(postedInbounds))
	for _, item := range postedInbounds {
		inboundMap, ok := toMap(item)
		if !ok {
			continue
		}
		tag := getString(inboundMap["tag"])
		if shouldPersistInTemplate(tag, templateInboundTags, inboundMap, true) {
			templateInbounds = append(templateInbounds, inboundMap)
		} else {
			managedInboundMaps = append(managedInboundMaps, inboundMap)
		}
	}

	templateOutbounds := make([]any, 0, len(postedOutbounds))
	for _, item := range postedOutbounds {
		outboundMap, ok := toMap(item)
		if !ok {
			continue
		}
		tag := getString(outboundMap["tag"])
		if shouldPersistInTemplate(tag, templateOutboundTags, outboundMap, false) {
			templateOutbounds = append(templateOutbounds, outboundMap)
		} else {
			managedOutboundMaps = append(managedOutboundMaps, outboundMap)
		}
	}

	templateMap["inbounds"] = templateInbounds
	templateMap["outbounds"] = templateOutbounds

	for key, value := range posted {
		if key == "inbounds" || key == "outbounds" {
			continue
		}
		templateMap[key] = value
	}

	sanitized, err := json.MarshalIndent(templateMap, "", "  ")
	if err != nil {
		return err
	}
	if err := s.SaveXraySetting(string(sanitized)); err != nil {
		return err
	}

	if err := syncManagedInbounds(managedInboundMaps, inboundService); err != nil {
		return err
	}
	return syncManagedOutbounds(managedOutboundMaps, outboundService)
}

func decodeJSONPayload(data string) (map[string]any, error) {
	if strings.TrimSpace(data) == "" {
		return map[string]any{}, nil
	}
	decoder := json.NewDecoder(strings.NewReader(data))
	decoder.UseNumber()
	result := map[string]any{}
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func collectTemplateTags(value any) map[string]struct{} {
	tags := make(map[string]struct{})
	for _, item := range toSlice(value) {
		if m, ok := toMap(item); ok {
			tag := getString(m["tag"])
			if tag != "" {
				tags[tag] = struct{}{}
			}
		}
	}
	return tags
}

func shouldPersistInTemplate(tag string, templateTags map[string]struct{}, payload map[string]any, isInbound bool) bool {
	if tag == "" {
		return true
	}
	if _, ok := templateTags[tag]; ok {
		return true
	}
	if isInbound {
		_, hasPort := payload["port"]
		_, hasProtocol := payload["protocol"]
		if !hasPort || !hasProtocol {
			return true
		}
	} else {
		if _, hasProtocol := payload["protocol"]; !hasProtocol {
			return true
		}
	}
	return false
}

func toSlice(value any) []any {
	if value == nil {
		return []any{}
	}
	if slice, ok := value.([]any); ok {
		return slice
	}
	return []any{}
}

func toMap(value any) (map[string]any, bool) {
	if value == nil {
		return map[string]any{}, false
	}
	m, ok := value.(map[string]any)
	return m, ok
}

func getString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	case fmt.Stringer:
		return v.String()
	default:
		return ""
	}
}

func marshalSection(value any) (string, error) {
	if value == nil {
		return "", nil
	}
	switch v := value.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return "", nil
		}
		return v, nil
	default:
		bs, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return "", err
		}
		return string(bs), nil
	}
}

func parsePort(value any) (int, error) {
	switch v := value.(type) {
	case nil:
		return 0, nil
	case json.Number:
		i, err := v.Int64()
		return int(i), err
	case float64:
		return int(v), nil
	case float32:
		return int(v), nil
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case uint64:
		return int(v), nil
	case string:
		if strings.TrimSpace(v) == "" {
			return 0, nil
		}
		i, err := strconv.Atoi(v)
		return i, err
	default:
		return 0, fmt.Errorf("invalid port type %T", value)
	}
}

func syncManagedInbounds(items []map[string]any, inboundService *InboundService) error {
	db := database.GetDB()
	existing, err := inboundService.GetAllInbounds()
	if err != nil {
		return err
	}

	existingByTag := make(map[string]*model.Inbound, len(existing))
	for _, inbound := range existing {
		existingByTag[inbound.Tag] = inbound
	}

	seen := map[string]struct{}{}
	for _, payload := range items {
		tag := getString(payload["tag"])
		inboundModel, err := buildInboundFromPayload(payload, existingByTag[tag])
		if err != nil {
			return err
		}
		if inboundModel == nil || inboundModel.Tag == "" {
			continue
		}

		if err := db.Save(inboundModel).Error; err != nil {
			return err
		}
		seen[inboundModel.Tag] = struct{}{}
	}

	for tag, inbound := range existingByTag {
		if _, ok := seen[tag]; ok {
			continue
		}
		if err := db.Delete(&model.Inbound{}, inbound.Id).Error; err != nil {
			return err
		}
	}
	return nil
}

func buildInboundFromPayload(payload map[string]any, existing *model.Inbound) (*model.Inbound, error) {
	if existing == nil {
		existing = &model.Inbound{}
		existing.Enable = true
	} else {
		copy := *existing
		existing = &copy
	}

	existing.Tag = getString(payload["tag"])
	if listen, ok := payload["listen"].(string); ok {
		existing.Listen = listen
	} else if payload["listen"] == nil {
		existing.Listen = ""
	}
	port, err := parsePort(payload["port"])
	if err != nil {
		return nil, err
	}
	existing.Port = port
	if protocol, ok := payload["protocol"].(string); ok {
		existing.Protocol = model.Protocol(protocol)
	}

	settings, err := marshalSection(payload["settings"])
	if err != nil {
		return nil, err
	}
	existing.Settings = settings

	streamSettings, err := marshalSection(payload["streamSettings"])
	if err != nil {
		return nil, err
	}
	existing.StreamSettings = streamSettings

	sniffing, err := marshalSection(payload["sniffing"])
	if err != nil {
		return nil, err
	}
	existing.Sniffing = sniffing

	existing.Remark = defaultRemark(existing.Remark, existing.Tag)
	return existing, nil
}

func syncManagedOutbounds(items []map[string]any, outboundService *OutboundService) error {
	db := database.GetDB()
	existing, err := outboundService.GetAllOutbounds()
	if err != nil {
		return err
	}

	existingByTag := make(map[string]*model.Outbound)
	for _, outbound := range existing {
		if outbound.Enable {
			existingByTag[outbound.Tag] = outbound
		}
	}

	seen := map[string]struct{}{}
	for _, payload := range items {
		tag := getString(payload["tag"])
		outboundModel, err := buildOutboundFromPayload(payload, existingByTag[tag])
		if err != nil {
			return err
		}
		if outboundModel == nil || outboundModel.Tag == "" {
			continue
		}

		if err := saveOutbound(db, outboundModel).Error; err != nil {
			return err
		}
		seen[outboundModel.Tag] = struct{}{}
	}

	for tag, outbound := range existingByTag {
		if _, ok := seen[tag]; ok {
			continue
		}
		if err := db.Delete(&model.Outbound{}, outbound.Id).Error; err != nil {
			return err
		}
	}
	return nil
}

func buildOutboundFromPayload(payload map[string]any, existing *model.Outbound) (*model.Outbound, error) {
	if existing == nil {
		existing = &model.Outbound{}
		existing.Enable = true
		existing.CreatedAt = time.Now().Unix()
	} else {
		copy := *existing
		existing = &copy
	}

	existing.Tag = getString(payload["tag"])
	if protocol, ok := payload["protocol"].(string); ok {
		existing.Protocol = protocol
	}

	settings, err := marshalSection(payload["settings"])
	if err != nil {
		return nil, err
	}
	existing.Settings = settings

	streamSettings, err := marshalSection(payload["streamSettings"])
	if err != nil {
		return nil, err
	}
	existing.StreamSettings = streamSettings

	proxySettings, err := marshalSection(payload["proxySettings"])
	if err != nil {
		return nil, err
	}
	existing.ProxySettings = proxySettings

	mux, err := marshalSection(payload["mux"])
	if err != nil {
		return nil, err
	}
	existing.Mux = mux

	existing.Remark = defaultRemark(existing.Remark, existing.Tag)
	existing.UpdatedAt = time.Now().Unix()
	return existing, nil
}

func saveOutbound(db *gorm.DB, outbound *model.Outbound) *gorm.DB {
	if outbound.Id == 0 {
		return db.Create(outbound)
	}
	return db.Save(outbound)
}

func defaultRemark(current, tag string) string {
	if strings.TrimSpace(current) != "" {
		return current
	}
	return tag
}

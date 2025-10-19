package controller

import (
	"encoding/json"
	"errors"

	"github.com/mhsanaei/3x-ui/v2/database/model"
	"github.com/mhsanaei/3x-ui/v2/web/service"

	"github.com/gin-gonic/gin"
)

// XraySettingController handles Xray configuration and settings operations.
type XraySettingController struct {
	XraySettingService service.XraySettingService
	SettingService     service.SettingService
	InboundService     service.InboundService
	OutboundService    service.OutboundService
	XrayService        service.XrayService
	WarpService        service.WarpService
}

// NewXraySettingController creates a new XraySettingController and initializes its routes.
func NewXraySettingController(g *gin.RouterGroup) *XraySettingController {
	a := &XraySettingController{}
	a.initRouter(g)
	return a
}

// initRouter sets up the routes for Xray settings management.
func (a *XraySettingController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/xray")
	g.GET("/getDefaultJsonConfig", a.getDefaultXrayConfig)
	g.GET("/getOutboundsTraffic", a.getOutboundsTraffic)
	g.GET("/getXrayResult", a.getXrayResult)

	g.POST("/", a.getXraySetting)
	g.POST("/warp/:action", a.warp)
	g.POST("/update", a.updateSetting)
	g.POST("/resetOutboundsTraffic", a.resetOutboundsTraffic)
}

// getXraySetting retrieves the Xray configuration template and inbound tags.
func (a *XraySettingController) getXraySetting(c *gin.Context) {
	templateJSON, err := a.SettingService.GetXrayConfigTemplate()
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.getSettings"), err)
		return
	}
	var template map[string]any
	if err := json.Unmarshal([]byte(templateJSON), &template); err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.getSettings"), err)
		return
	}

	// Ensure template structure exists
	if template == nil {
		template = map[string]any{}
	}

	template["inbounds"], err = a.appendManagedInbounds(template["inbounds"])
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.getSettings"), err)
		return
	}

	template["outbounds"], err = a.appendManagedOutbounds(template["outbounds"])
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.getSettings"), err)
		return
	}

	combinedTemplate, err := json.Marshal(template)
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.getSettings"), err)
		return
	}
	inboundTags, err := a.InboundService.GetInboundTags()
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.getSettings"), err)
		return
	}
	effectiveConfig, err := a.XrayService.GetXrayConfig()
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.getSettings"), err)
		return
	}

	effectiveConfigJSON, err := json.Marshal(effectiveConfig)
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.getSettings"), err)
		return
	}

	response := map[string]json.RawMessage{
		"xraySetting":     json.RawMessage(combinedTemplate),
		"inboundTags":     json.RawMessage(inboundTags),
		"effectiveConfig": json.RawMessage(effectiveConfigJSON),
	}

	jsonObj(c, response, nil)
}

// updateSetting updates the Xray configuration settings.
func (a *XraySettingController) updateSetting(c *gin.Context) {
	xraySetting := c.PostForm("xraySetting")
	if xraySetting == "" {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.modifySettings"), errors.New("empty xraySetting payload"))
		return
	}

	err := a.XraySettingService.ApplyAdvancedSetting(xraySetting, &a.InboundService, &a.OutboundService)
	if err == nil {
		a.XrayService.SetToNeedRestart()
	}
	jsonMsg(c, I18nWeb(c, "pages.settings.toasts.modifySettings"), err)
}

func (a *XraySettingController) appendManagedInbounds(base any) (any, error) {
	var templateInbounds []any
	if baseSlice, ok := base.([]any); ok {
		templateInbounds = append(templateInbounds, baseSlice...)
	} else if base != nil {
		if raw, ok := base.(json.RawMessage); ok {
			if err := json.Unmarshal(raw, &templateInbounds); err != nil {
				return nil, err
			}
		}
	}

	inbounds, err := a.InboundService.GetAllInbounds()
	if err != nil {
		return nil, err
	}

	for _, inbound := range inbounds {
		inboundMap, err := convertInboundModel(inbound)
		if err != nil {
			return nil, err
		}
		templateInbounds = append(templateInbounds, inboundMap)
	}

	return templateInbounds, nil
}

func (a *XraySettingController) appendManagedOutbounds(base any) (any, error) {
	var templateOutbounds []any
	if baseSlice, ok := base.([]any); ok {
		templateOutbounds = append(templateOutbounds, baseSlice...)
	} else if base != nil {
		if raw, ok := base.(json.RawMessage); ok {
			if err := json.Unmarshal(raw, &templateOutbounds); err != nil {
				return nil, err
			}
		}
	}

	outbounds, err := a.OutboundService.GetEnabledOutbounds()
	if err != nil {
		return nil, err
	}

	for _, outbound := range outbounds {
		outboundMap, err := convertOutboundModel(outbound)
		if err != nil {
			return nil, err
		}
		templateOutbounds = append(templateOutbounds, outboundMap)
	}

	return templateOutbounds, nil
}

func convertInboundModel(inbound *model.Inbound) (map[string]any, error) {
	inboundMap := map[string]any{
		"tag":      inbound.Tag,
		"protocol": string(inbound.Protocol),
		"port":     inbound.Port,
	}
	if inbound.Listen != "" {
		inboundMap["listen"] = inbound.Listen
	}
	if inbound.Settings != "" {
		var settings any
		if err := json.Unmarshal([]byte(inbound.Settings), &settings); err != nil {
			return nil, err
		}
		inboundMap["settings"] = settings
	}
	if inbound.StreamSettings != "" {
		var streamSettings any
		if err := json.Unmarshal([]byte(inbound.StreamSettings), &streamSettings); err != nil {
			return nil, err
		}
		inboundMap["streamSettings"] = streamSettings
	}
	if inbound.Sniffing != "" {
		var sniffing any
		if err := json.Unmarshal([]byte(inbound.Sniffing), &sniffing); err != nil {
			return nil, err
		}
		inboundMap["sniffing"] = sniffing
	}
	return inboundMap, nil
}

func convertOutboundModel(outbound *model.Outbound) (map[string]any, error) {
	outboundMap := map[string]any{
		"tag":      outbound.Tag,
		"protocol": outbound.Protocol,
	}
	if outbound.Settings != "" {
		var settings any
		if err := json.Unmarshal([]byte(outbound.Settings), &settings); err != nil {
			return nil, err
		}
		outboundMap["settings"] = settings
	}
	if outbound.StreamSettings != "" {
		var streamSettings any
		if err := json.Unmarshal([]byte(outbound.StreamSettings), &streamSettings); err != nil {
			return nil, err
		}
		outboundMap["streamSettings"] = streamSettings
	}
	if outbound.ProxySettings != "" {
		var proxySettings any
		if err := json.Unmarshal([]byte(outbound.ProxySettings), &proxySettings); err != nil {
			return nil, err
		}
		outboundMap["proxySettings"] = proxySettings
	}
	if outbound.Mux != "" {
		var mux any
		if err := json.Unmarshal([]byte(outbound.Mux), &mux); err != nil {
			return nil, err
		}
		outboundMap["mux"] = mux
	}
	return outboundMap, nil
}

// getDefaultXrayConfig retrieves the default Xray configuration.
func (a *XraySettingController) getDefaultXrayConfig(c *gin.Context) {
	defaultJsonConfig, err := a.SettingService.GetDefaultXrayConfig()
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.getSettings"), err)
		return
	}
	jsonObj(c, defaultJsonConfig, nil)
}

// getXrayResult retrieves the current Xray service result.
func (a *XraySettingController) getXrayResult(c *gin.Context) {
	jsonObj(c, a.XrayService.GetXrayResult(), nil)
}

// warp handles Warp-related operations based on the action parameter.
func (a *XraySettingController) warp(c *gin.Context) {
	action := c.Param("action")
	var resp string
	var err error
	switch action {
	case "data":
		resp, err = a.WarpService.GetWarpData()
	case "del":
		err = a.WarpService.DelWarpData()
	case "config":
		resp, err = a.WarpService.GetWarpConfig()
	case "reg":
		skey := c.PostForm("privateKey")
		pkey := c.PostForm("publicKey")
		resp, err = a.WarpService.RegWarp(skey, pkey)
	case "license":
		license := c.PostForm("license")
		resp, err = a.WarpService.SetWarpLicense(license)
	}

	jsonObj(c, resp, err)
}

// getOutboundsTraffic retrieves the traffic statistics for outbounds.
func (a *XraySettingController) getOutboundsTraffic(c *gin.Context) {
	outboundsTraffic, err := a.OutboundService.GetOutboundsTraffic()
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.getOutboundTrafficError"), err)
		return
	}
	jsonObj(c, outboundsTraffic, nil)
}

// resetOutboundsTraffic resets the traffic statistics for the specified outbound tag.
func (a *XraySettingController) resetOutboundsTraffic(c *gin.Context) {
	tag := c.PostForm("tag")
	err := a.OutboundService.ResetOutboundTraffic(tag)
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.resetOutboundTrafficError"), err)
		return
	}
	jsonObj(c, "", nil)
}

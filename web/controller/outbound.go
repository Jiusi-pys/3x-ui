package controller

import (
	"strconv"

	"github.com/mhsanaei/3x-ui/v2/database/model"
	"github.com/mhsanaei/3x-ui/v2/web/service"
	"github.com/mhsanaei/3x-ui/v2/web/session"

	"github.com/gin-gonic/gin"
)

// OutboundController handles HTTP requests related to Xray outbounds management.
type OutboundController struct {
	outboundService service.OutboundService
	xrayService     service.XrayService
}

// NewOutboundController creates a new OutboundController and sets up its routes.
func NewOutboundController(g *gin.RouterGroup) *OutboundController {
	a := &OutboundController{}
	a.initRouter(g)
	return a
}

// initRouter initializes the routes for outbound-related operations.
func (a *OutboundController) initRouter(g *gin.RouterGroup) {
	g.GET("/list", a.getOutbounds)
	g.GET("/get/:id", a.getOutbound)
	g.GET("/tags", a.getOutboundTags)

	g.POST("/add", a.addOutbound)
	g.POST("/del/:id", a.delOutbound)
	g.POST("/update/:id", a.updateOutbound)
}

// getOutbounds retrieves the list of outbounds for the logged-in user.
func (a *OutboundController) getOutbounds(c *gin.Context) {
	user := session.GetLoginUser(c)
	outbounds, err := a.outboundService.GetOutbounds(user.Id)
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.outbounds.toasts.obtain"), err)
		return
	}
	jsonObj(c, outbounds, nil)
}

// getOutbound retrieves a specific outbound by its ID.
func (a *OutboundController) getOutbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, I18nWeb(c, "get"), err)
		return
	}
	outbound, err := a.outboundService.GetOutbound(id)
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.outbounds.toasts.obtain"), err)
		return
	}
	jsonObj(c, outbound, nil)
}

// getOutboundTags retrieves all outbound tags.
func (a *OutboundController) getOutboundTags(c *gin.Context) {
	tags, err := a.outboundService.GetOutboundTags()
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.outbounds.toasts.obtain"), err)
		return
	}
	jsonObj(c, tags, nil)
}

// addOutbound creates a new outbound.
func (a *OutboundController) addOutbound(c *gin.Context) {
	outbound := &model.Outbound{}
	err := c.ShouldBind(outbound)
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.outbounds.toasts.obtain"), err)
		return
	}

	user := session.GetLoginUser(c)
	outbound.UserId = user.Id

	err = a.outboundService.AddOutbound(outbound)
	jsonMsgObj(c, I18nWeb(c, "pages.outbounds.addOutbound"), outbound, err)
	if err == nil {
		a.xrayService.SetToNeedRestart()
	}
}

// delOutbound deletes an outbound by ID.
func (a *OutboundController) delOutbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, I18nWeb(c, "delete"), err)
		return
	}

	err = a.outboundService.DelOutbound(id)
	jsonMsgObj(c, I18nWeb(c, "pages.outbounds.delOutbound"), id, err)
	if err == nil {
		a.xrayService.SetToNeedRestart()
	}
}

// updateOutbound updates an existing outbound by ID.
func (a *OutboundController) updateOutbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.outbounds.toasts.obtain"), err)
		return
	}

	outbound := &model.Outbound{
		Id: id,
	}
	err = c.ShouldBind(outbound)
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.outbounds.toasts.obtain"), err)
		return
	}

	err = a.outboundService.UpdateOutbound(outbound)
	jsonMsgObj(c, I18nWeb(c, "pages.outbounds.updateOutbound"), outbound, err)
	if err == nil {
		a.xrayService.SetToNeedRestart()
	}
}

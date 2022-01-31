package controllers

import (
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"merchant/common"
	"merchant/models"
	"strings"
	"time"
)

type BaseController struct {
	beego.Controller
	MemberId   string
	BodyObject interface{}
}

// Prepare is checking the authorization
func (o *BaseController) Prepare() {
	if o.Ctx.Request.Header.Get("Authorization") != "" {
		acecessToken := o.Ctx.Input.Header("Authorization")
		v, err := models.GetAuthenticationByToken(acecessToken)
		if err == nil && v.ExpiryTime > time.Now().Unix() {
			o.MemberId = v.MemberId
			o.CheckAuthorization()
			return
		}
	}
	o.Ctx.Output.SetStatus(401)
	o.Data["json"] = common.CreateErrorMessageWithCode(common.ErrorCodeUnauthrized)
	o.ServeJSON()
	return
}

// CheckAuthorization is checking the authorization
func (o *BaseController) CheckAuthorization() {
	var isAuthorized = true
	var err error
	idStr := o.Ctx.Input.Param(":id")
	member, err := models.GetMemberById(o.MemberId)
	if err != nil {
		o.Ctx.Output.SetStatus(401)
		o.Data["json"] = common.CreateErrorMessageWithCode(common.ErrorCodeUnauthrized)
		o.ServeJSON()
	}
	var updatedMember *models.Member
	if strings.Contains(o.Ctx.Request.URL.Path, "/member") {
		var v models.Member
		if o.Ctx.Request.Method == "PUT" || o.Ctx.Request.Method == "POST" {
			if err = json.Unmarshal(o.Ctx.Input.RequestBody, &v); err == nil {
				if idStr != "" {
					v.Id = idStr
				}
				if v.Merchant == nil {
					v.Merchant = member.Merchant
				}
				o.BodyObject = v
			}
		}
		if err == nil && idStr != "" {
			updatedMember, err = models.GetMemberById(idStr)
		}
	} else if strings.Contains(o.Ctx.Request.URL.Path, "/merchant") {
		var v models.Merchant
		if o.Ctx.Request.Method == "PUT" || o.Ctx.Request.Method == "POST" {
			if err = json.Unmarshal(o.Ctx.Input.RequestBody, &v); err == nil {
				if idStr != "" {
					v.Id = idStr
				}
				o.BodyObject = v
			}
		}
	}
	if err != nil {
		o.Ctx.Output.SetStatus(400)
		o.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUndefined, err.Error())
		o.ServeJSON()
		return
	}
	keyPath := o.Ctx.Request.Method + o.Ctx.Request.URL.Path
	switch {
	case strings.HasPrefix(keyPath, "GET/v1/member/"), strings.HasPrefix(keyPath, "PUT/v1/member/"):
		var v models.Member
		if o.BodyObject != nil {
			v = o.BodyObject.(models.Member)
		}
		if idStr != "" {
			if (v.Merchant != nil && v.Merchant.Id != member.Merchant.Id) ||
				member.Role == common.RoleUser && idStr != o.MemberId ||
				(member.Role == common.RoleAdministrator && updatedMember.Merchant.Id != member.Merchant.Id) {
				isAuthorized = false
			}
		}
	case strings.HasPrefix(keyPath, "DELETE/v1/member/"):
		if member.Role == common.RoleUser ||
			(member.Role == common.RoleAdministrator && updatedMember.Merchant.Id != member.Merchant.Id) {
			isAuthorized = false
		}
	case strings.HasPrefix(keyPath, "POST/v1/member/"):
		var v = o.BodyObject.(models.Member)
		if member.Role == common.RoleUser ||
			(member.Role == common.RoleAdministrator && v.Merchant != nil && v.Merchant.Id != member.Merchant.Id) {
			isAuthorized = false
		}
	case strings.HasPrefix(keyPath, "GET/v1/merchant/"):
		if member.Role != common.RoleSuperAdmin && member.Merchant.Id != idStr {
			isAuthorized = false
		}
	case strings.HasPrefix(keyPath, "PUT/v1/merchant/"), strings.HasPrefix(keyPath, "DELETE/v1/merchant/"):
		if member.Role == common.RoleUser ||
			(member.Role == common.RoleAdministrator && member.Merchant.Id != idStr) {
			isAuthorized = false
		}
	case strings.HasPrefix(keyPath, "POST/v1/merchant/"):
		if member.Role != common.RoleSuperAdmin {
			isAuthorized = false
		}
	}
	if !isAuthorized {
		o.Ctx.Output.SetStatus(401)
		o.Data["json"] = common.CreateErrorMessageWithCode(common.ErrorCodeUnauthrized)
		o.ServeJSON()
	}
	return
}

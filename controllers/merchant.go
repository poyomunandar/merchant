package controllers

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"pace_merchant/common"
	"pace_merchant/models"
)

// MerchantController operations for Merchant
type MerchantController struct {
	BaseController
}

// Prepare ...
func (c *MerchantController) Prepare() {
	c.BaseController.Prepare()
}

// URLMapping ...
func (c *MerchantController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Merchant, only member with role superadmin can do this. It will create the merchant along with default administrator member for the new merchant with email administator@[merchantCode].com and default password Merchant!234
// @Param	Authorization	header	string	false	"Access Token, for example: 27dd4391dfc6428496e198150c0f6e1e ..."
// @Param	body		body 	common.MerchantRequest	true		"body for Merchant content"
// @Success 201 {object} models.Merchant
// @Failure 401 unauthorized
// @Failure 400 bad request
// @router / [post]
func (c *MerchantController) Post() {
	var v = c.BodyObject.(models.Merchant)
	var m models.Member
	if _, err := models.AddMerchant(&v); err == nil {
		m.Role = common.RoleAdministrator
		m.EmailAddress = fmt.Sprintf("%s@%s.com", common.RoleAdministrator, v.MerchantCode)
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(common.DefaultMerchantPassword), bcrypt.DefaultCost)
		m.Password = string(hashedPassword)
		m.Merchant = &models.Merchant{MerchantCode: v.MerchantCode}
		if _, err := models.AddMember(&m); err != nil {
			models.DeleteMerchant(v.MerchantCode)
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUndefined, err.Error())
			c.ServeJSON()
			return
		}
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = v
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUndefined, err.Error())
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Merchant by id, can only view the merchant of the member. Superadmin can view any merchants
// @Param	Authorization	header	string	false	"Access Token, for example: 27dd4391dfc6428496e198150c0f6e1e ..."
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Merchant
// @Failure 400 bad request
// @router /:id [get]
func (c *MerchantController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	v, err := models.GetMerchantById(idStr)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUndefined, err.Error())
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Merchant, only role administrator can do this and only can do to its merchant. Superadmin can update any merchants
// @Param	Authorization	header	string	false	"Access Token, for example: 27dd4391dfc6428496e198150c0f6e1e ..."
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	common.MerchantRequest	true		"body for Merchant content"
// @Success 200 {object} common.ErrorMessage
// @Failure 401 unauthorized
// @Failure 400 bad request
// @router /:id [put]
func (c *MerchantController) Put() {
	var v = c.BodyObject.(models.Merchant)
	if err := models.UpdateMerchantById(&v); err == nil {
		c.Data["json"] = common.CreateSuccessMessage()
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUndefined, err.Error())
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Merchant, all members belong to the merchant will be deleted as well. It's soft delete. Only administrator can delete their merchant. Superadmin can delete any merchants
// @Param	Authorization	header	string	false	"Access Token, for example: 27dd4391dfc6428496e198150c0f6e1e ..."
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} common.ErrorMessage
// @Failure 401 unauthorized
// @Failure 400 bad request
// @router /:id [delete]
func (c *MerchantController) Delete() {
	member, _ := models.GetMemberById(c.MemberId)
	members, err := models.GetMemberByMerchantId(member.Merchant.MerchantCode)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUndefined, err.Error())
		c.ServeJSON()
		return
	}
	// delete all members
	for _, item := range *members {
		item.IsDeleted = 1
		err = models.UpdateMemberById(&item)
	}
	if err := models.DeleteMerchant(member.Merchant.MerchantCode); err == nil {
		c.Data["json"] = common.CreateSuccessMessage()
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUndefined, err.Error())
	}
	c.ServeJSON()
}

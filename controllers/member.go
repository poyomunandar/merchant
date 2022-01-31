package controllers

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"merchant/common"
	"merchant/models"
	"strings"
)

// MemberController operations for Member
type MemberController struct {
	BaseController
}

// Prepare ...
func (c *MemberController) Prepare() {
	c.BaseController.Prepare()
}

// URLMapping ...
func (c *MemberController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Member only for role administrator, and can only create member under administrator's merchant (no need to define merchant in the payload request). Superadmin can create the member under any merchants
// @Param	Authorization	header	string	false	"Access Token, for example: 27dd4391dfc6428496e198150c0f6e1e ..."
// @Param	body		body 	common.MemberRequest	true		"body for Member content"
// @Success 201 {object} models.Member
// @Failure 401 unauthorized
// @Failure 400 bad request
// @router / [post]
func (c *MemberController) Post() {
	member, _ := models.GetMemberById(c.MemberId)
	var v = c.BodyObject.(models.Member)
	if v.EmailAddress == "" || v.Password == "" || v.Role == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = common.CreateErrorMessageWithCode(common.ErrorCodeRequiredParamEmpty, "Password", "EmailAddress", "Role")
		c.ServeJSON()
		return
	}
	if !common.IsEmailValid(v.EmailAddress) {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = common.CreateErrorMessageWithCode(common.ErrorEmailFormatInvalid)
		c.ServeJSON()
		return
	}
	if !common.MapRoles[v.Role] {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = common.CreateErrorMessageWithCode(common.ErrorRoleNotAvailable)
		c.ServeJSON()
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(v.Password), bcrypt.DefaultCost)
	v.Password = string(hashedPassword)
	if v.Merchant == nil {
		v.Merchant = &models.Merchant{Id: member.Merchant.Id}
	}
	if _, err := models.AddMember(&v); err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = v
	} else {
		if strings.Contains(err.Error(), "member_email_address_uindex") {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = common.CreateErrorMessageWithCode(common.ErrorCodeEmailExist)
			c.ServeJSON()
			return
		} else if strings.Contains(err.Error(), "member_merchant_id_fk") {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = common.CreateErrorMessageWithCode(common.ErrorMerchantIdNotExist)
			c.ServeJSON()
			return
		} else {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUndefined, err.Error())
		}
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Member by id, for role user can only see their self, for role administrator can see anyone in the same merchant. Superadmin can view any members
// @Param	Authorization	header	string	false	"Access Token, for example: 27dd4391dfc6428496e198150c0f6e1e ..."
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Member
// @Failure 401 unauthorized
// @Failure 400 bad request
// @router /:id [get]
func (c *MemberController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	updatedMember, err := models.GetMemberById(idStr)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUndefined, err.Error())
		c.ServeJSON()
		return
	}
	c.Data["json"] = updatedMember
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Member, for role user can only see their self, for role administrator can see all members in the same merchant. Superadmin can view any members
// @Param	Authorization	header	string	false	"Access Token, for example: 27dd4391dfc6428496e198150c0f6e1e ..."
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 - all columns are in lowercase"
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} []models.Member
// @Failure 401 unauthorized
// @Failure 400 bad request
// @router / [get]
func (c *MemberController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Ctx.Output.SetStatus(400)
				err := errors.New("Error: invalid query key/value pair")
				c.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUndefined, err.Error())
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}
	member, _ := models.GetMemberById(c.MemberId)
	if member.Role == common.RoleUser {
		query["id"] = c.MemberId
	} else if member.Role == common.RoleAdministrator {
		query["Merchant.Id"] = member.Merchant.Id
	}

	l, err := models.GetAllMember(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUndefined, err.Error())
	} else {
		if l == nil {
			c.Data["json"] = []models.Member{}
		} else {
			c.Data["json"] = l
		}
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Member, for role user, can only update their self, for role administrator can update only the member with the same merchant. Superadmin can update any members
// @Param	Authorization	header	string	false	"Access Token, for example: 27dd4391dfc6428496e198150c0f6e1e ..."
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	common.MemberRequest	true		"body for Member content"
// @Success 200 {object} common.ErrorMessage
// @Failure 401 unauthorized
// @Failure 400 bad request
// @router /:id [put]
func (c *MemberController) Put() {
	v := c.BodyObject.(models.Member)
	if !common.IsEmailValid(v.EmailAddress) {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = common.CreateErrorMessageWithCode(common.ErrorEmailFormatInvalid)
		c.ServeJSON()
		return
	}
	if err := models.UpdateMemberById(&v); err == nil {
		c.Data["json"] = common.CreateSuccessMessage()
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUndefined, err.Error())
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Member, only administrator can do this for the member of the same merchant. Superadmin can delete any members
// @Param	Authorization	header	string	false	"Access Token, for example: 27dd4391dfc6428496e198150c0f6e1e ..."
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {object} common.ErrorMessage
// @Failure 401 unauthorized
// @Failure 400 bad request
// @router /:id [delete]
func (c *MemberController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	if err := models.DeleteMember(idStr); err == nil {
		c.Data["json"] = common.CreateSuccessMessage()
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUndefined, err.Error())
	}
	c.ServeJSON()
}

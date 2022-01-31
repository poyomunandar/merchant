package controllers

import (
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
	"merchant/common"
	"merchant/models"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AuthenticationController operations for Authentication
type AuthenticationController struct {
	beego.Controller
}

// URLMapping ...
func (c *AuthenticationController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// Post ...
// @Title Post
// @Description Authenticate credential to get acess token
// @Param	body		body 	common.AuthenticationRequest	true		"body for Authentication content"
// @Success 201 {object} models.Authentication
// @Failure 401 Unauthorized
// @Failure 400 Bad Request
// @router / [post]
func (c *AuthenticationController) Post() {
	var v common.AuthenticationRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if v.EmailAddress == "" || v.Password == "" {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = common.CreateErrorMessageWithCode(common.ErrorCodeRequiredParamEmpty, "password", "emailAddress")
			c.ServeJSON()
			return
		}
		member, err := models.GetMemberByEmail(v.EmailAddress)
		if err != nil {
			c.Ctx.Output.SetStatus(401)
			c.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUnauthrized, err.Error())
			c.ServeJSON()
			return
		}
		if err = bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(v.Password)); err != nil {
			c.Ctx.Output.SetStatus(401)
			c.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUnauthrized, err.Error())
			c.ServeJSON()
			return
		}
		token := strings.ReplaceAll(uuid.New().String(), "-", "")
		authentication := models.Authentication{
			ExpiryTime: time.Now().Add(30 * time.Minute).Unix(),
			Token:      token,
			MemberId:   member.Id,
		}
		if _, err := models.AddAuthentication(&authentication); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = authentication
		} else {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUndefined, err.Error())
		}
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = common.CreateErrorMessage(common.ErrorCodeUndefined, err.Error())
	}
	c.ServeJSON()
}

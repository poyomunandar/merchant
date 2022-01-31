// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"merchant/controllers"
)

func init() {
	beego.Router("/v1/authentication", &controllers.AuthenticationController{}, "post:Post")
	beego.Router("/v1/member", &controllers.MemberController{}, "get:GetAll")
	beego.Router("/v1/member/:id", &controllers.MemberController{}, "get:GetOne")
	beego.Router("/v1/member/", &controllers.MemberController{}, "post:Post")
	beego.Router("/v1/member/:id", &controllers.MemberController{}, "put:Put")
	beego.Router("/v1/member/:id", &controllers.MemberController{}, "delete:Delete")
	beego.Router("/v1/merchant/:id", &controllers.MerchantController{}, "get:GetOne")
	beego.Router("/v1/merchant/", &controllers.MerchantController{}, "post:Post")
	beego.Router("/v1/merchant/:id", &controllers.MerchantController{}, "put:Put")
	beego.Router("/v1/merchant/:id", &controllers.MerchantController{}, "delete:Delete")
}

package routers

import (
	"merchant/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/authentication",
			beego.NSInclude(
				&controllers.AuthenticationController{},
			),
		),

		beego.NSNamespace("/member",
			beego.NSInclude(
				&controllers.MemberController{},
			),
		),

		beego.NSNamespace("/merchant",
			beego.NSInclude(
				&controllers.MerchantController{},
			),
		),
	)
	beego.AddNamespace(ns)
}

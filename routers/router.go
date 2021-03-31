// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	// topic "gobee/controllers/v1/topic"
	controllers "gobee/controllers/v1/topic"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	/* ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/object",
			beego.NSInclude(
				&controllers.ObjectController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	) */

	ns := web.NewNamespace("/v1",
		web.NSNamespace("/topic",
			web.NSInclude(
				&controllers.TopicController{},
			),
			// beego.NSRouter("/GetTopic", &controllers.TopicController{}),
		),
	)
	web.AddNamespace(ns)
}

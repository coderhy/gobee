package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["gobee/controllers:TopicController"] = append(beego.GlobalControllerRouter["gobee/controllers:TopicController"],
        beego.ControllerComments{
            Method: "GetTopic",
            Router: "/GetTopic",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

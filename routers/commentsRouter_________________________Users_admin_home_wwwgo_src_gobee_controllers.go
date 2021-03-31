package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["gobee/controllers/v1/topic:TopicController"] = append(beego.GlobalControllerRouter["gobee/controllers/v1/topic:TopicController"],
        beego.ControllerComments{
            Method: "GetTopic",
            Router: "/GetTopic",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gobee/controllers/v1/topic:TopicController"] = append(beego.GlobalControllerRouter["gobee/controllers/v1/topic:TopicController"],
        beego.ControllerComments{
            Method: "GetTopicAll",
            Router: "/GetTopicAll",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["gobee/controllers/v1/topic:TopicController"] = append(beego.GlobalControllerRouter["gobee/controllers/v1/topic:TopicController"],
        beego.ControllerComments{
            Method: "GetTopicPanic",
            Router: "/GetTopicPanic",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

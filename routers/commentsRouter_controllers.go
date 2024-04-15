package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/agenda_personal_parametros_api_mid/controllers:MainController"] = append(beego.GlobalControllerRouter["github.com/udistrital/agenda_personal_parametros_api_mid/controllers:MainController"],
        beego.ControllerComments{
            Method: "Get",
            Router: "/get/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

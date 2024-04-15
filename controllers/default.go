package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

// MainController ...
// @Title Get
// @Description create Get
// @Param id query string true "id"
// @Success 201 {object} []models.Default
// @Failure 403 body is empty
// @router /get/:id [get]
func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

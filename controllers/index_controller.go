/*
 * @Description:
 * @Author: gphper
 * @Date: 2021-09-19 20:47:35
 */
package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"goredismanager/comment"
	"goredismanager/global"
	"goredismanager/model"
	"goredismanager/service"

	"github.com/gin-gonic/gin"
)

type indexController struct {
	BaseController
}

var Ic = indexController{}

func (con *indexController) Show(c *gin.Context) {

	dbSlice := make([]struct{}, 16)
	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"services":     global.RedisServiceStorage,
		"service_name": global.UseClient.ConnectName,
		"db":           global.UseClient.Db,
		"db_slice":     dbSlice,
	})

}

func (con *indexController) SaveConfig(c *gin.Context) {

	var serviceReq model.ServiceConfigReq
	var err error

	err = con.FormBind(c, &serviceReq)
	if err != nil {
		con.Error(c, err.Error())
	}

	err = service.AddServiceConf(serviceReq)
	if err != nil {
		con.Error(c, err.Error())
		return
	}
	con.Success(c, "", "添加连接成功")
}

func (con *indexController) Switch(c *gin.Context) {

	var (
		err    error
		serReq model.ServiceSwitchReq
	)

	err = con.FormBind(c, &serReq)
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	err = service.ServiceSwitch(serReq)
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	con.Success(c, "", "切换连接成功")
}

func (con *indexController) SearchKey(c *gin.Context) {

	var (
		err    error
		req    model.RedisKeyReq
		result []model.RedisKey
	)

	err = con.FormBind(c, &req)
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	result, _ = service.SearchKeyType(req)

	gen := comment.NewTrie()

	for _, v := range result {

		if v.Type == "none" {
			con.AjaxReturn(c, AJAXFAIL, "Key值不存在")
			return
		}

		v.Key += "_" + v.Type

		stringSlice := strings.Split(v.Key, ":")

		gen.Insert(stringSlice)

	}

	resultSlice := comment.GetOne(gen.Root, "")

	con.AjaxReturn(c, AJAXSUCCESS, resultSlice)
}

func (con *indexController) AddKey(c *gin.Context) {

	var (
		req model.AddKeyReq
		err error
	)

	err = con.FormBind(c, &req)

	fmt.Printf("%+v", req)

	if err != nil {
		con.Error(c, "添加失败")
	}

	err = service.AddKey(req)
	if err != nil {
		con.Error(c, "添加失败")
	}

	con.Success(c, "", "添加成功")
}

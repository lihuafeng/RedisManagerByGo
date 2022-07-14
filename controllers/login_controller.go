/*
 * @Description:
 * @Author: gphper
 * @Date: 2021-11-07 17:20:54
 */
package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goredismanager/global"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/dchest/captcha"
)

type loginController struct {
	BaseController
}

var LoginC = loginController{}

func (con *loginController) ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login/login.html", gin.H{})
}

func (con *loginController) Login(c *gin.Context) {

	username := c.PostForm("account")
	password := c.PostForm("password")
	code := c.PostForm("code")

	if !con.CaptchaVerify(c, code) {
		con.Error(c, "验证码错误")
		return
	}
	sql_query, _ := global.Db.Prepare("select * from admin_user where user_name=?")
	user_info, err := sql_query.Query(username)
	user_info.Scan()
	fmt.Print(user_info)
	fmt.Print(err)
	if _, ok := global.Accounts[username]; ok {
		if global.Accounts[username] == password {
			userInfo := make(map[string]interface{})
			userInfo["username"] = username
			//session 存储数据
			session := sessions.Default(c)
			userstr, _ := json.Marshal(userInfo)

			session.Set("userInfo", string(userstr))
			session.Save()

			con.Success(c, "/index", "登录成功")
		} else {
			con.Error(c, "账号密码错误")
		}
	} else {
		con.Error(c, "账号密码错误")
		return
	}

}

func (con *loginController) LoginOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("userInfo")
	session.Save()
	con.Success(c, "login", "退出成功")
}

func (con *loginController) Captcha(c *gin.Context) {
	w, h := 107, 36
	captchaId := captcha.NewLen(3)
	session := sessions.Default(c)
	session.Set("captcha", captchaId)
	_ = session.Save()
	_ = Serve(c.Writer, c.Request, captchaId, ".png", "zh", false, w, h)
}
func (con *loginController) CaptchaVerify(c *gin.Context, code string) bool {
	session := sessions.Default(c)
	if captchaId := session.Get("captcha"); captchaId != nil {
		session.Delete("captcha")
		_ = session.Save()
		if captcha.VerifyString(captchaId.(string), code) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func Serve(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool, width, height int) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		_ = captcha.WriteImage(&content, id, width, height)
	case ".wav":
		w.Header().Set("Content-Type", "audio/x-wav")
		_ = captcha.WriteAudio(&content, id, lang)
	default:
		return captcha.ErrNotFound
	}

	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}

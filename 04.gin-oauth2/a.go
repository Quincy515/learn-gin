package main

import (
	"fmt"
	"gin-oauth2/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

const authServerURL = "http://127.0.0.1:8081"

var (
	oauth2Config = oauth2.Config{
		ClientID:     "clienta",
		ClientSecret: "123",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:8080/github/getcode",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/auth",  // 获取授权码地址
			TokenURL: authServerURL + "/token", // 获取 token 地址
		},
	}
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("public/*")
	//codeUrl, _ := url.ParseRequestURI("http://localhost:8080/getcode")
	//loginUrl := "http://127.0.0.1:8081/auth?" +
	//	"response_type=code&client_id=clienta&redirect_uri=" +
	//	codeUrl.String()
	// state 参数，传递给服务端，验证通过会原封不动的传回客户端，
	// 在 /getcode 里理论上需要对 state 进行判断，防止被串改。
	loginUrl := oauth2Config.AuthCodeURL("myclient")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "a-index.html", map[string]string{
			"loginUrl": loginUrl,
		})
	})
	// 注册用户
	r.Any("/user/reg", func(c *gin.Context) {
		data := map[string]string{
			"error":   "",
			"message": "",
		}
		if c.Request.Method == "POST" {
			userID := "" // 第三方账号 ID
			if c.Query("token") != "" {
				passport := utils.GetUserInfo(authServerURL+"/info", c.Query("token"), true)
				userID = passport.UserID
			}
			source := c.Query("source")
			uname, upass, upass2 := c.PostForm("userName"), c.PostForm("userPass"), c.PostForm("userPass2")
			user, err := utils.AddNewUser(uname, upass, upass2, userID, source)
			if err != nil {
				data["error"] = err.Error()
			} else {
				if userID != "" {
					data["message"] = fmt.Sprintf("账号绑定成功,您的用户名是%s,第三方账号是:%s", user.UserName, userID)
				} else {
					data["message"] = fmt.Sprintf("账号创建成功,您的用户名是%s", user.UserName)
				}
			}
		}
		c.HTML(200, "reg.html", data)
	})
	r.GET("/github/getcode", func(c *gin.Context) {
		//source := c.Param("source")                  // 来源
		source := "github"
		code, _ := c.GetQuery("code")                // 得到授权码
		token, err := oauth2Config.Exchange(c, code) // 请求 token
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
		} else {
			//c.JSON(200, token)
			// 通过第三方登录平台提交 token，获取用户在第三方平台的用户 ID
			passport := utils.GetUserInfo(authServerURL+"/info", token.AccessToken, true)
			user := utils.GetUserName(source, passport.UserID)
			if user == nil {
				//c.String(200, "您需要注册并绑定账号")
				c.Redirect(302, "/user/reg?token="+token.AccessToken+"&source=github")
			} else {
				c.String(200, "您在本站的用户名是: %s", user.UserName)
			}
		}
	})
	//r.GET("/info", func(c *gin.Context) {
	//	token := c.Query("token")
	//	ret := utils.GetUserInfo(authServerURL+"/info", token, true)
	//	c.Writer.Header().Add("Content-type", "application/json")
	//	c.String(200, ret)
	//})
	r.Run(":8080")
}

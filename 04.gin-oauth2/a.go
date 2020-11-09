package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

const authServerURL = "http://127.0.0.1:8081"

var (
	oauth2Config = oauth2.Config{
		ClientID:     "clienta",
		ClientSecret: "123",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:8080/getcode",
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
	r.GET("/getcode", func(c *gin.Context) {
		code, _ := c.GetQuery("code")
		//c.JSON(200, gin.H{"code": code})
		token, err := oauth2Config.Exchange(c, code)
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
		} else {
			c.JSON(200, token)
		}
	})
	r.Run(":8080")
}

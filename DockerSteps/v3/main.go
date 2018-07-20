package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
	"github.com/appleboy/gin-jwt"
	"net/http"
	"github.com/spf13/viper"
)

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	c.JSON(200, gin.H{
		"userID": claims["id"],
		"text":   "Hello World.",
	})
}

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

func main() {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath("./Configs")
	viper.AddConfigPath("Images/dockerfile-learn/Configs")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Can't find config")
	}

	port := viper.GetString("port")
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if port == "" {
		port = "8000"
	}

	adminID := viper.GetString("admin.name")
	adminFirst := viper.GetString("admin.first-name")
	adminLast := viper.GetString("admin.last-name")
	adminPwd := viper.GetString("admin.password")
	log.Info(adminID, adminFirst, adminLast, adminPwd)
	// the jwt middleware
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "docker test",
		Key:        []byte(viper.GetString("secret-key")),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(userId string, password string, c *gin.Context) (interface{}, bool) {
			if (userId == adminID && password == adminPwd) || (userId == "test" && password == "test") {
				return &User{
					UserName:  adminID,
					LastName:  adminLast,
					FirstName: adminFirst,
				}, true
			}

			return nil, false
		},
		Authorizator: func(user interface{}, c *gin.Context) bool {
			if v, ok := user.(string); ok && v == adminID {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header:Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}

	r.POST("/login", authMiddleware.LoginHandler)

	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", helloHandler)
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
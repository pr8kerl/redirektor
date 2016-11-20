package main

import (
	"fmt"
	"github.com/pr8kerl/redirektor"
	"gopkg.in/gin-gonic/gin.v1"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	currentUser          string = "luser"
	currentUserFirstName string = "luser"
	basename             string
	cfgfile              string
)

func init() {

	basename = filepath.Base(os.Args[0])
	cfgfile = filepath.Abs(os.Args[0]) + ".ini"

	// setup config
	var cfg = Config{}
	err := cfg.LoadFromFile(cfgfile)
	if err != nil {
		fmt.Printf("error reading config: %s\n", err)
		os.Exit(1)
	}
	bindaddress = fmt.Sprintf("%s:%d", cfg.Server.BindAddr, cfg.Server.BindPort)
}

func main() {

	r, err := redirektor.NewRedirektor(cfg)
	if err != nil {
		fmt.Printf("error initialising redirektor: %s\n", err)
		os.Exit(1)
	}

	// Creates a router without any middleware by default
	//gin.SetMode(gin.ReleaseMode)
	s := gin.New()
	// Global middlewares
	s.Use(MyLogger(gin.DefaultWriter))
	s.Use(gin.Recovery())
	s.Use(SetJellyBeans())
	//r.Use(GetUser)

	//r.GET("/", index)
	s.StaticFS(appRoot, http.Dir(cfg.Server.HtmlPath))

	api := s.Group("/api")
	{
		api.GET("/dbs", r.GetDBs)
		api.GET("/redirekts", r.GetAll)
		api.GET("/redirekts/:dbname", r.Get)
		api.PUT("/redirekts/:dbname", r.Set)
	}

	s.Run(bindaddress)

}

func index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": 200, "message": "hello"})
}

func SetJellyBeans() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Powered-By", "Black Jelly Beans")
		c.Next()
	}
}

func GetUser(c *gin.Context) {
	luser := c.Request.Header.Get("X-Remote-User")
	if len(luser) > 0 {
		currentUser = luser
		fmt.Printf("current user: %s\n", currentUser)
	}
	c.Next()
}

func MyLogger(out io.Writer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		fmt.Fprintf(out, "[GIN] %v | %3d | %13v | %s | %s | %-7s %s\n%s",
			end.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			currentUser,
			method,
			path,
			comment,
		)
	}
}

// @contact.name API Support
// @contact.url http://github.com/axiaoxin-com/pink-lady
// @contact.email 254606826@qq.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name apikey

package routes

import (
	"net/http"

	"github.com/axiaoxin-com/logging"
	"github.com/axiaoxin-com/pink-lady/routes/docs"
	"github.com/axiaoxin-com/pink-lady/statics"
	"github.com/axiaoxin-com/pink-lady/webserver"
	"github.com/gin-contrib/pprof"

	// docs is generated by Swag CLI, you have to import it.
	_ "github.com/axiaoxin-com/pink-lady/routes/docs"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	// DisableGinSwaggerEnvkey 设置该环境变量时关闭 swagger 文档
	DisableGinSwaggerEnvkey = "DISABLE_GIN_SWAGGER"
)

// Register 在 gin engine 上注册 url 对应的 HandlerFunc
func Register(httpHandler http.Handler) {
	app, ok := httpHandler.(*gin.Engine)
	if !ok {
		panic("HTTP handler must be *gin.Engine")
	}

	// api 文档变量设置，注意这里依赖 viper 读配置，需要保证在 main 中已预先加载这些配置项
	docs.SwaggerInfo.Title = viper.GetString("apidocs.title")
	docs.SwaggerInfo.Description = viper.GetString("apidocs.desc")
	docs.SwaggerInfo.Version = Version
	docs.SwaggerInfo.Host = viper.GetString("apidocs.host")
	docs.SwaggerInfo.BasePath = viper.GetString("apidocs.basepath")
	docs.SwaggerInfo.Schemes = viper.GetStringSlice("apidocs.schemes")

	// Group x 默认 url 路由
	x := app.Group("/x", webserver.GinBasicAuth())
	{
		if viper.GetBool("server.pprof") {
			pprof.RouteRegister(x, "/pprof")
		}
		if viper.GetBool("server.metrics") {
			x.GET("/metrics", webserver.PromExporterHandler())
		}
		// ginSwagger 生成的在线 API 文档路由
		x.GET("/apidocs/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, DisableGinSwaggerEnvkey))
		// 默认的 ping 方法，返回 server 相关信息
		x.Any("/ping", Ping)
	}
	// 注册 favicon.ico 和 robots.txt
	app.GET("/favicon.ico", func(c *gin.Context) {
		file, err := statics.Files.ReadFile("favicon.ico")
		if err != nil {
			logging.Error(c, "read favicon file error:"+err.Error())
		}
		c.Data(http.StatusOK, "image/x-icon", file)
		return
	})

	app.GET("/robots.txt", func(c *gin.Context) {
		file, err := statics.Files.ReadFile("robots.txt")
		if err != nil {
			logging.Error(c, "read robots file error:"+err.Error())
		}
		c.Data(http.StatusOK, "text/plain", file)
		return
	})

	// 注册其他 gin HandlerFunc
	Routes(app)
}

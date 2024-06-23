package routes

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Unleash/unleash-client-go/v3"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/imedcl/manager-api/docs"
	"github.com/imedcl/manager-api/pkg/config"
	"github.com/imedcl/manager-api/pkg/data"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	db             *data.DB
	client         *unleash.Client
	connectionAuth data.ConnectionAuth
	currentUser    *data.User
)

func setTraceID(ctx *gin.Context) {
	if ctx.GetHeader("trace-ID") != "" {
		ctx.Set("trace-ID", ctx.GetHeader("trace-ID"))
	} else {
		id := uuid.New()
		ctx.Set("trace-ID", id)
	}
	ctx.Next()
}

func logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.Keys["trace-ID"],
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

func Create(cfg *config.Config, clt *unleash.Client) {
	client = clt
	// force colors in terminal
	gin.ForceConsoleColor()

	// Register a route.
	router := gin.New()

	// CORS Configuration
	origin := strings.TrimSuffix(cfg.AppUrl(), "/")
	if gin.Mode() == "debug" {
		origin = "*"
	}
	// Temporal old domain accept: Remove when migration is finished
	oldDomain := "https://autentia-admin-dev.autentia.io"

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{origin, oldDomain},
		AllowMethods:     []string{"PUT", "OPTIONS", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "accept", "origin", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add Trace ID to context
	router.Use(setTraceID)

	// Logger
	router.Use(logger())

	// Recovery from server error
	router.Use(gin.Recovery())

	// Migrations of DB
	migrateDB(cfg)
	logrus.Println("Migration complete!")
	// Add health check route
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	logrus.Println("Healthz registered!")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public Routes
	LoginRoute(router, db, client)
	LoginGuestRoute(router, db, client)
	RecoveryRoute(router, db, client)
	ReconfirmationRoute(router, db, client)
	ConfirmRoute(router, db, client)
	ValidateRoute(router, db, client)
	ActivateRoute(router, db, client)

	router.Static("/upload", "./upload")

	// Private Routes
	authorized := router.Group("/", VerifyToken)

	//user
	UsersRolesRoute(authorized, db, client)
	UserPostRoute(authorized, db, client)
	UserInstrospectionRoute(authorized, db, client)
	UserGetRoute(authorized, db, client)
	UserPutRoute(authorized, db, client)
	UserFeedbackRoute(authorized, db, client)
	LogoutRoute(authorized, db, client)
	DeleteUserRoute(authorized, db, client)

	//module
	ModulesGetRoute(authorized, db, client)
	ModulesDeleteRoute(authorized, db, client)

	//log
	LogGetRoute(authorized, db, client)
	LogPostRoute(authorized, db, client)

	logrus.Println("Autentia routes!")

	// Start the server
	logrus.Println("Start server!")
	if err := router.Run(cfg.Port()); err != nil {
		logrus.Fatalf("%v", err)
	}
}

func migrateDB(cfg *config.Config) {
	connectionAuth.Database = cfg.DbName()
	connectionAuth.UserName = cfg.DbUserName()
	connectionAuth.Password = cfg.DbPassword()
	connectionAuth.Port = cfg.DbPort()
	connectionAuth.Host = cfg.DbHost()
	connectionAuth.SSL = cfg.DbSSL()
	connectionAuth.SSLCa = cfg.DbSSLCa()
	connectionAuth.SSLCert = cfg.DbSSLCert()
	connectionAuth.SSLKey = cfg.DbSSLKey()
	connectionAuth.TimeZone = cfg.DbTimeZone()
	fmt.Println(cfg)
	db, _ = connectionAuth.Connect()
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	err := data.Migrate(*db)
	if err != nil {
		logrus.Print("Error al migrar", err.Error())
	}
	InitUserSuperAdminManager()

}

type Claims struct {
	User *data.User `json:"user"`
	jwt.StandardClaims
}

type MessageResponse struct {
	Code    int    `json:"code"`
	Details string `json:"details"`
}

type PasswordMessageResponse struct {
	Code         int                       `json:"code"`
	Details      string                    `json:"details"`
	Requirements config.PasswordValidation `json:"requirements"`
}

func VerifyToken(c *gin.Context) {
	currentUser = &data.User{}
	r := c.Request
	cfg := config.New()
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		c.JSON(http.StatusBadRequest, MessageResponse{
			Details: "Usuario no autorizado",
			Code:    http.StatusBadRequest,
		})
		c.Abort()
		return
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) <= 1 {
		c.JSON(http.StatusBadRequest, MessageResponse{
			Details: "Usuario no autorizado",
			Code:    http.StatusBadRequest,
		})
		c.Abort()
		return
	}
	reqToken = splitToken[1]

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(reqToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.SignKey()), nil
	})
	if err == nil && token.Valid {
		currentUser = claims.User
		if db.IsActive(currentUser.ID) {
			var userData, _ = db.GetUserByID(currentUser.ID)
			if userData.LastTokenSession != reqToken {
				c.JSON(http.StatusUnauthorized, MessageResponse{
					Details: "Usuario no autorizado",
					Code:    http.StatusUnauthorized,
				})
				c.Abort()
				return
			}
			db.ExtendToken(currentUser)
			c.Next()
			return
		} else {
			logrus.Print("user inactive")
			c.JSON(http.StatusUnauthorized, config.SetError("Sesión expirada"))
			c.Abort()
			return
		}
	} else {
		logrus.Print("invalid token")
		c.JSON(http.StatusUnauthorized, config.SetError("Token inválido"))
		c.Abort()
		return
	}
}
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func InitUserSuperAdminManager() {

	db.CreateDefaultUser("sam", "Super Admin", "ti@ghh.cl", "admin")

	/* 	result := db.Exec(`
	   	        DELETE FROM users
	   	        WHERE id NOT IN (
	   	            SELECT MIN(id)
	   	            FROM users
	   	            GROUP BY email
	   	        )
	   	    `)
	   	if result.Error != nil {
	   		panic(result.Error)
	   	}
	   	fmt.Println("Duplicated emails removed successfully.") */
}

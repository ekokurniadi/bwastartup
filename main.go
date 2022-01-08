package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/payment"
	"bwastartup/transaction"
	"bwastartup/user"
	webHandler "bwastartup/web/handler"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	s := godotenv.Load()
	if s != nil {
		fmt.Println(s)
	}

	// call repository
	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	// call service
	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)
	authService := auth.NewService()

	// call handler
	userHandler := handler.NewUserHandler(userService, authService, userRepository)
	campaignHandler := handler.NewCampaignHandler(campaignService, campaignRepository)
	transactionHandler := handler.NewTransactionHandler(transactionService, transactionRepository)

	// web handler
	userWebHandler := webHandler.NewUserHandler(userService)
	campaignWebHandler := webHandler.NewCampaignHandler(campaignService)
	transactionWebHandler := webHandler.NewTransactionHandler(transactionService)
	sessionWebHandler := webHandler.NewSessionHandler(userService)

	router := gin.Default()
	router.Use(cors.Default())
	cookieStore := cookie.NewStore([]byte(auth.SECRET_KEY))
	router.Use(sessions.Sessions("bwastartup", cookieStore))

	router.LoadHTMLGlob("web/templates/**/*")
	router.HTMLRender = loadTemplates("./web/templates")
	router.Static("images", "./images")
	router.Static("css", "./web/assets/css")
	router.Static("js", "./web/assets/js")
	router.Static("webfonts", "./web/assets/webfonts")
	api := router.Group("/api/v1")

	// endpoint all about users
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.POST("/fetch_data", userHandler.GetAllUserOnWeb)

	// endpoint all about campaigns
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)
	api.POST("/fetch_campaign", campaignHandler.GetAllCampaignsOnWeb)

	// endpoint all about transactions
	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)
	api.POST("/fetch_transactions", transactionHandler.GetAllTransactionsOnWeb)

	// web router to user page
	router.GET("/users", authAdminMiddleWare(), userWebHandler.Index)
	router.GET("/users/new", authAdminMiddleWare(), userWebHandler.New)
	router.POST("/users", authAdminMiddleWare(), userWebHandler.Create)
	router.GET("/users/edit/:id", authAdminMiddleWare(), userWebHandler.Edit)
	router.POST("/users/update/:id", authAdminMiddleWare(), userWebHandler.Update)
	router.GET("/users/avatar/:id", authAdminMiddleWare(), userWebHandler.NewAvatar)
	router.POST("/users/avatar/:id", authAdminMiddleWare(), userWebHandler.CreateAvatar)
	router.GET("/users/delete/:id", authAdminMiddleWare(), userWebHandler.Delete)

	// web router to campaign page
	router.GET("/campaigns", authAdminMiddleWare(), campaignWebHandler.Index)

	// web router to transactions page
	router.GET("/transactions", authAdminMiddleWare(), transactionWebHandler.Index)

	router.GET("/login", sessionWebHandler.New)
	router.POST("/session", sessionWebHandler.Create)

	router.Run()

}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.ApiResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")

		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			response := helper.ApiResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.ApiResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)

		if err != nil {
			response := helper.ApiResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)

	}

}

func authAdminMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		userIDSession := session.Get("userID")
		if userIDSession == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
	}
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.html")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
		panic(err.Error())
	}

	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}

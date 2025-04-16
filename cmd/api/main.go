package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"trello-backend/docs"
	"trello-backend/graph/generated"
	"trello-backend/graph/resolver"
	"trello-backend/internal/app"
	"trello-backend/internal/config"
	"trello-backend/internal/repositories"
	"trello-backend/internal/routes"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
)

// @title Trello 後端 API
// @version 1.0
// @description Trello 後端 API 文件
// @host localhost:8080
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func initDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.GetDBConnString()), &gorm.Config{})
	if err != nil {
		log.Fatal("無法連線到資料庫:", err)
	}

	// Run migrations
	app.Migrate(db)

	log.Println("資料庫遷移完成")

	return db
}

func setupSwagger() {
	docs.SwaggerInfo.Title = "Trello 後端 API"
	docs.SwaggerInfo.Description = "Trello 後端 API 文件"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}
}

func main() {
	cfg := config.LoadConfig()
	db := initDB(cfg)

	// 初始化 Swagger 文件
	setupSwagger()

	// 使用 wire 進行相依性注入
	api, err := app.InitializeAPI(db, cfg.JWTSecret)
	if err != nil {
		log.Fatal("無法初始化 API:", err)
	}

	// 設定路由
	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router := routes.NewRouter(engine, cfg.JWTSecret, cfg)

	// 註冊所有 handlers
	for name, handler := range api.GetHandlers() {
		router.RegisterHandler(name, handler)
	}

	// 設定所有路由
	router.SetupRoutes()

	// ===== GraphQL 與 Playground 設定 =====
	userRepo := repositories.NewUserRepository(db)
	gqlResolver := &resolver.Resolver{UserRepo: userRepo}
	gqlHandler := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: gqlResolver}))
	gqlHandler.Use(extension.Introspection{})
	gqlHandler.AddTransport(transport.POST{})
	gqlHandler.AddTransport(transport.GET{})
	gqlHandler.AddTransport(transport.Options{})
	gqlHandler.AddTransport(transport.MultipartForm{})

	engine.Any("/graphql", func(c *gin.Context) {
		gqlHandler.ServeHTTP(c.Writer, c.Request)
	})
	engine.GET("/playground", func(c *gin.Context) {
		playground.Handler("GraphQL Playground", "/graphql").ServeHTTP(c.Writer, c.Request)
	})
	// ===== End =====

	// 啟動伺服器
	if err := engine.Run(":8080"); err != nil {
		log.Fatal("伺服器啟動失敗:", err)
	}
}

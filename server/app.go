package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/algakz/banner_service/pkg/auth"
	authhttp "github.com/algakz/banner_service/pkg/auth/delivery/http"
	authpostgres "github.com/algakz/banner_service/pkg/auth/repository/postgres"
	authusecase "github.com/algakz/banner_service/pkg/auth/usecase"
	"github.com/algakz/banner_service/pkg/banner"

	bnhttp "github.com/algakz/banner_service/pkg/banner/delivery/http"
	bnpostgres "github.com/algakz/banner_service/pkg/banner/repository/postgres"
	bnusecase "github.com/algakz/banner_service/pkg/banner/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

type App struct {
	httpServer *http.Server

	authUC   auth.UseCase
	bannerUC banner.UseCase
}

func NewApp() *App {
	db := InitDB()

	authRepo := authpostgres.NewUserRepository(db)
	bannerRepo := bnpostgres.NewBannerPostgres(db)

	return &App{
		bannerUC: bnusecase.NewBannerUseCase(bannerRepo),
		authUC: authusecase.NewAuthUseCase(
			authRepo,
			viper.GetString("auth.salt"),
			[]byte(viper.GetString("auth.signed_key")),
			viper.GetDuration("auth.token_ttl"),
		),
	}
}

func (a *App) Run(port string) error {
	logrus.SetFormatter(new(logrus.JSONFormatter))
  logrus.SetOutput(os.Stdout)
  logrus.SetLevel(logrus.DebugLevel)
  logrus.Info("configs for logrus setted")
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	// Set up http handlers
	authhttp.RegisterHTTPEndpoints(router, a.authUC)
	authMiddleware := authhttp.NewAuthMiddleware(a.authUC)
	api := router.Group("/api", authMiddleware)
	bnhttp.RegisterHTTPEndpoints(api, a.bannerUC)

	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    10 * time.Second,
	}
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			logrus.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func InitDB() *pgxpool.Pool {
	if err := gotenv.Load(); err != nil {
		logrus.Fatalf("failed while loading .end: %s", err.Error())
	}
	dbpool, err := pgxpool.New(
		context.Background(),
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			viper.GetString("db.username"),
			os.Getenv("DB_PASSWORD"),
			viper.GetString("db.host"),
			viper.GetString("db.port"),
			viper.GetString("db.db"),
			viper.GetString("db.sslmode"),
		),
	)
	if err != nil {
		logrus.Fatalf("failed while creating connection pool: %s", err.Error())
		os.Exit(1)
	}

	return dbpool
}

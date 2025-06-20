package main

import (
	"context"
	"github.com/muriloFlores/StoreManager/infrastructure/reset_token"
	"github.com/muriloFlores/StoreManager/infrastructure/security"
	"github.com/muriloFlores/StoreManager/infrastructure/security/uuid_generator"
	"github.com/muriloFlores/StoreManager/infrastructure/workers/email"
	"github.com/muriloFlores/StoreManager/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/muriloFlores/StoreManager/pkg/config"

	"github.com/muriloFlores/StoreManager/infrastructure/db"
	"github.com/muriloFlores/StoreManager/infrastructure/db/postgres_repository"
	"github.com/muriloFlores/StoreManager/infrastructure/notifications"
	"github.com/muriloFlores/StoreManager/infrastructure/queue"
	"github.com/muriloFlores/StoreManager/infrastructure/security/jwt_manager"
	"github.com/muriloFlores/StoreManager/infrastructure/templates"
	web_http "github.com/muriloFlores/StoreManager/infrastructure/web/http"

	"github.com/muriloFlores/StoreManager/infrastructure/web/router"

	"github.com/muriloFlores/StoreManager/internal/core/use_case/auth"
	"github.com/muriloFlores/StoreManager/internal/core/use_case/user"
)

func main() {
	log.Println("Iniciando a aplicação Store Manager...")

	appLogs := logger.NewLogger()

	cfg, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalf("FATAL: error in load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbpool, err := db.NewDBPool(ctx, cfg)
	if err != nil {
		log.Fatalf("FATAL: dabase initialization failed: %v", err)
	}
	defer dbpool.Close()

	redisConnectionOpt := asynq.RedisClientOpt{Addr: cfg.RedisAddress}

	userRepo := postgres_repository.NewPostgresUserRepository(dbpool)
	actionTokenRepo := postgres_repository.NewActionTokenRepository(dbpool, appLogs)
	passwordHasher := security.NewPasswordHasher()
	idGenerator := uuid_generator.NewUUIDGenerator()
	tokenManager := jwt_manager.NewJWTGenerator(cfg.JWTSecret)

	templateManager, err := templates.NewHTMLTemplateManager()
	if err != nil {
		log.Fatalf("FATAL: Falha ao carregar templates HTML: %v", err)
	}

	cryptToken := reset_token.NewCryptoTokenGenerator()
	if err != nil {
		log.Fatalf("FATAL: Falha ao carregar templates de email: %v", err)
	}
	emailSender := notifications.NewSmtpSender(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpSenderEmail, cfg.SmtpAppPassword)
	taskEnqueuer := queue.NewTaskEnqueuer(redisConnectionOpt, appLogs)

	userUseCases := user.NewUserUseCases(
		userRepo,
		passwordHasher,
		idGenerator,
		cryptToken,
		taskEnqueuer,
		actionTokenRepo,
		appLogs,
	)

	authUseCases := auth.NewAuthUseCases(
		userRepo,
		passwordHasher,
		tokenManager,
		actionTokenRepo,
		cryptToken,
		taskEnqueuer,
		appLogs,
	)

	emailProcessor := email.NewEmailProcessor(emailSender, templateManager)
	go email.RunTaskServer(redisConnectionOpt, emailProcessor, appLogs)

	userHandler := web_http.NewUserHandler(userUseCases)
	authHandler := web_http.NewAuthHandler(authUseCases, appLogs)
	mainRouter := router.NewRouter(userHandler, authHandler, tokenManager)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: mainRouter,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Servidor HTTP escutando na porta %s", cfg.ServerPort)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("FATAL: Erro ao iniciar o servidor: %v", err)
	case sig := <-shutdown:
		log.Printf("INFO: Sinal de desligamento recebido: %v. Iniciando graceful shutdown...", sig)
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("ERRO: Desligamento forçado do servidor: %v", err)
			server.Close()
		}
		log.Println("Aplicação finalizada.")
	}
}

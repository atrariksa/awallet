package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atrariksa/awallet/blackbox"
	"github.com/atrariksa/awallet/configs"
	"github.com/atrariksa/awallet/constants"
	"github.com/atrariksa/awallet/drivers"
	"github.com/atrariksa/awallet/handlers"
	"github.com/atrariksa/awallet/middlewares"
	"github.com/atrariksa/awallet/migrations"
	"github.com/atrariksa/awallet/repos"
	"github.com/atrariksa/awallet/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	cmdMessage :=
		`
	Please use following commands :
	1. use "migrate up" to migrate tables
	2. use "server" to run service
	3. use "blackbox" to run test cases
	`
	if len(os.Args) == 1 {
		log.Fatalln(cmdMessage)
	}
	command := os.Args[1]
	switch command {
	case "server":
		server()
	case "migrate":
		migrate(os.Args)
	case "blackbox":
		runBlackbox()
	default:
		log.Println(fmt.Sprintf(`Unknown command "%v". %v`, command, cmdMessage))
	}
}

func server() {
	cfg := configs.Get()
	addr := cfg.APP.HOST + ":" + cfg.APP.PORT

	// The HTTP Server
	server := &http.Server{Addr: addr, Handler: setupService(cfg)}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

func setupService(cfg *configs.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	c := drivers.GetRedisClient(cfg)
	cacheRepo := repos.NewCache(cfg, c)

	dbRead := drivers.NewDBClientRead(cfg)
	dbWrite := drivers.NewDBClientWrite(cfg)

	userRepoRead := repos.UserRepoRead{DBRead: dbRead, Cache: cacheRepo}
	userRepoWrite := repos.UserRepoWrite{DBWrite: dbWrite, Cache: cacheRepo}

	userService := services.UserService{
		UserRepoRead:  &userRepoRead,
		UserRepoWrite: &userRepoWrite,
	}

	tokenService := services.NewTokenService(cfg)

	registerHandler := handlers.RegisterHandler{
		UserService:  &userService,
		TokenService: &tokenService,
	}

	r.Post(constants.CREATE_USER_PATH, registerHandler.Handle)

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddlewareHandler(cfg))

		userBalanceWrite := repos.UserBalanceRepoWrite{DBWrite: dbWrite, Cache: cacheRepo}
		userBalanceRead := repos.UserBalanceRepoRead{DBRead: dbRead, Cache: cacheRepo}
		userBalanceService := services.UserBalanceService{
			UserRepoRead:     &userRepoRead,
			UserBalanceWrite: &userBalanceWrite,
			UserBalanceRead:  &userBalanceRead,
		}

		readBalanceHandler := handlers.ReadBalanceHandler{UserBalanceService: &userBalanceService}
		r.Get(constants.READ_BALANCE_PATH, readBalanceHandler.Handle)

		topupBalanceHandler := handlers.TopupBalanceHandler{UserBalanceService: &userBalanceService}
		r.Post(constants.TOPUP_BALANCE_PATH, topupBalanceHandler.Handle)

		transferHandler := handlers.TransferHandler{UserBalanceService: &userBalanceService}
		r.Post(constants.TRANSFER_PATH, transferHandler.Handle)

		topTransactionsPerUserHandler := handlers.TopTransactionsPerUserHandler{UserBalanceService: &userBalanceService}
		r.Get(constants.TOP_TRANSACTIONS_PER_USER_PATH, topTransactionsPerUserHandler.Handle)

		topUserHandler := handlers.ListTopUserHandler{UserService: &userService}
		r.Get(constants.TOP_USERS_PATH, topUserHandler.Handle)
	})

	return r
}

func migrate([]string) {
	cfg := configs.Get()
	dbWrite := drivers.NewDBClientWrite(cfg)
	m := migrations.Migrator{DB: dbWrite}
	m.MigrateUp()
}

func runBlackbox() {
	cfg := configs.Get()
	bb := blackbox.NewBlackBox(cfg)
	bb.Run()
}

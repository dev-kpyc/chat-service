package application

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/dev-kpyc/chat-service/pkg/chatroom"
	"github.com/dev-kpyc/chat-service/pkg/clients/grpc"
	"github.com/dev-kpyc/chat-service/pkg/messaging"
	repository "github.com/dev-kpyc/chat-service/pkg/repository/postgres"
	"github.com/dev-kpyc/chat-service/pkg/user"
	"github.com/heptiolabs/healthcheck"
	_ "github.com/lib/pq" // postgres driver
)

// ChatServer ...
type ChatServer struct {
	db     *sql.DB
	repo   *repository.Postgres
	config map[string]string
	checks map[string]healthcheck.Check
}

// Start ...
func Start() {
	log.Println("Starting chat server")

	chat := ChatServer{}
	chat.configure()
	chat.initializeDB()
	chat.initChecks()
	chat.waitUntilReady()

	chat.initializeRepository()
	chat.registerAPIs()
	chat.registerHealthChecks()
}

func (chat *ChatServer) initializeDB() {
	config := chat.config

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config["postgres.hostname"],
		config["postgres.port"],
		config["postgres.username"],
		config["postgres.password"],
		config["postgres.database"],
	)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	chat.db = db
}

func (chat *ChatServer) initializeRepository() {
	chat.repo = repository.NewPostgresRepository(chat.db)

	chat.repo.Initialize()
}

func (chat *ChatServer) registerAPIs() {

	usersvc := user.New(chat.repo)
	chatsvc := chatroom.New(chat.repo, usersvc)
	messagesvc := messaging.New(chat.repo, usersvc)
	clients := grpc.New(&chatsvc, &messagesvc, &usersvc)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "0.0.0.0", chat.config["server.port"]))
	if err != nil {
		panic(err)
	}

	clients.RegisterEndpoints(lis)
}

func (chat *ChatServer) initChecks() {

	chat.checks = make(map[string]healthcheck.Check)
	dbcheck := healthcheck.DatabasePingCheck(chat.db, 1*time.Second)
	chat.checks["database"] = dbcheck
}

func (chat *ChatServer) waitUntilReady() {
	for _, check := range chat.checks {
		check()
	}
	log.Println("Server is ready")
}

func (chat *ChatServer) registerHealthChecks() {
	health := healthcheck.NewHandler()

	for name, check := range chat.checks {
		health.AddReadinessCheck(name, check)
	}

	http.ListenAndServe(":8081", health)
}

package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rsilraf/pos_goexpert/desafios/clean_architecture/configs"
	"github.com/rsilraf/pos_goexpert/desafios/clean_architecture/internal/event/handler"
	"github.com/rsilraf/pos_goexpert/desafios/clean_architecture/internal/infra/graph"
	"github.com/rsilraf/pos_goexpert/desafios/clean_architecture/internal/infra/grpc/pb"
	"github.com/rsilraf/pos_goexpert/desafios/clean_architecture/internal/infra/grpc/service"
	"github.com/rsilraf/pos_goexpert/desafios/clean_architecture/internal/infra/web/webserver"
	"github.com/rsilraf/pos_goexpert/desafios/clean_architecture/pkg/events"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

var cfg *configs.Conf

func init() {
	c, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	cfg = c
}

func main() {

	db, err := sql.Open(cfg.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	// EVENT DISPATCHER
	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", handler.NewOrderCreatedHandler(rabbitMQChannel))
	eventDispatcher.Register("OrdersListed", handler.NewOrdersListedHandler(rabbitMQChannel))

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db, eventDispatcher)

	// WEB
	webserver := webserver.NewWebServer(cfg.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("POST", "/order", webOrderHandler.Create)
	webserver.AddHandler("GET", "/order", webOrderHandler.ListAll)

	fmt.Println("Starting web server on port", cfg.WebServerPort)
	go webserver.Start()

	// GRPC
	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", cfg.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	// GRAPHQL
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", cfg.GraphQLServerPort)
	http.ListenAndServe(":"+cfg.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.RabbitMQUser, cfg.RabbitMQPassword, cfg.RabbitMQHost, cfg.RabbitMQPort))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

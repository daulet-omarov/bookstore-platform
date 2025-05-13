package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/daulet-omarov/book-service/models"
	"github.com/daulet-omarov/bookstore-platform/your-module-path/bookpb"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config config
	logger *log.Logger
	books  models.BookModel
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4040, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("BOOKSTORE_DB_DSN"), "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle time")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Printf("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		books:  models.BookModel{DB: db},
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		logger.Printf("starting %s server on %d", cfg.env, cfg.port)
		err = srv.ListenAndServe()
		logger.Fatal(err)
	}()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	bookpb.RegisterBookServiceServer(grpcServer, &bookServer{books: models.BookModel{DB: db}})

	log.Println("BookService gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type bookServer struct {
	bookpb.UnimplementedBookServiceServer
	books models.BookModel
}

func (s *bookServer) CheckBook(ctx context.Context, req *bookpb.BookRequest) (*bookpb.BookResponse, error) {
	id, err := strconv.ParseInt(req.BookId, 10, 64)
	if err != nil {
		return nil, err
	}
	book, err := s.books.Get(id)
	if err != nil {
		return nil, err
	}
	return &bookpb.BookResponse{
		Available: book.Stock > 0,
		Title:     book.Title,
		Quantity:  int32(book.Stock),
	}, nil
}

func (s *bookServer) UpdateBook(ctx context.Context, req *bookpb.UpdateRequest) (*bookpb.UpdateResponse, error) {
	id, err := strconv.ParseInt(req.BookId, 10, 64)
	if err != nil {
		return nil, err
	}
	book, err := s.books.Get(id)
	if err != nil {
		return nil, err
	}
	book.Stock = book.Stock - int64(req.Delta)
	success := s.books.Update(book) == nil
	return &bookpb.UpdateResponse{Success: success}, nil
}

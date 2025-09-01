package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/blyundup/microservices_course/week_1/grpc/internal/config"
	"github.com/blyundup/microservices_course/week_1/grpc/internal/config/env"
	desc "github.com/blyundup/microservices_course/week_1/grpc/pkg/note_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedNoteV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	builderInsert := sq.Insert("note").
		PlaceholderFormat(sq.Dollar).
		Columns("title", "body").
		Values(gofakeit.City(), gofakeit.Address().Street).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var noteID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&noteID)
	if err != nil {
		log.Fatalf("failed to insert note: %v", err)
	}

	log.Printf("inserted note with id: %v", noteID)

	return &desc.CreateResponse{
		Id: noteID,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	builderSelectOne := sq.Select("id", "title", "body", "created_at", "updated_at").
		From("note").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var id int64
	var title, body string
	var created_at time.Time
	var updated_at sql.NullTime

	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &title, &body, &created_at, &updated_at)
	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
	}

	log.Printf("id: %v, title: %v, body: %v, created_at: %v, updated_at: %v", id, title, body, created_at, updated_at)

	var updatedAtTime *timestamppb.Timestamp
	if updated_at.Valid {
		updatedAtTime = timestamppb.New(updated_at.Time)
	}

	return &desc.GetResponse{
		Note: &desc.Note{
			Id:        id,
			Info:      &desc.NoteInfo{Title: title, Content: body},
			CreatedAt: timestamppb.New(created_at),
			UpdatedAt: updatedAtTime,
		},
	}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

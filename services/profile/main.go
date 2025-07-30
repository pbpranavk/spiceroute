package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"

	pb "spiceroute/proto"

	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
)

func pqStringArray(slice []string) interface{} {
	return slice
}

type server struct {
	pb.UnimplementedProfileServiceServer
	db *sql.DB
}

func (s *server) UpsertPreference(ctx context.Context, p *pb.Preference) (*pb.Preference, error) {
	_, err := s.db.ExecContext(ctx, `
	  INSERT INTO preferences (user_id, cuisines, allergies, budget_week, spicy)
	  VALUES ($1, $2, $3, $4, $5)
	  ON CONFLICT (user_id) DO UPDATE SET cuisines=$2, allergies=$3, budget_week=$4, spicy=$5`,
		p.UserId, pqStringArray(p.Cuisines), pqStringArray(p.Allergies), p.BudgetWeek, p.Spicy)
	return p, err
}

func main() {
	db, err := sql.Open("pgx", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	lis, _ := net.Listen("tcp", ":50051")
	grpcServer := grpc.NewServer()
	pb.RegisterProfileServiceServer(grpcServer, &server{db: db})
	log.Fatal(grpcServer.Serve(lis))
}

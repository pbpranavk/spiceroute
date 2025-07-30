package main

import (
	"context"
	"log"
	"net"

	"spiceroute/pkg/database"
	"spiceroute/pkg/models"
	pb "spiceroute/proto"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedProfileServiceServer
	db *gorm.DB
}

func (s *server) UpsertPreference(ctx context.Context, p *pb.Preference) (*pb.Preference, error) {
	preference := models.Preference{
		UserID:     p.UserId,
		Cuisines:   p.Cuisines,
		Allergies:  p.Allergies,
		BudgetWeek: p.BudgetWeek,
		Spicy:      p.Spicy,
	}

	// Use Upsert (Create or Update)
	result := s.db.WithContext(ctx).Where("user_id = ?", p.UserId).
		Assign(preference).
		FirstOrCreate(&preference)

	if result.Error != nil {
		return nil, result.Error
	}

	// Convert back to protobuf
	return &pb.Preference{
		UserId:     preference.UserID,
		Cuisines:   preference.Cuisines,
		Allergies:  preference.Allergies,
		BudgetWeek: preference.BudgetWeek,
		Spicy:      preference.Spicy,
	}, nil
}

func (s *server) GetPreference(ctx context.Context, p *pb.Preference) (*pb.Preference, error) {
	var preference models.Preference

	result := s.db.WithContext(ctx).Where("user_id = ?", p.UserId).First(&preference)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil if not found
		}
		return nil, result.Error
	}

	// Convert to protobuf
	return &pb.Preference{
		UserId:     preference.UserID,
		Cuisines:   preference.Cuisines,
		Allergies:  preference.Allergies,
		BudgetWeek: preference.BudgetWeek,
		Spicy:      preference.Spicy,
	}, nil
}

func main() {
	// Initialize database connection
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProfileServiceServer(grpcServer, &server{db: db})

	log.Println("Profile service starting on :50051")
	log.Fatal(grpcServer.Serve(lis))
}

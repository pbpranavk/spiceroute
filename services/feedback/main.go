package main

import (
	"context"
	"log"
	"net"
	"time"

	"spiceroute/pkg/database"
	"spiceroute/pkg/models"
	pb "spiceroute/proto"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type server struct {
	db *gorm.DB
	pb.UnimplementedFeedbackServiceServer
}

func (s *server) SubmitFeedback(ctx context.Context, batch *pb.FeedbackBatch) (*emptypb.Empty, error) {
	// Use a transaction for batch operations
	return &emptypb.Empty{}, s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, f := range batch.Entries {
			// Parse cooked_at time
			cookedAt, err := time.Parse(time.RFC3339, f.CookedAt)
			if err != nil {
				// If parsing fails, use current time
				cookedAt = time.Now()
			}

			feedback := models.Feedback{
				UserID:          f.UserId,
				DishID:          f.DishId,
				Rating:          f.Rating,
				Skipped:         f.Skipped,
				SubstitutedWith: f.SubstitutedWith,
				Comment:         f.Comment,
				CookedAt:        cookedAt,
			}

			// Use Upsert (Create or Update) based on user_id, dish_id, and cooked_at
			result := tx.Where("user_id = ? AND dish_id = ? AND cooked_at = ?",
				f.UserId, f.DishId, cookedAt).
				Assign(feedback).
				FirstOrCreate(&feedback)

			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
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
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterFeedbackServiceServer(grpcServer, &server{db: db})

	log.Println("Feedback service starting on :50054")
	log.Fatal(grpcServer.Serve(lis))
}

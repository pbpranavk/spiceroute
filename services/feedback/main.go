package feedback

import (
	"context"
	"database/sql"

	pb "spiceroute/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	db *sql.DB
	pb.UnimplementedFeedbackServiceServer
}

func (s *server) SubmitFeedback(ctx context.Context, batch *pb.FeedbackBatch) (*emptypb.Empty, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	stmt, err := tx.PrepareContext(ctx, `
	  INSERT INTO feedback (user_id, dish_id, rating, skipped, substituted_with, comment, cooked_at)
	  VALUES ($1,$2,$3,$4,$5,$6,$7)
	  ON CONFLICT (user_id, dish_id, cooked_at) DO UPDATE SET
		rating = EXCLUDED.rating,
		skipped = EXCLUDED.skipped,
		substituted_with = EXCLUDED.substituted_with,
		comment = EXCLUDED.comment`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for _, f := range batch.Entries {
		_, err := stmt.ExecContext(ctx, f.UserId, f.DishId, f.Rating, f.Skipped, f.SubstitutedWith, f.Comment, f.CookedAt)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	return &emptypb.Empty{}, tx.Commit()
}

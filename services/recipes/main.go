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
	db *gorm.DB
	pb.UnimplementedRecipeServiceServer
}

func (s *server) CreateRecipe(ctx context.Context, r *pb.Recipe) (*pb.RecipeID, error) {
	recipe := models.Recipe{
		Name:          r.Name,
		Cuisine:       r.Cuisine,
		PrepMinutes:   r.PrepMinutes,
		Calories:      r.Calories,
		Ingredients:   r.Ingredients,
		Cost:          r.Cost,
		ShelfLifeDays: r.ShelfLifeDays,
		Tags:          r.Tags,
		Nutrition:     r.Nutrition,
	}

	result := s.db.WithContext(ctx).Create(&recipe)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pb.RecipeID{Id: recipe.ID}, nil
}

func (s *server) ListRecipes(ctx context.Context, q *pb.RecipeQuery) (*pb.RecipeList, error) {
	var recipes []models.Recipe

	query := s.db.WithContext(ctx).Model(&models.Recipe{})

	// Apply filters if provided
	if len(q.Cuisines) > 0 {
		query = query.Where("cuisine IN ?", q.Cuisines)
	}

	// Note: Spicy filter would need to be implemented based on recipe tags or a separate field
	// For now, we'll just limit the results
	query = query.Limit(100)

	result := query.Find(&recipes)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert to protobuf
	var pbRecipes []*pb.Recipe
	for _, recipe := range recipes {
		pbRecipes = append(pbRecipes, &pb.Recipe{
			Id:            recipe.ID,
			Name:          recipe.Name,
			Cuisine:       recipe.Cuisine,
			PrepMinutes:   recipe.PrepMinutes,
			Calories:      recipe.Calories,
			Ingredients:   recipe.Ingredients,
			Cost:          recipe.Cost,
			ShelfLifeDays: recipe.ShelfLifeDays,
			Tags:          recipe.Tags,
			Nutrition:     recipe.Nutrition,
		})
	}

	return &pb.RecipeList{Recipes: pbRecipes}, nil
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
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRecipeServiceServer(grpcServer, &server{db: db})

	log.Println("Recipe service starting on :50053")
	log.Fatal(grpcServer.Serve(lis))
}

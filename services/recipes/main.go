package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"

	pb "spiceroute/proto"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
)

func pqStringArray(slice []string) interface{} {
	return slice
}

type server struct {
	db *sql.DB
	pb.UnimplementedRecipeServiceServer
}

func (s *server) CreateRecipe(ctx context.Context, r *pb.Recipe) (*pb.RecipeID, error) {
	id := uuid.New()
	_, err := s.db.ExecContext(ctx, `
	  INSERT INTO recipes (id, name, cuisine, prep_minutes, calories, ingredients, cost, shelf_life_days, tags, nutrition)
	  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		id, r.Name, r.Cuisine, r.PrepMinutes, r.Calories, pqStringArray(r.Ingredients),
		r.Cost, r.ShelfLifeDays, pqStringArray(r.Tags), r.Nutrition)
	return &pb.RecipeID{Id: id.String()}, err
}

func (s *server) ListRecipes(ctx context.Context, q *pb.RecipeQuery) (*pb.RecipeList, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, cuisine, prep_minutes, calories, ingredients, cost, shelf_life_days, tags, nutrition FROM recipes LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []*pb.Recipe
	for rows.Next() {
		var r pb.Recipe
		var ingredients, tags []string
		err := rows.Scan(&r.Id, &r.Name, &r.Cuisine, &r.PrepMinutes, &r.Calories, &ingredients, &r.Cost, &r.ShelfLifeDays, &tags, &r.Nutrition)
		if err != nil {
			return nil, err
		}
		r.Ingredients = ingredients
		r.Tags = tags
		recipes = append(recipes, &r)
	}
	return &pb.RecipeList{Recipes: recipes}, nil
}

func main() {
	db, err := sql.Open("pgx", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	lis, _ := net.Listen("tcp", ":50053")
	grpcServer := grpc.NewServer()
	pb.RegisterRecipeServiceServer(grpcServer, &server{db: db})
	log.Fatal(grpcServer.Serve(lis))
}

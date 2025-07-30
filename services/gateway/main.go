package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	pb "spiceroute/proto"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Initialize gRPC connections
	profileConn, _ := grpc.NewClient("profile:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	plannerConn, _ := grpc.NewClient("planner:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	recipesConn, _ := grpc.NewClient("recipes:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	feedbackConn, _ := grpc.NewClient("feedback:50056", grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Initialize service clients
	profile := pb.NewProfileServiceClient(profileConn)
	planner := pb.NewPlannerServiceClient(plannerConn)
	recipes := pb.NewRecipeServiceClient(recipesConn)
	feedback := pb.NewFeedbackServiceClient(feedbackConn)

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// User Preferences & Profile APIs
	r.Route("/preferences", func(r chi.Router) {
		r.Post("/onboarding", func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var pref pb.Preference
			json.Unmarshal(body, &pref)

			result, err := profile.UpsertPreference(context.Background(), &pref)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		})

		r.Get("/{user_id}", func(w http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "user_id")
			pref := &pb.Preference{UserId: userID}

			result, err := profile.GetPreference(context.Background(), pref)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		})

		r.Put("/{user_id}", func(w http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "user_id")
			body, _ := io.ReadAll(r.Body)
			var pref pb.Preference
			json.Unmarshal(body, &pref)
			pref.UserId = userID

			result, err := profile.UpsertPreference(context.Background(), &pref)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		})
	})

	// Recipe Management APIs
	r.Route("/recipes", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			cuisines := r.URL.Query()["cuisine"]
			spicyStr := r.URL.Query().Get("spicy")
			spicy := false
			if spicyStr == "true" {
				spicy = true
			}

			query := &pb.RecipeQuery{
				Cuisines: cuisines,
				Spicy:    spicy,
			}

			result, err := recipes.ListRecipes(context.Background(), query)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		})

		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			recipeID := chi.URLParam(r, "id")
			// Note: This would need a GetRecipe method in the proto
			// For now, we'll return a placeholder
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"id": recipeID, "message": "Recipe details endpoint"})
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var recipe pb.Recipe
			json.Unmarshal(body, &recipe)

			result, err := recipes.CreateRecipe(context.Background(), &recipe)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		})

		r.Get("/search", func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query().Get("q")
			cuisines := r.URL.Query()["cuisine"]

			// This would integrate with the Vector service for semantic search
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"query":    query,
				"cuisines": cuisines,
				"message":  "Semantic search endpoint",
			})
		})

		r.Get("/recommendations", func(w http.ResponseWriter, r *http.Request) {
			userID := r.URL.Query().Get("user_id")
			// This would integrate with the Vector service for recommendations
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"user_id": userID,
				"message": "Personalized recommendations endpoint",
			})
		})
	})

	// Meal Planning APIs
	r.Route("/plans", func(r chi.Router) {
		r.Post("/generate", func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var req pb.PlanRequest
			json.Unmarshal(body, &req)

			result, err := planner.GeneratePlan(context.Background(), &req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		})

		r.Get("/{user_id}", func(w http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "user_id")
			// This would need a GetPlans method in the proto
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"user_id": userID,
				"message": "Get user meal plans endpoint",
			})
		})

		r.Get("/{user_id}/{plan_id}", func(w http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "user_id")
			planID := chi.URLParam(r, "plan_id")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"user_id": userID,
				"plan_id": planID,
				"message": "Get specific meal plan endpoint",
			})
		})

		r.Post("/{user_id}/{plan_id}/regenerate", func(w http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "user_id")
			planID := chi.URLParam(r, "plan_id")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"user_id": userID,
				"plan_id": planID,
				"message": "Regenerate meal plan endpoint",
			})
		})
	})

	// Shopping & Ordering APIs
	r.Route("/shopping", func(r chi.Router) {
		r.Get("/{user_id}", func(w http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "user_id")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"user_id": userID,
				"message": "Get shopping lists endpoint",
			})
		})

		r.Post("/{user_id}/generate", func(w http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "user_id")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"user_id": userID,
				"message": "Generate shopping list endpoint",
			})
		})

		r.Post("/{user_id}/{list_id}/order", func(w http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "user_id")
			listID := chi.URLParam(r, "list_id")
			// This would integrate with the Orderer service
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"user_id": userID,
				"list_id": listID,
				"message": "Place grocery order endpoint",
			})
		})

		r.Get("/{user_id}/orders", func(w http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "user_id")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"user_id": userID,
				"message": "Get order history endpoint",
			})
		})
	})

	// Feedback APIs
	r.Route("/feedback", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var feedbackBatch pb.FeedbackBatch
			json.Unmarshal(body, &feedbackBatch)

			result, err := feedback.SubmitFeedback(context.Background(), &feedbackBatch)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		})

		r.Get("/{user_id}", func(w http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "user_id")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"user_id": userID,
				"message": "Get user feedback history endpoint",
			})
		})

		r.Get("/recipes/{recipe_id}", func(w http.ResponseWriter, r *http.Request) {
			recipeID := chi.URLParam(r, "recipe_id")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"recipe_id": recipeID,
				"message":   "Get recipe ratings endpoint",
			})
		})
	})

	// Analytics APIs
	r.Route("/analytics", func(r chi.Router) {
		r.Get("/{user_id}/nutrition", func(w http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "user_id")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"user_id": userID,
				"message": "Get nutrition insights endpoint",
			})
		})

		r.Get("/{user_id}/spending", func(w http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "user_id")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"user_id": userID,
				"message": "Get spending analytics endpoint",
			})
		})

		r.Get("/{user_id}/cooking-time", func(w http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "user_id")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"user_id": userID,
				"message": "Get cooking time insights endpoint",
			})
		})
	})

	// Health & Dietary APIs
	r.Route("/nutrition", func(r chi.Router) {
		r.Get("/calculator", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Nutrition calculator endpoint",
			})
		})
	})

	r.Route("/allergies", func(r chi.Router) {
		r.Get("/check", func(w http.ResponseWriter, r *http.Request) {
			recipeID := r.URL.Query().Get("recipe_id")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"recipe_id": recipeID,
				"message":   "Check recipe for allergens endpoint",
			})
		})
	})

	log.Printf("Gateway service starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	pb "spiceroute/proto"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	profileConn, _ := grpc.NewClient("profile:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	plannerConn, _ := grpc.NewClient("planner:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))

	profile := pb.NewProfileServiceClient(profileConn)
	planner := pb.NewPlannerServiceClient(plannerConn)

	r := chi.NewRouter()

	r.Post("/onboarding", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var pref pb.Preference
		json.Unmarshal(body, &pref)

		result, err := profile.UpsertPreference(context.Background(), &pref)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result)
	})

	r.Post("/plan/generate", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req pb.PlanRequest
		json.Unmarshal(body, &req)

		result, err := planner.GeneratePlan(context.Background(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}

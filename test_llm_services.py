#!/usr/bin/env python3
"""
Test script for SpiceRoute LLM-powered services
Demonstrates the enhanced vector and planner services with real LLM integration
"""

import asyncio
import httpx
import json
import os
from typing import Dict, Any

# Service URLs (adjust as needed for your deployment)
VECTOR_SERVICE_URL = "http://localhost:8080"  # Vector service port
PLANNER_SERVICE_URL = "http://localhost:8000"  # Planner service port

async def test_vector_service():
    """Test the enhanced vector service with LLM embeddings"""
    print("🧪 Testing Vector Service with LLM Integration")
    print("=" * 50)
    
    async with httpx.AsyncClient() as client:
        # Test health check
        try:
            response = await client.get(f"{VECTOR_SERVICE_URL}/health")
            print(f"✅ Health check: {response.json()}")
        except Exception as e:
            print(f"❌ Health check failed: {e}")
            return
        
        # Test embedding creation
        test_recipes = [
            {
                "id": "recipe_001",
                "text": "Spicy Chicken Tikka Masala with basmati rice and naan bread. A rich, creamy Indian curry with tender chicken pieces in a tomato-based sauce with aromatic spices.",
                "metadata": {"cuisine": "indian", "spice_level": "medium", "prep_time": 45}
            },
            {
                "id": "recipe_002", 
                "text": "Classic Margherita Pizza with fresh mozzarella, basil, and tomato sauce on a crispy thin crust. Traditional Italian comfort food.",
                "metadata": {"cuisine": "italian", "spice_level": "mild", "prep_time": 30}
            },
            {
                "id": "recipe_003",
                "text": "Vegetarian Pad Thai with rice noodles, tofu, bean sprouts, and peanuts in a tangy tamarind sauce. Authentic Thai street food.",
                "metadata": {"cuisine": "thai", "spice_level": "medium", "prep_time": 25}
            }
        ]
        
        print("\n📝 Creating embeddings...")
        for recipe in test_recipes:
            try:
                response = await client.post(
                    f"{VECTOR_SERVICE_URL}/embed",
                    json=recipe
                )
                result = response.json()
                print(f"✅ Embedded {recipe['id']}: {result['message']}")
            except Exception as e:
                print(f"❌ Failed to embed {recipe['id']}: {e}")
        
        # Test semantic search
        print("\n🔍 Testing semantic search...")
        search_queries = [
            "I want something spicy and flavorful",
            "Looking for Italian comfort food",
            "Need a quick vegetarian meal",
            "Show me Asian cuisine options"
        ]
        
        for query in search_queries:
            try:
                response = await client.post(
                    f"{VECTOR_SERVICE_URL}/query",
                    json={"text": query, "k": 3, "threshold": 0.5}
                )
                result = response.json()
                print(f"\n🔎 Query: '{query}'")
                print(f"📊 Found {result['total_found']} results:")
                for i, item in enumerate(result['results'][:3], 1):
                    print(f"   {i}. {item['id']} (similarity: {item['similarity']:.3f})")
            except Exception as e:
                print(f"❌ Search failed for '{query}': {e}")
        
        # List all embeddings
        try:
            response = await client.get(f"{VECTOR_SERVICE_URL}/embeddings")
            embeddings = response.json()
            print(f"\n📋 Total embeddings stored: {embeddings['total_embeddings']}")
        except Exception as e:
            print(f"❌ Failed to list embeddings: {e}")

async def test_planner_service():
    """Test the enhanced planner service with LLM recommendations"""
    print("\n\n🧪 Testing Planner Service with LLM Integration")
    print("=" * 50)
    
    async with httpx.AsyncClient() as client:
        # Test health check
        try:
            response = await client.get(f"{PLANNER_SERVICE_URL}/health")
            print(f"✅ Health check: {response.json()}")
        except Exception as e:
            print(f"❌ Health check failed: {e}")
            return
        
        # Test meal suggestions
        print("\n💡 Testing AI-powered meal suggestions...")
        suggestion_request = {
            "user_id": "user_123",
            "preferences": {
                "cuisines": ["italian", "indian", "mediterranean"],
                "spice_tolerance": "medium",
                "cooking_skill": "intermediate"
            },
            "dietary_restrictions": ["vegetarian"],
            "budget": 150.0,
            "days": 7
        }
        
        try:
            response = await client.post(
                f"{PLANNER_SERVICE_URL}/suggest",
                json=suggestion_request
            )
            result = response.json()
            print(f"✅ Generated {len(result['suggestions'])} meal suggestions")
            print(f"🎯 Cuisine variety: {', '.join(result['cuisine_variety'])}")
            print(f"💭 Reasoning: {result['reasoning'][:100]}...")
            print(f"📊 Nutrition goals: {result['nutrition_goals']}")
            
            print("\n🍽️  Top 5 suggestions:")
            for i, suggestion in enumerate(result['suggestions'][:5], 1):
                print(f"   {i}. {suggestion}")
                
        except Exception as e:
            print(f"❌ Meal suggestions failed: {e}")
        
        # Test meal planning
        print("\n📅 Testing AI-enhanced meal planning...")
        sample_dishes = [
            {
                "id": "dish_001",
                "name": "Chicken Tikka Masala",
                "cuisine": "indian",
                "prep_minutes": 45,
                "calories": 450,
                "ingredients": ["chicken", "yogurt", "spices", "tomato", "cream"],
                "cost": 12.50,
                "shelf_life_days": 3,
                "tags": ["spicy", "protein", "curry"],
                "nutrition": "High protein, moderate carbs"
            },
            {
                "id": "dish_002",
                "name": "Margherita Pizza",
                "cuisine": "italian", 
                "prep_minutes": 30,
                "calories": 350,
                "ingredients": ["dough", "mozzarella", "tomato", "basil"],
                "cost": 8.75,
                "shelf_life_days": 2,
                "tags": ["vegetarian", "comfort", "italian"],
                "nutrition": "Moderate protein, high carbs"
            },
            {
                "id": "dish_003",
                "name": "Greek Salad",
                "cuisine": "mediterranean",
                "prep_minutes": 15,
                "calories": 200,
                "ingredients": ["lettuce", "tomato", "cucumber", "olives", "feta"],
                "cost": 6.25,
                "shelf_life_days": 4,
                "tags": ["vegetarian", "healthy", "salad"],
                "nutrition": "Low calorie, high fiber"
            }
        ]
        
        plan_request = {
            "user_id": "user_123",
            "days": 3,
            "dishes": sample_dishes,
            "daily_calories": 2000.0,
            "budget_week": 100.0,
            "preferences": {
                "cuisines": ["italian", "indian", "mediterranean"],
                "spice_tolerance": "medium"
            },
            "dietary_restrictions": []
        }
        
        try:
            response = await client.post(
                f"{PLANNER_SERVICE_URL}/plan",
                json=plan_request
            )
            result = response.json()
            
            print(f"✅ Generated meal plan for {len(result['schedule'])} days")
            print(f"💰 Estimated savings: ${result['estimated_savings']:.2f}")
            print(f"🛒 Shopping list items: {len(result['shopping_list'])}")
            
            if result['recommendations']:
                print(f"\n💡 AI Recommendations:")
                for i, rec in enumerate(result['recommendations'][:3], 1):
                    print(f"   {i}. {rec}")
            
            print(f"\n📊 Nutrition Summary:")
            nutrition = result['nutrition_summary']
            print(f"   Total calories: {nutrition['total_calories']:.0f}")
            print(f"   Total cost: ${nutrition['total_cost']:.2f}")
            print(f"   Avg calories/day: {nutrition['avg_calories_per_day']:.0f}")
            
        except Exception as e:
            print(f"❌ Meal planning failed: {e}")

async def main():
    """Run all tests"""
    print("🚀 SpiceRoute LLM Services Test Suite")
    print("Testing enhanced vector and planner services with OpenAI integration")
    print("=" * 60)
    
    # Check if OpenAI API key is available
    if not os.getenv("OPENAI_API_KEY"):
        print("⚠️  Warning: OPENAI_API_KEY not set. Some features may use fallback implementations.")
    
    await test_vector_service()
    await test_planner_service()
    
    print("\n" + "=" * 60)
    print("✅ Test suite completed!")
    print("\n💡 To run with full LLM features, set your OPENAI_API_KEY environment variable")

if __name__ == "__main__":
    asyncio.run(main()) 
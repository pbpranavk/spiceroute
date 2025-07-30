import os
import logging
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from ortools.sat.python import cp_model
from openai import OpenAI
import asyncio
from typing import List, Optional, Dict, Any
import json

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(title="SpiceRoute Planner Service", version="1.0.0")

# Initialize OpenAI client
client = None
try:
    api_key = os.getenv("OPENAI_API_KEY")
    if api_key:
        client = OpenAI(api_key=api_key)
        logger.info("OpenAI client initialized successfully")
    else:
        logger.warning("OPENAI_API_KEY not found, LLM features will be limited")
except Exception as e:
    logger.error(f"Failed to initialize OpenAI client: {e}")

class Dish(BaseModel):
    id: str
    name: str
    cuisine: str
    prep_minutes: int
    calories: int
    ingredients: List[str]
    cost: float
    shelf_life_days: int
    tags: Optional[List[str]] = None
    nutrition: Optional[str] = None

class PlanRequest(BaseModel):
    user_id: str
    days: int
    dishes: List[Dish]
    daily_calories: float
    budget_week: float
    preferences: Optional[Dict[str, Any]] = None
    dietary_restrictions: Optional[List[str]] = None

class Day(BaseModel):
    day_index: int
    dish_ids: List[str]
    servings: List[int]
    total_calories: float
    total_cost: float

class PlanResponse(BaseModel):
    schedule: List[Day]
    cook_days: List[int]
    shopping_list: List[str]
    nutrition_summary: Optional[Dict[str, Any]] = None
    recommendations: Optional[List[str]] = None
    estimated_savings: Optional[float] = None

class SuggestionRequest(BaseModel):
    user_id: str
    preferences: Dict[str, Any]
    dietary_restrictions: List[str]
    budget: float
    days: int

class SuggestionResponse(BaseModel):
    suggestions: List[str]
    reasoning: str
    cuisine_variety: List[str]
    nutrition_goals: Dict[str, Any]

@app.get("/health")
def health_check():
    return {
        "status": "healthy",
        "openai_available": client is not None,
        "service": "planner"
    }

@app.post("/plan", response_model=PlanResponse)
async def plan(req: PlanRequest):
    try:
        # Filter dishes based on dietary restrictions
        filtered_dishes = filter_dishes_by_restrictions(req.dishes, req.dietary_restrictions)
        
        if not filtered_dishes:
            raise HTTPException(status_code=400, detail="No dishes available after applying dietary restrictions")
        
        # Generate optimized plan using constraint programming
        schedule, cook_days, shopping_list = generate_optimized_plan(
            filtered_dishes, req.days, req.daily_calories, req.budget_week
        )
        
        # Calculate nutrition summary
        nutrition_summary = calculate_nutrition_summary(schedule, filtered_dishes)
        
        # Generate AI-powered recommendations
        recommendations = None
        if client:
            recommendations = await generate_ai_recommendations(
                req.user_id, schedule, filtered_dishes, req.preferences
            )
        
        # Calculate estimated savings
        estimated_savings = calculate_savings(schedule, req.budget_week)
        
        return PlanResponse(
            schedule=schedule,
            cook_days=cook_days,
            shopping_list=shopping_list,
            nutrition_summary=nutrition_summary,
            recommendations=recommendations,
            estimated_savings=estimated_savings
        )
        
    except Exception as e:
        logger.error(f"Error generating plan: {e}")
        raise HTTPException(status_code=500, detail=f"Failed to generate plan: {str(e)}")

@app.post("/suggest", response_model=SuggestionResponse)
async def suggest_meals(req: SuggestionRequest):
    try:
        if not client:
            # Fallback suggestions without LLM
            return generate_fallback_suggestions(req)
        
        # Generate AI-powered meal suggestions
        suggestions, reasoning = await generate_ai_meal_suggestions(req)
        
        # Analyze cuisine variety
        cuisine_variety = analyze_cuisine_variety(suggestions)
        
        # Generate nutrition goals
        nutrition_goals = generate_nutrition_goals(req.preferences, req.dietary_restrictions)
        
        return SuggestionResponse(
            suggestions=suggestions,
            reasoning=reasoning,
            cuisine_variety=cuisine_variety,
            nutrition_goals=nutrition_goals
        )
        
    except Exception as e:
        logger.error(f"Error generating suggestions: {e}")
        raise HTTPException(status_code=500, detail=f"Failed to generate suggestions: {str(e)}")

def filter_dishes_by_restrictions(dishes: List[Dish], restrictions: Optional[List[str]]) -> List[Dish]:
    """Filter dishes based on dietary restrictions"""
    if not restrictions:
        return dishes
    
    filtered = []
    for dish in dishes:
        # Check if dish contains restricted ingredients
        dish_text = f"{dish.name} {' '.join(dish.ingredients)} {dish.nutrition or ''}".lower()
        
        should_exclude = False
        for restriction in restrictions:
            if restriction.lower() in dish_text:
                should_exclude = True
                break
        
        if not should_exclude:
            filtered.append(dish)
    
    return filtered

def generate_optimized_plan(dishes: List[Dish], days: int, daily_calories: float, budget_week: float):
    """Generate optimized meal plan using constraint programming"""
    model = cp_model.CpModel()
    D = len(dishes)
    T = days

    # Decision variables
    x = {}  # servings of dish i on day t
    for i in range(D):
        for t in range(T):
            x[i, t] = model.NewIntVar(0, 5, f"x_{i}_{t}")

    cook = {}  # whether dish i is cooked on day t
    for i in range(D):
        for t in range(T):
            cook[i, t] = model.NewBoolVar(f"cook_{i}_{t}")
            
            # Link cook and servings
            model.Add(x[i, t] > 0).OnlyEnforceIf(cook[i, t])
            model.Add(x[i, t] == 0).OnlyEnforceIf(cook[i, t].Not())

    # Calorie constraints per day
    for t in range(T):
        model.Add(sum(int(dishes[i].calories) * x[i, t] for i in range(D)) >= int(daily_calories * 0.9))
        model.Add(sum(int(dishes[i].calories) * x[i, t] for i in range(D)) <= int(daily_calories * 1.1))

    # Budget constraint
    model.Add(sum(int(dishes[i].cost * 100) * sum(x[i, t] for t in range(T)) for i in range(D)) <= int(budget_week * 100))

    # Minimize cooking sessions and prep time
    cook_sessions = model.NewIntVar(0, D * T, "cook_sessions")
    model.Add(cook_sessions == sum(cook[i, t] for i in range(D) for t in range(T)))
    model.Minimize(cook_sessions + sum(int(dishes[i].prep_minutes) * cook[i, t] for i in range(D) for t in range(T)))

    # Solve the model
    solver = cp_model.CpSolver()
    solver.parameters.max_time_in_seconds = 10.0
    status = solver.Solve(model)

    if status != cp_model.OPTIMAL and status != cp_model.FEASIBLE:
        raise Exception("No feasible solution found")

    # Extract results
    schedule = []
    for t in range(T):
        dish_ids = []
        servings = []
        total_calories = 0
        total_cost = 0
        
        for i in range(D):
            sv = solver.Value(x[i, t])
            if sv > 0:
                dish_ids.append(dishes[i].id)
                servings.append(sv)
                total_calories += dishes[i].calories * sv
                total_cost += dishes[i].cost * sv
        
        schedule.append(Day(
            day_index=t,
            dish_ids=dish_ids,
            servings=servings,
            total_calories=total_calories,
            total_cost=total_cost
        ))

    cook_days = [t for t in range(T) if any(solver.Value(cook[i, t]) == 1 for i in range(D))]
    
    # Generate shopping list
    shopping_list = generate_shopping_list(schedule, dishes)
    
    return schedule, cook_days, shopping_list

def generate_shopping_list(schedule: List[Day], dishes: List[Dish]) -> List[str]:
    """Generate consolidated shopping list from meal plan"""
    ingredient_counts = {}
    
    for day in schedule:
        for i, dish_id in enumerate(day.dish_ids):
            dish = next(d for d in dishes if d.id == dish_id)
            servings = day.servings[i]
            
            for ingredient in dish.ingredients:
                if ingredient in ingredient_counts:
                    ingredient_counts[ingredient] += servings
                else:
                    ingredient_counts[ingredient] = servings
    
    # Convert to shopping list format
    shopping_list = []
    for ingredient, count in ingredient_counts.items():
        if count > 1:
            shopping_list.append(f"{ingredient} (x{count})")
        else:
            shopping_list.append(ingredient)
    
    return shopping_list

def calculate_nutrition_summary(schedule: List[Day], dishes: List[Dish]) -> Dict[str, Any]:
    """Calculate nutrition summary for the meal plan"""
    total_calories = sum(day.total_calories for day in schedule)
    total_cost = sum(day.total_cost for day in schedule)
    avg_calories_per_day = total_calories / len(schedule)
    avg_cost_per_day = total_cost / len(schedule)
    
    return {
        "total_calories": total_calories,
        "total_cost": total_cost,
        "avg_calories_per_day": avg_calories_per_day,
        "avg_cost_per_day": avg_cost_per_day,
        "days_planned": len(schedule)
    }

async def generate_ai_recommendations(user_id: str, schedule: List[Day], dishes: List[Dish], preferences: Optional[Dict[str, Any]]) -> List[str]:
    """Generate AI-powered recommendations for meal planning"""
    try:
        # Create context for the AI
        context = f"""
        User ID: {user_id}
        Meal Plan Summary:
        - Days planned: {len(schedule)}
        - Total dishes: {sum(len(day.dish_ids) for day in schedule)}
        - Average calories per day: {sum(day.total_calories for day in schedule) / len(schedule):.0f}
        
        Dishes in plan: {[dishes[next(i for i, d in enumerate(dishes) if d.id == dish_id)].name for day in schedule for dish_id in day.dish_ids]}
        
        User preferences: {preferences or 'Not specified'}
        
        Please provide 3-5 specific, actionable recommendations to improve this meal plan.
        Focus on nutrition, variety, cost-effectiveness, and user satisfaction.
        """
        
        response = client.chat.completions.create(
            model="gpt-3.5-turbo",
            messages=[
                {"role": "system", "content": "You are a nutritionist and meal planning expert. Provide helpful, specific recommendations."},
                {"role": "user", "content": context}
            ],
            max_tokens=300,
            temperature=0.7
        )
        
        recommendations = response.choices[0].message.content.split('\n')
        return [rec.strip() for rec in recommendations if rec.strip() and not rec.startswith('-')]
        
    except Exception as e:
        logger.error(f"Error generating AI recommendations: {e}")
        return ["Consider adding more vegetables to your meals", "Try to include a variety of protein sources"]

async def generate_ai_meal_suggestions(req: SuggestionRequest) -> tuple[List[str], str]:
    """Generate AI-powered meal suggestions"""
    try:
        context = f"""
        User preferences: {req.preferences}
        Dietary restrictions: {req.dietary_restrictions}
        Budget: ${req.budget}
        Days to plan: {req.days}
        
        Please suggest 10 meal ideas that would work well for this user.
        Consider their preferences, restrictions, and budget.
        """
        
        response = client.chat.completions.create(
            model="gpt-3.5-turbo",
            messages=[
                {"role": "system", "content": "You are a culinary expert. Suggest diverse, delicious meals that fit the user's requirements."},
                {"role": "user", "content": context}
            ],
            max_tokens=500,
            temperature=0.8
        )
        
        suggestions = response.choices[0].message.content.split('\n')
        meal_suggestions = [s.strip() for s in suggestions if s.strip() and not s.startswith('-')]
        
        # Generate reasoning
        reasoning_response = client.chat.completions.create(
            model="gpt-3.5-turbo",
            messages=[
                {"role": "system", "content": "Explain why these meal suggestions are appropriate for the user."},
                {"role": "user", "content": f"Explain why these meals are good for someone with preferences {req.preferences} and restrictions {req.dietary_restrictions}"}
            ],
            max_tokens=200,
            temperature=0.7
        )
        
        reasoning = reasoning_response.choices[0].message.content
        
        return meal_suggestions, reasoning
        
    except Exception as e:
        logger.error(f"Error generating AI meal suggestions: {e}")
        return generate_fallback_suggestions(req)

def generate_fallback_suggestions(req: SuggestionRequest) -> SuggestionResponse:
    """Generate fallback suggestions without LLM"""
    suggestions = [
        "Grilled chicken with quinoa and roasted vegetables",
        "Lentil curry with brown rice",
        "Salmon with sweet potato and green beans",
        "Vegetarian stir-fry with tofu",
        "Mediterranean salad with chickpeas",
        "Turkey meatballs with whole wheat pasta",
        "Black bean tacos with avocado",
        "Greek yogurt parfait with berries",
        "Vegetable soup with whole grain bread",
        "Grilled fish with steamed vegetables"
    ]
    
    return SuggestionResponse(
        suggestions=suggestions,
        reasoning="These are balanced, nutritious meals that work well for most dietary preferences.",
        cuisine_variety=["Mediterranean", "Asian", "Mexican", "American"],
        nutrition_goals={"protein": "20-30g per meal", "fiber": "5-10g per meal", "calories": "400-600 per meal"}
    )

def analyze_cuisine_variety(suggestions: List[str]) -> List[str]:
    """Analyze cuisine variety in suggestions"""
    cuisines = set()
    for suggestion in suggestions:
        suggestion_lower = suggestion.lower()
        if any(word in suggestion_lower for word in ['curry', 'indian', 'masala']):
            cuisines.add('Indian')
        elif any(word in suggestion_lower for word in ['taco', 'mexican', 'salsa']):
            cuisines.add('Mexican')
        elif any(word in suggestion_lower for word in ['mediterranean', 'greek', 'olive']):
            cuisines.add('Mediterranean')
        elif any(word in suggestion_lower for word in ['stir-fry', 'asian', 'tofu']):
            cuisines.add('Asian')
        else:
            cuisines.add('American')
    
    return list(cuisines)

def generate_nutrition_goals(preferences: Dict[str, Any], restrictions: List[str]) -> Dict[str, Any]:
    """Generate nutrition goals based on preferences and restrictions"""
    goals = {
        "protein": "20-30g per meal",
        "fiber": "5-10g per meal",
        "calories": "400-600 per meal"
    }
    
    if "vegetarian" in restrictions:
        goals["protein"] = "15-25g per meal (plant-based)"
    
    if "low-carb" in preferences:
        goals["carbohydrates"] = "20-30g per meal"
    
    return goals

def calculate_savings(schedule: List[Day], budget_week: float) -> float:
    """Calculate estimated savings compared to eating out"""
    total_cost = sum(day.total_cost for day in schedule)
    estimated_eating_out_cost = len(schedule) * 15  # Assume $15 per meal eating out
    return estimated_eating_out_cost - total_cost

@app.on_event("startup")
async def startup_event():
    logger.info("Starting SpiceRoute Planner Service")

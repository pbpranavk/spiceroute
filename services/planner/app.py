from fastapi import FastAPI
from pydantic import BaseModel
from ortools.sat.python import cp_model

app = FastAPI()

class Dish(BaseModel):
    id: str
    prep_minutes: int
    calories: int
    shelf_life_days: int
    cost: float

class PlanRequest(BaseModel):
    days: int
    dishes: list[Dish]
    daily_calories: float
    budget_week: float

class Day(BaseModel):
    day_index: int
    dish_ids: list[str]
    servings: list[int]

class PlanResponse(BaseModel):
    schedule: list[Day]
    cook_days: list[int]
    shopping_list: list[str]

@app.post("/plan", response_model=PlanResponse)
def plan(req: PlanRequest):
    model = cp_model.CpModel()
    D = len(req.dishes)
    T = req.days

    x = {}
    for i in range(D):
        for t in range(T):
            x[i, t] = model.NewIntVar(0, 5, f"x_{i}_{t}")  # servings

    cook = {}
    for i in range(D):
        for t in range(T):
            cook[i, t] = model.NewBoolVar(f"cook_{i}_{t}")

            # link cook and servings
            model.Add(x[i, t] > 0).OnlyEnforceIf(cook[i, t])
            model.Add(x[i, t] == 0).OnlyEnforceIf(cook[i, t].Not())

    # calories per day
    for t in range(T):
        model.Add(sum(int(req.dishes[i].calories) * x[i, t] for i in range(D)) >= int(req.daily_calories * 0.9))
        model.Add(sum(int(req.dishes[i].calories) * x[i, t] for i in range(D)) <= int(req.daily_calories * 1.1))

    # budget
    model.Add(sum(int(req.dishes[i].cost * 100) * sum(x[i, t] for t in range(T)) for i in range(D))
              <= int(req.budget_week * 100))

    # minimize number of cook sessions and total prep time
    cook_sessions = model.NewIntVar(0, D * T, "cook_sessions")
    model.Add(cook_sessions == sum(cook[i, t] for i in range(D) for t in range(T)))
    model.Minimize(cook_sessions + sum(int(req.dishes[i].prep_minutes) * cook[i, t] for i in range(D) for t in range(T)))

    solver = cp_model.CpSolver()
    solver.parameters.max_time_in_seconds = 10.0
    status = solver.Solve(model)

    schedule = []
    for t in range(T):
        dish_ids = []
        servings = []
        for i in range(D):
            sv = solver.Value(x[i, t])
            if sv > 0:
                dish_ids.append(req.dishes[i].id)
                servings.append(sv)
        schedule.append(Day(day_index=t, dish_ids=dish_ids, servings=servings))

    cook_days = [t for t in range(T) if any(solver.Value(cook[i, t]) == 1 for i in range(D))]

    return PlanResponse(schedule=schedule, cook_days=cook_days, shopping_list=[])

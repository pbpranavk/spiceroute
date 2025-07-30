from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from playwright.async_api import async_playwright
import uuid

app = FastAPI()

SESSIONS = {}

class OrderReq(BaseModel):
    retailer: str
    items: list[str]
    confirm_token: str = None

@app.post("/order/initiate")
async def initiate_order(req: OrderReq):
    session_id = str(uuid.uuid4())
    SESSIONS[session_id] = {
        "retailer": req.retailer,
        "items": req.items,
        "status": "pending"
    }
    return {
        "session_id": session_id,
        "review_url": f"/order/review/{session_id}"
    }

@app.get("/order/review/{session_id}")
async def review_order(session_id: str):
    if session_id not in SESSIONS:
        raise HTTPException(status_code=404, detail="Session not found")

    order = SESSIONS[session_id]
    return {
        "retailer": order["retailer"],
        "items": order["items"],
        "message": "Click confirm to proceed with automated checkout",
        "confirm_url": f"/order/confirm/{session_id}"
    }

@app.post("/order/confirm/{session_id}")
async def confirm_order(session_id: str):
    if session_id not in SESSIONS:
        raise HTTPException(status_code=404, detail="Session not found")

    order = SESSIONS[session_id]
    retailer = order["retailer"]
    items = order["items"]

    async with async_playwright() as pw:
        browser = await pw.chromium.launch(headless=True)
        page = await browser.new_page()

        if retailer == "doordash":
            await page.goto("https://www.doordash.com")
            await page.screenshot(path=f"checkout_preview_{session_id}.png")

            for item in items:
                await page.keyboard.type(item)
                await page.keyboard.press("Enter")

        await browser.close()

    order["status"] = "submitted"
    return {
        "status": "order submitted",
        "screenshot": f"/screenshots/checkout_preview_{session_id}.png"
    }

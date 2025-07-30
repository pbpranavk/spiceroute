from fastapi import FastAPI
from pydantic import BaseModel
from playwright.async_api import async_playwright
import asyncio

app = FastAPI()

class OrderReq(BaseModel):
    retailer: str
    items: list[str]
    confirm_token: str

@app.post("/order")
async def order(req: OrderReq):
    # check token came from UI confirm
    async with async_playwright() as pw:
        browser = await pw.chromium.launch(headless=False)
        page = await browser.new_page()
        if req.retailer == "doordash":
            await page.goto("https://www.doordash.com")
            # ... fill cart. Use selectors you record and store
        await browser.close()
    return {"status": "submitted"}

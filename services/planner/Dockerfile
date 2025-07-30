FROM python:3.11-slim
RUN pip install fastapi uvicorn ortools pydantic
COPY app.py /app/app.py
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000"]

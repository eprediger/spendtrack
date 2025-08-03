from fastapi import FastAPI

from app.infrastructure.http.routers.health_router import health_router

app = FastAPI(
    title="Spendtrack backend",
    version="0.1.0",
)

app.include_router(health_router, prefix="/api/v1")

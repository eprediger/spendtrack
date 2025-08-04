from fastapi import FastAPI

from app.infrastructure.config.settings import get_settings
from app.infrastructure.http.routers.health_router import health_router

settings = get_settings()
app = FastAPI(
    title="Spendtrack backend",
    version=settings.version,
)

app.include_router(health_router, prefix="/api/v1")

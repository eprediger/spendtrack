from fastapi import APIRouter

from app.application.domain.entities.health import HealthStatus
from app.application.services.health_service import HealthService
from app.infrastructure.adapters.health_adapter import HealthAdapter
from app.infrastructure.config.settings import get_settings

health_router = APIRouter(tags=["health"])

service = HealthService(get_settings())
adapter = HealthAdapter(service)


@health_router.get("/health", response_model=HealthStatus)
def health_check():
    """Health check endpoint handler."""
    return adapter.get_health_status()

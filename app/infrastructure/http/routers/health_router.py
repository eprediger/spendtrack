from fastapi import APIRouter

from app.application.domain.entities.health import HealthStatus
from app.application.services.health_service import HealthService
from app.infrastructure.adapters.health_adapter import HealthAdapter

health_router = APIRouter(tags=["health"])

service = HealthService()
adapter = HealthAdapter(service)


@health_router.get("/health", response_model=HealthStatus)
def health_check():
    """Health check endpoint handler."""
    return adapter.get_health_status()

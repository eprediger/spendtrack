from app.application.domain.entities.health import HealthStatus
from app.application.ports.health_port import HealthPort


class HealthAdapter:
    """Application adapter for health check"""

    def __init__(self, service: HealthPort):
        self._service = service

    def get_health_status(self) -> HealthStatus:
        """Get current health status"""
        return self._service.get_health_status()

from datetime import UTC, datetime

from app.application.domain.entities.health import HealthStatus
from app.application.ports.health_port import HealthPort
from app.infrastructure.config.settings import Settings


class HealthService(HealthPort):
    """Implementation of the HealthPort"""

    def __init__(self, settings: Settings):
        self._settings = settings

    def get_health_status(self) -> HealthStatus:
        """Get current health status"""
        return HealthStatus(
            healthy=True,
            timestamp=datetime.now(UTC),
            version=self._settings.version,
            message="Application healthy",
        )

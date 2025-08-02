from datetime import UTC, datetime

from app.application.domain.entities.health import HealthStatus
from app.application.ports.health_port import HealthPort


class HealthService(HealthPort):
    """Implementation of the HealthPort"""

    def get_health_status(self) -> HealthStatus:
        """Get current health status"""
        return HealthStatus(
            healthy=True,
            timestamp=datetime.now(UTC),
            version="0.1.0",
            message="Application healthy",
        )

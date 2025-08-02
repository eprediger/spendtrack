from abc import ABC, abstractmethod

from app.application.domain.entities.health import HealthStatus


class HealthPort(ABC):
    """Port for health check operations."""

    @abstractmethod
    def get_health_status(self) -> HealthStatus:
        """Get health status."""

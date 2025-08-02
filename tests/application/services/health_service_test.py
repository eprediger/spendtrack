from datetime import datetime, UTC

import pytest

from app.application.domain.entities.health import HealthStatus
from app.application.services.health_service import HealthService


class TestHealthService:
    """Test cases for HealthService."""

    @pytest.fixture
    def health_service(self) -> HealthService:
        return HealthService()

    def test_get_health_status_returns_healthy(self, health_service):
        """Test that check_health calls the adapter."""
        result = health_service.get_health_status()

        assert result.healthy == True
        assert isinstance(result.timestamp, datetime)
        assert result.version == "0.1.0"
        assert result.message == "Application healthy"

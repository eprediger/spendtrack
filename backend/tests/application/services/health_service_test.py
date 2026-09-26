from datetime import datetime
from unittest.mock import Mock

import pytest

from app.application.services.health_service import HealthService


class TestHealthService:
    """Test cases for HealthService."""

    @pytest.fixture
    def mock_settings(self):
        """Mock settings."""
        return Mock()

    @pytest.fixture
    def health_service(self, mock_settings) -> HealthService:
        return HealthService(mock_settings)

    def test_get_health_status_returns_healthy(self, health_service, mock_settings):
        expected_version = "0.1.0"
        mock_settings.version = expected_version

        """Test that check_health calls the adapter."""
        result = health_service.get_health_status()

        assert result.healthy == True
        assert isinstance(result.timestamp, datetime)
        assert result.version == expected_version
        assert result.message == "Application healthy"

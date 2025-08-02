from datetime import datetime
from unittest.mock import Mock

import pytest

from app.application.domain.entities.health import HealthStatus
from app.infrastructure.adapters.health_adapter import HealthAdapter


class TestHealthAdapter:
    """Test cases for HealthAdapter."""

    @pytest.fixture
    def mock_health_port(self):
        """Mock health adapter fixture."""
        return Mock()

    @pytest.fixture
    def health_adapter(self, mock_health_port):
        """Health service fixture."""
        return HealthAdapter(mock_health_port)

    def test_check_health_calls_service(self, health_adapter, mock_health_port):
        """Test that get_health_status returns healthy status."""
        expected_status = HealthStatus(
            healthy=True,
            timestamp=datetime.now(),
            version="1.0.0",
            message="Test message"
        )
        mock_health_port.get_health_status.return_value = expected_status

        result = health_adapter.get_health_status()

        mock_health_port.get_health_status.assert_called_once()
        assert result.healthy == True
        assert result.version == "1.0.0"
        assert result.message == "Test message"
        assert isinstance(result.timestamp, datetime)

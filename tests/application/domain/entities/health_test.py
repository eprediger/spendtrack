from datetime import datetime

from app.application.domain.entities.health import HealthStatus


class TestHealthStatus:
    """Test cases for HealthStatus entity."""

    def test_health_status_creation(self):
        """Test HealthStatus entity creation."""
        timestamp = datetime.now()
        status = HealthStatus(
            healthy=True,
            timestamp=timestamp,
            version="1.0.0",
            message="Test message"
        )

        assert status.healthy == True
        assert status.timestamp == timestamp
        assert status.version == "1.0.0"
        assert status.message == "Test message"

    def test_health_status_without_message(self):
        """Test HealthStatus entity creation without message."""
        timestamp = datetime.now()
        status = HealthStatus(
            healthy=True,
            timestamp=timestamp,
            version="1.0.0"
        )

        assert status.healthy == True
        assert status.timestamp == timestamp
        assert status.version == "1.0.0"
        assert status.message is None

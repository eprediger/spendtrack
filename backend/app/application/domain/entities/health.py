from datetime import datetime

from pydantic import BaseModel


class HealthStatus(BaseModel):
    """Health status entity."""

    healthy: bool
    timestamp: datetime
    version: str
    message: str | None = None

    class ConfigDict:
        from_attributes = True

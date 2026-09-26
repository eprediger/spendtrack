from functools import lru_cache

from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    version: str = "0.1.0"


@lru_cache
def get_settings() -> Settings:
    return Settings()

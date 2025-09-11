from pydantic import PostgresDsn
from pydantic_settings import BaseSettings
from typing import Optional

class Settings(BaseSettings):
    APP_NAME: str
    DEBUG: bool
    PORT: int

    POSTGRES_USER: str
    POSTGRES_PASSWORD: str
    POSTGRES_DB: str
    POSTGRES_PORT: int
    POSTGRES_HOST: str

    DATABASE_URL: Optional[PostgresDsn] = None

    class Config:
        env_file = ".env"

settings = Settings()

# Build DATABASE_URL if not explicitly provided
if not settings.DATABASE_URL:
    settings.DATABASE_URL = (
        f"postgresql+asyncpg://{settings.POSTGRES_USER}:{settings.POSTGRES_PASSWORD}"
        f"@{settings.POSTGRES_HOST}:{settings.POSTGRES_PORT}/{settings.POSTGRES_DB}"
    )

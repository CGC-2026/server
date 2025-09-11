from pydantic import field_validator
from pydantic_settings import BaseSettings, SettingsConfigDict

class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        extra="ignore",
    )

    APP_NAME: str
    DEBUG: bool = False
    PORT: int = 8000

    POSTGRES_USER: str
    POSTGRES_PASSWORD: str
    POSTGRES_DB: str
    POSTGRES_PORT: int = 5432
    POSTGRES_HOST: str = "localhost"

    DATABASE_URL: str | None = None

    @field_validator("DATABASE_URL", mode="before")
    @classmethod
    def assemble_db_url(cls, v, info):
        if v:  # already provided via env
            return v
        d = info.data
        return (
            f"postgresql+asyncpg://{d['POSTGRES_USER']}:{d['POSTGRES_PASSWORD']}"
            f"@{d['POSTGRES_HOST']}:{d['POSTGRES_PORT']}/{d['POSTGRES_DB']}"
        )

    @field_validator("DATABASE_URL")
    @classmethod
    def validate_scheme(cls, v):
        if not v.startswith(("postgresql+asyncpg://", "postgresql://", "postgres://")):
            raise ValueError(
                "DATABASE_URL must start with one of: postgresql+asyncpg://, postgresql://, postgres://"
            )
        return v

settings = Settings()

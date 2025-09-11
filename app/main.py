from contextlib import asynccontextmanager
from fastapi import FastAPI
from app.db import engine, Base
from app import settings

app = FastAPI(title=settings.APP_NAME)

@asynccontextmanager
async def lifespan(app: FastAPI):
    # start up
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    yield
    # Clean up here (after)
    

@app.get("/")
async def root():
    return {"message": "Hello World"}
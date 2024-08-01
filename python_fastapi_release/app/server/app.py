from fastapi import FastAPI
from app.configs.config import settings
from starlette.middleware.cors import CORSMiddleware


async def startup(): ...


async def shutdown(): ...


def create_app() -> FastAPI:
    _app = FastAPI(title=settings.APP_NAME, version=settings.API_VERSION)
    _app.add_middleware(
        CORSMiddleware,
        allow_origins=settings.BACKEND_CORS_ORIGINS,
        allow_methods=["*"],
        allow_headers=["*"],
    )
    _app.add_event_handler("startup", startup)
    _app.add_event_handler("shutdown", shutdown)
    return _app

from asyncio import current_task
from collections.abc import AsyncIterator
from contextlib import asynccontextmanager

from sqlalchemy.ext.asyncio import (
    AsyncSession, 
    async_sessionmaker,
    async_scoped_session, 
    create_async_engine
)
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()


class Database:
    def __init__(self, dns: str) -> None:
        self._engine = create_async_engine(dns, pool_pre_ping=True, pool_size=1, max_overflow=0)

        session_factory = async_sessionmaker(
            bind=self._engine,
            autocommit=False,
            autoflush=False,
            expire_on_commit=False,
            class_=AsyncSession,
        )
        self._session_factory = async_scoped_session(session_factory, scopefunc=current_task)

    @asynccontextmanager
    async def session(self) -> AsyncIterator[AsyncSession]:
        session: AsyncSession = self._session_factory()
        try:
            yield session
        except Exception:
            await session.rollback()
            raise
        finally:
            await session.close()

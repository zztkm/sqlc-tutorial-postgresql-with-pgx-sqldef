import os

from fastapi import FastAPI
from fastapi.responses import JSONResponse

from gen.sqlc.query import AsyncQuerier
from gen.sqlc.models import Author
from db import Database

app = FastAPI()

postgres_uri = os.environ["DNS"].replace("postgresql", "postgresql+asyncpg")
db = Database(postgres_uri)


@app.get("/")
async def root():
    return {"message": "Hello World"}


@app.get("/authors/{id}")
async def user(id: int) -> Author | JSONResponse:
    async with db.session() as session:
        conn = await session.connection()
        querier = AsyncQuerier(conn=conn)

        author = await querier.get_author(id=id)
        if author:
            return author
        else:
            return JSONResponse(status_code=404, content={"message": "author not found"})


@app.get("/authors")
async def users() -> list[Author]:
    async with db.session() as session:
        conn = await session.connection()
        querier = AsyncQuerier(conn=conn)

        res: list[Author] = []
        async for author in querier.list_authors():
            res.append(author)
        return res


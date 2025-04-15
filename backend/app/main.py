from fastapi import FastAPI, Depends
from fastapi.middleware.cors import CORSMiddleware
import psycopg
from sqlalchemy.orm import Session

from app.models.database import create_tables, get_db
from app.routes import auth, hotels, rooms, bookings

app = FastAPI(title="Hotel Extranet API")

# Disable CORS. Do not remove this for full-stack development.
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Allows all origins
    allow_credentials=True,
    allow_methods=["*"],  # Allows all methods
    allow_headers=["*"],  # Allows all headers
)

app.include_router(auth.router, prefix="/auth", tags=["authentication"])
app.include_router(hotels.router, prefix="/hotels", tags=["hotels"])
app.include_router(rooms.router, prefix="/rooms", tags=["rooms"])
app.include_router(bookings.router, prefix="/bookings", tags=["bookings"])

create_tables()

@app.get("/healthz", tags=["health"])
async def healthz():
    return {"status": "ok"}

@app.get("/", tags=["root"])
async def root():
    return {
        "message": "Welcome to Hotel Extranet API",
        "version": "1.0.0",
        "docs": "/docs"
    }

@app.on_event("startup")
async def startup_event():
    pass

from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from typing import List, Optional
from pydantic import BaseModel
from datetime import date, datetime

from app.models.database import Room, Hotel, User, get_db
from app.services.auth import get_current_user, get_current_admin

router = APIRouter()

class RoomBase(BaseModel):
    name: str
    description: Optional[str] = None
    price_weekday: float
    price_weekend: float

class RoomCreate(RoomBase):
    hotel_id: int

class RoomUpdate(RoomBase):
    pass

class RoomResponse(RoomBase):
    id: int
    hotel_id: int
    created_at: datetime
    updated_at: datetime

    class Config:
        from_attributes = True

class BulkPriceUpdate(BaseModel):
    room_ids: List[int]
    start_date: date
    end_date: date
    price_weekday: Optional[float] = None
    price_weekend: Optional[float] = None

@router.get("", response_model=List[RoomResponse])
async def get_rooms(
    skip: int = 0, 
    limit: int = 100,
    hotel_id: Optional[int] = None,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    query = db.query(Room)
    
    if hotel_id:
        query = query.filter(Room.hotel_id == hotel_id)
    
    if current_user.role != "admin":
        hotels = db.query(Hotel).filter(Hotel.user_id == current_user.id).all()
        hotel_ids = [hotel.id for hotel in hotels]
        query = query.filter(Room.hotel_id.in_(hotel_ids))
    
    rooms = query.offset(skip).limit(limit).all()
    return rooms

@router.put("/bulk-price-update", response_model=dict)
async def bulk_update_prices(
    update_data: BulkPriceUpdate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    updated_count = 0
    
    for room_id in update_data.room_ids:
        room = db.query(Room).filter(Room.id == room_id).first()
        
        if not room:
            continue
        
        if current_user.role != "admin":
            hotel = db.query(Hotel).filter(Hotel.id == room.hotel_id).first()
            if hotel.user_id != current_user.id:
                continue
        
        if update_data.price_weekday is not None:
            room.price_weekday = update_data.price_weekday
        
        if update_data.price_weekend is not None:
            room.price_weekend = update_data.price_weekend
        
        updated_count += 1
    
    db.commit()
    
    return {
        "message": f"Successfully updated {updated_count} rooms",
        "updated_count": updated_count
    }


@router.get("/{room_id}", response_model=RoomResponse)
async def get_room(
    room_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    room = db.query(Room).filter(Room.id == room_id).first()
    
    if not room:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Room not found"
        )
    
    if current_user.role != "admin":
        hotel = db.query(Hotel).filter(Hotel.id == room.hotel_id).first()
        if hotel.user_id != current_user.id:
            raise HTTPException(
                status_code=status.HTTP_403_FORBIDDEN,
                detail="Not enough permissions to access this room"
            )
    
    return room

@router.post("", response_model=RoomResponse, status_code=status.HTTP_201_CREATED)
async def create_room(
    room_data: RoomCreate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    hotel = db.query(Hotel).filter(Hotel.id == room_data.hotel_id).first()
    
    if not hotel:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Hotel not found"
        )
    
    if current_user.role != "admin" and hotel.user_id != current_user.id:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not enough permissions to add rooms to this hotel"
        )
    
    db_room = Room(
        hotel_id=room_data.hotel_id,
        name=room_data.name,
        description=room_data.description,
        price_weekday=room_data.price_weekday,
        price_weekend=room_data.price_weekend
    )
    
    db.add(db_room)
    db.commit()
    db.refresh(db_room)
    
    return db_room

@router.put("/{room_id}", response_model=RoomResponse)
async def update_room(
    room_id: int,
    room_data: RoomUpdate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    room = db.query(Room).filter(Room.id == room_id).first()
    
    if not room:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Room not found"
        )
    
    if current_user.role != "admin":
        hotel = db.query(Hotel).filter(Hotel.id == room.hotel_id).first()
        if hotel.user_id != current_user.id:
            raise HTTPException(
                status_code=status.HTTP_403_FORBIDDEN,
                detail="Not enough permissions to update this room"
            )
    
    room.name = room_data.name
    room.description = room_data.description
    room.price_weekday = room_data.price_weekday
    room.price_weekend = room_data.price_weekend
    
    db.commit()
    db.refresh(room)
    
    return room

@router.delete("/{room_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_room(
    room_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    room = db.query(Room).filter(Room.id == room_id).first()
    
    if not room:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Room not found"
        )
    
    if current_user.role != "admin":
        hotel = db.query(Hotel).filter(Hotel.id == room.hotel_id).first()
        if hotel.user_id != current_user.id:
            raise HTTPException(
                status_code=status.HTTP_403_FORBIDDEN,
                detail="Not enough permissions to delete this room"
            )
    
    db.delete(room)
    db.commit()
    
    return None

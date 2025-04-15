from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from typing import List, Optional
from pydantic import BaseModel, Field
from datetime import datetime

from app.models.database import Hotel, User, get_db
from app.services.auth import get_current_user, get_current_admin

router = APIRouter()

class HotelBase(BaseModel):
    name: str
    location: str
    description: Optional[str] = None

class HotelCreate(HotelBase):
    pass

class HotelUpdate(HotelBase):
    pass

class HotelResponse(HotelBase):
    id: int
    user_id: int
    created_at: datetime
    updated_at: datetime

    class Config:
        from_attributes = True

@router.get("", response_model=List[HotelResponse])
async def get_hotels(
    skip: int = 0, 
    limit: int = 100, 
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    if current_user.role == "admin":
        hotels = db.query(Hotel).offset(skip).limit(limit).all()
    else:
        hotels = db.query(Hotel).filter(Hotel.user_id == current_user.id).offset(skip).limit(limit).all()
    
    return hotels

@router.get("/{hotel_id}", response_model=HotelResponse)
async def get_hotel(
    hotel_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    hotel = db.query(Hotel).filter(Hotel.id == hotel_id).first()
    
    if not hotel:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Hotel not found"
        )
    
    if current_user.role != "admin" and hotel.user_id != current_user.id:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not enough permissions to access this hotel"
        )
    
    return hotel

@router.post("", response_model=HotelResponse, status_code=status.HTTP_201_CREATED)
async def create_hotel(
    hotel_data: HotelCreate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    db_hotel = Hotel(
        user_id=current_user.id,
        name=hotel_data.name,
        location=hotel_data.location,
        description=hotel_data.description
    )
    
    db.add(db_hotel)
    db.commit()
    db.refresh(db_hotel)
    
    return db_hotel

@router.put("/{hotel_id}", response_model=HotelResponse)
async def update_hotel(
    hotel_id: int,
    hotel_data: HotelUpdate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    hotel = db.query(Hotel).filter(Hotel.id == hotel_id).first()
    
    if not hotel:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Hotel not found"
        )
    
    if current_user.role != "admin" and hotel.user_id != current_user.id:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not enough permissions to update this hotel"
        )
    
    hotel.name = hotel_data.name
    hotel.location = hotel_data.location
    hotel.description = hotel_data.description
    
    db.commit()
    db.refresh(hotel)
    
    return hotel

@router.delete("/{hotel_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_hotel(
    hotel_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    hotel = db.query(Hotel).filter(Hotel.id == hotel_id).first()
    
    if not hotel:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Hotel not found"
        )
    
    if current_user.role != "admin" and hotel.user_id != current_user.id:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not enough permissions to delete this hotel"
        )
    
    db.delete(hotel)
    db.commit()
    
    return None

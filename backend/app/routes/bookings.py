from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from typing import List, Optional
from pydantic import BaseModel
from datetime import date, datetime, timedelta

from app.models.database import Booking, Room, Hotel, User, get_db
from app.services.auth import get_current_user, get_current_admin

router = APIRouter()

class BookingBase(BaseModel):
    room_id: int
    checkin_date: date
    checkout_date: date

class BookingCreate(BookingBase):
    pass

class BookingUpdate(BaseModel):
    checkin_date: Optional[date] = None
    checkout_date: Optional[date] = None

class BookingResponse(BookingBase):
    id: int
    user_id: int
    total_price: float
    created_at: datetime
    updated_at: datetime

    class Config:
        from_attributes = True

def calculate_total_price(room, checkin_date, checkout_date):
    delta = checkout_date - checkin_date
    days = delta.days
    
    if days <= 0:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Checkout date must be after checkin date"
        )
    
    weekday_count = 0
    weekend_count = 0
    
    current_date = checkin_date
    while current_date < checkout_date:
        if current_date.weekday() >= 5:
            weekend_count += 1
        else:
            weekday_count += 1
        current_date += timedelta(days=1)
    
    total_price = (weekday_count * room.price_weekday) + (weekend_count * room.price_weekend)
    
    return total_price

@router.get("", response_model=List[BookingResponse])
async def get_bookings(
    skip: int = 0, 
    limit: int = 100,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    if current_user.role == "admin":
        bookings = db.query(Booking).offset(skip).limit(limit).all()
    else:
        hotels = db.query(Hotel).filter(Hotel.user_id == current_user.id).all()
        hotel_ids = [hotel.id for hotel in hotels]
        
        rooms = db.query(Room).filter(Room.hotel_id.in_(hotel_ids)).all()
        room_ids = [room.id for room in rooms]
        
        bookings = db.query(Booking).filter(Booking.room_id.in_(room_ids)).offset(skip).limit(limit).all()
    
    return bookings

@router.get("/{booking_id}", response_model=BookingResponse)
async def get_booking(
    booking_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    booking = db.query(Booking).filter(Booking.id == booking_id).first()
    
    if not booking:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Booking not found"
        )
    
    if current_user.role != "admin" and booking.user_id != current_user.id:
        room = db.query(Room).filter(Room.id == booking.room_id).first()
        hotel = db.query(Hotel).filter(Hotel.id == room.hotel_id).first()
        
        if hotel.user_id != current_user.id:
            raise HTTPException(
                status_code=status.HTTP_403_FORBIDDEN,
                detail="Not enough permissions to access this booking"
            )
    
    return booking

@router.post("", response_model=BookingResponse, status_code=status.HTTP_201_CREATED)
async def create_booking(
    booking_data: BookingCreate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    room = db.query(Room).filter(Room.id == booking_data.room_id).first()
    
    if not room:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Room not found"
        )
    
    total_price = calculate_total_price(
        room, 
        booking_data.checkin_date, 
        booking_data.checkout_date
    )
    
    db_booking = Booking(
        user_id=current_user.id,
        room_id=booking_data.room_id,
        checkin_date=booking_data.checkin_date,
        checkout_date=booking_data.checkout_date,
        total_price=total_price
    )
    
    db.add(db_booking)
    db.commit()
    db.refresh(db_booking)
    
    return db_booking

@router.put("/{booking_id}", response_model=BookingResponse)
async def update_booking(
    booking_id: int,
    booking_data: BookingUpdate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    booking = db.query(Booking).filter(Booking.id == booking_id).first()
    
    if not booking:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Booking not found"
        )
    
    if current_user.role != "admin" and booking.user_id != current_user.id:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not enough permissions to update this booking"
        )
    
    if booking_data.checkin_date:
        booking.checkin_date = booking_data.checkin_date
    
    if booking_data.checkout_date:
        booking.checkout_date = booking_data.checkout_date
    
    if booking_data.checkin_date or booking_data.checkout_date:
        room = db.query(Room).filter(Room.id == booking.room_id).first()
        booking.total_price = calculate_total_price(
            room, 
            booking.checkin_date, 
            booking.checkout_date
        )
    
    db.commit()
    db.refresh(booking)
    
    return booking

@router.delete("/{booking_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_booking(
    booking_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    booking = db.query(Booking).filter(Booking.id == booking_id).first()
    
    if not booking:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Booking not found"
        )
    
    if current_user.role != "admin" and booking.user_id != current_user.id:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not enough permissions to delete this booking"
        )
    
    db.delete(booking)
    db.commit()
    
    return None

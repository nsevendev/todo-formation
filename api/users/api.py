from ninja import Router
from .models import User
from django.shortcuts import get_object_or_404
from pydantic import BaseModel
import uuid

router = Router()

# Schéma de validation Pydantic pour User
class UserSchema(BaseModel):
    id: uuid.UUID
    pseudo: str

class UserCreateSchema(BaseModel):
    pseudo: str
    password: str

@router.get("/", response=list[UserSchema])
def list_users(request):
    return list(User.objects.all())

@router.post("/", response=UserSchema)
def create_user(request, payload: UserCreateSchema):
    user = User.objects.create(**payload.dict())
    return user

@router.get("/{user_id}", response=UserSchema)
def get_user(request, user_id: uuid.UUID):
    user = get_object_or_404(User, id=user_id)
    return user

@router.delete("/{user_id}")
def delete_user(request, user_id: uuid.UUID):
    user = get_object_or_404(User, id=user_id)
    user.delete()
    return {"message": "User deleted"}
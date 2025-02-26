from ninja import Router
from .models import Task
from users.models import User
from django.shortcuts import get_object_or_404
from pydantic import BaseModel
import uuid

router = Router()

# Schéma de validation Pydantic pour Task
class TaskSchema(BaseModel):
    id: uuid.UUID
    libelle: str
    done: bool
    user_id: uuid.UUID

class TaskCreateSchema(BaseModel):
    libelle: str
    user_id: uuid.UUID

@router.get("/", response=list[TaskSchema])
def list_tasks(request):
    return list(Task.objects.all())

@router.post("/", response=TaskSchema)
def create_task(request, payload: TaskCreateSchema):
    user = get_object_or_404(User, id=payload.user_id)
    task = Task.objects.create(libelle=payload.libelle, user=user)
    return task

@router.get("/{task_id}", response=TaskSchema)
def get_task(request, task_id: uuid.UUID):
    task = get_object_or_404(Task, id=task_id)
    return task

@router.delete("/{task_id}")
def delete_task(request, task_id: uuid.UUID):
    task = get_object_or_404(Task, id=task_id)
    task.delete()
    return {"message": "Task deleted"}
from ninja import Router
from .models import ParamApp
from django.shortcuts import get_object_or_404
from pydantic import BaseModel
import uuid

router = Router()

# Schéma de validation Pydantic pour ParamApp
class ParamAppSchema(BaseModel):
    id: uuid.UUID
    title_app: str
    copyright: str

@router.get("/", response=list[ParamAppSchema])
def list_config(request):
    return list(ParamApp.objects.all())

@router.get("/{config_id}", response=ParamAppSchema)
def get_config(request, config_id: uuid.UUID):
    config = get_object_or_404(ParamApp, id=config_id)
    return config
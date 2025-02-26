import uuid
from django.db import models

class ParamApp(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    title_app = models.CharField(max_length=10)
    copyright = models.CharField(max_length=10)

    def __str__(self):
        return self.title_app
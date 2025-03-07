"""
URL configuration for config project.

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/5.1/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
from django.contrib import admin
from django.urls import path
from ninja import NinjaAPI

api = NinjaAPI()

# Import des routes de chaque app
from users.api import router as users_router
from tasks.api import router as tasks_router
from paramsApp.api import router as params_app_router

# Ajout des route dans ninja
api.add_router("/users/", users_router)
api.add_router("/tasks/", tasks_router)
api.add_router("/param-app/", params_app_router)

urlpatterns = [
    path('admin/', admin.site.urls),
    path('api/', api.urls),
]

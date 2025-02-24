# Dockerfile
FROM python:3.13.2-slim
RUN apt-get update && apt-get install -y gcc libpq-dev
WORKDIR /app/api
COPY ./requirements/base.txt /app/api/requirements/base.txt
RUN pip install --upgrade pip && pip install -r /app/api/requirements/base.txt
COPY . .
CMD ["sh", "-c", "python manage.py migrate && python manage.py runserver 0.0.0.0:80"]
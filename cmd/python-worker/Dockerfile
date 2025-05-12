FROM python:3.10-slim
# FROM tensorflow/tensorflow:2.15.0-py3

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

RUN apt-get update && \
    apt-get install -y build-essential libglib2.0-0 libsm6 libxext6 libxrender-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . /app
COPY ./models /app/models

RUN pip install --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt

ENV TF_CPP_MIN_LOG_LEVEL=2

RUN ls -l /app/models

CMD ["python", "sqs_consumer.py"]

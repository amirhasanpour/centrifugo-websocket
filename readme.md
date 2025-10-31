# Real-time Chat Application with Microservices Architecture
Distributed real-time chat application built with Go microservices architecture using gRPC for inter-service communication. Features Centrifugo for WebSocket-based real-time messaging, PostgreSQL for persistent data storage, and Fiber for high-performance HTTP routing. Containerized all services with Docker and orchestrated deployments using Docker Compose for development and testing.

---

## Features and Tools

- **Real-time messaging with WebSocket connections powered by [Centrifugo](https://github.com/centrifugal/centrifugo) for instant message delivery**
- **Microservices architecture with independent, scalable services for authentication and chat functionality**
- **[gRPC](https://github.com/grpc/grpc-go) communication for fast, efficient, and type-safe inter-service messaging between microservices**
- **RESTful API Gateway using [Fiber](https://github.com/gofiber/fiber) framework for high-performance HTTP routing and request handling**
- **[JWT-based](https://github.com/golang-jwt/jwt) authentication service for secure user management, registration, and token validation**
- **Room-based chat system supporting multiple chat rooms with membership management**
- **Message persistence with full chat history storage and retrieval capabilities**
- **[PostgreSQL](https://github.com/postgres/postgres) as the primary relational database for user data, rooms, and message storage**
- **[Centrifugo](https://github.com/centrifugal/centrifugo) as the real-time engine for WebSocket connections and publish-subscribe messaging**
- **[Docker](https://github.com/docker/compose) containerization for all microservices ensuring isolated and reproducible environments**
- **[Docker Compose](https://github.com/docker/compose) for orchestration and management of all containerized services**
- **[Fiber](https://github.com/gofiber/fiber) framework for building high-performance HTTP APIs with minimal memory footprint**
- **[GORM](https://github.com/go-gorm/gorm) as the ORM for database operations with PostgreSQL**
- **[Protocol Buffers](https://github.com/protocolbuffers/protobuf) for defining service contracts and generating gRPC client/server code**

---

## How to run?

### Using Docker Compose

cd to the project directory and run this command:

```bash
docker-compose up --build -d
```

to stop all services:

```bash
docker-compose down
```

to stop services and remove database volumes:

```bash
docker-compose down -v
```

---

## Testing the Application

### Register User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login User

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Validate Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/validate \
  -H "Content-Type: application/json" \
  -d '{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

### Get User Profile

```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Create Room

```bash
curl -X POST http://localhost:8080/api/v1/chat/rooms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "name": "General Chat",
    "description": "General discussion room for everyone"
  }'
```

### Get All Rooms

```bash
curl -X GET http://localhost:8080/api/v1/chat/rooms \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Join Room

```bash
curl -X POST http://localhost:8080/api/v1/chat/rooms/join \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "room_id": "room-1234-5678-90ab-cdef12345678"
  }'
```

### Send Message

```bash
curl -X POST http://localhost:8080/api/v1/chat/messages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "room_id": "room-1234-5678-90ab-cdef12345678",
    "content": "Hello everyone! How are you doing today?"
  }'
```

### Get Room Messages

```bash
curl -X GET "http://localhost:8080/api/v1/chat/rooms/room-1234-5678-90ab-cdef12345678/messages?limit=20" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Health Check

```bash
curl -X GET http://localhost:8080/api/v1/health
```
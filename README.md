# Spy Cat Agency

A CRUD application for managing spy cats, missions, targets, and notes for the Spy Cat Agency. This project demonstrates building RESTful APIs in Go using the Gin framework, interacting with a PostgreSQL database via GORM, integrating with third-party services (TheCatAPI for breed validation), and containerizing the database with Docker Compose.

## Overview

The Spy Cat Agency system allows you to:
- **Manage Spy Cats**:  
  Create, list, retrieve, update, and delete spy cats.  
  Each cat is described by:
    - **Name**
    - **Years of Experience**
    - **Breed** (validated via [TheCatAPI](https://api.thecatapi.com/v1/breeds))
    - **Salary**
- **Manage Missions & Targets**:
    - **Missions**: Create a mission for a spy cat, including 1–3 targets.  
      Each mission stores the assigned cat, target details, and its completion state.
    - **Targets**: Each target (unique to a mission) includes:
        - **Name**
        - **Country**
        - **Notes**
        - **Status** ("ONGOING" or "COMPLETED")

  Additional rules:
    - A mission cannot be deleted if it is assigned to a cat.
    - New targets cannot be added to a completed mission.
    - A target cannot be deleted if it is completed.
    - Completing all targets in a mission automatically marks the mission as completed.
    - A cat can only have one ongoing mission at a time.
- **Manage Notes**:
    - Create and update notes for targets.
    - Note updates are disallowed if the target or its associated mission is completed.

- **General Features**:
    - Uses **Gin** as the web framework.
    - Uses **GORM** for database operations (PostgreSQL, dockerized).
    - Validates request payloads and returns appropriate HTTP status codes.
    - Integrates TheCatAPI for breed validation.
    - Includes logging middleware (via Gin).

## Directory Structure

```plaintext
spy_cats/
├── cmd/
│   ├── main.go              # Application entry point
│   └── docs/                     # Documentation files (or additional command tools)
│              
├── config
│   └── config.go            # Configuration and environment variables
├── internal
│   ├── cat                  # Spy Cat domain
│   │   ├── cat.go
│   │   ├── cat_handler.go
│   │   ├── cat_repository.go
│   │   └── cat_service.go
│   ├── mission              # Mission domain
│   │   ├── mission.go
│   │   ├── mission_handler.go
│   │   ├── mission_repository.go
│   │   └── mission_service.go
│   ├── note                 # Note domain
│   │   ├── note.go
│   │   ├── note_handler.go
│   │   ├── note_repository.go
│   │   └── note_service.go
│   └── target               # Target domain
│       ├── target.go
│   │   ├── target_handler.go
│       ├── target_repository.go
│       └── target_service.go
├── pkg
│   ├── database             # Database connection & migration logic
│   │   ├── migrations       # SQL migration files
│   │   └── db.go
│   └── router               # Route setup
│       └── router.go
├── docker-compose.yml
├── Dockerfile
├── README.md
├── go.mod
└── go.sum
```



## Getting Started
### Prerequisites

    Go (1.24.1)
    Docker & Docker Compose
    Git


## Setup Instructions
### 1. Clone the Repository:
```
git clone https://github.com/genryusaishigikuni/spy_cats.git
cd spy_cats
```

### 2. Start the Database and the Application with Docker Compose:
```
docker-compose up -d
```
or depending on your docker compose version

```
docker compose up -d
```



## API Documentation
### It's provided as four generated json files from postman collections in cmd/docs directory 

## Environment Variables
```
DB_HOST – Database host (default: localhost)
DB_PORT – Database port (default: 5432) (5431:5432 in docker-compose file)
DB_USER – Database user (e.g., catadmin)
DB_PASSWORD – Database password
DB_NAME – Database name (e.g., spycatsdb)
SERVER_PORT – API port (default: :8080)
THECATAPI_KEY (optional) – API key for TheCatAPI (if required)
```
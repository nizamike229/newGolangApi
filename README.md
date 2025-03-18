# REST API Project :rocket:
This is a RESTful API built with Go, featuring CRUD operations, authentication, authorization, and a custom logger. The project is designed to manage Person entities and secure access with JWT-based authentication.

# Features :sparkles:
- CRUD Operations :wrench: - Create, Read, and Delete operations for the Person model.
- Authentication :key: - User login with JWT tokens stored in cookies.
- Authorization :shield: - Role-based access control with middleware.
- Custom Logger :loudspeaker: - Custom logging system with INFO, WARNING, and ERROR levels, including timestamps.
- Database :floppy_disk: - PostgreSQL with GORM for ORM.
- Framework :gear: - Built using go-chi for routing.
# Tech Stack :computer:
- Language: Go :golang:
- outer: go-chi/chi/v5 :arrow_right:
- ORM: gorm.io/gorm with gorm.io/driver/postgres :database:
- JWT: github.com/golang-jwt/jwt/v5 :lock:
- Logger: Custom implementation over fmt.Println :speech_balloon:
- Database: PostgreSQL :elephant:

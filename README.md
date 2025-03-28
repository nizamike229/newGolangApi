# [![My Skills](https://skillicons.dev/icons?i=go)](https://skillicons.dev)  REST API Project :rocket: 
This is a RESTful TO-DO API built with Go, featuring CRUD operations, authentication, authorization, and a custom logger. The project is designed to manage Task entities and secure access with JWT-based authentication.

# Features :sparkles:
- CRUD Operations :wrench: - Create, Read, and Delete operations for the Task model.
- Authentication :key: - User login with JWT tokens stored in cookies.
- Authorization :shield: - Role-based access control with middleware.
- Custom Logger :loudspeaker: - Custom logging system with INFO, WARNING, and ERROR levels, including timestamps.
- Database :floppy_disk: - PostgreSQL with GORM for ORM.
- Framework :gear: - Built using go-chi for routing.
# Tech Stack :computer:
- Language: Go
- outer: go-chi/chi/v5
- ORM: gorm.io/gorm with gorm.io/driver/postgres
- JWT: github.com/golang-jwt/jwt/v5
- Logger: Custom implementation over fmt.Println
- Database: PostgreSQL

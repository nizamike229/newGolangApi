# REST API Project
This is a RESTful API built with Go, featuring CRUD operations, authentication, authorization, and a custom logger. The project is designed to manage Person entities and secure access with JWT-based authentication.

# Features
- CRUD Operations: Create, Read, and Delete operations for the Person model.
- Authentication: User login with JWT tokens stored in cookies.
- Authorization: Role-based access control with middleware.
- Custom Logger: Custom logging system with INFO, WARNING, and ERROR levels, including timestamps.
- Database: PostgreSQL with GORM for ORM.
- Framework: Built using go-chi for routing.
# Tech Stack
- Language: Go
- outer: go-chi/chi/v5
- ORM: gorm.io/gorm with gorm.io/driver/postgres
- JWT: github.com/golang-jwt/jwt/v5
- Logger: Custom implementation over fmt.Println
- Database: PostgreSQL

## Microservice based app in golang

-- The application demonstrates how to use microservices to create a scalable and maintainable architecture

# Architecture / Services

- API Gateway: The entry point for all client requests.
- Authentication Service: Handles user authentication.
- Logger Service: Page application events.
- Mailer Service: Sends emails.
- Listener Service: Listens for events from RabbitMQ and processes them.

# Tech stack

- Golang: The primary language.
- Docker: Used for containerizing the services.
- Docker Compose: Used for orchestrating the multi-container Docker application.
- RabbitMQ: Used for messaging between services.
- PostgreSQL: The database used by the Authentication Service.
- MongoDB: The database used by the Logger Service.
- gRPC: Communication method between api-gateway and other services where it needs async calls.

# BoulderLog

BoulderLog is a web application for climbers to log and track their bouldering sessions.

## Goals

- User can easily and quickly record each route between climbs
- User should see their progress visually
- Minimalist design, optimized for mobile use
- Possibly fun gamified elements

## Features

- User authentication
- Log boulder attempts with grade, difficulty, and additional details
- View daily progress and statistics with visualizations
- Download boulder logs as CSV
- Dark mode support

## Application Architecture

BoulderLog follows a layered architecture pattern:

1. **Presentation Layer**: 
   - Handlers (`handlers/`) handle HTTP requests and responses.
   - Templates (`components/`) render the UI using the templ templating engine.

2. **Service Layer** (`services/`):
   - Contains business logic and coordinates between handlers and data access.
   - Manages user authentication and boulder logging operations.

3. **Data Access Layer** (`db/`):
   - Implements data persistence using PostgreSQL.
   - Provides an interface for data operations, allowing for easy swapping of storage mechanisms.

4. **Model Layer** (`models/`):
   - Defines data structures used throughout the application.

5. **Utility Layer** (`utils/`):
   - Contains shared utilities like logging and JWT management.

## Technology Stack

- Go 1.22
- templ for HTML templating
- HTMX for dynamic content updates
- Tailwind CSS for styling
- PostgreSQL for data storage
- Docker for development and deployment

## Getting Started

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Set up environment variables (create a `.env` file based on `.env.example`)
4. Run the application in development mode with Docker: `docker compose up --build`

The compose file will run the application in development mode with hot reloading.
It Generates Tailwind CSS: `tailwindcss -i ./static/css/input.css -o ./static/css/tailwind.css --watch`
And runs the application in development mode with [air-verse/air](https://github.com/air-verse/air)
This also runs `templ generate` to generate the templ templates into Go code.


## License

This project is licensed under the MIT License - see the LICENSE file for details.
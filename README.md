# boulderlog
App for tracking my bouldering progress

## Goals

- User can easily and quickly record each route between routes
- User should see their progress visually
- Minimalist design, optimized for mobile use
- Possibly fun gamified elements

## Features

- Login with username and password
- Route logging
- Route history
- Graphs for grades topped
- Graphs for perceived difficulty of a grade

### Route logging

- Record each route, with date, location, topped (true/false), difficulty, and perceived difficulty
- Date should be automatically set to today
- Location should be set once per session
- Perceived difficulty should be set for each route
	- 1 - (Topped) flash very easily
	- 2 - (Topped) flash with some difficulty
	- 3 - (Topped) flash with a lot of difficulty
	- 4 - (Topped) topped with a lot of difficulty and multiple attempts
	- 5 - (Not topped) very close to topping, maybe next time
	- 6 - (Not topped) could do all moves separately but not together
	- 7 - (Not topped) some moves just could not be done
	- 8 - (Not topped) couldn't do any moves

## Tech stack

- Go (https://go.dev/)
- Templ (https://templ.guide/)
- htmx (https://htmx.org/)
- Tailwind (https://tailwindcss.com/)
- Postgres (https://www.postgresql.org/)
- Docker (https://www.docker.com/)

## Development

### Tailwind CSS

To generate the Tailwind CSS, run the following command:

```bash
tailwindcss -i ./static/css/input.css -o ./static/css/tailwind.css --watch
```

### Run the server

To run the server, run the following command:

```bash
air
```
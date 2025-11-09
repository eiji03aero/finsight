# Commands

When meets corresponding condition, follow the instructions attached.

## Create a scribble

- Condition: When you are asked to create a scribble.
- Instructions:
  - Get a timestamp in format of `YYYYMMDD-HHmm`
  - Create a file as a path of `docs/scribbles/${timestamp}-${general description user gave you}`
  - The content of the file must be empty

## Backend Development

- Condition: When working with the Go backend project.
- Instructions:
  - All Go commands (go get, go mod tidy, go run, etc.) must be executed inside the Docker container
  - Use `docker-compose exec api <command>` to run commands in the api container
  - Example: `docker-compose exec api go mod tidy`


## Frontend development

- Make sure to write import path in absolute path starting `@` alias.
- Make sure to run command inside `client` container when you need to run:
  - `npm ...` commands
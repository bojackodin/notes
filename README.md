# Auth
- Sign-up
- Sign-in

# Note
- Create note
- List notes

## make docker.image 

## docker compose -f deployment/docker-compose.yml up -d

## make migrate.up

curl -i -X POST \
-H "Authorization: Bearer INVALID" \
localhost:8080/notes

curl -X POST -i \
-H "Content-Type: application/json" \
-d '{"username":"your_username", "password": "your_password"}' \
localhost:8080/sign-up

curl -X POST -i \
-H "Content-Type: application/json" \
-d '{"username":"your_username", "password": "your_password"}' \
localhost:8080/sign-in

curl -i -H "Authorization: Bearer your_token" \
localhost:8080/notes

curl -i -X POST \
-H "Authorization: Bearer your_token" \
-H "Content-Type: application/json" \
-d '{"title":"This is a smple text with erors"}' \
localhost:8080/notes

curl -i -X POST \
-H "Authorization: Bearer your_token" \
-H "Content-Type: application/json" \
-d '{"title":"This is a simple text without errors"}' \
localhost:8080/notes

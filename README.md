# make docker.image 

# docker compose -f deployment/docker-compose.yml up -d

# docker run --rm -v $(pwd)/migrations:/migrations --network host migrate/migrate -path /migrations -database 'postgres://postgres:password@localhost:5432/?sslmode=disable' up

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

curl -i -X POST \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjY0NjkxOTcsImlhdCI6MTcyNjQ2NTU5NywidXNlcl9pZCI6NH0.BS5zfrVahEfxBmFEINEJYeYMoIZilfUnZvSjbyPS8kw" \
-H "Content-Type: application/json" \
-d '{"title":"This is a smple text with erors"}' \
localhost:8080/notes


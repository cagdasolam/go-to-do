run:
	go run cmd/server/main.go

docker-up:
	docker-compose up --build

k8s-deploy:
	kubectl apply -f k8s/

swagger:
	swag init -g cmd/server/main.go
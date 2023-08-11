run: 
	go run src/main.go

elastic:
	make docker-compose up -d
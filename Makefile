build:
	docker build -t go-file-server .
run:
	docker run --rm --name gfs -p 8080:8080 --volume ./data:/data go-file-server
dev:
	go run src/main.go
stop:
	docker stop gfs
clean:
	docker rmi go-file-server

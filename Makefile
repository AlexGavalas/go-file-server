build:
	docker build -t go-file-server .
compile:
	go build -o ./main ./src
run:
	docker run --rm --name gfs -p 8080:8080 --volume ./data:/data go-file-server
dev:
	ENV=development go run ./src
stop:
	docker stop gfs
clean:
	docker rmi go-file-server

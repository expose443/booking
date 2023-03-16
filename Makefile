build:
	docker build -t forum:1.0 .

run:
	docker run -d --name forum-app -p8080:8080 forum:1.0 && echo "server started at http://localhost:8080/"

docker:
	docker run  --rm -p 5432:5432 --name postgres -e POSTGRES_PASSWORD=verysecretpassowrd -e POSTGRES_USER=postgres -e POSTGRES_NAME=hotel postgres
stop:
	docker stop postgres
	docker container rm postgres
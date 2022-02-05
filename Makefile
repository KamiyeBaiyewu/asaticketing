default: start
project:=asa-tickets
version=${shell git rev-parse --short HEAD}



build:
	docker-compose -p ${project} build --build-arg APP_VERSION=${version}

up:
	docker-compose -p ${project} up server

start:
	docker-compose -p ${project} up -d 


logs: 
	docker-compose -p ${project} logs -f ${service}

stop: 
	docker-compose -p ${project} down


shell:
	docker-compose -p ${project} exec ${service} sh


start-server:
	docker-compose -p ${project} up -d server 

stop-server:
	docker-compose -p ${project} down server 

clean-server: build start-server

up-server: build up

# build-run: build up
clean: stop build start
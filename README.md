# Golang Training 2021

![image](https://user-images.githubusercontent.com/15091368/127484596-56c85468-9ba4-4ea1-a919-5b70815d72c5.png)

### Prerequisites

Create `secret` directory with `.token` file for external API.

### Make commands

`make build` builds docker images from `Dockerfile`

`make run` runs docker container

`make stop` stops running container

### Docker

`docker kill --signal=SIGINT stock-service` - send `SIGING` signal to running container. Observe graseful shutdown.

`docker ps` to get information about running containers

`docker inspect stock-service` to find information about container

`docker logs stock-service` to get container logs

#### Builder pattern for docker images

![docker](https://user-images.githubusercontent.com/25442973/128024871-ee1bcd2c-07f7-4bb1-87ef-5fe5f3f8e840.png)

### Docker swarm

Init swarm
```
docker swarm init
```

Deploy stack
```
docker stack deploy -c docker-stack.yaml xxx
```

Remove stack
```
docker stack rm xxx
```

Leave swarm
```
docker swarm leave -f
```

List of running services   
```
docker service ls
```

Update number of replicas
```
docker service update --replicas=3 xxx_stock
```

![replicas](https://user-images.githubusercontent.com/25442973/128883081-4771b174-4549-471e-a50c-717b01c384c2.png)

![multiple-containers](https://user-images.githubusercontent.com/25442973/128886978-33d93db4-7103-4806-98a6-e69ef496b1cf.png)


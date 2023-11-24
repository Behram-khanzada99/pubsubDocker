# pubsubDocker

# Go and Redis Application

This is a simple Go application that demonstrates publishing and subscribing to Redis using the go-redis library. The application generates and stores 100,000 objects in a Redis queue, and a subscriber consumes and processes these objects.
The application is then deployed into Docker container, which is then deployed onto the kubenetes environment.

## Prerequisites

Before deploying the application, ensure you have the following installed:
- [Redis] (https://redis.uptrace.dev/) installed on your machine
- [Docker] (https://www.docker.com/) installed on your machine.
- [Kubernetes] (https://www.kubernetes.com) installed on your machine.

## Getting Started

- The program has a producer that generates random objects.
- The struct object has three fields: ID, Name and Value
- The generateObject() generates objects. 
- Object ID is set to the number of iteration of the loop
- Object name is set to "Object%d" where %d is the id.
- Object Value is set to a random integer within range of 100.
- The objects are pushed to the redis in-memory list using LPush method.
- The subscriber uses multiple Goroutines to consume the objects produced by the publisher from the redis queue.

## Dockerfile

This Dockerfile starts from the official Golang image, sets the working directory, copies the current directory into the container, builds the Go application, exposes port 8080, and specifies the command to run the application.

## docker-compose.yml

This docker-compose.yml file defines two services: 'redis' and 'app'. The 'redis' service uses the official Redis image and exposes port 6379. The 'app' service builds the Go application using the Dockerfile in the current directory, exposes port 8080, and depends on the redis service. This ensures that the Redis service is started before the Go application.

## deployment.yaml

The Go application is deployed with one replicas, and it listens on port 8080.
Redis is deployed with a single replica, and it listens on port 6379.

## service.yaml

The Go application service is exposed on port 80 and forwards traffic to port 8080 on the pods labeled app: go-app.

## How to deploy Kubernetes service on Terraform
- Make sure you have Kubernetes and MiniKube Installed
- Start MiniKube by running `minikube start` in the terminal.
- Navigate to the terraform folder in the project directory and execute the following commands:
  ```
  terraform init
  terraform plan
  terraform apply
  ```
- Run the following commands to verify if the deployments and pods are made:
  - `kubectl get pods` to list all the pods and check their status.
  - `kubectl get deployments` to list all the deployments.

- Sometimes the status changes from "running" to something like "CrashLoopBackOff" after the application has executed successfully and exited. In order to check if the application has executed successfully run this command: `kubectl logs [pod-name]`


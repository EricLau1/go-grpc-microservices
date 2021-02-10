# Golang Microservices With GRPC And Kubernetes (Minikube)

## Requirements

- Docker
- Minikube
- Kubectl
- Golang
- Protoc

### Start Minikube

```bash
minikube start

eval $(minikube docker-env)
```

### Deploy MongoDB

```bash
cd k8s/mongodb

kubectl apply -f .
```

### Create One User on MongoDB

```bash
    mongo -u admin -p admin --authenticationDatabase admin

    use microservices

    db.createUser({user: 'user', pwd: 'password', roles:[{'role': 'readWrite', 'db': 'microservices'}]});

    show users;

    # testing authentication with new user
    mongo -u user -p password --authenticationDatabase microservices

    use microservices

    show collections
```

### Build Proto Files

```bash
# root directory
protoc -I=./messages --go_out=plugins=grpc:. ./messages/*.proto
```

### Build Services

```bash
# root directory
sh build.sh
```

### Build Image

```bash
cd k8s/docker

sh build.sh
```

### Deploy Services

```bash
cd k8s/services

kubectl apply -f .
```

### Generate EXTERNAL-IP

```bash
minikube tunnel
```

### Visualize EXTERNAL-IP

```bash
kubectl get svc api-service
```

### Visualize Logs

```bash
minikube dashboard
```

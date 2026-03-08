# legendary-succotash

legendary-succotash is a template REST API backend built with Go and PostgreSQL. It includes support for custom PostgreSQL images that you can create and preload with data. The project also comes with Kubernetes configuration files for easy deployment.

## Technologies and Libraries

1. [Go](https://go.dev/)
2. [PostgreSQL](https://www.postgresql.org/)
3. [Gofiber](https://gofiber.io/)
4. [GORM](https://gorm.io/)
5. [Docker](https://www.docker.com/)

## Building and Running the Application

You can deploy the application using either Docker or Kubernetes. Below are instructions for both methods.

### A. Create Custom PostgreSQL Docker Image

1. Modify the `database/init.sql` file to define your databases, tables, relations, and any initial data you want to preload.
2. During the build process, the Docker image will be created by executing the script inside the container.

#### Command to build the PostgreSQL Docker image:

```bash
docker build -f db.Dockerfile -t dp487/legendary-succotash:latest .
```

### Push the PostgreSQL Docker image to Docker Hub:

```bash
docker push dp487/legendary-succotashdb
```

### To run the Docker image locally, use the following command:

```bash
docker run -d --rm -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres --name legendary-succotashdb dp487/legendary-succotashdb:latest
```

### B. Docker Compose

You can run the application using Docker Compose with the following command:

```bash
docker-compose up --build
```

### C. Kubernetes

To deploy the application on a Kubernetes cluster, YAML files are provided for both the Go backend and PostgreSQL database.

Run the following commands to start the Kubernetes deployment and services:

```bash
kubectl apply -f deployment.yaml
kubectl apply -f postgres-deployment.yaml
```

### Future Tasks

1. Add tests
2. Add support for gRPC
3. Add support for GraphQL database
4. Add support for Auto build and push to dockerhub through Github actions

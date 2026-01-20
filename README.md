<div align="center"><h1>Containerized E-commerce API</h1></div>

## Introduction
This project is an ecommerce API developed using Golang and docker containers. Its purpose is to provide a hands-on experience in building a scalable and efficient ecommerce platform. By using Golang, we aim to leverage its performance and concurrency features. The project is designed to handle various aspects of an ecommerce platform, including product management, order processing, and customer management. With the use of containers, we ensure easy deployment and scalability of the application.

## Installation
1. **Requirements**: Install Docker (minimum 25.0.3), Docker Compose plugin, and Go 1.23+. These tools are used for local builds and testing.

2. **Clone the repository**: `git clone https://github.com/buliro/distributed_ecommerce.git`

3. **Navigate to the project directory**: `cd distributed_ecommerce`

4. **Configure environment variables**: Copy `services/public.env` (or the provided sample file) to `services/.env` and populate credentials such as Africa's Talking keys, database details, and API port.

5. **Start local stack with Docker Compose**:
   ```bash
   cd services
   docker compose up --build
   ```
   The compose file (and the `docker/Dockerfile` it references) will build the Go API and run the auxiliary services required for development.

## Usage
You can interact with the API through the command line in a bash terminal. See [Examples](#examples) for more.

## Features
1. **Dockerized Go Application**: The project is a Go application that is containerized using Docker, allowing for easy setup, consistent environments, and scalability.

2. **Kubernetes Integration**: The application is designed to be deployed on a Kubernetes cluster, providing robust orchestration capabilities such as automated rollouts, rollbacks, service discovery, and load balancing.

3. **ORY Hydra Integration**: The project integrates with ORY Hydra, an OAuth 2.0 and OpenID Connect provider, to handle authentication and authorization, ensuring secure access to your application.

4. **PostgreSQL Database**: The application uses a PostgreSQL database for data storage, providing a powerful, open-source object-relational database system with a strong reputation for reliability, data integrity, and correctness.

5. **Africa's Talking API Integration**: The project uses the Africa's Talking API to send SMS notifications to customers when they place an order, enhancing the user experience and providing real-time updates.

## Deployment Requirements
- **Container Registry**: Docker Hub, ECR, or another OCI-compliant registry with push access for the Jenkins pipeline.
- **Kubernetes Cluster**: Kubernetes ≥1.27 with `kubectl` configured locally. Create a namespace (default manifests target `ecommerce`).
- **Runtime Secrets**: Database credentials and any third-party API keys must be provisioned via Kubernetes secrets (see `k8s/secret.yaml`).
- **Image Pull Secret**: `registry-cred` Kubernetes secret configured with your registry credentials so the cluster can pull the published image.
- **Jenkins Agent Tooling**: Docker CLI, kubectl, and Go toolchain installed on the Jenkins build agent. Jenkins credentials named `dockerhub-credentials` and `kubeconfig` are expected by the pipeline.

## Deployment Steps
1. **Build and Push the Image**
   ```bash
   docker build -f docker/Dockerfile -t <registry>/<repo>/ecommerce-api:$(git rev-parse --short HEAD) services
   docker push <registry>/<repo>/ecommerce-api:$(git rev-parse --short HEAD)
   docker tag <registry>/<repo>/ecommerce-api:$(git rev-parse --short HEAD) <registry>/<repo>/ecommerce-api:latest
   docker push <registry>/<repo>/ecommerce-api:latest
   ```

2. **Prepare Kubernetes Resources**
   ```bash
   kubectl apply -f k8s/namespace.yaml
   kubectl create secret docker-registry registry-cred \
     --docker-server=<registry> \
     --docker-username=<user> \
     --docker-password=<token> \
     --namespace ecommerce
   kubectl apply -f k8s/config.yaml
   kubectl apply -f k8s/secret.yaml
   ```

3. **Deploy the API**
   ```bash
   kubectl apply -f k8s/deployment.yaml
   kubectl apply -f k8s/service.yaml
   kubectl -n ecommerce rollout status deploy/ecommerce-api
   ```

4. **Configure Jenkins CI/CD**
   - Create/verify credentials `dockerhub-credentials` (registry user/token) and `kubeconfig` (service-account kubeconfig scoped to `ecommerce`).
   - Create a Pipeline job pointing to this repository so it picks up the `Jenkinsfile` in the project root.
   - Configure the GitHub webhook to call `https://<jenkins-host>/github-webhook/` on push events.

5. **Monitor and Verify**
   ```bash
   kubectl -n ecommerce get pods
   kubectl -n ecommerce logs -f deploy/ecommerce-api
   kubectl -n ecommerce port-forward svc/ecommerce-api 8080:80
   curl http://localhost:8080/api/v1/status
   ```

Troubleshooting tips for common issues (image pull errors, failed probes, RBAC) are covered in the Jenkinsfile comments and Kubernetes manifests; ensure probes succeed before routing traffic.

## Post-Deployment Verification
- Check rollout history: `kubectl -n ecommerce rollout history deploy/ecommerce-api`
- Inspect pod details: `kubectl -n ecommerce describe pod <pod-name>`
- View application logs: `kubectl -n ecommerce logs -l app=ecommerce-api --tail=100`
- Smoke test API via port-forward (example for `/status` endpoint shown above).

## Troubleshooting
- **ImagePullBackOff**: Confirm `registry-cred` secret exists and Jenkins pushed the correct tag (inspect pod events via `kubectl describe`).
- **Failed Readiness/Liveness Probes**: Verify database connectivity and environment variables defined in `k8s/config.yaml` / `k8s/secret.yaml`. Adjust `initialDelaySeconds` if migrations take longer.
- **RBAC or Kubeconfig Issues**: Ensure the service account referenced by the Jenkins kubeconfig has permissions to `get`, `list`, `watch`, and `patch` deployments in the `ecommerce` namespace (`kubectl auth can-i patch deployment --as <service-account>`).
- **Database Connection Errors**: Confirm Postgres service is reachable (`kubectl -n ecommerce exec <pod> -- nc -zv postgres.ecommerce.svc.cluster.local 5432`).

## Examples
Check API status
```
curl -XGET 0.0.0.0:8080/api/v1/status
```

Get existing products
```
curl -XGET 0.0.0.0:8080/api/v1/products?page=1
```

Create a product
```
curl -X POST -H "Content-Type: application/json" -d '{"name": "Product 1", "price": 200, "stock": 5}' 0.0.0.0:8080/api/v1/products
```

Create a user - ***replace the phone number with your number to test SMS functionality***
```
curl -X POST -H "Content-Type: application/json" -d '{"name": "Customer 1", "phone": "+254700123456", "password": "secret"}' 0.0.0.0:8080/api/v1/customers
```

Get the currently logged in customer without authentication
```
curl -XGET 0.0.0.0:8080/api/v1/customers/me
```

Authenticate an existing customer
```
curl -X POST -H "Content-Type: application/json" -d '{"phone": "+254700123456", "password": "secret"}' 0.0.0.0:8080/api/v1/customers/login
```

## Contributing
1. **Fork the Repository**: Start by forking the project repository to your own GitHub account. This creates a copy of the repository under your account where you can make changes without affecting the original project.

2. **Clone the Forked Repository**: Clone the forked repository to your local machine. This allows you to work on the project locally.

3. **Create a New Branch**: Always create a new branch for each feature or bug fix you're working on. This keeps your changes organized and separated from the main project.

4. **Make Your Changes**: Make the changes you want to contribute. Be sure to follow the project's coding standards and conventions.

5. **Commit Your Changes**: Commit your changes to your branch. Write a clear and concise commit message describing what changes you made and why.

6. **Push Your Changes**: Push your changes to your forked repository on GitHub.

7. **Submit a Pull Request**: Go to the original project repository on GitHub and submit a pull request. In the pull request, describe the changes you made and why they should be included in the project.

For bug reports and feature requests, it's best to open an issue in the project's issue tracker. Describe the bug or feature in detail, including steps to reproduce (for bugs) or use cases (for features). Always check the issue tracker first to see if someone else has already reported the issue or requested the feature.

Remember, the key to a successful contribution is communication. Always be respectful and considerate of others, and remember that all contributions, no matter how small, are valued in an open source project.

## License
The project is distributed under the MIT license. Please see the [LICENSE](./LICENSE) file for more information.

## Authors
- [Leroy Buliro](http://github.com/leroysb)

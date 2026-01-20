<div align="center"><h1>Containerized E-commerce API</h1></div>

## Introduction
This project is an ecommerce API developed using Golang and docker containers. Its purpose is to provide a hands-on experience in building a scalable and efficient ecommerce platform. By using Golang, we aim to leverage its performance and concurrency features. The project is designed to handle various aspects of an ecommerce platform, including product management, order processing, and customer management. With the use of containers, we ensure easy deployment and scalability of the application.

## Installation
1. **Requirements**: Install Docker (minimum 25.0.3), Docker Compose plugin, and Go 1.23+. These tools are used for local builds and testing.

2. **Clone the repository**: `git clone https://github.com/buliro/distributed_ecommerce.git`

3. **Navigate to the project directory**: `cd distributed_ecommerce`

4. **Configure environment variables**: Copy `services/.env.example` to `services/.env` and populate credentials such as Africa's Talking keys, database details, Redis settings, and Hydra client secrets.

5. **Start local stack with Docker Compose**:
   ```bash
   docker compose -f docker/docker-compose.yml up --build --detach
   ```
   This builds the Go API from the `services` module and starts Hydra, PostgreSQL, Redis, and the API container. When you are done, run `docker compose -f docker/docker-compose.yml down` to stop everything.

6. **Run tests locally** (optional but recommended):
   ```bash
   cd services
   go test ./...
   ```

## Usage
You can interact with the API through the command line in a bash terminal. See [Examples](#examples) for more.

## Features
1. **Dockerized Go Application**: The project is a Go application that is containerized using Docker, allowing for easy setup, consistent environments, and scalability.

2. **Kubernetes Integration**: The application is designed to be deployed on a Kubernetes cluster, providing robust orchestration capabilities such as automated rollouts, rollbacks, service discovery, and load balancing.

3. **Hybrid OAuth + Redis Authorization**: The project integrates with ORY Hydra for token issuance and introspection, while Redis-backed sessions hydrate user context for the API to enforce fine-grained, customer-aware authorization.

4. **PostgreSQL Database**: The application uses a PostgreSQL database for data storage, providing a powerful, open-source object-relational database system with a strong reputation for reliability, data integrity, and correctness.

5. **Africa's Talking API Integration**: The project uses the Africa's Talking API to send SMS notifications to customers when they place an order, enhancing the user experience and providing real-time updates.

## Hybrid Authentication Overview

To keep the authorization rules simple for developers while retaining Hydra’s token lifecycle guarantees, the API combines Hydra with Redis-backed session storage:

```mermaid
flowchart LR
    Client[Client App] -->|1. Login request| API{Go API}
    API -->|2. Validate credentials
    (Postgres)| DB[(PostgreSQL)]
    API -->|3. Request access token| Hydra[[ORY Hydra]]
    Hydra -->|4. Issue token| API
    API -->|5. Store user+token| Redis[(Redis)]
    Client -->|6. Call private endpoint
    with Bearer token| API
    API -->|7. Introspect token| Hydra
    API -->|8. Hydrate user context| Redis
    API -->|9. Authorized response| Client
```

### Login & Request Flow

1. **User credentials validated** against PostgreSQL.
2. **Hydra issues an access token** via the client-credentials grant configured for the API.
3. **Redis session created** mapping token → customer profile (ID, phone, name) with a configurable TTL (`SESSION_TTL_SECONDS`).
4. **Auth middleware** introspects the incoming token with Hydra, then loads the user context from Redis before passing control to protected handlers.
5. **Logout or TTL expiry** removes the session entry from Redis, immediately revoking access for that token even if Hydra still considers it active.

### Operational Notes

- Redis connection details are configured through `REDIS_ADDR`, `REDIS_PASSWORD`, and `SESSION_TTL_SECONDS` in `services/.env`.
- Middleware exposes the hydrated customer on `c.Locals("user")`, and retains the token on `c.Locals("token")` for downstream handlers such as logout.
- Integration tests under `services/internal/tests` use `miniredis` to simulate the flow during `go test`.

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

Authenticate an existing customer (stores token + customer in Redis)
```
curl -X POST -H "Content-Type: application/json" -d '{"phone": "+254700123456", "password": "secret"}' 0.0.0.0:8080/api/v1/customers/login
```

Access a protected endpoint with the issued token
```
curl -H "Authorization: Bearer <access_token>" 0.0.0.0:8080/api/v1/customers/me
```

Revoke the session (removes Redis entry) and confirm access is lost
```
curl -X POST -H "Authorization: Bearer <access_token>" 0.0.0.0:8080/api/v1/customers/logout
curl -H "Authorization: Bearer <access_token>" 0.0.0.0:8080/api/v1/customers/me
# → 401 Session expired
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

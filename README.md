<div align="center"><h1>Containerized E-commerce API</h1></div>

## Introduction
This project is an ecommerce API developed using Golang and docker containers. Its purpose is to provide a hands-on experience in building a scalable and efficient ecommerce platform. By using Golang, we aim to leverage its performance and concurrency features. The project is designed to handle various aspects of an ecommerce platform, including product management, order processing, and customer management. With the use of containers, we ensure easy deployment and scalability of the application.

## Installation
1. **Requirements**: Before starting, ensure that you have Docker version 25.0.3 and Go version 1.22 installed on your machine. These are necessary to build and run the project.

2. **Clone the main repository**: Use the command `git clone http://github.com/leroysb/go_kubernetes.git` to clone the main project repository to your local machine.

3. **Navigate to the project directory**: Use the command `cd go_kubernetes` to change your current directory to the project's root directory.

4. **Clone the Hydra repository**: The project also requires the ORY Hydra OAuth2 and OpenID Connect server. Use the command `git clone https://github.com/ory/hydra.git` to clone this repository into the project's root directory.

5. **Update environment variables**: rename the file `public.env` to `.env`. Head over to [Africa's Talking](https://developers.africastalking.com) and log into your account. If you do not have an account, sign up then obtain the necessary credentials such as username, apiKey, and shortcode. Update the missing values in your `.env`.

6. **Build and run the project with Docker Compose**: Use the command `docker compose -f docker-compose.yml -f hydra/quickstart.yml up --build` to start the project. This command tells Docker Compose to build and run the Docker containers defined in the `docker-compose.yml`, `hydra/quickstart.yml`, and `hydra/quickstart-postgres.yml` files.

## Usage
You can interact with the API through the command line in a bash terminal. See [Examples](#examples) for more.

## Features
1. **Dockerized Go Application**: The project is a Go application that is containerized using Docker, allowing for easy setup, consistent environments, and scalability.

2. **Kubernetes Integration**: The application is designed to be deployed on a Kubernetes cluster, providing robust orchestration capabilities such as automated rollouts, rollbacks, service discovery, and load balancing.

3. **ORY Hydra Integration**: The project integrates with ORY Hydra, an OAuth 2.0 and OpenID Connect provider, to handle authentication and authorization, ensuring secure access to your application.

4. **PostgreSQL Database**: The application uses a PostgreSQL database for data storage, providing a powerful, open-source object-relational database system with a strong reputation for reliability, data integrity, and correctness.

5. **Africa's Talking API Integration**: The project uses the Africa's Talking API to send SMS notifications to customers when they place an order, enhancing the user experience and providing real-time updates.

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

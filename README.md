# PortsService

_This documentation explains how to use the PortsService, a service for managing Port information._

## Table of Contents

* [Building the Service](#building-the-service)
* [Running the Service](#running-the-service)
* [Storage Options](#storage-options)
* [Features](#features)
    + [Data Loading](#data-loading)
    + [Graceful Shutdown](#graceful-shutdown)
    + [HTTP Server](#http-server)
* [Environment Variables](#environment-variables)
    + [API_PORT](#apiport)
        - [Example](#example)
    + [STORAGE_PROVIDER](#storageprovider)
        - [Example](#example-1)
    + [REDIS_ADDRESS](#redisaddress)
        - [Example](#example-2)
* [Running in Docker](#running-in-docker)
* [HTTP API](#http-api)
    + [Status Codes](#status-codes)
    + [Example Request](#example-request)
* [Tests](#tests)
    + [Running Tests](#running-tests)
    + [Running Tests with Race Condition Check](#running-tests-with-race-condition-check)
* [Dummy Test Data](#dummy-test-data)
    + [Generating the File](#generating-the-file)
* [Generating Mocks](#generating-mocks)
    + [Installing gomock](#installing-gomock)
    + [Regenerating Mocks](#regenerating-mocks)
* [Code Quality and Linting](#code-quality-and-linting)
    + [Running golangci-lint with Docker](#running-golangci-lint-with-docker)
* [Third-Party Libraries](#third-party-libraries)
* [Running PortsService with Custom Configuration](#running-portsservice-with-custom-configuration)
* [Troubleshooting](#troubleshooting)
* [Future Improvements](#future-improvements)
  
## Building the Service

To build the PortsService, execute the following command:

```bash
go build -o portsService main.go
```

## Running the Service

To run the built binary with the ports.json file, execute:

```bash
./portsService ports.json
```

The service will launch an HTTP service on port 8080 by default. You can change the default port using the API_PORT
environment variable.

## Storage Options

By default, PortsService uses in-memory storage. To use Redis as the storage, set the STORAGE_PROVIDER environment
variable to "redis". When using Redis, provide the Redis address using the REDIS_ADDRESS environment variable.

## Features

### Data Loading

At startup, the service will load all port information from the provided file into the storage. The service uses an
upserting mechanism, meaning that if a port already exists in the storage, it will be updated; otherwise, it will be
inserted.

### Graceful Shutdown

The service is designed to shut down gracefully if it receives a termination signal (e.g., SIGINT or SIGTERM). This
ensures that any ongoing requests are completed before the service stops, preventing data corruption or loss.

### HTTP Server

The service exposes an HTTP server to retrieve port data. Clients can interact with the PortsService using
standard HTTP methods, such as GET, to fetch information about ports.
More info [here](#http-api)

## Environment Variables

This section provides an overview of the environment variables that can be modified to configure the PortsService.

### API_PORT

This variable allows you to change the default port on which the HTTP service listens. Default value: `8080`.

#### Example

```bash
export API_PORT=9000
```

### STORAGE_PROVIDER

This variable sets the storage provider for the service. Set its value to `redis` to use Redis as the storage provider,
or `memory` to use in-memory storage. If an invalid value is provided, the default value `memory` will be used. Default
value: `memory`.

#### Example

```bash
export STORAGE_PROVIDER=redis
```

### REDIS_ADDRESS

If using Redis as the storage provider, this variable must be set to the address of the Redis instance. There is no
default value.

#### Example

```bash 
export REDIS_ADDRESS=127.0.0.1:6379
```

To set these environment variables, use the export command before running the service, or include them in a .env file
when using Docker Compose.

Remember to adjust the values according to your requirements and infrastructure.

## Running in Docker

To execute the service inside a Docker container, and start a Redis instance inside a Docker container (using this Redis
as storage), run:

```bash
docker compose up
```

Or, to run the containers as a daemon:

```bash
docker compose up -d
```

## HTTP API

The HTTP API has a single endpoint: `GET /port/{portKey}`. The response has the following structure:

```json
{
  "name": "Ajman",
  "city": "Ajman",
  "country": "United Arab Emirates",
  "alias": [],
  "regions": [],
  "coordinates": [
    55.5136433,
    25.4052165
  ],
  "province": "Ajman",
  "timezone": "Asia/Dubai",
  "unlocs": [
    "AEAJM"
  ],
  "code": "52000"
}
```

### Status Codes

- 200: Port found in the storage.
- 204: Port not found in the storage.
- 400: Invalid request.
- 500: Server error.

### Example Request

To send an example request, use `curl`:

```bash
curl "http://localhost:8080/port/AEAJM" -s | jq
```

The response will be:

```json
{
  "name": "Ajman",
  "city": "Ajman",
  "country": "United Arab Emirates",
  "alias": [],
  "regions": [],
  "coordinates": [
    55.5136433,
    25.4052165
  ],
  "province": "Ajman",
  "timezone": "Asia/Dubai",
  "unlocs": [
    "AEAJM"
  ],
  "code": "52000"
}
```

## Tests

To ensure the reliability and correctness of the PortsService, it is essential to run tests. This section explains how
to execute tests for the PortsService.

### Running Tests

To run all tests, execute the following command:

```bash
go test -v ./...
```

This command will run all tests in the PortsService project and display the results, including passed and failed tests,
with verbose output.

### Running Tests with Race Condition Check

To run all tests with a race condition check, execute the following command:

```bash
go test -v -race ./...
```

This command will run all tests and check for race conditions in the PortsService project. The race detector is a
valuable tool for identifying data races in your code, which can lead to unexpected behavior and hard-to-debug issues.

By regularly running tests and checking for race conditions, you can ensure the quality and stability of the
PortsService.

## Dummy Test Data

If you need to test the PortsService with a larger dataset, you can generate a dummy test file called `data.json` with
approximately 200 MB of data. This section explains how to create the dummy test data.

### Generating the File

To generate the `data.json` file, execute the following command:

```bash
go run testdata/main.go
```

This command will create a `data.json` file in the current directory, containing approximately 200 MB of dummy port
data. You can use this file to test the PortsService's performance, storage capabilities, and handling of large
datasets.

Remember to replace the `ports.json` file with the newly generated `data.json` file when running the PortsService if you
want to use the larger dataset for testing.

## Generating Mocks

When developing and testing the PortsService, it's essential to use mocks to isolate components and simulate external
dependencies' behavior. To generate mocks for the PortsService, you need to have `gomock` installed. This section
explains how to install `gomock` and regenerate mocks.

### Installing gomock

To install `gomock`, run the following command:

```
go install github.com/golang/mock/mockgen@v1.6.0
```

This command installs the `gomock` package and the `mockgen` tool at the specified version (v1.6.0).

### Regenerating Mocks

To regenerate the mocks for the PortsService, run the following command:

```
go generate ./...
```

This command will traverse all directories within the project and regenerate mocks based on the source files and
interfaces defined in the project. The newly generated mocks can be used in your tests to isolate components and
simulate the behavior of external dependencies, making it easier to write effective and accurate test cases.

## Code Quality and Linting

To maintain high-quality code and ensure adherence to best practices, you can use the `golangci-lint` tool. This section
explains how to run `golangci-lint` using Docker.

### Running golangci-lint with Docker

To run `golangci-lint` with Docker, execute the following command:

```
docker run -t --rm -v $(pwd):/app -w /app golangci/golangci-lint golangci-lint run -v
```

This command will run `golangci-lint` in a Docker container, mounting the current directory (`$(pwd)`) to `/app` inside
the container, and setting the container's working directory to `/app`.

The `-v` flag enables verbose output, displaying the results of the linting process. If any issues are
found, `golangci-lint` will report them, allowing you to review and correct the problems to maintain high-quality code.

By regularly running `golangci-lint`, you can ensure that your code follows best practices and prevent potential issues
in the PortsService.

## Third-Party Libraries

The PortsService uses the following third-party libraries:

- `github.com/go-chi/chi`: A lightweight, idiomatic, and composable router for building Go HTTP services.
- `github.com/kelseyhightower/envconfig`: A Go library that processes environment variable-based configuration. It
  simplifies the use of environment variables for configuring applications.
- `github.com/redis/go-redis`: A Redis client for Go, providing a simple, general-purpose API for Redis commands.
- `go.uber.org/zap`: A fast, structured, leveled logging library for Go. It provides both low-level logging capabilities
  and high-level functionality, such as leveled logging and message fields.

## Running PortsService with Custom Configuration

To run the PortsService with a custom configuration, follow these steps:

Set the desired environment variables using the export command. For example:

```bash
export API_PORT=9000
export STORAGE_PROVIDER=redis
export REDIS_ADDRESS=127.0.0.1:6379
```

Alternatively, change the environment variables in docker-compose.yaml file when using Docker Compose.

Run the service:

```bash 
./portsService ports.json
```

By setting the environment variables, you can customize the service to fit your requirements and infrastructure.

## Troubleshooting

If you encounter issues while running the PortsService, consider the following troubleshooting steps:

Verify that the environment variables are set correctly.
Check the logs for any error messages or warnings.
Ensure that the service is running and listening on the specified port.
Confirm that the Redis instance (if used) is running and reachable.
If the issue persists, consult the documentation and source code for additional information.

## Future Improvements

While the current implementation of PortsService provides essential features and functionality, there is always room for
enhancement and optimization. This section outlines potential future improvements to the PortsService (in no particular
order).

1. **Error Wrapping**: Improve error handling by implementing error wrapping, which provides more context and detailed
   information about the cause of errors. This can help with debugging and identifying the root cause of issues more
   efficiently.

2. **Caching**: Implement caching mechanisms to improve performance when retrieving frequently accessed port data.

3. **Pagination**: Add support for pagination when returning large datasets. This will improve the performance and user
   experience when working with extensive collections of port data.

4. **API Versioning**: Introduce API versioning to allow for backward compatibility and smoother updates when new
   features are added or existing functionality is modified.

5. **Authentication and Authorization**: Add authentication and authorization mechanisms to restrict access to the API
   and protect sensitive data. This can be achieved using API keys, or other security protocols.

6. **Logging and Monitoring**: Enhance logging and monitoring capabilities to provide more detailed information about
   the service's performance, errors, and usage. Integrating with monitoring tools like Prometheus and Grafana can
   further improve observability.

7. **Rate Limiting**: Implement rate limiting to prevent abuse of the API and ensure fair resource allocation among
   users.

8. **Data Validation**: Enhance data validation on input and storage to ensure the consistency and correctness of the
   port data.

9. **Bulk Operations**: Add support for bulk operations, such as inserting or updating multiple ports at once, to
   improve efficiency and reduce the number of API calls needed for large-scale data management.

10. **Additional Endpoints**: Extend the API to include additional endpoints for more advanced querying and data
    manipulation, such as searching for ports by name, country, or other attributes.

11. **Documentation and Examples**: Continuously improve the documentation, providing more detailed explanations,
    examples, and tutorials to help users get the most out of the PortsService.

12. **Integration Tests**: Develop and implement integration tests to verify the correct interaction between different
    components of the PortsService, such as the API, storage providers, and external services. Integration tests can
    help identify potential issues and ensure the overall reliability and stability of the service.

By implementing these improvements, the PortsService can be further optimized, providing a more robust, efficient, and
user-friendly solution for managing port data.




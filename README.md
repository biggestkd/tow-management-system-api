# Car Management Tool API

## Project Overview

**Tow Management System API** is a Go API that enables connectivity between the UI application and downstream resources (e.g. databases, SaaS services, etc.).
This backend application will provide security between the frontend application and critical resources necessary to enable the tool's services.

Additional information can be found [here](https://app.clickup.com/9014465481/v/dc/8cmvmy9-1634)

### Getting Started

The API's codebase is organized according to a layered architecture to improve maintainability, scalability, and testability. Below is an overview of the main directories and their roles:

- `/handler`: Contains the
- `/model`: Contains the data structures and models used throughout the application.
- `/repository`: Contains the database implementation
- `/service`: Contains the business logic for the system.
- `/utils`: Contains utility functions that are shared across modules

This structure enables a clear separation of concerns and supports modular development across feature areas.

### Prerequisites

#### Install Go (version go1.23.1)
* Follow the instructions [here](https://go.dev/doc/install) or with [Homebrew](https://formulae.brew.sh/formula/go)

### Run Application

#### Local Instructions
1. Ensure that all dependencies are working (MongoDB)
2. Set environment variables.
```bash
MONGO_DB_USER="" # find in atlas console
MONGO_DB_PASSWORD="" # find in atlas console
MONGO_CLUSTER_HOSTNAME="" # find in atlas console
APP_NAME="tow-management-system-api"
ENVIRONMENT="local" # or dev/prod
PORT="8080"
```
3. Run command:
```bash
go run main.go
```

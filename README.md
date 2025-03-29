# cal-blog-service

This is a backend service I built in go to house all the data for my personal blog at https://calcorbin.com/blog.

## Table of Contents
* [cal-blog-service](#cal-blog-service)
  * [Table of Contents](#table-of-contents)
  * [Local Development](#local-development)
    * [Prerequisites](#prerequisites)
    * [Installation and Running](#installation-and-running)
    * [Linting](#linting)
    * [Testing](#testing)
    

## Local Development

### Prerequisites
- Go 1.24.0 or higher

### Installation and Running

1. Clone the repository
   ```bash
   git clone https://github.com/CalCorbin/cal-blog-service.git
   cd cal-blog-service
   ```
2. Install dependencies
   ```bash
   go mod download
   ```
3. Start the server
   ```bash
   make start
   ```

### Linting

1. Run linting locally
   ```bash
    make lint
    ```
2. Lint all changes before creating a pull request, otherwise changes will be rejected.

### Testing

1. Run tests locally
    ```bash
    make test
   ```
2. My goal with this repo is maintain both 100% test coverage, as well as useful tests. I want to exercise the code, and
feel confident in any changes. Pull requests will be automatically rejected if tests are not passing.
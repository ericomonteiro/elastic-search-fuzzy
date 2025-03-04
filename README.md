# Elasticsearch Go Example

This repository contains an example project demonstrating how to use Elasticsearch with Go. The project includes setting up an Elasticsearch client, creating an index, inserting documents, and querying the index through a simple HTTP server.

## Prerequisites

- Go 1.24 or later
- Elasticsearch 7.x or later

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/elasticsearch-go-example.git
   cd elasticsearch-go-example
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

## Usage

1. Start Elasticsearch on your local machine or ensure you have access to an Elasticsearch instance.

2. Run the application:

   ```sh
   go run main.go
   ```

3. The server will start on `http://localhost:8080`.

## Endpoints

### Insert a Document

- **URL:** `/brand`
- **Method:** `POST`
- **Payload:**

  ```json
  {
    "Name": "BrandName",
    "Terms": "Term1, Term2, Term3"
  }
  ```

- **Response:**

    - **201 Created:** Document inserted successfully
    - **400 Bad Request:** Invalid request payload

### Search Documents

- **URL:** `/search`
- **Method:** `GET`
- **Query Parameters:**
    - `terms`: The search terms

- **Response:**

    - **200 OK:** Returns the matching document
    - **400 Bad Request:** Missing 'terms' query parameter
    - **500 Internal Server Error:** Error during search

## Project Structure

- `main.go`: Main application file containing the server setup and Elasticsearch operations.
- `api.info.go`: Contains the `InfoRequest` struct and methods for configuring and executing the Info API request.
- `esapi.response.go`: Contains the `Response` struct and methods for handling API responses.

## License

This project is licensed under the Apache License 2.0. See the `LICENSE` file for more details.
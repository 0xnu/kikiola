## Kikiola

[![Release](https://img.shields.io/github/release/0xnu/kikiola.svg)](https://github.com/0xnu/kikiola/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xnu/kikiola)](https://goreportcard.com/report/github.com/0xnu/kikiola)
[![Go Reference](https://pkg.go.dev/badge/github.com/0xnu/kikiola.svg)](https://pkg.go.dev/github.com/0xnu/kikiola)
[![License](https://img.shields.io/github/license/0xnu/kikiola)](/LICENSE)

Kikiola is a high-performance vector database written in [Go](https://go.dev). It efficiently stores, indexes, and searches for vectors, making it suitable for similarity search, recommendation systems, artificial intelligence, and machine learning applications.

### Features

+ Fast and efficient vector storage and retrieval
+ Support for high-dimensional vectors
+ Indexing techniques for fast similarity search
+ Simple and intuitive API for easy integration
+ Scalable architecture for handling large datasets
+ Text embedding support for text-based queries

### Installation

To install Kikiola, ensure that you have Go installed on your system. Then, follow these steps:

1. Clone the Kikiola repository:

```sh
git clone https://github.com/0xnu/kikiola.git
```

2. Navigate to the project directory:

```sh
cd kikiola
```

3. Build the project:

```sh
go build ./...
```

4. Run the Kikiola server:

```sh
go run cmd/main.go
```

The Kikiola server will start running on `http://localhost:3400`.

### Usage - API Endpoints

Kikiola provides the following API endpoints:

+  `POST /vectors`: Insert a new vector
+  `GET /vectors/{id}`: Retrieve a vector by ID
+  `DELETE /vectors/{id}`: Delete a vector by ID
+  `POST /search`: Search for the nearest neighbors of a vector

#### cURL Examples

Here are some examples of how to use Kikiola with cURL:

1. Insert a new vector:

```sh
curl -X POST -H "Content-Type: application/json" -d '{
  "id": "vector1",
  "embedding": [0.1, 0.2, 0.3],
  "metadata": {
    "name": "Vector 1",
    "category": "sample"
  }
}' http://localhost:3400/vectors
```

2. Retrieve a vector by ID:

```sh
curl -X GET http://localhost:3400/vectors/vector1
```

3. Delete a vector by ID:

```sh
curl -X DELETE http://localhost:3400/vectors/vector1
```

4. Search for the nearest neighbors of a vector:

```sh
curl -X POST -H "Content-Type: application/json" -d '{
  "vector": {
    "id": "query_vector",
    "embedding": [0.5, 0.6, 0.7]
  },
  "k": 5
}' http://localhost:3400/search
```

### Integration with Other Applications or Systems

To use Kikiola in your Go applications or systems, follow these steps:

1. Install Kikiola as a dependency:

```sh
go get github.com/0xnu/kikiola
```

2. Import the necessary packages in your Go code:

```sh
import (
    "github.com/0xnu/kikiola/pkg/db"
    "github.com/0xnu/kikiola/pkg/server"
)
```

3. Create a new storage instance and server:

```sh
storage, err := db.NewStorage("data/vectors.db")
if err != nil {
    log.Fatal(err)
}
defer storage.Close()

server := server.NewServer(storage)
```

4. Start the Kikiola server:

```sh
log.Fatal(server.Start(":9090"))
```

Now, you can use Kikiola's functionality in your application or system by making HTTP requests to the appropriate endpoints.

### License

This project is licensed under the [MIT License](./LICENSE).

### Copyright

(c) 2024 [Finbarrs Oketunji](https://finbarrs.eu).
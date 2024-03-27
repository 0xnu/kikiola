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
+  `POST /search`: Search for the nearest neighbours of a vector

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

4. Search for the nearest neighbours of a vector:

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


### Kikiola Embedding Performance Benchmarks

You can execute `benchmark.go` to test the performance of Kikiola when embedding 1,000,000 records: `go run benchmark.go`

| Metric          | Value            |
|----------------|-------------------|
| Number of vectors | 1,000,000        |
| Vector dimension  | 128              |
| Embedding time    | 1m 41.964261347s |
| Embedding speed   | 9,807.36 vectors/sec |

The benchmark results provide insights into the performance of Kikiola when embedding vectors. Let's examine the results:

+ **Number of Vectors**: The benchmark tested the embedding of 1,000,000 vectors. It represents a significant dataset size and helps evaluate Kikiola's performance in handling large-scale vector embedding tasks.

+ **Vector Dimension**: Each vector in the benchmark has a dimension of 128. The vector dimension indicates the number of features or attributes associated with each vector. Using a dimension of 128 in many vector embedding scenarios is industry standard.

+ **Embedding Time**: It takes 1 minute and 41.964261347 seconds to embed all 1,000,000 vectors. The duration includes the time required to process and index each vector using Kikiola's embedding functionality.

+ **Embedding Speed**: Kikiola achieves an impressive embedding speed of 9,807.36 vectors per second. The metric shows the average number of vectors Kikiola embeds in one second. A higher embedding speed implies faster processing and improved efficiency in handling large datasets.

The benchmark results demonstrate that Kikiola is capable of efficiently embedding a large number of high-dimensional vectors. With an embedding speed of nearly 10,000 vectors per second, Kikiola showcases its performance capabilities in processing and indexing vectors at scale.

These results can be valuable for organisations and software engineers considering Kikiola for their vector embedding needs. The high embedding speed and the ability to handle a large dataset of 1,000,000 vectors with 128 dimensions highlight Kikiola's efficiency and scalability.

It's important to note that the actual performance may vary depending on factors such as hardware specifications, system configuration, and the nature of the vectors you are embedding. However, this benchmark strongly indicates Kikiola's performance characteristics and suitability for large-scale vector embedding tasks.

### License

This project is licensed under the [MIT License](./LICENSE).

### Copyright

(c) 2024 [Finbarrs Oketunji](https://finbarrs.eu).
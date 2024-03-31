## Usage - API Endpoints

Kikiola provides the following API endpoints:

+  `POST /vectors`: Insert a new vector
+  `GET /vectors/{id}`: Retrieve a vector by ID
+  `DELETE /vectors/{id}`: Delete a vector by ID
+  `GET /query/{id}`: Retrieve the original text content associated with an embedding ID
+  `POST /search`: Search for the nearest neighbours of a vector

#### cURL Examples

Here are some examples of how to use Kikiola with cURL:

1. Insert a new vector:

```sh
curl -X POST -H "Content-Type: application/json" -d '{
  "id": "badf35f6-e291-46cb-986b-01d57e6df80b",
  "embedding": [0.1, 0.2, 0.3],
  "metadata": {
    "name": "Vector 1",
    "category": "sample"
  }
}' http://localhost:3400/vectors
```

2. Retrieve a vector by ID:

```sh
curl -X GET http://localhost:3400/vectors/badf35f6-e291-46cb-986b-01d57e6df80b
```

3. Delete a vector by ID:

```sh
curl -X DELETE http://localhost:3400/vectors/badf35f6-e291-46cb-986b-01d57e6df80b
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

5. Tensor Compression:

```sh
curl -X POST -H "Content-Type: application/json" -d '{
  "id": "f5fbc794-c23a-4924-85af-775d4d34ecde",
  "embedding": [0.1, 0.2, 0.3],
  "metadata": {
    "name": "Tensor Compression",
    "category": "sample"
  },
  "compressed": true,
  "quantizationParams": {
    "min": -1.0,
    "max": 1.0,
    "bits": 8
  },
  "pruningMask": [false, true, false],
  "sparseIndices": [0, 2]
}' http://localhost:3400/vectors
```

6. Bioinformatics — Storing and Analysing Gene Sequences:

```sh
curl -X POST -H "Content-Type: application/json" -d '{
  "id": "badf35f6-e291-46cb-862d-01d57e6df80b",
  "embedding": [
    1, 0, 0, 0,
    0, 1, 0, 0,
    0, 0, 1, 0,
    0, 0, 0, 1,
    1, 0, 0, 0,
    0, 1, 0, 0,
    0, 0, 1, 0,
    0, 0, 0, 1
  ],  // Placeholder for the actual embedding
  "metadata": {
    "name": "Gene Sequence 1",
    "category": "gene",
    "sequence": "ACGTACGT" // Placeholder for the actual gene sequence
  }
}' http://localhost:3400/vectors
```

7. Retrieve the original text content associated with an embedding ID

```sh
curl -X GET "http://localhost:3400/query/83635f86-56b3-4bdd-a9bf-428dcebb8674"
```

### Integration with Other Applications or Systems

To use Kikiola in your Go applications or systems, follow these steps:

1. Install Kikiola as a dependency:

```sh
go get github.com/0xnu/kikiola
```

2. Import the necessary packages in your Go code:

```go
import (
    "github.com/0xnu/kikiola/pkg/db"
    "github.com/0xnu/kikiola/pkg/server"
)
```

3. Create a new storage instance and server:

```go
storage, err := db.NewStorage("data/vectors.db")
if err != nil {
    log.Fatal(err)
}
defer storage.Close()

server := server.NewServer(storage)
```

4. Start the Kikiola server:

```go
log.Fatal(server.Start(":9090"))
```

Now, you can use Kikiola's functionality in your application or system by making HTTP requests to the appropriate endpoints.

### Additional Integrations

Follow the tutorials below to integrate Kikiola and vector search into your current application.

[JavaScript](../tutorials/javascript.md) - [Python](../tutorials/python.md)
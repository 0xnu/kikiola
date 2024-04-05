## Usage - API Endpoints

Kikiola provides the following API endpoints:

+  `POST /vectors`: Insert a new vector
+  `GET /vectors/{id}`: Retrieve a vector by ID
+  `DELETE /vectors/{id}`: Delete a vector by ID
+  `PATCH vectors/{id}/metadata`: Update the metadata of a vector
+  `GET /query/{id}`: Retrieve the original text content associated with an embedding ID
+  `POST /search`: Search for the nearest neighbours of a vector
+  `POST /objects`: Insert a new object (e.g., document, image, audio, video, or any other file type)
+  `GET /objects/{id}`: Retrieve an object by ID
+  `DELETE /objects/{id}`: Delete an object by ID
+  `PATCH /objects/{id}/metadata`: Update the metadata of an object
+  `PATCH /objects/{id}/content`: Update the content of an object by uploading a new file

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

7. Retrieve the original text content associated with an embedding ID:

```sh
curl -X GET "http://localhost:3400/query/83635f86-56b3-4bdd-a9bf-428dcebb8674"
```

8. Update the metadata of a vector:

```sh
curl -X PATCH "http://localhost:3400/vectors/83635f86-56b3-4bdd-a9bf-428dcebb8674/metadata" -H "Content-Type: application/json" -d '{"metadata": {"name": "PDF Embeddings", "category": "pdf"}}'
```

9. Hybrid search (combining sparse embedding and traditional keyword search):

```sh
curl -X POST -H "Content-Type: application/json" -d '{
  "vector": {
    "ID": "query_vector",
    "Embedding": [0.1, 0.2, 0.3],
    "Metadata": {"text": "search keywords"}
  },
  "k": 10
}' http://localhost:3400/search
```

10. Reranking of search results:

```sh
curl -X POST -H "Content-Type: application/json" -d '{
  "vector": {
    "ID": "query_vector",
    "Embedding": [0.1, 0.2, 0.3]
  },
  "k": 10,
  "rerank": true
}' http://localhost:3400/search
```

11. Setting Alpha value (influence factor) for hybrid search:

```sh
curl -X POST -H "Content-Type: application/json" -d '{
  "vector": {
    "ID": "query_vector",
    "Embedding": [0.1, 0.2, 0.3],
    "Metadata": {"text": "search keywords"}
  },
  "k": 10,
  "alpha": 0.7
}' http://localhost:3400/search
```

12. Support for more distance/similarity metrics beyond just cosine similarity:

```sh
curl -X POST -H "Content-Type: application/json" -d '{
  "vector": {
    "ID": "query_vector",
    "Embedding": [0.1, 0.2, 0.3]
  },
  "k": 10,
  "metric": "euclidean"
}' http://localhost:3400/search
```

13. ANN (Approximate Nearest Neighbor) indexing for faster search on large datasets:

```sh
curl -X POST -H "Content-Type: application/json" -d '{
  "vector": {
    "ID": "query_vector",
    "Embedding": [0.02329646609723568, -0.044301047921180725, -0.014636795036494732]
  },
  "k": 10,
  "indexType": "annoy",
  "params": {
    "n_trees": 10
  }
}' http://localhost:3400/search
```

14. Ability to explain individual search results:

```sh
curl -X POST -H "Content-Type: application/json" -d '{
  "vector": {
    "ID": "query_vector",
    "Embedding": [0.1, 0.2, 0.3]
  },
  "k": 10,
  "explain": true
}' http://localhost:3400/search
```

15. Customizable relevance tuning and boosting of results:

```sh
curl -X POST -H "Content-Type: application/json" -d '{
  "vector": {
    "ID": "query_vector",
    "Embedding": [0.1, 0.2, 0.3]
  },
  "k": 10,
  "boost": {
    "category": "sample",
    "weight": 2.0
  }
}' http://localhost:3400/search
```

16. Insert a new object:

```sh
curl -X POST -H "Content-Type: multipart/form-data" -F "data={
  \"id\": \"0539f0ac-6771-47c6-8f5e-2cdf272a6de0\",
  \"metadata\": {
    \"name\": \"Oxford\",
    \"category\": \"Image\"
  }
};type=application/json" -F "object=@oxford.jpg" http://localhost:3400/objects
```

17. Retrieve an object by ID:

```sh
curl -X GET http://localhost:3400/objects/0539f0ac-6771-47c6-8f5e-2cdf272a6de0
```

18. Delete an object by ID:

```sh
curl -X DELETE http://localhost:3400/objects/0539f0ac-6771-47c6-8f5e-2cdf272a6de0
```

19. PATCH an object metadata by ID:

```sh
curl -X PATCH -H "Content-Type: application/json" -d '{
  "metadata": {
    "name": "Oxford High Street",
    "category": "Images"
  }
}' http://localhost:3400/objects/0539f0ac-6771-47c6-8f5e-2cdf272a6de0/metadata
```

20. PATCH an object content with a new one by ID:

```sh
curl -X PATCH -H "Content-Type: multipart/form-data" -F "object=@oxford_high_street.webp" http://localhost:3400/objects/0539f0ac-6771-47c6-8f5e-2cdf272a6de0/content
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
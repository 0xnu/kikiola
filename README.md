## Kikiola

[![Release](https://img.shields.io/github/release/0xnu/kikiola.svg)](https://github.com/0xnu/kikiola/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xnu/kikiola)](https://goreportcard.com/report/github.com/0xnu/kikiola)
[![Go Reference](https://pkg.go.dev/badge/github.com/0xnu/kikiola.svg)](https://pkg.go.dev/github.com/0xnu/kikiola)
[![License](https://img.shields.io/github/license/0xnu/kikiola)](/LICENSE)

Kikiola is a high-performance vector database written in [Go](https://go.dev). It efficiently stores, indexes, and searches for vectors, making it suitable for similarity search, recommendation systems, artificial intelligence, and machine learning applications.

### Features

+ Tensor Compression
+ Support for high-dimensional vectors
+ Handles concurrency and multiple writes
+ Text embedding support for text-based queries
+ Simple and intuitive API for easy integration
+ Indexing techniques for fast similarity search
+ Fast and efficient vector storage and retrieval
+ Scalable architecture for handling large datasets
+ Distributed Storage: multiple nodes or shards for scalability
+ Objects (e.g., document, image, audio, video, or any other file type)

### Run

To run Kikiola, ensure that you have Go installed on your system. Then, follow these steps:

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

### Test

To test Kikiola, ensure that you have Go installed on your system. Then, follow these steps:

```sh
go test ./...
```

### Usage, Use Cases, and Benchmark

+ [Usage](./docs/USAGE.md)
+ [Docker](./docs/DOCKER.md)
+ [Benchmark](./docs/BENCHMARK.md)
+ [Quay](./docs/QUAY.md)
+ [JFrog](./docs/JFROG.md)
+ [GitLab](./docs/GITLAB.md)
+ [Microsoft Azure](./docs/AZURE.md)
+ [Amazon Web Services (AWS)](./docs/AWS.md)
+ [Google Cloud Platform (GCP)](./docs/GCP.md)
+ [Generate and Store Embeddings](./tutorials/generate_and_store_embeddings_docs_images.md) - Documents and Images
+ [Generate and Store Embeddings](./tutorials/generate_and_store_embeddings_genome.md) - Genome Sequence
+ [Generate and Store Embeddings](./tutorials/generate_and_store_embeddings_genome_huggingface.ipynb) - Hugging Face 🤗
+ [Generate and Store Embeddings](./tutorials/generate_and_store_embeddings_10k.md) - SEC Form 10-K
+ [Generate and Store Embeddings](./tutorials/generate_and_store_embeddings_videos.md) - Videos
+ [Generate and Store Embeddings](./tutorials/generate_and_store_embeddings_audios.md) - Audios

### License

This project is licensed under the [MIT License](./LICENSE).

### Copyright

(c) 2024 [Finbarrs Oketunji](https://finbarrs.eu).
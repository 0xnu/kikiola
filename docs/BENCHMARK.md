## Kikiola Embedding Performance Benchmarks

You can execute `benchmark.go` to test the performance of Kikiola when embedding 1,000,000 records: `go run benchmark.go`

| Metric            | Value                |
|-------------------|----------------------|
|                   |                      |
27-03-2024
| Number of vectors | 1,000,000            |
| Vector dimension  | 128                  |
| Embedding time    | 1m 41.964261347s     |
| Embedding speed   | 9,807.36 vectors/sec |
28-03-2024
| Number of vectors | 1,000,000            |
| Vector dimension  | 128                  |
| Embedding time    | 1m 53.342241584s     |
| Embedding speed   | 8,822.84 vectors/sec |

The benchmark results (27-03-2024) provide insights into the performance of Kikiola when embedding vectors. Let's examine the results:

+ **Number of Vectors**: The benchmark tested the embedding of 1,000,000 vectors. It represents a significant dataset size and helps evaluate Kikiola's performance in handling large-scale vector embedding tasks.

+ **Vector Dimension**: Each vector in the benchmark has a dimension of 128. The vector dimension indicates the number of features or attributes associated with each vector. Using a dimension of 128 in many vector embedding scenarios is industry standard.

+ **Embedding Time**: It takes 1 minute and 41.964261347 seconds to embed all 1,000,000 vectors. The duration includes the time required to process and index each vector using Kikiola's embedding functionality.

+ **Embedding Speed**: Kikiola achieves an impressive embedding speed of 9,807.36 vectors per second. The metric shows the average number of vectors Kikiola embeds in one second. A higher embedding speed implies faster processing and improved efficiency in handling large datasets.

The benchmark results demonstrate that Kikiola is capable of efficiently embedding a large number of high-dimensional vectors. With an embedding speed of nearly 10,000 vectors per second, Kikiola showcases its performance capabilities in processing and indexing vectors at scale.

These results can be valuable for organisations and software engineers considering Kikiola for their vector embedding needs. The high embedding speed and the ability to handle a large dataset of 1,000,000 vectors with 128 dimensions highlight Kikiola's efficiency and scalability.

It's important to note that the actual performance may vary depending on factors such as hardware specifications, system configuration, and the nature of the vectors you are embedding. However, this benchmark strongly indicates Kikiola's performance characteristics and suitability for large-scale vector embedding tasks.
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
+ Tensor Compression

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

### Usage and Benchmark

+ [Usage](./USAGE.md)
+ [Benchmark](./BENCHMARK.md)
+ [Generate and Store Embeddings](./tutorials/generate_and_store_embeddings_docs_images.md) - Documents and Images:

    - <span title="Portable Document Format"><img src="https://img.icons8.com/color/48/000000/pdf.png" width="18" height="18"/></span> .pdf
    - <span title="Comma-Separated Values"><img src="https://img.icons8.com/color/48/000000/csv.png" width="18" height="18"/></span> .csv
    - <span title="Microsoft Excel Open XML Spreadsheet"><img src="https://img.icons8.com/color/48/000000/xls.png" width="18" height="18"/></span> .xlsx
    - <span title="Joint Photographic Experts Group"><img src="https://img.icons8.com/color/48/000000/jpg.png" width="18" height="18"/></span> .jpg
    - <span title="Portable Network Graphics"><img src="https://img.icons8.com/color/48/000000/png.png" width="18" height="18"/></span> .png
    - <span title="Web Picture Format"><img src="https://cdn-icons-png.flaticon.com/512/8263/8263085.png" width="18" height="18"/></span> .webp
    - <span title="Graphics Interchange Format"><img src="https://img.icons8.com/color/48/000000/gif.png" width="18" height="18"/></span> .gif

+ [Generate and Store Embeddings](./tutorials/generate_and_store_embeddings_videos.md) - Videos:

    - <span title="MPEG-4 Part 14"><img src="https://cdn-icons-png.flaticon.com/512/136/136545.png" width="18" height="18"/></span> .mp4
    - <span title="Audio Video Interleave"><img src="https://img.icons8.com/color/48/000000/avi.png" width="18" height="18"/></span> .avi
    - <span title="Matroska Multimedia Container"><img src="https://img.icons8.com/color/48/000000/mkv.png" width="18" height="18"/></span> .mkv
    - <span title="QuickTime Movie"><img src="https://img.icons8.com/color/48/000000/mov.png" width="18" height="18"/></span> .mov
    - <span title="Windows Media Video"><img src="https://cdn-icons-png.freepik.com/512/8300/8300652.png" width="18" height="18"/></span> .wmv
    - <span title="Flash Video"><img src="https://img.icons8.com/color/48/000000/flv.png" width="18" height="18"/></span> .flv
    - <span title="Web Media"><img src="https://cdn-icons-png.flaticon.com/512/8300/8300667.png" width="18" height="18"/></span> .webm
    - <span title="Moving Picture Experts Group"><img src="https://cdn-icons-png.flaticon.com/512/9645/9645792.png" width="18" height="18"/></span> .mpeg
    - <span title="Third Generation Partnership Project"><img src="https://cdn-icons-png.flaticon.com/512/8744/8744396.png" width="18" height="18"/></span> .3gp
    - <span title="High Efficiency Video Coding"><img src="https://img.icons8.com/fluency/48/000000/file.png" width="18" height="18"/></span> .hevc
    - <span title="Advanced Video Coding High Definition"><img src="https://img.icons8.com/fluency/48/000000/file.png" width="18" height="18"/></span> .avchd
    - <span title="Ogg Video"><img src="https://img.icons8.com/color/48/000000/ogg.png" width="18" height="18"/></span> .ogg
    - <span title="DVD Video Object"><img src="https://cdn-icons-png.flaticon.com/512/8300/8300651.png" width="18" height="18"/></span> .vob
    - <span title="MPEG-4 Part 14"><img src="https://cdn-icons-png.flaticon.com/512/9704/9704669.png" width="18" height="18"/></span> .m4v
    - <span title="MPEG Transport Stream"><img src="https://cdn-icons-png.flaticon.com/512/9405/9405038.png" width="18" height="18"/></span> .ts

### License

This project is licensed under the [MIT License](./LICENSE).

### Copyright

(c) 2024 [Finbarrs Oketunji](https://finbarrs.eu).
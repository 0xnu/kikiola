## Kikiola - Generate and Store Embeddings

The below example demonstrates how to generate embeddings from gene sequences using the OpenAI API. It retrieves gene sequence data from the [Ensembl](https://www.ensembl.org/index.html) REST API by specifying the gene ID. The code processes the API response to extract the relevant gene sequences from the nested data structure. It then uses the [text-embedding-3-small](https://platform.openai.com/docs/guides/embeddings) model to generate embeddings from the extracted gene sequences and stores them in Kikiola. You can use the stored embeddings for various use cases, such as [Retrieval Augmented Generation](https://blogs.nvidia.com/blog/what-is-retrieval-augmented-generation/) (RAG) or gene sequence analysis.

1. Install the required libraries if you haven't already:

```sh
pip install openai requests
```

2. Set OpenAI API Key as an environment variable: `export OPENAI_API_KEY="sk-**********************"`

3. Specify the Gene ID in the code below:

```python
import requests
import sys
import os
import json
from openai import OpenAI

class KikiolaGenomeEmbedding:
    def __init__(self):
        self.api_key = os.environ.get("OPENAI_API_KEY")
        if self.api_key is None:
            raise ValueError("OPENAI_API_KEY environment variable not set.")
        self.client = OpenAI(api_key=self.api_key)
        self.document_text = ""
        self.embeddings = []

    def load_genome_sequences(self, gene_id):
        server = "http://rest.ensembl.org"
        ext = f"/genetree/member/id/{gene_id}?"
        r = requests.get(server + ext, headers={"Content-Type": "application/json"})
        if not r.ok:
            r.raise_for_status()
            sys.exit()
        decoded = r.json()
        self.extract_genome_sequences(decoded)

    def extract_genome_sequences(self, data):
        if isinstance(data, dict):
            if 'sequence' in data and 'mol_seq' in data['sequence']:
                genome_sequence = data['sequence']['mol_seq']['seq']
                self.document_text += genome_sequence + '\n'
            else:
                for value in data.values():
                    self.extract_genome_sequences(value)
        elif isinstance(data, list):
            for item in data:
                self.extract_genome_sequences(item)

    def generate_embeddings(self):
        chunk_size = 8000
        chunks = [self.document_text[i:i + chunk_size] for i in range(0, len(self.document_text), chunk_size)]
        self.embeddings = []

        for chunk in chunks:
            response = self.client.embeddings.create(model="text-embedding-3-small", input=chunk)
            self.embeddings.append(response.data[0].embedding)
    def store_embeddings(self):
        for i, embedding in enumerate(self.embeddings):
            vector_data = {
                "id": f"genome_sequence_{i}",
                "embedding": embedding,
                "metadata": {
                    "name": f"Genome Sequence Embeddings - Part {i+1}",
                    "category": "genome_sequence"
                }
            }
            print(f"Vector data for Part {i+1}: {vector_data}")
            response = requests.post("http://localhost:3400/vectors", json=vector_data)
            if response.status_code == 200:
                print(f"Embeddings stored for Part {i+1}. Status code: {response.status_code}")
            else:
                print(f"Error storing embeddings for Part {i+1}. Status code: {response.status_code}")
                print(f"Error response: {response.text}")

if __name__ == "__main__":
    gene_id = "ENSG00000157764"
    embeddings_generator = KikiolaGenomeEmbedding()
    embeddings_generator.load_genome_sequences(gene_id)
    embeddings_generator.generate_embeddings()
    embeddings_generator.store_embeddings()
    print("Kikiola Embeddings Completed.")
```

> Before running the Python code, ensure the Kikiola server is running and accessible at the specified URL and port (e.g., `http://localhost:3400`).

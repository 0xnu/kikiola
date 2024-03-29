## Kikiola - Generate and Store Embeddings

Below is an example of generating embeddings from a PDF document using the OpenAI API. It extracts text from a PDF, generates embeddings using the [text-embedding-3-small](https://platform.openai.com/docs/guides/embeddings) model, and stores the embeddings in Kikiola. You can use the stored embeddings for various use cases, such as [Retrieval Augmented Generation](https://blogs.nvidia.com/blog/what-is-retrieval-augmented-generation/) (RAG).

1. Install the required libraries if you haven't already:

```sh
pip install openai PyPDF2 requests
```

2. Set OpenAI API Key as an environment variable: `export OPENAI_API_KEY="sk-**********************"`

3. Specify the location of the PDF locally or the URL of a PDF file in the code below.

```python
import io
import os
import requests
import PyPDF2
from openai import OpenAI

class KikiolaEmbedding:
    def __init__(self):
        self.api_key = os.environ.get("OPENAI_API_KEY")
        if self.api_key is None:
            raise ValueError("OPENAI_API_KEY environment variable is not set.")
        self.client = OpenAI(api_key=self.api_key)
        self.pdf_text = ""
        self.embeddings = []

    def load_pdf(self, pdf_path_or_url):
        if pdf_path_or_url.startswith("http"):
            # Download PDF from URL
            response = requests.get(pdf_path_or_url)
            pdf_file = io.BytesIO(response.content)
        else:
            # Open local PDF file
            pdf_file = open(pdf_path_or_url, 'rb')

        reader = PyPDF2.PdfReader(pdf_file)
        self.pdf_text = ""
        for page in reader.pages:
            self.pdf_text += page.extract_text()

        pdf_file.close()

    def generate_embeddings(self):
        response = self.client.embeddings.create(
            model="text-embedding-3-small",
            input=self.pdf_text
        )
        self.embeddings = response.data[0].embedding

    def store_embeddings(self):
        vector_data = {
            "id": "83635f86-56b3-4bdd-a9bf-428dcebb8674",
            "embedding": self.embeddings,
            "metadata": {
                "name": "PDF Embeddings",
                "category": "document"
            }
        }
        response = requests.post("http://localhost:3400/vectors", json=vector_data)
        print(f"Embeddings stored. Status code: {response.status_code}")

pdf_path_or_url = "company_manual.pdf"  # or a URL to a PDF file

embeddings_generator = KikiolaEmbedding()
embeddings_generator.load_pdf(pdf_path_or_url)
embeddings_generator.generate_embeddings()
embeddings_generator.store_embeddings()

print("Kikiola Embeddings Completed.")
```

> Before running the Python code, ensure the Kikiola server is running and accessible at the specified URL and port (e.g., `http://localhost:3400`).

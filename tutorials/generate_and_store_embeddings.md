## Kikiola - Generate and Store Embeddings

The below example demonstrates how to generate embeddings from a document or images using OpenAI API. It supports formats such as `.PDF`, `.CSV`, `.XLSX`, `.JPG`, `.PNG`, `.WEBP`, and `.GIF`. It takes the document or images, extracts text, and uses the [text-embedding-3-small](https://platform.openai.com/docs/guides/embeddings) model to generate embeddings before storing them in Kikiola. You can use the stored embeddings for various use cases, such as [Retrieval Augmented Generation](https://blogs.nvidia.com/blog/what-is-retrieval-augmented-generation/) (RAG).

1. Install the required libraries if you haven't already:

```sh
pip install openai PyPDF2 Pillow pandas pytesseract requests
```

2. Set OpenAI API Key as an environment variable: `export OPENAI_API_KEY="sk-**********************"`

3. Specify the location of the document or image locally or the URL in the code below:

```python
import io
import os
import requests
from PyPDF2 import PdfReader
import pandas as pd
from openai import OpenAI
from PIL import Image
import pytesseract
from csv import reader
from urllib.parse import urlparse

class KikiolaEmbedding:

    SUPPORTED_FILE_TYPES = (".pdf", ".csv", ".xlsx", ".jpg", ".jpeg", ".png", ".webp", ".gif")

    def __init__(self):
        self.api_key = os.environ.get("OPENAI_API_KEY")
        if self.api_key is None:
            raise ValueError("OPENAI_API_KEY environment variable not set.")
        self.client = OpenAI(api_key=self.api_key)
        self.document_text = ""
        self.embeddings = []

    def load_document(self, document_path_or_url):
        _urlparse = urlparse(document_path_or_url)
        if _urlparse.scheme in ["http", "https"]:
            response = requests.get(document_path_or_url)
            document_file = io.BytesIO(response.content)
        else:
            document_file = open(document_path_or_url, 'rb')

        if document_path_or_url.lower().endswith(".pdf"):
            self.load_pdf(document_file)
        elif document_path_or_url.lower().endswith(".csv"):
            self.load_csv(document_file)
        elif document_path_or_url.lower().endswith(".xlsx"):
            self.load_xlsx(document_file)
        elif document_path_or_url.lower().endswith(self.SUPPORTED_FILE_TYPES):
            self.load_image(document_file)
        else:
            raise ValueError("Unsupported file type.")
        if isinstance(document_file, io.IOBase):
            document_file.close()

    def load_pdf(self, pdf_file):
        pdf_reader = PdfReader(pdf_file)
        self.document_text = ""
        for page in pdf_reader.pages:
            self.document_text += page.extract_text()

    def load_csv(self, csv_file):
        csv_reader = reader(csv_file)
        self.document_text = " \n".join(', '.join(row) for row in csv_reader)

    def load_xlsx(self, xlsx_file):
        xlsx_data = pd.read_excel(xlsx_file)
        self.document_text = xlsx_data.to_csv(index=False)

    def load_image(self, image_file):
        self.document_text = pytesseract.image_to_string(Image.open(image_file))

    def generate_embeddings(self):
        response = self.client.embeddings.create(model="text-embedding-3-small", input=self.document_text)
        self.embeddings = response.data[0].embedding

    def store_embeddings(self):
        vector_data = {
            "id": "6da53ec8-5bc8-42ad-8069-9cd0b3815152",
            "embedding": self.embeddings,
            "metadata": {
                "name": "Document Embeddings",
                "category": "document"
            }
        }
        response = requests.post("http://localhost:3400/vectors", json=vector_data)
        print(f"Embeddings stored. Status code: {response.status_code}")

if __name__ == "__main__":
    document_path_or_url = "nano_bytes.gif"
    embeddings_generator = KikiolaEmbedding()
    embeddings_generator.load_document(document_path_or_url)
    embeddings_generator.generate_embeddings()
    embeddings_generator.store_embeddings()

    print("Kikiola Embeddings Completed.")
```

> Before running the Python code, ensure the Kikiola server is running and accessible at the specified URL and port (e.g., `http://localhost:3400`).

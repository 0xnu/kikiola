## Kikiola - Generate and Store Embeddings

The code demonstrates how to generate embeddings from a specific SEC filing, specifically a 10-K report, using the OpenAI API. It takes the URL of the SEC filing, extracts the text, and splits it into chunks. It then uses the `gpt-3.5-turbo` model to classify each chunk into predefined items like `Item 1. Business`, `Item 1A. Risk Factors`, etc.

Next, it generates embeddings for each item using the [text-embedding-3-large](https://platform.openai.com/docs/guides/embeddings) model before storing them in Kikiola.

You can use the stored embeddings for various use cases, such as similarity search, clustering, or [Retrieval Augmented Generation](https://blogs.nvidia.com/blog/what-is-retrieval-augmented-generation/) (RAG) tasks specific to financial documents.

1. Install the required libraries if you haven't already:

```sh
pip install openai pydantic requests uuid
```

2. Set OpenAI API Key as an environment variable: `export OPENAI_API_KEY="sk-**********************"`

3. Specify the [Form 10-K](https://www.investopedia.com/terms/1/10-k.asp) URL in the code below:

```python
import os
import time
import uuid
import requests
from pydantic import BaseModel
from typing import List, Optional
from openai import OpenAI

client = OpenAI()

class Item(BaseModel):
    name: Optional[str]

class ItemResponse(BaseModel):
    items: List[Item]

class Kikiola10KEmbedding:
    def __init__(self, url: str):
        self.api_key = os.environ.get("OPENAI_API_KEY")
        if self.api_key is None:
            raise ValueError("OPENAI_API_KEY environment variable not set.")
        self.url = url
        self.sec_10k = self.load_10k()
        self.prompt = """
        You are an expert at extracting text data from SEC filings like 10-Ks, 10-Qs, etc.
        Your sole job is to take a chunk of text and classify which Item in the SEC filing it belongs to.
        The goal is to identify and extract the text corresponding to the following Items:
        - Item 1. Business
        - Item 1A. Risk Factors
        - Item 2. Properties
        - Item 3. Legal Proceedings
        - Item 7. Management's Discussion and Analysis of Financial Condition and Results of Operations
        - Item 7A. Quantitative and Qualitative Disclosures About Market Risk
        - Item 8. Financial Statements and Supplementary Data

        Please adhere to the following guidelines:
        1. Be accurate, factual, and concise with your responses.
        2. Your response MUST be grounded in truth and the data that is present in the text. Do not make any assumptions or inferences beyond what is explicitly stated in the text.
        3. If the chunk of text does not clearly belong to any of the specified Items, return an empty list []. Do not attempt to guess or assign an Item if there is insufficient evidence in the text.
        4. If the chunk of text belongs to multiple Items, return a list of all applicable Items, like ["Item 1", "Item 1A"].
        5. Do not provide any explanations or additional information in your response. Only return a list of Items or an empty list.

        Your response should be formatted as a valid JSON list, like ["Item 1", "Item 1A"] or [].
        """
        self.items_map = {
            "Item 1": [],
            "Item 1A": [],
            "Item 2": [],
            "Item 3": [],
            "Item 7": [],
            "Item 7A": [],
            "Item 8": [],
            "unknown": [],
        }
        self.embeddings = []

    def load_10k(self) -> str:
        response = requests.get(self.url)
        return response.text

    def split_text_into_chunks(self, text: str, chunk_size: int) -> List[str]:
        chunks = []
        for i in range(0, len(text), chunk_size):
            chunks.append(text[i:i + chunk_size])
        return chunks

    def classify_items(self, chunk: str) -> ItemResponse:
        completion = client.chat.completions.create(
            model="gpt-3.5-turbo",
            messages=[
                {
                    "role": "system",
                    "content": self.prompt,
                },
                {
                    "role": "user",
                    "content": f"Extract the items from the following chunk: {chunk}",
                },
            ],
        )
        return completion.choices[0].message.content

    def extract_items(self) -> None:
        filing_text = self.sec_10k.replace("\n", " ")
        chunk_size = 4096
        chunks = self.split_text_into_chunks(filing_text, chunk_size)

        start_time = time.time()

        for index, chunk in enumerate(chunks):
            response = self.classify_items(chunk)
            items = eval(response)
            print(f"Chunk {index}. Items: {items}")

            for item in items:
                if item in self.items_map:
                    self.items_map[item].append(chunk)

            time.sleep(1)

        end_time = time.time()
        print(f"Time taken: {end_time - start_time} seconds")

    def generate_embeddings(self) -> None:
        self.embeddings = []
        document_uuid = str(uuid.uuid4())

        for item_name, chunks in self.items_map.items():
            document_text = "\n".join(chunks)
            response = client.embeddings.create(
                model="text-embedding-3-large",
                input=document_text
            )
            self.embeddings.append({
                "id": f"{document_uuid}_{item_name.replace(' ', '')}",
                "embedding": response.data[0].embedding,
                "metadata": {
                    "name": "sec_filing",
                    "category": "securities"
                }
            })

    def store_embeddings(self) -> None:
        server_url = "http://localhost:3400/vectors"
        for vector_data in self.embeddings:
            print(f"Vector data: {vector_data}")
            try:
                response = requests.post(server_url, json=vector_data)
                if response.status_code == 200:
                    print(f"Embeddings stored for {vector_data['id']}. Status code: {response.status_code}")
                else:
                    print(f"Error storing embeddings for {vector_data['id']}. Status code: {response.status_code}")
                    print(f"Error response: {response.text}")
            except requests.exceptions.RequestException as e:
                print(f"Error storing embeddings for {vector_data['id']}: {e}")

if __name__ == "__main__":
    url = "https://www.sec.gov/Archives/edgar/data/789019/000095017023035122/msft-20230630.htm"
    embeddings_generator = Kikiola10KEmbedding(url)
    embeddings_generator.extract_items()
    embeddings_generator.generate_embeddings()
    embeddings_generator.store_embeddings()
    print("Kikiola 10-K Embeddings Completed.")
```

> Before running the Python code, ensure the Kikiola server is running and accessible at the specified URL and port (e.g., `http://localhost:3400`).

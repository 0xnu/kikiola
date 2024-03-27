## Kikiola - Python Tutorial

Python applications can interact with Kikiola by making HTTP requests to the API endpoints exposed by the server. Here's how you can use Kikiola in a Python application:

1. Install the `requests` library if you haven't already:

```sh
pip install requests
```

2. Make HTTP requests to the Kikiola API endpoints using the `requests` library:

```python
import requests

# Insert a new vector
vector_data = {
    "id": "bd8a4438-3c82-4ce2-a78a-c619e27e4b36",
    "embedding": [0.1, 0.2, 0.3],
    "metadata": {
        "name": "Vector 1",
        "category": "sample"
    }
}
response = requests.post("http://localhost:3400/vectors", json=vector_data)
print(response.status_code)

# Retrieve a vector by ID
response = requests.get("http://localhost:3400/vectors/bd8a4438-3c82-4ce2-a78a-c619e27e4b36")
vector = response.json()
print(vector)

# Search for nearest neighbors
search_data = {
    "vector": {
        "id": "query_vector",
        "embedding": [0.5, 0.6, 0.7]
    },
    "k": 5
}
response = requests.post("http://localhost:3400/search", json=search_data)
results = response.json()
print(results)
```

In the above examples, the `requests.post()` and `requests.get()` functions send `POST` and `GET` requests to the Kikiola API endpoints.

> Before running the Python code, ensure the Kikiola server is running and accessible at the specified URL and port (e.g., `http://localhost:3400`).

You can adapt these examples to fit your specific use case and integrate Kikiola into your Python applications by making the appropriate HTTP requests to the Kikiola API endpoints.
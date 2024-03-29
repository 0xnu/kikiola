## Kikiola - Generate and Store Embeddings

The below example demonstrates how to generate embeddings from a video using OpenAI API. It supports formats such as `.mp4`, `.avi`, `.mkv`, `.mov`, `.wmv`, `.flv`, `.webm`, `.mpeg`, `.3gp`, `.hevc`, `.avchd`, `.ogg`, `.vob`, `.m4v`, and `.ts`. Upon specifying the location or URL of a video, it processes the video frames using the Python Imaging Library (PIL) to convert them into a format suitable for Pytesseract. Pytesseract extracts text from these frames before the [text-embedding-3-small](https://platform.openai.com/docs/guides/embeddings) model uses the extracted text to generate embeddings stored in Kikiola. You can use the stored embeddings for various use cases, such as [Retrieval Augmented Generation](https://blogs.nvidia.com/blog/what-is-retrieval-augmented-generation/) (RAG).

1. Install the required libraries if you haven't already:

```sh
pip install openai Pillow numpy opencv-python pytesseract requests
```

2. Set OpenAI API Key as an environment variable: `export OPENAI_API_KEY="sk-**********************"`

3. Specify the location of the document or image locally or the URL in the code below:

```python
import io
import os
import requests
import cv2
import numpy as np
from PIL import Image
import pytesseract
from openai import OpenAI
from urllib.parse import urlparse

class KikiolaVideoEmbedding:

    SUPPORTED_FILE_TYPES = (".mp4", ".avi", ".mkv", ".mov", ".wmv", ".flv", ".webm", ".mpeg", ".3gp", ".hevc", ".avchd", ".ogg", ".vob", ".m4v", ".ts")

    def __init__(self):
        self.api_key = os.environ.get("OPENAI_API_KEY")
        if self.api_key is None:
            raise ValueError("OPENAI_API_KEY environment variable not set.")
        self.client = OpenAI(api_key=self.api_key)

    def load_and_process_video(self, video_file):
        video = cv2.VideoCapture(video_file)
        frame_counter = 0

        while video.isOpened():
            ret, frame = video.read()
            if not ret:
                break

            frame = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
            pil_image = Image.fromarray(frame)
            document_text = pytesseract.image_to_string(pil_image)
            embeddings = self.generate_embeddings(document_text)

            self.store_embeddings(embeddings)

            frame_counter += 1

        video.release()

        print("Processed {} frames".format(frame_counter))

    def run(self, video_path_or_url):
        video_file = self.get_video_file(video_path_or_url)
        if video_path_or_url.lower().endswith(self.SUPPORTED_FILE_TYPES):
            self.load_and_process_video(video_file)
        else:
            raise ValueError("Unsupported file type.")
        
    def generate_embeddings(self, document_text):
        response = self.client.embeddings.create(model="text-embedding-3-small", input=document_text)
        return response.data[0].embedding

    def store_embeddings(self, embeddings):
        vector_data = {
            "id": "fc614fd4-59c5-4df9-bd1e-57f98d811976",
            "embedding": embeddings,
            "metadata": {
                "name": "Video Embeddings",
                "category": "video"
            }
        }
        response = requests.post("http://localhost:3400/vectors", json=vector_data)
        print(f"Embeddings stored. Status code: {response.status_code}")

    def get_video_file(self, video_path_or_url):
        _urlparse = urlparse(video_path_or_url)
        if _urlparse.scheme in ["http", "https"]:
            response = requests.get(video_path_or_url)
            video_file = io.BytesIO(response.content)
        else:
            video_file = video_path_or_url
        return video_file

if __name__ == "__main__":
    video_path_or_url = "vigro_deep_all_about_you.mp4"
    embeddings_generator = KikiolaVideoEmbedding()
    embeddings_generator.run(video_path_or_url)
    print("Kikiola Embeddings Completed.")
```

> Before running the Python code, ensure the Kikiola server is running and accessible at the specified URL and port (e.g., `http://localhost:3400`).

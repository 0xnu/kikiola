## Kikiola - Generate and Store Embeddings

The example below demonstrates how to generate embeddings from audio using OpenAI API. It supports formats such as `.wav`, `.mp3`, `.aac`, `.flac`, `.aiff`, `.ogg`, `.wma`, `.pcm`, `.alac`, and `.ac3`.

To get started, simply provide the local file path or URL of the audio you wish to process. The code uses the `pydub` library, which efficiently loads and converts the audio into a suitable format. It then proceeds to generate embeddings, which are stored in Kikiola, from the audio data using the OpenAI API's [text-embedding-3-small](https://platform.openai.com/docs/guides/embeddings) model.

The stored embeddings can be used for various use cases, such as [Retrieval Augmented Generation](https://blogs.nvidia.com/blog/what-is-retrieval-augmented-generation/) (RAG), speech synthesis, and music composition.

1. Install the required libraries if you haven't already:

```sh
pip install openai pydub requests
```

2. Set OpenAI API Key as an environment variable: `export OPENAI_API_KEY="sk-**********************"`

3. Specify the location of the document or image locally or the URL in the code below:

```python
import io
import os
import requests
from pydub import AudioSegment
from openai import OpenAI
from urllib.parse import urlparse

class KikiolaAudioEmbedding:

    SUPPORTED_FILE_TYPES = (".wav", ".mp3", ".aac", ".flac", ".aiff", ".ogg", ".wma", ".pcm", ".alac", ".ac3")

    def __init__(self):
        self.api_key = os.environ.get("OPENAI_API_KEY")
        if self.api_key is None:
            raise ValueError("OPENAI_API_KEY environment variable not set.")
        self.client = OpenAI(api_key=self.api_key)

    def load_and_process_audio(self, audio_file):
        audio = AudioSegment.from_file(audio_file)

        audio_data = io.BytesIO()
        audio.export(audio_data, format="wav")
        audio_data.seek(0)

        embeddings = self.generate_embeddings(audio_data)
        self.store_embeddings(embeddings)

    def run(self, audio_path_or_url):
        audio_file = self.get_audio_file(audio_path_or_url)
        if audio_path_or_url.lower().endswith(self.SUPPORTED_FILE_TYPES):
            self.load_and_process_audio(audio_file)
        else:
            raise ValueError("Unsupported file type.")

    def generate_embeddings(self, audio_data):
        response = self.client.embeddings.create(
            model="text-embedding-3-small",
            input=audio_data.read()
        )
        embeddings = response.data[0].embedding
        return embeddings

    def store_embeddings(self, embeddings):
        vector_data = {
            "id": "fc614fd4-59c5-4df9-bd1e-57f98d811976",
            "embedding": embeddings,
            "metadata": {
                "name": "Audio Embeddings",
                "category": "audio"
            }
        }
        response = requests.post("http://localhost:3400/vectors", json=vector_data)
        print(f"Embeddings stored. Status code: {response.status_code}")

    def get_audio_file(self, audio_path_or_url):
        _urlparse = urlparse(audio_path_or_url)
        if _urlparse.scheme in ["http", "https"]:
            response = requests.get(audio_path_or_url)
            audio_file = io.BytesIO(response.content)
        else:
            audio_file = audio_path_or_url
        return audio_file

if __name__ == "__main__":
    audio_path_or_url = "wande_coal_olamide_kpe_paso.mp3"
    embeddings_generator = KikiolaAudioEmbedding()
    embeddings_generator.run(audio_path_or_url)
    print("Kikiola Embeddings Completed.")
```

> Before running the Python code, ensure the Kikiola server is running and accessible at the specified URL and port (e.g., `http://localhost:3400`).

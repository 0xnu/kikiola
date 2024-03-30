## Using Kikiola Docker Image

Below is how you can utilise [Kikiola Docker Image](https://hub.docker.com/r/0xnu20/kikiola) in your project.

### Prerequisites

- Docker installed on your system

### Usage

1. Pull the Docker image: `docker pull 0xnu20/kikiola:latest`.


2. Run a container from the image: `docker run -p 3400:3400 0xnu20/kikiola:latest`.

This command starts a container from the `0xnu20/kikiola:latest` image and maps port 3400 from the container to port 3400 on the host machine.

3. Access the application:

The Golang Endpoints inside the container will be accessible at `http://localhost:3400`.

### Additional Information

- The `0xnu20/kikiola:latest` image is built from a Golang package and exposes port 3400.
- The image is configured to start the Golang application when the container runs automatically.
- Refer to the Docker documentation for more advanced usage options to customise the container's behaviour or pass environment variables.
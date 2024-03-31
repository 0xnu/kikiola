## Kikiola Docker Image x GitLab Container Registry

To push the [Kikiola Docker image](https://hub.docker.com/r/0xnu20/kikiola/) to [GitLab Container Registry](https://docs.gitlab.com/ee/user/packages/container_registry/), follow these steps:

1. Authenticate with the GitLab Container Registry:

   + Log in to your GitLab account.
   + Navigate to your project's repository.
   + Go to `Settings` > `Repository` and expand the `Container Registry` section.
   + Copy the authentication command provided, which should look like:

     ```sh
     docker login registry.gitlab.com -u <your-username> -p <your-access-token>
     ```

   + Run the authentication command in your terminal to log in to the GitLab Container Registry.

2. Tag your Docker image with the GitLab Container Registry URL:

   ```sh
   docker tag 0xnu20/kikiola:latest registry.gitlab.com/your-namespace/your-project-name/your-image-name:latest
   ```

3. Push the tagged image to the GitLab Container Registry:

   ```sh
   docker push registry.gitlab.com/your-namespace/your-project-name/your-image-name:latest
   ```

4. Configure deployment settings (optional):

   + If you want to use the Docker image in your GitLab CI/CD pipeline, you can configure the deployment settings in your `.gitlab-ci.yml` file.
   + Specify the image name and tag in the appropriate job or stage.
   + GitLab will automatically pull the image from the Container Registry during deployment.

5. Pull and run the Docker image from the GitLab Container Registry:

   + To pull the Docker image from the GitLab Container Registry, use the following command:

     ```sh
     docker pull registry.gitlab.com/your-namespace/your-project-name/your-image-name:latest
     ```

   + To run the Docker container, use the `docker run` command:

     ```sh
     docker run -d -p 3400:3400 registry.gitlab.com/your-namespace/your-project-name/your-image-name:latest
     ```

Please replace `your-namespace`, `your-project-name`, and `your-image-name` with your actual GitLab namespace, project name, and desired image name, respectively.

Following these steps, you can push your Kikiola Docker image to the GitLab Container Registry and integrate it seamlessly with your GitLab projects and CI/CD pipelines.
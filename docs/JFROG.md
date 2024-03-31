## Kikiola Docker Image x JFrog Container Registry

To push the [Kikiola Docker image](https://hub.docker.com/r/0xnu20/kikiola/) to [JFrog Container Registry](https://jfrog.com/container-registry/), follow these steps:

1. Create a new repository in JFrog Container Registry:

   + Log in to your JFrog Platform account.
   + Navigate to the `Artifactory` section.
   + Click on `Repositories` in the left sidebar.
   + Click on `New Repository` and select `Remote Repository`.
   + Choose `Docker` as the repository type and provide a name for your repository.
   + Configure the necessary settings, such as the remote URL (e.g., `https://registry-1.docker.io`), and save the repository.

2. Authenticate with JFrog Container Registry:

   + In the JFrog Platform, go to `Administration` > `Identity and Access` > `Access Tokens`.
   + Click `Generate Admin Token` and copy the generated access token.
   + Run the following command in your terminal to log in to JFrog Container Registry:

     ```sh
     docker login your-jfrog-instance.jfrog.io -u your-username -p your-access-token
     ```

3. Tag your Docker image with the JFrog Container Registry URL:

   ```sh
   docker tag 0xnu20/kikiola:latest your-jfrog-instance.jfrog.io/your-repo-name/your-image-name:latest
   ```

4. Push the tagged image to JFrog Container Registry:

   ```sh
   docker push your-jfrog-instance.jfrog.io/your-repo-name/your-image-name:latest
   ```

5. Configure build integration (optional):

   + You can create a new pipeline in the JFrog Platform if you want to automate the build and push process using JFrog Pipelines.
   + Define the steps in your pipeline YAML file to build the Docker image and push it to the JFrog Container Registry.
   + Configure the pipeline to trigger specific events, such as code commits or schedules.

6. Pull and run the Docker image from JFrog Container Registry:

   + To pull the Docker image from JFrog Container Registry, use the following command:

     ```sh
     docker pull your-jfrog-instance.jfrog.io/your-repo-name/your-image-name:latest
     ```

   + To run the Docker container, use the `docker run` command:

     ```sh
     docker run -d -p 3400:3400 your-jfrog-instance.jfrog.io/your-repo-name/your-image-name:latest
     ```

Please replace `your-jfrog-instance`, `your-username`, `your-access-token`, `your-repo-name`, and `your-image-name` with your actual JFrog instance URL, username, access token, repository name, and desired image name, respectively.

Following these steps, you can push your Kikiola Docker image to JFrog Container Registry and leverage JFrog's features for artifact management, security scanning, and pipeline automation.
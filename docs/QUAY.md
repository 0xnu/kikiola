## Kikiola Docker Image x Quay

To push the [Kikiola Docker image](https://hub.docker.com/r/0xnu20/kikiola/) to [Quay](https://quay.io/), follow these steps:

1. Create a new repository on Quay:

   + Log in to your Quay account.
   + Click on the `+` icon in the top right corner and select `New Repository`.
   + Provide a name for your repository and configure the necessary settings.
   + Click `Create Public Repository` or `Create Private Repository` based on your preference.

2. Tag your Docker image with the Quay repository URL:

   ```sh
   docker tag 0xnu20/kikiola:latest quay.io/your-username/your-repository-name:latest
   ```

3. Push the tagged image to Quay:

   ```sh
   docker push quay.io/your-username/your-repository-name:latest
   ```

4. Configure webhook notifications (optional):

   + You can set up webhook notifications if you want to receive notifications when you push new images to your Quay repository.
   + In your Quay repository settings, navigate to the `Webhooks` section.
   + Click on `Create Webhook` and specify the required information, including the webhook URL and the events for which you wish to receive notifications.

5. Pull and run the Docker image from Quay:

   + To pull the Docker image from Quay, use the following command:

     ```sh
     docker pull quay.io/your-username/your-repository-name:latest
     ```

   + To run the Docker container, use the `docker run` command:

     ```sh
     docker run -d -p 80:80 quay.io/your-username/your-repository-name:latest
     ```

Please replace `your-username` and `your-repository-name` with your Quay username and the desired repository name.

Following these steps, you can push your Kikiola Docker image to Quay and manage and distribute the Kikiola Docker image using Quay's features and integrations.
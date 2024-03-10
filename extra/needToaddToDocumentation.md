# section on finding env names inside pods to expose to containers

finding env vars that run in your container on pods:

kubectl exec -it tpm-frontend-8dfb5484b-z5v22 -- /bin/sh

then type in: env

from here you an find tpm backend service and tpm backend service port
and call the vars in using os.getenv thing in javascript, then can do the following:

http://$TPM_BACKEND_SERVICE:$TPM_BACKEND_SERVICE_PORT/api/v1/createUser

but need to call this on server side, not client side, as it is a server call not a client call.

# Configure Docker to Push to and Pull from the Registry

% docker push registry.digitalocean.com/tpm-containers-test/tpm-frontend:0.0.01

To interact with your registry using the docker command-line interface (CLI), you need to configure docker using the DigitalOcean command-line tool, doctl. Install doctl and authenticate it with an API token.

Then, call the registry login command:

`doctl registry login`

This command adds credentials to docker so that pull and push commands to your DigitalOcean registry will be authenticated.

You can then use the docker tag command to tag your image with the fully qualified destination path, and docker push to upload it:

`docker tag <my-image> registry.digitalocean.com/<my-registry>/<my-image>`
`docker push registry.digitalocean.com/<my-registry>/<my-image>`

## How to Free Up Space in Your Container Registry

Over time after pushing many updates to an image, it can result in unused layers and blobs taking up space. Running garbage collection will remove these unused layers and blobs.

Validated on 3 Dec 2021 • Last edited on 27 Sep 2023
The DigitalOcean Container Registry (DOCR) is a private Docker image registry with additional tooling support that enables integration with your Docker environment and DigitalOcean Kubernetes clusters. DOCR registries are private and co-located in the datacenters where DigitalOcean Kubernetes clusters are operated for secure, stable, and performant rollout of images to your clusters.

Manage Images and Tags
To manually delete images and tags in your registry:

Navigate to Container Registry in the control panel.

If you have a live registry and you have pushed images to it, the images are listed here.

Click the plus, +, next to a repository to see its image versions and each version’s tags. You can also see untagged images here, which we recommend deleting manually or through garbage collection.

# section on switching between clusters in kubectl using config

`kubectl config get-contexts`
`kubectl config use-context <context name>`

# digital ocean uses amd64 architecture

this means that your containers need to be built for the same architecture, so if on mac, build it with that architecture specifically as it cannot use arm architecture to run the container on the digitial ocean clusters.
push the versions of your images with new tags to container registry but build using:
docker build --platform linux/amd64

tpm-backend:0.0.02
tpm-frontend:0.0.02

# section on exposing frontend to public

check frontend-service yaml, no ingress needed for external IP as Loadbalancer provisioning will do this
through the following command:
apply frontend service
then use external Ip and port exposed by service to see frontend: 159.65.212.242:3000

![alt text](image-6.png)

# move get logs command to debugging section

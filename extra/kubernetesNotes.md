kubernetes uses yaml files to create infrastructure

--OVERVIEW--

use KIND for local testing, this is Kubernetes In Docker

creating a single cluster with kind (starts cluster in docker container):
`kind create cluster --image kindest/node:v1.23.5`

see the running nodes :
`kubectl get nodes`

--NAMESPACES--

conceptually compartmentalise your resources into virtual clusters, that way you can allocate a set amount of resource to each namespace. can set namespace for employee type, department, products etc.

create namespace:
`kubectl create namespace <namespace name>`

see all namespaces:
`kubectl get ns`

once our namespace has been created we can go ahead and add all our resources, including things such as config maps and secrets

--CONFIGMAPS--

ConfigMaps is an object used to store non-confidential data in key value pairs, pods can consume configMaps as env variables, command line args or as config files in a volume.
The advantage of configMaps is that it allows you to decouple environmnet specific configurations from your applications, making the app easily portable.

creating configmaps that pass key value pairs:
`kubectl -n <namespace name> create configmap <configmap name> --from-literal ENV_VAR_NAME=value`
e.g:
![alt text](image-1.png)

list out configmaps:
`kubectl -n <namespace name> get cm`

output configmap in yaml format:
`kubectl -n <namespace name> get cm <configmap name> -o yaml`

--SECRETS--
work similar to configmaps but are encrypted at rest.
It's important to keep secrets and config maps separate as we can use RBAC to give access to secrets to some people but others can only see configMaps

creating secrets:
`kubectl -n <namespace name> create secret generic <secret name> --from-literal SECRET_KEY=secretValue`
![alt text](image-2.png)

in the picture below we can see how the secret value is set as an environment variable in the deployment yaml:
![alt text](image-3.png)

list all secrets:
`kubectl -n <namespace name> get secret`

--DEPLOYMENTS--

Use kubernetes deployment to define how we want kubernetes to run the containers, we do this through the spec in the yaml file.

in the deployments yaml file you can state pods and number of replicas (more pods to allow failover)

in yaml file you set kind: Deployment for deployment files,
can set instance number under replicas in spec area

in the spec you also define the containers you want to run with:
name, image, ports, you can also add env variables and pass in secret values that you have stored already:
![alt text](image.png)

the env variables can also be mapped to configmaps

deploy application to k8s:
`kubectl -n <namespace name> apply -f <deploy.yaml file>`

once the deployment happens you can view the pods:
`kubectl -n <namespace name> get pods`

## Services

https://kubernetes.io/docs/concepts/services-networking/service/

Once we expose a port through our deployment yaml, we need to send traffic to that port, to do this we use services.

A service is a construct that defines how we want traffic to flow to our pods. It is an abstract way to expose an application running on a set of pods as a network service. You can define a service with service.yaml file:
![alt text](image-4.png)

In the image above, the targetPort needs to match the port of the selector. The selector also needs to be the same name as the label for the app that you have set in the deployment.yaml file.

difference between port and target port:
the target port is what the service is exposing and listening for external traffic. Once it receives traffic, it then forwards it to the port, which is the port that the pods inside the cluster are listening on.

port is the port on which the Service itself listens for traffic, while targetPort is the port to which traffic is forwarded to the backend Pods. They are separate because the Service may expose a different port to clients than what the Pods are listening on internally. This abstraction allows for flexibility in managing networking and port assignments within Kubernetes deployments.

Services can provide a public or private endpoint.

to apply the service yaml file:
`kubectl -n <namespace name> -f <path to service.yaml>`

to list the service:
`kubectl -n <namespace name> get svc`

running list service will give us an internal IP address, we can access our deployment on, as well as show the port thats exposed. This internal IP address allows other applications deployed in your cluster to access the pod running on the service.
![alt text](image-5.png)

Note that this internal cluster IP address is not accessible from outside the cluster by default, it allows different parts of the application within the cluster to communicate to each other. If you need to make your service accessible from outside the cluster, you can use things such as an ingress controller or LoadBalancer.

If you have two different deployments running different images within the same namespace and each deployment exposes a Service, the Pods managed by these deployments can communicate with each other using their internal Cluster IP addresses and exposed ports.

Here's how it works:

1.You have two Deployments, each deploying Pods running different container images.

2.Each Deployment exposes a Service, which creates a stable endpoint for accessing the Pods.

3.The Pods within the same namespace can communicate with each other using the DNS name of the Service, which resolves to the internal Cluster IP address of the Service. They can use the internal Cluster IP address and exposed ports to send and receive data.

--EXTRA NOTES--

command to get cluster names:
`kind get clusters`

command to load docker images into kind:
`kind load docker-image <Image name> -n <cluster name>`

`kubectl describe deployment tpm-backend  -n tpm`

`kubectl describe pod/tpm-backend-86c48f8b46-6t4m7 -n tpm`

`docker image tag tpm-backend:latest tpm-backend:0.0.1`

https://documentation.breadnet.co.uk/kubernetes/kb/kubectl-commands/

because pods that are in the same namespace share the same network, they both work off the same localhost, therefore images deployed in the same namespace can talk to each other i.e. cloud sql proxy and your db connection vals:

https://github.com/GoogleCloudPlatform/cloud-sql-proxy/blob/main/examples/k8s-sidecar/proxy_with_sa_key.yaml

--update code and push to kind pod --

### when creating a dockerfile/code change

1. rebuild your image with a new tag e.g.

```diff
-    docker build -t tpm-backend:0.0.2
+    docker build -t tpm-backend:0.0.2
```

2. change deployment file to point image section to new tag

```diff
containers:
        - name: tpm-backend
-          image: tpm-backend:0.0.1
+          image: tpm-backend:0.0.2
```

3. redeploy the deployment file using the following:

`kubectl apply -f <path to deployment.yaml>`

### View pod status

`kubectl get pods`

`kubectl get pods/<pod name>`

### View container logs from the pod its running in

`kubectl logs -f pod/tpm-backend-c5587bf9b-89pds -c cloud-sql-proxy`

-c stands for container
`kubectl logs -f pod/<pod name> -c <container name that is set in deployment.yaml>`



## Debugging

First see what is running within the pod and see if the descriptors make sense:
`kubectl describe pods/<pod name>`

### Debugging CrashLoopBackOff

### questions

what is the command to list namespaces and select one as default one? can only find config related commands online

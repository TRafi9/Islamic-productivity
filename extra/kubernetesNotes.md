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

--DEPLOYMENTS--

in the deployments yaml file you can state pods and number of replicas (more pods to allow failover)

in yaml file you set kind: Deployment for deployment files,
can set instance number under replicas in spec area

in the spec you also define the containers you want to run with:
name, image, ports, you can also add env variables and pass in secret values that you have stored already:
![alt text](image.png)

the env variables can also be mapped to configmaps

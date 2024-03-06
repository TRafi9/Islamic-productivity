# Kubernetes Component Documentation

## worker nodes

When you deploy kubernetes you get a cluster, within a cluster you will get a set of machines, which are called worker nodes, which will run your application.

These worker nodes will run pods inside them, which are your containers at run time. A worker node can have multiple applications running on it.

There are three components each worker node is comprised of:

- **kubelet** - An agent that runs on each worker node in a cluster, it makes sure that containers are running on a pod and are healthy, its primary responsibility is to manage the containers scheduled on the pods. It also is responsible to schedule to containers to run on the pods.
  Because kubelets are responsible for starting the pods, they have to interface with both the worker node resources e.g. cpu/ram and the container.
- **container runtime** - each node needs a container runtime installed so containers are able to be run there, there are different types of container runtimes e.g. containered/ docker engine.
- **kube-proxy** - this part includes a quick explanation of services and how they interact with kube-proxy. Because pods are ephemeral, their connection details can and will change. Due to this, we need a stable way to connect to the pods that run our containers. Enter Services. Services are logical abstractions which create stable endpoints to connect to the pods, we can define which pods we want the endpoints to expose connectivity to by setting them in `selector` part of the yaml. Doing so will expose those pods to be connected to from anywhere in the cluster. When network packets from the client hit the service, the service then distributes the network packet to one of the nodes which hold a pod that is tagged with the same `selector` value e.g. our `tpm-backend`. It does this based on an algorithm as to not cause any networking bottlenecks, e.g. it might use round-robin to send different network packets to different nodes. Once the network packet is at the node, the kube-proxy then chooses which pod to distribute that data to (as we could have multiple pods spun up for the same application to increase availability), it will make this decision based on multiple factors e.g. pod health and load balancing algorithms in place.



## Ingress
- ingress controllers run the ingress yaml rules, they are not defaultly contained and started within the cluster, your preferred ingress controller (e.g. nginx) needs to be added to the cluster, so it can run the ingress rules defined
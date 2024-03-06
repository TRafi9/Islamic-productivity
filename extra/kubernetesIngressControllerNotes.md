## Understanding Ingress Controllers

# I dont.

### deploy Nginx ingress controller through kubectl

The first thing is deploying your ingress controller to manage ingress as per your ingress manifest.

Applies a yaml manifest
`kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.10.0/deploy/static/provider/cloud/deploy.yaml`

The ingress controller normally goes into its own namespace to make it easier to manage its resources and also separate its own resources so it has the capability to manage its own traffic, the above command deploys nginx into its own namespace `ingress-nginx`

Even though the ingress controller is in its own namespace, by default it listens out for all the ingress resources deployed in our cluster across all namespaces.

### How to get IP address of ingress controller

The next step is to get the IP address of the ingress controller, specifically the loadBalancers IP address. This is the IP address you need to point external traffic to as an entry point to your cluster, i.e. you can set this IP address in the DNS settings of your cloudflare website.
To find this IP address you can run the following command:
`kubectl -n <ingress-controller-namespace> get svc`

Note that this external IP address will be `pending` in the kind cluster, this is because LoadBalancer services are designed to use the load-balancer infrastructure your cloud provider offers, and since we aren't running in the cloud, there is none.

This also means that when you apply your ingress manifest, there will be no address assigned, as the loadBalancer's external IP address will be the one assigned to your ingress.

### Kubernetes NodePort vs LoadBalancer vs Ingress? When should I use what?

https://medium.com/google-cloud/kubernetes-nodeport-vs-loadbalancer-vs-ingress-when-should-i-use-what-922f010849e0


## section on how to run nginx on kind cluster to get external IP address

finding env vars that run in your container on pods:
kubectl exec -it tpm-frontend-8dfb5484b-z5v22 -- /bin/sh

then you can call using os.env version in javascript or whatever

http://$TPM_BACKEND_SERVICE:$TPM_BACKEND_SERVICE_PORT/

### Ingress

    helm install stable/nginx-ingress \
    --namespace hotrod --name nginx \
    -f ingress/values.yaml


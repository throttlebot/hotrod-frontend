### Ingress

    helm install stable/nginx-ingress \
    --tiller-namespace gitlab-managed-apps \
    --namespace hotrod --name nginx \
    -f ingress/values.yaml

    kubectl apply -f ingress/frontend-ingress.yaml \
    --namespace hotrod

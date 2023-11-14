# kubazulo
Kubeconfig Authentication Helper for Kubernetes API-Server

## Configure kubeconfig

This is a sample configuration.

Please export the new file for `KUBECONFIG` variable:
>export KUBECONFIG=$HOME/.kube/kubazulo.config

Then copy the configuration below into the kubeconfig file:
```

users:
- name: kubazulo-azuread
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      args:
      - --client-id
      - {{ YOUR_CLIENT_ID }}
      - --tenant-id
      - {{ YOUR_TENANT_ID}}
      - --force-login
      - {{ false || true}}
      command: kubazulo
      env: null
      interactiveMode: IfAvailable
      provideClusterInfo: false
```



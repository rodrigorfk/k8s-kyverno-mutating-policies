# Kyverno Policies: Nginx Sidecar Injection

## Overview

This set of Kyverno policies demonstrates the capability of automatically injecting a sidecar container (in this case, an Nginx proxy) into Kubernetes pods and managing its associated TLS certificate. It's important to note that this example is intended to showcase Kyverno's features and is not designed as a production-ready solution or a real-world use case.

The policies serve as an educational tool to illustrate how Kyverno can:
1. Dynamically modify pod specifications
2. Generate additional resources based on pod creation
3. Integrate with other Kubernetes components like cert-manager

While using Nginx as a sidecar is the chosen example, the concepts demonstrated here can be applied to various other scenarios where sidecar injection and certificate management are required.

The example consists of two main policies:

1. `nginx-inject-sidecar`: Injects an Nginx sidecar container into pods.
2. `nginx-inject-sidecar-generate-certificate`: Generates a TLS certificate for the Nginx sidecar.

These policies showcase how Kyverno can enhance application deployments by adding configurable components and managing their security requirements automatically.

## Policy Details

### 1. nginx-inject-sidecar

This policy injects an Nginx sidecar container into pods when specific conditions are met.

Key features:
- Triggers on pod creation when the label `nginx.kyverno.io/inject: "true"` is present.
- Uses a ConfigMap to store sidecar configuration.
- Allows customization of CPU and memory resources via annotations.
- Mounts a TLS certificate for secure communication.

### 2. nginx-inject-sidecar-generate-certificate

This policy generates a TLS certificate for the Nginx sidecar using cert-manager.

Key features:
- Creates a Certificate resource for each pod with the Nginx sidecar.
- Configurable certificate properties (duration, renewal, issuer) via a ConfigMap.
- Automatically sets appropriate DNS names and IP addresses for the certificate.
- Carries over the ownerReferences from the pod to the generated Certificate resource, ensuring automatic cleanup when the pod's owner (e.g., a Deployment/ReplicaSet) is deleted.

## How It Works

1. When a pod with the label `nginx.kyverno.io/inject: "true"` is created, the `nginx-inject-sidecar` policy is triggered.
2. The policy injects an Nginx sidecar container with configurations from a specified ConfigMap.
3. Simultaneously, the `nginx-inject-sidecar-generate-certificate` policy creates a Certificate resource for the pod.
4. cert-manager processes the Certificate resource and generates a TLS certificate.
5. The Nginx sidecar uses this certificate for secure communication.
6. When the pod's owner (such as a Deployment/ReplicaSet) is deleted, the Certificate resource is automatically deleted due to the inherited ownerReferences.

## Configuration

### ConfigMap

Both policies use a ConfigMap named `nginx-inject-sidecar-<revision>` in the `kube-system` namespace. This ConfigMap should contain:

- `image`: The Nginx image to use for the sidecar.
- `certificateDuration`: Duration of the TLS certificate (default: 43800h0m0s).
- `certificateRenewBefore`: When to renew the certificate (default: 8760h0m0s).
- `issuerRefGroup`, `issuerRefKind`, `issuerRefName`: References to the cert-manager issuer.

### Custom Resource Requirements

You can customize the Nginx sidecar's resource requirements using the following annotations on the pod:

- `nginx.kyverno.io/proxyCPU`: CPU request
- `nginx.kyverno.io/proxyCPULimit`: CPU limit
- `nginx.kyverno.io/proxyMemory`: Memory request
- `nginx.kyverno.io/proxyMemoryLimit`: Memory limit

## Testing the Policies

1. Apply the policies to your cluster:
   ```bash
   kubectl apply -f policy.yaml
   ```

2. Create the necessary ConfigMap:
   ```yaml
   apiVersion: v1
   kind: ConfigMap
   metadata:
     name: nginx-inject-sidecar-default
     namespace: kube-system
   data:
     image: nginx:latest
     certificateDuration: "43800h0m0s"
     certificateRenewBefore: "8760h0m0s"
   ```

3. Create a test Deployment with the injection label:
   ```yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: test-nginx-sidecar
     labels:
       app: myapp
   spec:
     replicas: 1
     selector:
       matchLabels:
         app: myapp
     template:
       metadata:
         labels:
           app: myapp
           nginx.kyverno.io/inject: "true"
       spec:
         containers:
         - name: main-container
           image: busybox
           command: ['sh', '-c', 'echo The app is running! && sleep 3600']
   ```

4. Apply the Deployment:
   ```bash
   kubectl apply -f test-deployment.yaml
   ```

5. Verify the sidecar injection:
   ```bash
   kubectl get pod -l app=myapp -o yaml
   ```

   You should see:
   - An additional `nginx-sidecar` container in the pod spec.
   - TLS volume and volumeMount configurations.
   - Resource limits and requests as specified or defaulted.

6. Check for the generated Certificate resource:
   ```bash
   kubectl get certificate
   ```

   You should see a certificate named `nginx-myapp`.

7. Verify the ownerReferences:
   ```bash
   kubectl get certificate nginx-myapp -o yaml
   ```

   In the output, you should see an `ownerReferences` section that points to the ReplicaSet. This ensures that the certificate will be deleted when the Deployment is removed.

8. Test cleanup by deleting the Deployment:
   ```bash
   kubectl delete deployment test-nginx-sidecar
   ```

   Verify that both the pods and the certificate are deleted:
   ```bash
   kubectl get pods -l app=myapp
   kubectl get certificate nginx-myapp
   ```

   Both commands should return no resources, confirming that the cleanup worked as expected.

## Important Notes

- Ensure cert-manager is installed and configured in your cluster.
- The policies assume the existence of a secret named `nginx-<instance>` for environment variables. Make sure this secret exists or modify the policy accordingly.
- The `instance` value is derived from various labels. Ensure at least one of the specified labels is present on your pods.
- The policies use `admission: true` and `background: true`, meaning they apply to both new and existing resources.
- The generated Certificate resources inherit the ownerReferences from the pods. This means that when you delete a Deployment or other resource that owns the pod, the associated Certificate will also be automatically deleted, preventing resource leaks.

By using these policies, you can automatically inject Nginx sidecars and manage their TLS certificates, enhancing the security and functionality of your applications with minimal manual configuration.
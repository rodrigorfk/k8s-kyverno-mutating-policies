# Kyverno Policy: Mutate Public Registry Images

## What This Policy Does

This Kyverno policy, named "mutate-public-registry-images", provides automation around image registries in Kubernetes. It has two primary functions:

1. Replace public image registries with private registries for public images.
2. Add an imagePullSecret to the pod spec for private images.

The policy uses a ConfigMap to allow customizable mapping between public registries and their private counterparts, which is particularly useful for implementing ECR pull-through cache.

## Key Features

1. **Public Registry Replacement**: Automatically replaces specified public registries with private alternatives.
2. **ECR Integration**: Handles Amazon ECR registries, including cross-region scenarios.
3. **Image Existence Check**: Verifies if an image exists in ECR before applying the mutation.
4. **Flexible Mapping**: Uses a ConfigMap for easy updates to registry mappings.
5. **Support for Containers and InitContainers**: Applies mutations to both types of containers in a pod.
6. **ImagePullSecret Addition**: Automatically adds the necessary imagePullSecret for private registries.

## How It Works

The policy consists of several rules:

### 1. replace-image-registry-pod-containers and replace-image-registry-pod-init-containers

These rules handle the replacement of public registries with their private counterparts for both regular containers and init containers.

- They use the `public-registry-mapping` ConfigMap to determine the replacement registry.
- They handle special cases for Docker Hub, including adding the `library/` prefix when necessary.
- They support both tag-based and digest-based image references.

### 2. replace-image-registry-pod-containers-ecr and replace-image-registry-pod-init-containers-ecr

These rules specifically handle ECR registries:

- They extract the registry ID from the original ECR URL.
- They make an API call to an `ecr-image-checker` service to verify if the image exists in the target ECR registry.
- If the image exists, they replace the registry with the one in the current region.

### 3. add-imagepullsecret-containers and add-imagepullsecret-init-containers

These rules add an imagePullSecret to pods that use images from the private registry proxy.

## Configuration

The policy relies on two ConfigMaps:

1. `public-registry-mapping` in the `kube-system` namespace:
   - Maps public registries to their private counterparts.
   - Example mapping:
     ```yaml
     docker.io: 000000000000.dkr.ecr.eu-west-1.amazonaws.com/docker-hub
     gcr.io: my-registry-proxy.example.com/gcr
     ghcr.io: 000000000000.dkr.ecr.eu-west-1.amazonaws.com/github
     public.ecr.aws: 000000000000.dkr.ecr.eu-west-1.amazonaws.com/ecr-public
     quay.io: 000000000000.dkr.ecr.eu-west-1.amazonaws.com/quay
     registry.k8s.io: 000000000000.dkr.ecr.eu-west-1.amazonaws.com/kubernetes
     # ... (other mappings)
     ```

2. `cluster-config` in the `kube-system` namespace:
   - Provides the current AWS region for the cluster.
   - Example:
     ```yaml
     region: eu-west-1
     ```

## Testing the Policy

To test this policy:

1. Apply the policy to your cluster:
   ```bash
   kubectl apply -f policy.yaml
   ```

2. Ensure the required ConfigMaps are in place:
   ```bash
   kubectl apply -f configs.yaml
   ```

3. Test with different pod configurations:

   a. Public Registry (Docker Hub):
   ```yaml
   apiVersion: v1
   kind: Pod
   metadata:
     name: pod-docker-hub
     namespace: default
   spec:
     containers:
     - name: nginx
       image: nginx
   ```

   b. Google Container Registry (GCR):
   ```yaml
   apiVersion: v1
   kind: Pod
   metadata:
     name: test-pod-gcr
   spec:
     containers:
     - name: nginx
       image: gcr.io/nginx
   ```

   c. ECR in the Current Region:
   ```yaml
   apiVersion: v1
   kind: Pod
   metadata:
     name: test-pod-ecr-current-region
   spec:
     containers:
     - name: nginx
       image: 000000000000.dkr.ecr.eu-west-1.amazonaws.com/nginx
   ```

   d. ECR in a Different Region (Image Exists):
   ```yaml
   apiVersion: v1
   kind: Pod
   metadata:
     name: test-pod-ecr-different-region-exists
   spec:
     containers:
     - name: nginx
       image: 000000000000.dkr.ecr.us-east-1.amazonaws.com/nginx
   ```

   e. ECR in a Different Region (Image Does Not Exist):
   ```yaml
   apiVersion: v1
   kind: Pod
   metadata:
     name: test-pod-ecr-different-region-not-exists
   spec:
     containers:
     - name: nginx
       image: 111111111111.dkr.ecr.us-east-1.amazonaws.com/nginx
   ```

4. Apply each pod configuration:
   ```bash
   kubectl apply -f <pod-config-file>.yaml
   ```

5. Verify the mutations for each pod:
   ```bash
   kubectl get pod <pod-name> -o yaml
   ```

Expected Results:

a. For the Docker Hub image:
   - The image should be changed to the ECR equivalent specified in the `public-registry-mapping` ConfigMap.

b. For the GCR image:
   - The image should be changed to use the proxy registry specified in the `public-registry-mapping` ConfigMap.
   - An imagePullSecret should be added.

c. For the ECR image in the current region:
   - No changes should be made to the image.

d. For the ECR image in a different region (image exists):
   - The image registry should be changed to the current region's ECR.
   - The image path and tag should remain the same.

e. For the ECR image in a different region (image does not exist):
   - No changes should be made to the image.

## Important Notes

- This policy only affects newly created pods.
- It requires the `ecr-image-checker` service to be running in the `kube-system` namespace for ECR image verification.
- The policy assumes that the necessary permissions are in place for pulling images from the private registries.
- Regular updates to the `public-registry-mapping` ConfigMap may be necessary as new registries are added or changed.

By using this policy, you can ensure that all images used in your cluster are pulled from approved, private registries, enhancing security and potentially improving pull performance through caching mechanisms like ECR pull-through cache.
# Kyverno Policy: Enforce Default Hardware Architecture for Pods

## What This Policy Does

This Kyverno policy, named "mutate-pod-hardware-arch", addresses a common challenge in Kubernetes clusters with mixed hardware architectures. It automatically assigns a default architecture to newly created pods when one isn't explicitly specified. This helps prevent unintended scheduling issues that can occur when deploying workloads across nodes with different architectures.

Key features:
1. **Default Architecture Assignment**: Adds a nodeSelector for a default architecture (configurable, defaults to arm64) when not specified.
2. **Namespace-Level Configuration**: Allows setting a default architecture per namespace.
3. **Flexible Application**: Can be disabled for specific namespaces if needed.
4. **DaemonSet Awareness**: Skips enforcement for DaemonSet pods to allow them to run on all node types.

# Kyverno Policy: Enforce Default Hardware Architecture for Pods

## Why This Policy is Useful

In heterogeneous Kubernetes clusters (e.g., mixed amd64 and arm64 nodes), ensuring pods are scheduled on compatible nodes is crucial. This policy addresses a significant limitation in Kubernetes scheduling:

1. **Scheduler Limitation**: The Kubernetes scheduler is not aware of the hardware architecture required by the container images specified in a pod. It will schedule pods based on available resources, without considering architecture compatibility.

2. **Silent Failures**: Without proper architecture specification, pods might be scheduled on incompatible nodes. This leads to containers failing to start, often with cryptic error messages that don't clearly indicate an architecture mismatch.

3. **Image Pull Failures**: Even if a pod is scheduled on an incompatible node, the kubelet will attempt to pull and run the image. This results in runtime errors that can be difficult to diagnose, especially in large, diverse clusters.

4. **Resource Waste**: Incorrectly scheduled pods consume cluster resources (like IP addresses and node capacity) without providing any utility, potentially leading to resource exhaustion and decreased cluster efficiency.

5. **Deployment Delays**: In auto-scaling scenarios or during large deployments, architecture mismatches can cause significant delays as the system repeatedly tries to start pods on incompatible nodes.

By automatically adding architecture constraints, this policy helps prevent these issues, ensuring that pods are only scheduled on nodes with compatible hardware architectures. This leads to more predictable deployments, faster problem resolution, and better utilization of cluster resources.

## How It Works

The policy applies to the creation of new Pods and consists of one main rule:

### Rule: mutate-pod-hardware-arch

This rule applies when creating a new Pod:

1. It checks if the Pod is not part of a DaemonSet and doesn't already have an architecture specified (via nodeSelector or affinity).
2. If these conditions are met, it adds a nodeSelector for `kubernetes.io/arch`.
3. The architecture value is determined in this order:
   a. From the namespace label `policies.kyverno.io/default-arch` if present.
   b. Falls back to 'arm64' if not specified.

The policy can be disabled for specific namespaces by adding the label `policies.kyverno.io/disable-default-arch-enforcement: "true"`.

## Testing the Policy

To test this policy, you'll need a Kubernetes cluster with Kyverno installed. Here are steps to test it:

1. **Apply the Policy**:
   Apply the policy:
   ```bash
   kubectl apply -f policy.yaml
   ```

2. **Create a Test Namespace**:
   ```bash
   kubectl create namespace test-arch
   ```

3. **Create a Pod without Architecture Specification**:
   ```yaml
   apiVersion: v1
   kind: Pod
   metadata:
     name: test-pod
     namespace: test-arch
   spec:
     containers:
     - name: nginx
       image: nginx
   ```
   Save this as `test-pod.yaml` and apply:
   ```bash
   kubectl apply -f test-pod.yaml
   ```

4. **Verify the Architecture Assignment**:
   ```bash
   kubectl get pod test-pod -n test-arch -o yaml
   ```
   You should see a nodeSelector added with `kubernetes.io/arch: arm64`.

5. **Test Namespace-Level Configuration**:
   Label the namespace with a different architecture:
   ```bash
   kubectl label namespace test-arch policies.kyverno.io/default-arch=amd64
   ```
   Create another pod and verify it gets the amd64 architecture.
   ```bash
   kubectl delete -f test-pod.yaml
   kubectl apply -f test-pod.yaml
   kubectl get pod test-pod -n test-arch -o yaml
   ```

6. **Test Policy Disable**:
   Create a new namespace and disable the policy:
   ```bash
   kubectl create namespace no-arch-enforce
   kubectl label namespace no-arch-enforce policies.kyverno.io/disable-default-arch-enforcement=true
   ```
   Create a pod in this namespace and verify no architecture is added.
   ```bash
   kubectl delete -f test-pod.yaml
   kubectl apply -f test-pod.yaml
   kubectl get pod test-pod -n test-arch -o yaml
   ```

## Important Notes

- This policy only affects newly created Pods that don't already specify an architecture.
- It doesn't modify existing Pods or Pods created by DaemonSets.
- The policy uses `admission: true` and `background: false`, meaning it only applies during resource creation, not to existing resources.
- Ensure your cluster has nodes with the default architecture (arm64) or the architecture specified in namespace labels.

By using this policy, you can ensure that Pods in your cluster are always scheduled on nodes with compatible architectures, reducing the risk of deployment failures in heterogeneous environments.
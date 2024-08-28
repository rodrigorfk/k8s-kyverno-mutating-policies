# Kyverno Policy: Set KEDA Prometheus Scaler ServerAddress

## What This Policy Does

This Kyverno policy, named "keda-prometheus-serveraddress", addresses a common challenge when using KEDA (Kubernetes Event-driven Autoscaling) with Prometheus Scalers. It centralizes the management of the Prometheus server address, making it easier to update across multiple ScaledObjects.

Key features:
1. **Centralized Configuration**: Stores the Prometheus server address in a ConfigMap.
2. **Automatic Updates**: Modifies ScaledObjects to use the centralized server address.
3. **Flexible Application**: Only affects ScaledObjects that opt-in via an annotation.

## Why This Policy is Useful

Without this policy, updating the Prometheus server address requires modifying each individual ScaledObject. This can be time-consuming and error-prone, especially in environments with many microservices. This policy automates the process, reducing operational overhead and ensuring consistency.

## How It Works

The policy consists of two main rules:

### Rule 1: scaled-object

This rule applies when creating or updating a ScaledObject:

1. It checks for the annotation `keda.prometheus.sh/use-central-serveraddress: "true"`.
2. If present, it reads the server address from the `keda-prometheus-serveraddress` ConfigMap in the `observability` namespace.
3. It then adds or updates the `serverAddress` field in the Prometheus trigger of the ScaledObject.

### Rule 2: configmap

This rule applies when the `keda-prometheus-serveraddress` ConfigMap is created or updated:

1. It finds all ScaledObjects with the annotation `keda.prometheus.sh/use-central-serveraddress: "true"`.
2. For each matching ScaledObject, it updates the `serverAddress` in the Prometheus trigger to match the new value in the ConfigMap.

## Testing the Policy

To test this policy, you'll need a Kubernetes cluster with Kyverno and KEDA installed. Here are steps to test it:

1. **Apply the Policy**:
   Apply the policy:
   ```bash
   kubectl apply -f policy.yaml
   ```

2. **Create the ConfigMap**:
   ```bash
   kubectl create namespace observability
   kubectl create configmap keda-prometheus-serveraddress -n observability --from-literal=main=http://prometheus.observability:9090
   ```

3. **Create a ScaledObject**:
   ```yaml
   apiVersion: keda.sh/v1alpha1
   kind: ScaledObject
   metadata:
     name: prometheus-scaledobject
     annotations:
       keda.prometheus.sh/use-central-serveraddress: "true"
   spec:
     scaleTargetRef:
       deploymentName: my-deployment
     triggers:
      - type: prometheus
        metadata:
          serverAddress: http://old-address:9090
          query: sum(rate(http_requests_total{job="app"}[2m]))
   ```
   Save this as `scaledobject.yaml` and apply:
   ```bash
   kubectl apply -f scaledobject.yaml
   ```

4. **Verify the ServerAddress**:
   ```bash
   kubectl get scaledobject prometheus-scaledobject -o yaml
   ```
   You should see the `serverAddress` updated to the value from the ConfigMap.

5. **Update the ConfigMap**:
   ```bash
   kubectl patch configmap keda-prometheus-serveraddress -n observability --type merge -p '{"data":{"main":"http://new-prometheus.monitoring:9090"}}'
   ```

6. **Verify the Update**:
   Check the ScaledObject again:
   ```bash
   kubectl get scaledobject prometheus-scaledobject -o yaml
   ```
   The `serverAddress` should now reflect the new value.

## Important Notes

- The policy only affects ScaledObjects with the specific annotation. This allows for gradual adoption or exceptions if needed.
- Ensure the `keda-prometheus-serveraddress` ConfigMap exists in the `observability` namespace before creating ScaledObjects that depend on it.

By using this policy, you centralize the management of the Prometheus server address for KEDA scalers, simplifying operations and reducing the potential for configuration drift across your ScaledObjects.
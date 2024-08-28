# Kyverno Policy: Record Creation Details

## What This Policy Does

This Kyverno policy, named "record-creation-details", enhances the traceability of ConfigMap resources in your Kubernetes cluster. It does two main things:

1. **Adds Creator Information**: When a ConfigMap is created, it automatically adds an annotation `kyverno.io/created-by` containing details about who or what created the resource.

2. **Protects Creator Information**: Once added, it prevents this annotation from being modified or deleted, ensuring the integrity of the creation information.

## How It Works

### Rule 1: add-userinfo

This rule activates when a ConfigMap is being created. It adds an annotation `kyverno.io/created-by` to the ConfigMap's metadata. The value of this annotation is a string representation of the `userInfo` field from the admission request. This could include details like the username, user ID, groups the user belongs to, and extra information about the requester.

### Rule 2: prevent-updates-deletes-userinfo-annotations

This rule kicks in when someone tries to update a ConfigMap. It checks if the `kyverno.io/created-by` annotation exists. If it does, the rule ensures that:

- The annotation is not removed
- The value of the annotation is not changed

If either of these conditions is violated, the update is denied.

## Testing the Policy

To test this policy, you'll need a Kubernetes cluster with Kyverno installed. Here are some steps to test it:

1. **Apply the Policy**:
   Apply the policy:
   ```bash
   kubectl apply -f policy.yaml
   ```

2. **Test with Dry Run**:
   You can use `--dry-run=server` to test changes without actually applying them:
   ```bash
   kubectl create configmap test-cm --from-literal=key1=value1 --dry-run=server -o yaml --as=rodrigo --as-group=kubeadm:cluster-admins
   ```
   This will show you what the ConfigMap would look like if created, including the new annotation, the `--as` argument allow you to override the API server username and simulate different users. Here is an output example:
   ```yaml
   apiVersion: v1
   data:
     key1: value1
   kind: ConfigMap
   metadata:
     annotations:
       kyverno.io/created-by: '{"groups":["kubeadm:cluster-admins","system:authenticated"],"username":"rodrigo"}'
     creationTimestamp: "2024-08-27T14:01:18Z"
     name: test-cm
     namespace: default
     uid: 2e348e00-c7fe-4bfa-8958-04a4dc27c263
  ```


3. **Create a ConfigMap**:
   ```bash
   kubectl create configmap test-cm --from-literal=key1=value1
   ```

4. **Verify the Annotation**:
   ```bash
   kubectl get configmap test-cm -o yaml
   ```
   You should see the `kyverno.io/created-by` annotation in the output.

5. **Try to Modify the Annotation**:
   ```bash
   kubectl annotate configmap test-cm kyverno.io/created-by="new-value" --overwrite
   ```
   This should be denied.

6. **Try to Remove the Annotation**:
   ```bash
   kubectl annotate configmap test-cm kyverno.io/created-by-
   ```
   This should also be denied.

## Important Notes

- This policy only applies to ConfigMaps. To extend it to other resource types, modify the `kinds` field in both rules.
- The policy uses `validationFailureAction: Enforce`, meaning it will block operations that violate the policy. To make it advisory, change this to `Audit`.
- The `background: false` setting means this policy only applies to new resources, not existing ones.

By using this policy, you gain better visibility into who is creating ConfigMaps in your cluster, which can be crucial for auditing and security purposes.
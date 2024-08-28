# Kyverno's Mutating Webhooks playground

Welcome to the Kyverno Mutating Webhooks Playground! This repository provides hands-on examples of Kyverno policies, demonstrating both simple and complex logic for automating Kubernetes workflows using Dynamic Admission Control.

## 🎯 What's Inside
- Practical examples of Kyverno policies
- Simple to complex scenarios for Kubernetes workflow automation
- Ready-to-use development environment

## 🚀 Quick Start

### Prerequisites
- Docker installed and running on your machine
- Visual Studio Code with the "Remote - Containers" extension installed

### Setting Up the Development Environment
1. Open Visual Studio Code
2. Press `F1` or `cmd + shift + p` to open the command palette
3. Type and select "Dev Containers: Clone Repository in Container Volume..."
4. Enter the Git URL of this repository when prompted
5. Wait for the container to build and start (this may take a few minutes)

Once the container is ready, you'll have a fully configured environment with all necessary tools pre-installed!

## 🏃‍♂️ Running the Playground
Our playground uses a local Kubernetes cluster managed by [Kind](https://kind.sigs.k8s.io/). To set everything up:
1. Open the integrated terminal in VSCode (Ctrl+` or View > Terminal)
2. Run the following command:
```bash
make setup
```
This command will:
- Create a local Kubernetes cluster using Kind
- Install all necessary dependecies
- Set up Kyverno in the cluster

Once it is done, you can run `kubectl get pods -A` to list all pods and the output should be something similar to the following:
```
NAMESPACE            NAME                                                       READY   STATUS      RESTARTS   AGE
kube-system          coredns-7db6d8ff4d-gd5vs                                   1/1     Running     0          94m
kube-system          coredns-7db6d8ff4d-qz5sc                                   1/1     Running     0          94m
kube-system          etcd-my-cluster-control-plane                              1/1     Running     0          94m
kube-system          kindnet-7652s                                              1/1     Running     0          94m
kube-system          kindnet-gvvlp                                              1/1     Running     0          94m
kube-system          kindnet-qk9dw                                              1/1     Running     0          94m
kube-system          kube-apiserver-my-cluster-control-plane                    1/1     Running     0          94m
kube-system          kube-controller-manager-my-cluster-control-plane           1/1     Running     0          94m
kube-system          kube-proxy-7mkbq                                           1/1     Running     0          94m
kube-system          kube-proxy-t6snm                                           1/1     Running     0          94m
kube-system          kube-proxy-xpg6z                                           1/1     Running     0          94m
kube-system          kube-scheduler-my-cluster-control-plane                    1/1     Running     0          94m
kyverno              kyverno-admission-controller-776987899-gcmws               1/1     Running     0          80m
kyverno              kyverno-background-controller-86b9f95c96-vgcxn             1/1     Running     0          80m
kyverno              kyverno-cleanup-admission-reports-28746100-5jq47           0/1     Completed   0          9m20s
kyverno              kyverno-cleanup-cluster-admission-reports-28746100-klns2   0/1     Completed   0          9m20s
kyverno              kyverno-cleanup-cluster-ephemeral-reports-28746100-ckxgh   0/1     Completed   0          9m20s
kyverno              kyverno-cleanup-controller-7bbfc97569-h7tc4                1/1     Running     0          80m
kyverno              kyverno-cleanup-ephemeral-reports-28746100-dp769           0/1     Completed   0          9m20s
kyverno              kyverno-cleanup-update-requests-28746100-r779s             0/1     Completed   0          9m20s
kyverno              kyverno-reports-controller-665ccb5b65-fffzw                1/1     Running     0          80m
local-path-storage   local-path-provisioner-7d4d9bdcc5-zgr52                    1/1     Running     0          94m
```

## 📚 Exploring the Examples
After setup, you're ready to explore the examples:

1. Navigate to the `src/examples/` directory
2. Each subdirectory contains a specific scenario or use case
3. Follow the README in each example directory for specific instructions

## 🛠️ Troubleshooting
If you encounter any issues:

1. Ensure Docker is running and has enough resources allocated
2. Try rebuilding the dev container: F1 > "Dev Containers: Rebuild Container"
3. Check the "Problems" tab in VS Code for any error messages

For more help, please open an issue in this repository.
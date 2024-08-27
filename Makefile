create-cluster:
	@kind get clusters | grep -q my-cluster || kind create cluster --name my-cluster --config src/00-setup/kind-cluster.yaml --wait 5m

setup: create-cluster
	@helmfile -f src/00-setup/helmfile.yaml sync

delete-cluster:
	@kind delete cluster --name my-cluster

cleanup: delete-cluster
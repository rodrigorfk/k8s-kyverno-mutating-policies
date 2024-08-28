create-cluster:
	@kind get clusters | grep -q my-cluster || kind create cluster --name my-cluster --config src/setup/kind-cluster.yaml --wait 5m

install-dependencies:
	@helmfile -f src/setup/helmfile.yaml sync
	@kubectl apply -f https://raw.githubusercontent.com/kedacore/keda/v2.15.1/config/crd/bases/keda.sh_scaledobjects.yaml

deploy-ecr-image-checker:
	@(cd src/apps/ecr-image-checker && make build-container)
	@kind load docker-image 000000000000.dkr.ecr.us-east-1.amazonaws.com/ecr-image-checker:latest --name my-cluster
	@kubectl apply -f src/apps/ecr-image-checker/deploy/manifests.yaml
	@kubectl -n kube-system rollout restart deploy/ecr-image-checker

setup: create-cluster install-dependencies deploy-ecr-image-checker

delete-cluster:
	@kind delete cluster --name my-cluster

cleanup: delete-cluster
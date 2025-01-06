#!/bin/bash
set -x

if [ "$(uname -m)" == "aarch64" ]; then
  echo "ARM64"
  export ARCH=arm64
else
  echo "AMD64"
  export ARCH=amd64
fi

curl -Lo ./kind https://kind.sigs.k8s.io/dl/latest/kind-linux-$ARCH
chmod +x ./kind
mv ./kind /usr/local/bin/kind

curl -L -o kubebuilder https://go.kubebuilder.io/dl/latest/linux/$ARCH
chmod +x kubebuilder
mv kubebuilder /usr/local/bin/

KUBECTL_VERSION=$(curl -L -s https://dl.k8s.io/release/stable.txt)
curl -LO "https://dl.k8s.io/release/$KUBECTL_VERSION/bin/linux/$ARCH/kubectl"
chmod +x kubectl
mv kubectl /usr/local/bin/kubectl

docker network create -d=bridge --subnet=172.19.0.0/24 kind

kind version
kubebuilder version
docker --version
go version
kubectl version --client

echo "alias k=\"kubectl\"" >>~/.bashrc

# Update the system
sudo apt-get update

# Install Docker
sudo apt-get install -y docker.io

# Install apt-transport-https and curl
sudo apt-get install -y apt-transport-https curl

# Add Kubernetes apt repository
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list

# Update the system
sudo apt-get update

# Install Kubernetes
sudo apt-get install -y kubelet kubeadm kubectl

# Disable swap to allow Kubernetes to work properly
sudo swapoff -a

# Initialize the Kubernetes cluster with a specific version, in this case v1.21.0
sudo kubeadm init --kubernetes-version=v1.21.0

# Set up local kubeconfig
mkdir -p "$HOME/.kube"
sudo cp -i /etc/kubernetes/admin.conf "$HOME/.kube/config"
sudo chown $(id -u):$(id -g) "$HOME/.kube/config"

# Apply Flannel CNI network overlay
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml

# Allow the master node to run pods
kubectl taint nodes --all node-role.kubernetes.io/master-

echo "Kubernetes cluster setup completed."


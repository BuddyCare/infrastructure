REPO_NAME=$1
SERVICE_NAME=$2
INFRASTRUCTURE_NAME=$3
pwd
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
ansible-playbook -i ./$REPO_NAME/$INFRASTRUCTURE_NAME/ansible/inventory.ini ./$REPO_NAME/$INFRASTRUCTURE_NAME/ansible/playbooks/$SERVICE_NAME.yaml

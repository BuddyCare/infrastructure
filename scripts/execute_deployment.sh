SERVICE=$1
cd /home/ubuntu/ansible
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
ansible-playbook -i inventory.ini ./playbooks/$SERVICE.yaml

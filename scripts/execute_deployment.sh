REPO_NAME=$1
SERVICE_NAME=$2
INFRASTRUCTURE_NAME=$3
AWS_ACCESS_KEY_ID=$4
AWS_SECRET_ACCESS_KEY=$5
APP_VERSION=$6

pwd
echo "desde script execute deployment $AWS_ACCESS_KEY_ID $AWS_SECRET_ACCESS_KEY"
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
ansible-playbook -i ./$REPO_NAME/$INFRASTRUCTURE_NAME/ansible/inventory.ini ./$REPO_NAME/$INFRASTRUCTURE_NAME/ansible/playbooks/$SERVICE_NAME.yaml \
--extra-vars "aws_access_key_id=${AWS_ACCESS_KEY_ID} aws_secret_access_key=${AWS_SECRET_ACCESS_KEY} app_version=${APP_VERSION}"

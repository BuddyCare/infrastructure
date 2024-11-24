#!/bin/bash

# Actualizar los repositorios del sistema
sudo apt update

# Instalar dependencias necesarias
sudo apt install -y software-properties-common apt-transport-https ca-certificates curl gnupg

# Agregar el repositorio oficial de Ansible
sudo add-apt-repository --yes --update ppa:ansible/ansible

# Instalar Ansible
sudo apt install -y ansible

# Verificar la instalación de Ansible
ansible --version

# Crear el directorio de configuración de Ansible si no existe
sudo mkdir -p /etc/ansible

# Copiar ansible.cfg a /etc/ansible
if [ -f "/home/ubuntu/ansible/ansible.cfg" ]; then
    sudo cp /home/ubuntu/ansible/ansible.cfg /etc/ansible/ansible.cfg
    echo "ansible.cfg copiado a /etc/ansible/ansible.cfg"
else
    echo "El archivo ansible.cfg no se encontró en /home/ubuntu"
fi

# Copiar inventory.ini a /etc/ansible
if [ -f "/home/ubuntu/ansible/inventory.ini" ]; then
    sudo cp /home/ubuntu/ansible/inventory.ini /etc/ansible/inventory.ini
    echo "inventory.ini copiado a /etc/ansible/inventory.ini"
else
    echo "El archivo inventory.ini no se encontró en /home/ubuntu"
fi

# Cambiar permisos para el directorio y archivos de configuración
sudo chmod 644 /etc/ansible/ansible.cfg
sudo chmod 644 /etc/ansible/inventory.ini

# Instalar Helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Verificar la instalación de Helm
helm version

# Instalar la colección Kubernetes con ansible-galaxy
ansible-galaxy collection install kubernetes.core

# Mensaje final
echo "Ansible and Helm installed and configurated correctly."

echo "installing python libs for k8s...."
sudo apt-get install -y python3-kubernetes
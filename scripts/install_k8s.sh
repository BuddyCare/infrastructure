#!/bin/bash

# Variables
K3S_VERSION="v1.24.4+k3s1"
INSTALL_K3S_EXEC="server"

# Actualizar los repositorios del sistema
sudo apt update && sudo apt upgrade -y

# Instalar dependencias necesarias
sudo apt install -y curl bash

# Descargar e instalar K3s
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=${K3S_VERSION} sh -

# Habilitar y verificar el servicio de K3s
sudo systemctl enable k3s
sudo systemctl start k3s
sudo systemctl status k3s

# Mostrar el estado del clúster
echo "K3s ha sido instalado exitosamente. Aquí está el estado del clúster:"
sudo k3s kubectl get nodes

# Mensaje final
echo "La instalación de K3s ha finalizado. Puedes ejecutar 'sudo k3s kubectl get nodes' para verificar el estado del clúster."

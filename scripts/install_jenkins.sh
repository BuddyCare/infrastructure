#!/bin/bash

# Actualizar los repositorios
sudo apt update

# Instalar Java 17 y herramientas necesarias
sudo apt install -y openjdk-17-jre wget gnupg

# Verificar la instalación de Java 17
java -version

# Importar la clave GPG de Jenkins y agregar el repositorio
wget -q -O - https://pkg.jenkins.io/debian/jenkins.io-2023.key | sudo tee /etc/apt/trusted.gpg.d/jenkins.asc > /dev/null
echo "deb https://pkg.jenkins.io/debian-stable binary/" | sudo tee /etc/apt/sources.list.d/jenkins.list

# Actualizar los repositorios nuevamente para incluir Jenkins
sudo apt update

# Instalar Jenkins
sudo apt install -y jenkins

# Detener Jenkins para evitar que se ejecute antes de configurar
sudo systemctl stop jenkins

# Configurar Jenkins para que use Java 17
sudo sed -i 's|^JAVA_HOME=.*|JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64|' /etc/default/jenkins

# Crear el directorio init.groovy.d si no existe
if [ ! -d "/usr/share/jenkins/ref/init.groovy.d" ]; then
    sudo mkdir -p /usr/share/jenkins/ref/init.groovy.d
    sudo chmod 777 /usr/share/jenkins/ref/init.groovy.d
fi

# to create an admin user automatically
sudo cp /home/ubuntu/jenkins/pipelines/security.groovy /usr/share/jenkins/ref/init.groovy.d/security.groovy
# to create the pipeline for services
sudo cp /home/ubuntu/jenkins/pipelines/admin_service.groovy /usr/share/jenkins/ref/init.groovy.d/admin_service.groovy

echo 2.0 | sudo tee /usr/share/jenkins/ref/jenkins.install.UpgradeWizard.state > /dev/null

# Iniciar y habilitar el servicio de Jenkins
sudo systemctl start jenkins
sudo systemctl enable jenkins

# Verificar el estado de Jenkins
sudo systemctl status jenkins
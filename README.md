# infrastructure
to implement all the infrastructure

# current folder structure
```
infrastructure/
├── ec2/                             # Creación y manejo de instancias EC2
│   ├── create/                      # Código Go para crear instancias EC2
│   │   ├── main.go                  # Punto de entrada para la creación de instancias
│   │   ├── ec2_utils.go             # Funciones auxiliares para creación y manejo
│   │   ├── ssm_utils.go             # Funciones para ejecutar comandos SSM en instancias
│   │   └── README.md                # Documentación de la creación de instancias
│   ├── install/                     # Código Go para instalar servicios en instancias EC2
│   │   ├── main.go                  # Punto de entrada para instalaciones
│   │   ├── install_utils.go         # Funciones para manejar instalaciones
│   │   ├── scripts/                 # Scripts Bash para instalación directa
│   │   │   ├── install-jenkins.sh
│   │   │   ├── install-ansible.sh
│   │   │   ├── install-k8s.sh
│   │   │   └── common-setup.sh      # Configuraciones comunes (ej. actualizaciones)
│   │   └── README.md                # Documentación de las instalaciones
│   └── templates/                   # Plantillas para user-data
│       ├── userdata-jenkins.sh
│       ├── userdata-k8s.sh
│       └── userdata-common.sh
│
├── ansible/                         # Ansible para configuraciones avanzadas
│   ├── playbooks/                   # Playbooks para manejar Kubernetes
│   │   ├── deploy-service.yml       # Playbook para desplegar servicios en K8s
│   │   ├── configure-k8s.yml        # Playbook para configurar Kubernetes
│   │   ├── common.yml               # Playbook con tareas comunes
│   └── ansible.cfg                  # Configuración de Ansible
│
├── jenkins/                         # Configuración de Jenkins
│   ├── pipelines/                   # Definición de pipelines
│   │   ├── my-service-pipeline.groovy
│   │   └── Jenkinsfile              # Pipeline para CI/CD del servicio
│   ├── jobs/                        # Configuración de jobs Jenkins
│   │   ├── my-service-job.xml       # Configuración XML para un job
│   │   └── other-job.xml
│   └── dockerfiles/                 # Dockerfile para Jenkins (si es necesario)
│       └── Dockerfile.jenkins
│
└── pipeline/                        # Scripts auxiliares para CI/CD
├── build-and-deploy.sh          # Script de construcción y despliegue
└── test-runner.groovy           # Script para ejecutar pruebas
```
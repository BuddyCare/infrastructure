import jenkins.model.*
import hudson.model.*
import org.jenkinsci.plugins.workflow.job.WorkflowJob
import org.jenkinsci.plugins.workflow.cps.CpsFlowDefinition

// Crear o actualizar un job tipo Pipeline
def jobName = 'user-server'
def jenkins = Jenkins.instance
def myJob = jenkins.getItem(jobName) ?: jenkins.createProject(WorkflowJob, jobName)
def pipelineScript = '''
pipeline {
    agent any
    environment {
        REPO_NAME = 'core-services'
        SERVICE_NAME = 'admin_service'
        DOCKER_USER = 'alexrondon89'
        INFRASTRUCTURE_NAME = 'infrastructure'
    }
    stages {
        stage('Git Checkout core services') {
            steps {
                script {
                    sh "mkdir -p $WORKSPACE/$REPO_NAME"
                }
                dir("$WORKSPACE/$REPO_NAME") {
                    git url: 'https://github.com/BuddyCare/core-services.git', branch: 'main'
                }
            }
        }
        stage('Git Checkout infrastructure') {
            steps {
                script {
                    sh "mkdir -p $WORKSPACE/$REPO_NAME/$INFRASTRUCTURE_NAME"
                }
                dir("$WORKSPACE/$REPO_NAME/$INFRASTRUCTURE_NAME") {
                    git url: 'https://github.com/BuddyCare/infrastructure.git', branch: 'main'
                }
            }
        }
        stage('Building, tagging and pushing the Docker image'){
            steps{
                echo 'building image...'
                sh 'sudo docker image build -t $JOB_NAME:v1.$BUILD_ID -f $WORKSPACE/$REPO_NAME/$SERVICE_NAME/Dockerfile $WORKSPACE/$REPO_NAME/$SERVICE_NAME'
                echo 'tagging image...'
                sh 'sudo docker image tag $JOB_NAME:v1.$BUILD_ID $DOCKER_USER/$JOB_NAME:v1.$BUILD_ID'
                sh 'sudo docker image tag $JOB_NAME:v1.$BUILD_ID $DOCKER_USER/$JOB_NAME:latest'
                echo 'pushing docker image...'
                withCredentials([string(credentialsId: 'dockerhub_access1', variable: 'dockerhub_access1')]) {
                    sh "sudo docker login -u $DOCKER_USER -p $dockerhub_access1"
                    sh 'sudo docker image push $DOCKER_USER/$JOB_NAME:v1.$BUILD_ID'
                    sh 'sudo docker image push $DOCKER_USER/$JOB_NAME:latest'
                }
            }
        }
        stage('deployment k8s'){
            steps{
                echo 'Copying Ansible playbook from EC2...'
                sh '''
                    sudo chown ubuntu:ubuntu $WORKSPACE/$REPO_NAME/$INFRASTRUCTURE_NAME/scripts/execute_deployment.sh && \
                    sudo -u ubuntu chmod +x $WORKSPACE/$REPO_NAME/$INFRASTRUCTURE_NAME/scripts/execute_deployment.sh  && \
                    sudo -u ubuntu $WORKSPACE/$REPO_NAME/$INFRASTRUCTURE_NAME/scripts/execute_deployment.sh $REPO_NAME $SERVICE_NAME $INFRASTRUCTURE_NAME
                '''
            }
        }
    }
}
'''
myJob.definition = new CpsFlowDefinition(pipelineScript, true)
myJob.save()
jenkins.reload()

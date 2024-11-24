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
    }
    stages {
        stage('Git Checkout') {
            steps {
                git url: 'https://github.com/BuddyCare/core-services.git', branch: 'main'
            }
        }
        stage('Building, tagging and pushing the Docker image'){
            steps{
                echo 'building image...'
                sh 'sudo docker image build -t $JOB_NAME:v1.$BUILD_ID -f $WORKSPACE/$SERVICE_NAME/Dockerfile $WORKSPACE/$SERVICE_NAME'
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
                    sudo chown ubuntu:ubuntu /home/ubuntu/jenkins/scripts/execute_deployment.sh && \
                    sudo -u ubuntu chmod +x /home/ubuntu/jenkins/scripts/execute_deployment.sh  && \
                    sudo -u ubuntu /home/ubuntu/jenkins/scripts/execute_deployment.sh $SERVICE_NAME
                '''
            }
        }
    }
}
'''
myJob.definition = new CpsFlowDefinition(pipelineScript, true)
myJob.save()
jenkins.reload()

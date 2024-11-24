package internal

import (
	"context"
	"fmt"
	"github.com/BuddyCare/infrastructure/ec2/install/config"
	"github.com/BuddyCare/infrastructure/ec2/install/util"
	"golang.org/x/crypto/ssh"
)

type InstallSvc struct {
	sshClient    *ssh.Client
	instanceInfo config.InstanceInfo
	keyInfo      config.KeyInfo
}

func NewInstallSvc(ctx context.Context, instanceInfo config.InstanceInfo, keyInfo config.KeyInfo) (InstallSvc, error) {
	srv := InstallSvc{
		instanceInfo: instanceInfo,
		keyInfo:      keyInfo,
	}
	signer, err := util.GetSigner(srv.keyInfo.Location)
	if err != nil {
		return InstallSvc{}, err
	}

	sshClient, err := util.GetSshClient(signer, srv.instanceInfo.User, srv.instanceInfo.PublicIp)
	if err != nil {
		return InstallSvc{}, err
	}
	srv.sshClient = sshClient
	return srv, nil
}

func (installSvc InstallSvc) executeCommands(commands []string) error {
	for _, command := range commands {
		session, err := installSvc.sshClient.NewSession()
		if err != nil {
			return fmt.Errorf("failed to create session: %s", err)
		}

		output, err := session.CombinedOutput(command)
		if err != nil {
			return fmt.Errorf("failed to run command: %s", err)
		}
		fmt.Printf(string(output))
		session.Close()
	}
	return nil
}

func (installSvc InstallSvc) InstallDocker(ctx context.Context) error {
	fmt.Printf("installing docker...")
	dockerCommands := []string{
		"sudo apt-get update",
		"sudo apt-get install docker.io -y",
		"sudo systemctl start docker",
		"sudo systemctl enable docker",
		"sudo groupadd docker || true",            // use '|| true' in order to ignore error if group already exists
		"sudo usermod -aG docker ubuntu || true",  // make sure user is defined
		"sudo usermod -aG docker jenkins || true", // make sure user is defined
		"sudo systemctl restart docker",           // restart docker to set changes
	}
	err := installSvc.executeCommands(dockerCommands)
	if err != nil {
		return fmt.Errorf("failed to install docker: %s", err)
	}
	return nil
}

func (installSvc InstallSvc) InstallJenkins(ctx context.Context) error {
	fmt.Printf("installing jenkins...")
	fileCommands := []string{
		"mkdir -p /home/ubuntu/jenkins/pipelines",
		"mkdir -p /home/ubuntu/jenkins/scripts",
	}
	err := installSvc.executeCommands(fileCommands)
	if err != nil {
		return fmt.Errorf("failed to create jenkins folders: %s", err)
	}

	conn := util.RemoteConnInfo{
		User:        installSvc.instanceInfo.User,
		PublicIp:    installSvc.instanceInfo.PublicIp,
		KeyLocation: installSvc.keyInfo.Location,
	}
	//copy groovy files
	remotePathForGroovyFiles := "/home/ubuntu/jenkins/pipelines"
	localPathForGroovyFiles := "./jenkins/*"
	err = util.CopyFileToRemoteHost(conn, remotePathForGroovyFiles, localPathForGroovyFiles)
	if err != nil {
		return fmt.Errorf("failed to copy groovy files in remote server: %s", err)
	}

	// copy jenkins scripts
	remotePathForInstall := "/home/ubuntu/jenkins/scripts"
	localPathForInstall := []string{
		"./scripts/install_jenkins.sh",
		"./scripts/execute_deployment.sh",
	}
	for _, file := range localPathForInstall {
		err = util.CopyFileToRemoteHost(conn, remotePathForInstall, file)
		if err != nil {
			return fmt.Errorf("failed to copy scripts in remote server: %s", err)
		}
	}

	scriptCommands := []string{
		"chmod +x /home/ubuntu/jenkins/scripts/install_jenkins.sh && /home/ubuntu/jenkins/scripts/install_jenkins.sh",
	}
	err = installSvc.executeCommands(scriptCommands)
	if err != nil {
		return fmt.Errorf("failed to execute install_jenkins script: %s", err)
	}
	return nil
}

func (installSvc InstallSvc) InstallAnsible(ctx context.Context) error {
	fmt.Printf("installing ansible...")
	folderCommands := []string{
		"mkdir -p /home/ubuntu/ansible/charts",
		"mkdir -p /home/ubuntu/ansible/playbooks",
		"mkdir -p /home/ubuntu/ansible/scripts",
	}

	err := installSvc.executeCommands(folderCommands)
	if err != nil {
		return fmt.Errorf("failed to create ansible folders: %s", err)
	}

	conn := util.RemoteConnInfo{
		User:        installSvc.instanceInfo.User,
		PublicIp:    installSvc.instanceInfo.PublicIp,
		KeyLocation: installSvc.keyInfo.Location,
	}

	//copy ansible script files
	remotePathForScriptFiles := "/home/ubuntu/ansible/scripts"
	localPathForScriptFiles := "./scripts/install_ansible.sh"
	err = util.CopyFileToRemoteHost(conn, remotePathForScriptFiles, localPathForScriptFiles)
	if err != nil {
		return fmt.Errorf("failed to copy ansible install script in remote server: %s", err)
	}

	scriptCommands := []string{
		"chmod +x /home/ubuntu/ansible/scripts/install_ansible.sh && /home/ubuntu/ansible/scripts/install_ansible.sh",
	}
	err = installSvc.executeCommands(scriptCommands)
	if err != nil {
		return fmt.Errorf("failed to execute install_ansible script: %s", err)
	}

	//copy ansible files
	remotePathForAnsibleFiles := "/home/ubuntu/ansible"
	localPathForAnsibleFiles := []string{
		"./ansible/inventory.ini",
		"./ansible/ansible.cfg",
	}
	for _, file := range localPathForAnsibleFiles {
		err := util.CopyFileToRemoteHost(conn, remotePathForAnsibleFiles, file)
		if err != nil {
			return fmt.Errorf("failed to copy ansible files in remote server: %s", err)
		}
	}

	//copy playbooks files
	remotePathForPlaybookFiles := "/home/ubuntu/ansible/playbooks"
	localPathForPlaybookFiles := "./ansible/playbooks/*"
	err = util.CopyFileToRemoteHost(conn, remotePathForPlaybookFiles, localPathForPlaybookFiles)
	if err != nil {
		return fmt.Errorf("failed to copy ansible playbooks files in remote server: %s", err)
	}

	//copy charts files
	remotePathForChartsFiles := "/home/ubuntu/ansible/charts"
	localPathForChartsFiles := "./charts/*"
	err = util.CopyFileToRemoteHost(conn, remotePathForChartsFiles, localPathForChartsFiles)
	if err != nil {
		return fmt.Errorf("failed to copy ansible chart files in remote server: %s", err)
	}

	return nil
}

func (installSvc InstallSvc) InstallK8s(ctx context.Context) error {
	folderCommands := []string{
		"mkdir -p /home/ubuntu/k8s/scripts",
	}
	err := installSvc.executeCommands(folderCommands)
	if err != nil {
		return fmt.Errorf("failed to create k8s folder: %s", err)
	}

	conn := util.RemoteConnInfo{
		User:        installSvc.instanceInfo.User,
		PublicIp:    installSvc.instanceInfo.PublicIp,
		KeyLocation: installSvc.keyInfo.Location,
	}

	// copy k8s install script
	remotePathForInstall := "/home/ubuntu/k8s/scripts"
	localPathForInstall := "./scripts/install_k8s.sh"
	err = util.CopyFileToRemoteHost(conn, remotePathForInstall, localPathForInstall)
	if err != nil {
		return fmt.Errorf("failed to copy k8s install script in remote server: %s", err)
	}

	scriptCommand := []string{
		"chmod +x /home/ubuntu/k8s/scripts/install_k8s.sh && /home/ubuntu/k8s/scripts/install_k8s.sh",
	}
	err = installSvc.executeCommands(scriptCommand)
	if err != nil {
		return fmt.Errorf("failed to install k8s: %s", err)
	}

	return nil
}

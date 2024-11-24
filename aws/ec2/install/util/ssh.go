package util

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type RemoteConnInfo struct {
	KeyLocation string
	User        string
	PublicIp    string
}

func GetSigner(keyAddress string) (ssh.Signer, error) {
	err := os.Chmod(keyAddress, 0400)
	if err != nil {
		return nil, fmt.Errorf("failed to change file permissions: %v", err)
	}

	key, err := os.ReadFile(keyAddress)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	return signer, nil
}

func GetSshClient(signer ssh.Signer, user string, ip string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //delete for production
	}

	client, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect: %v", err)
	}
	return client, nil
}

func CopyFileToRemoteHost(conn RemoteConnInfo, remotePath, localPath string) error {
	matches, err := filepath.Glob(localPath)
	if err != nil {
		return fmt.Errorf("failed to glob %s: %v", matches, err)
	}
	if len(matches) == 0 {
		return fmt.Errorf("no files match the pattern %s", matches)
	}

	for _, file := range matches {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			log.Fatalf("file %s does not exist", file)
		}
		command := fmt.Sprintf("scp -r -o StrictHostKeyChecking=no -i %s %s %s@%s:%s",
			conn.KeyLocation, file, conn.User, conn.PublicIp, remotePath)
		cmd := exec.Command("/bin/sh", "-c", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to copy file %s: %s", file, err)
		}
		fmt.Printf("Output: %s\n", string(output))
	}
	return nil
}

package core

import (
	"log"
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func InitPublicKey() *ssh.PublicKeys {
	config := ReadConfig() // [TODO] use viper
	privateKeyFile := ExpandHome(config.PrivateKeyFile)

	_, err := os.Stat(privateKeyFile)
	if err != nil {
		log.Panicf("read file %s failed %s\n", privateKeyFile, err.Error())
	}

	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, "")
	if err != nil {
		log.Panicf("generate publickeys failed: %s\n", err.Error())
	}

	return publicKeys
}

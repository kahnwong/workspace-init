package core

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	cliBase "github.com/kahnwong/cli-base"
)

func initPublicKey() (*ssh.PublicKeys, error) {
	privateKeyFile := ExpandHome(config.PrivateKeyFile)
	_, err := cliBase.CheckIfConfigExists(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("private key doesn't exist at %s: %w", privateKeyFile, err)
	}

	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, "")
	if err != nil {
		return nil, fmt.Errorf("failed to generate public keys: %w", err)
	}

	return publicKeys, nil
}

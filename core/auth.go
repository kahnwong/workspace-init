package core

import (
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	cliBase "github.com/kahnwong/cli-base"
	"github.com/rs/zerolog/log"
)

func initPublicKey() *ssh.PublicKeys {
	privateKeyFile := ExpandHome(config.PrivateKeyFile)
	_, err := cliBase.CheckIfConfigExists(privateKeyFile)
	if err != nil {
		log.Fatal().Msgf("Private key doesn't exist at: %s", privateKeyFile)
	}

	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, "")
	if err != nil {
		log.Fatal().Msg("Generate publickeys failed")
	}

	return publicKeys
}

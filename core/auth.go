package core

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func initPublicKey() *ssh.PublicKeys {
	privateKeyFile := ExpandHome(viper.GetString("privateKeyFile"))

	_, err := os.Stat(privateKeyFile)
	if err != nil {
		log.Panic().Err(err).Msgf("Read file %s failed", privateKeyFile)
	}

	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, "")
	if err != nil {
		log.Panic().Err(err).Msg("Generate publickeys failed")
	}

	return publicKeys
}

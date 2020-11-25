package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/AlejandroAM91/vault/internal/app/vault"
)

const (
	encryptCmd = "encrypt"
	decryptCmd = "decrypt"
)

func parseCommand(cmd string, args []string) (config vault.Config) {
	flagSet := flag.NewFlagSet(cmd, flag.ExitOnError)
	flagSet.BoolVar(&config.Keep, "k", false, "Keeps original file")
	flagSet.StringVar(&config.Outfile, "o", "", "Output file")
	flagSet.Parse(args)

	if len(flagSet.Args()) < 1 {
		return
	}
	config.Infile = flagSet.Args()[0]
	return config
}

func encrypt(args []string) {
	config := parseCommand(encryptCmd, args)
	if config.Outfile == "" {
		config.Outfile = config.Infile + ".evf"
	}
	vault.EncryptFile(config)
}

func decrypt(args []string) {
	config := parseCommand(decryptCmd, args)
	if config.Outfile == "" {
		infile := config.Infile
		config.Outfile = infile[0 : len(infile)-len(filepath.Ext(infile))]
	}
	vault.DecryptFile(config)
}

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		return
	}

	switch os.Args[1] {
	case encryptCmd:
		encrypt(os.Args[2:])
	case decryptCmd:
		decrypt(os.Args[2:])
	}
}

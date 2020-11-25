package vault

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/AlejandroAM91/vault/pkg/vault"
	"golang.org/x/crypto/ssh/terminal"
)

type Config struct {
	Infile  string
	Outfile string
	Keep    bool
}

func EncryptFile(config Config) {

	pass, err := readPassword(true)
	if err != nil {
		panic(err)
	}

	ifile, err := os.Open(config.Infile)
	if err != nil {
		panic(err)
	}
	defer ifile.Close()

	ofile, err := os.Create(config.Outfile)
	if err != nil {
		panic(err)
	}
	defer ofile.Close()

	vfile, err := vault.NewWriter(ofile, pass)
	if err != nil {
		panic(err)
	}
	defer vfile.Close()

	_, err = io.Copy(vfile, ifile)
	if err != nil {
		panic(err)
	}
}

func DecryptFile(config Config) {
	pass, err := readPassword(false)
	if err != nil {
		panic(err)
	}

	ifile, err := os.Open(config.Infile)
	if err != nil {
		panic(err)
	}
	defer ifile.Close()

	ofile, err := os.Create(config.Outfile)
	if err != nil {
		panic(err)
	}
	defer ofile.Close()

	vfile, err := vault.NewReader(ifile, pass)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(ofile, vfile)
	if err != nil {
		panic(err)
	}
}

func readPassword(confirm bool) ([]byte, error) {
	fmt.Print("Enter Password: ")
	pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println("")
	if err != nil {
		return nil, err
	}

	if confirm {
		fmt.Print("Confirm Password: ")
		cpass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println("")
		if err != nil {
			return nil, err
		}

		if string(pass) != string(cpass) {
			return nil, errors.New("Passwords don´t match")
		}
	}

	return pass, err
}

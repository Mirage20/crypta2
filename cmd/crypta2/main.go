package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/nacl/box"
	"io/ioutil"
	"os"
)

var (
	versionString string
)

func main() {
	rootCmd := &cobra.Command{
		Use:           "crypta2",
		Short:         "Just another cryptographic tool",
		SilenceErrors: true,
		Version:       versionString,
	}
	rootCmd.AddCommand(makeGenKeyCmd())
	rootCmd.AddCommand(makeEncryptCmd())
	rootCmd.AddCommand(makeDecryptCmd())
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func makeGenKeyCmd() *cobra.Command {
	cmdGenKey := &cobra.Command{
		Use:          "genkey filename",
		Short:        "Generates a key pair",
		Long:         "Generates a key pair and save as filename.pub and filename.pvt to working directory",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fileNamePub := fmt.Sprintf("%s.pub", args[0])
			fileNamePvt := fmt.Sprintf("%s.pvt", args[0])

			if _, err := os.Stat(fileNamePvt); err == nil {
				return fmt.Errorf("private key file %q already exists", fileNamePvt)
			}

			if _, err := os.Stat(fileNamePub); err == nil {
				return fmt.Errorf("public key file %q already exists", fileNamePub)
			}

			pubKey, pvtKey, err := box.GenerateKey(rand.Reader)
			if err != nil {
				return fmt.Errorf("cannot generate key pair: %v", err)
			}
			pubEncoded := []byte(base64.StdEncoding.EncodeToString(pubKey[:]))
			pvtEncoded := []byte(base64.StdEncoding.EncodeToString(pvtKey[:]))
			fmt.Printf("Writing public key %q\n", fileNamePub)

			err = ioutil.WriteFile(fileNamePub, pubEncoded, 0644)
			if err != nil {
				return fmt.Errorf("cannot write public key %q: %v", fileNamePub, err)
			}

			fmt.Printf("Writing private key %q\n", fileNamePvt)
			err = ioutil.WriteFile(fileNamePvt, pvtEncoded, 0600)
			if err != nil {
				return fmt.Errorf("cannot write private key %q: %v", fileNamePvt, err)
			}
			return nil
		},
	}
	return cmdGenKey
}

func makeEncryptCmd() *cobra.Command {
	var input string
	var pubKeyFile string
	cmdEncrypt := &cobra.Command{
		Use:          "encrypt [(-f|--file=) input] (-p|--public-key=) public-key file",
		Short:        "Encrypt input with the given public key",
		Long:         "Encrypt input file or stdin with the given public key",
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			inBytes, err := readInput(input)
			if err != nil {
				return fmt.Errorf("cannot read input: %v", err)
			}
			pubEnc, err := ioutil.ReadFile(pubKeyFile)
			if err != nil {
				return fmt.Errorf("cannot read public key file %q: %v", pubKeyFile, err)
			}
			pubBytes, err := base64.StdEncoding.DecodeString(string(pubEnc))
			if err != nil {
				return fmt.Errorf("cannot decode public key: %v", err)
			}

			var pubKey32 [32]byte
			copy(pubKey32[:], pubBytes)

			value, err := box.SealAnonymous(nil, inBytes, &pubKey32, nil)
			if err != nil {
				return fmt.Errorf("cannot encrypt the input with provided public key")
			}
			fmt.Println(base64.StdEncoding.EncodeToString(value))
			return nil
		},
	}

	cmdEncrypt.Flags().StringVarP(&input, "file", "f", "", "Input file path. Uses standard input if not provided")
	cmdEncrypt.Flags().StringVarP(&pubKeyFile, "public-key", "p", "", "Public key file path for decrypting")
	_ = cmdEncrypt.MarkFlagRequired("public-key")
	return cmdEncrypt
}

func makeDecryptCmd() *cobra.Command {
	var input string
	var pvtKeyFile string
	var pubKeyFile string
	cmdDecrypt := &cobra.Command{
		Use:          "decrypt [(-f|--file=) input] (-k|--private-key=) private-key file (-p|--public-key=) public-key file",
		Short:        "Decrypt input with the given private and public keys",
		Long:         "Decrypt input file or stdin with the given private and public keys",
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			inEncodedBytes, err := readInput(input)
			if err != nil {
				return fmt.Errorf("cannot read input: %v", err)
			}
			inBytes, err := base64.StdEncoding.DecodeString(string(inEncodedBytes))
			if err != nil {
				return fmt.Errorf("cannot decode input: %v", err)
			}
			pubEnc, err := ioutil.ReadFile(pubKeyFile)
			if err != nil {
				return fmt.Errorf("cannot read public key file %q: %v", pubKeyFile, err)
			}
			pubBytes, err := base64.StdEncoding.DecodeString(string(pubEnc))
			if err != nil {
				return fmt.Errorf("cannot decode public key: %v", err)
			}
			pvtEnc, err := ioutil.ReadFile(pvtKeyFile)
			if err != nil {
				return fmt.Errorf("cannot read private key file %q: %v", pvtKeyFile, err)
			}
			pvtBytes, err := base64.StdEncoding.DecodeString(string(pvtEnc))
			if err != nil {
				return fmt.Errorf("cannot decode private key: %v", err)
			}
			var pubKey32 [32]byte
			var pvtKey32 [32]byte
			copy(pubKey32[:], pubBytes)
			copy(pvtKey32[:], pvtBytes)

			value, ok := box.OpenAnonymous(nil, inBytes, &pubKey32, &pvtKey32)
			if !ok {
				return fmt.Errorf("cannot decrypt the input with provided key pair")
			}
			fmt.Println(string(value))
			return nil
		},
	}

	cmdDecrypt.Flags().StringVarP(&input, "file", "f", "", "Base64 input file path. Uses standard input if not provided")
	cmdDecrypt.Flags().StringVarP(&pubKeyFile, "public-key", "p", "", "Public key file path for decrypting")
	cmdDecrypt.Flags().StringVarP(&pvtKeyFile, "private-key", "k", "", "Private key file path for decrypting")
	_ = cmdDecrypt.MarkFlagRequired("public-key")
	_ = cmdDecrypt.MarkFlagRequired("private-key")
	return cmdDecrypt
}

func readInput(file string) ([]byte, error) {
	if len(file) > 0 {
		return ioutil.ReadFile(file)
	} else {
		return ioutil.ReadAll(os.Stdin)
	}
}

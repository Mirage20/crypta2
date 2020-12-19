# Crypta2

Crypta2 uses anonymous encryption and decryption from [golang.org/x/crypto/nacl/box](https://godoc.org/golang.org/x/crypto/nacl/box) library to securely transfer your data.

## Getting Started

### Installation

#### Linux

1. Download

    ```bash
    curl -LO https://github.com/Mirage20/crypta2/releases/latest/download/crypta2-linux-x64.tar.gz
    ```
2. Extract

    ```bash
    tar -xzvf crypta2-linux-x64.tar.gz
    ```
3. Install

    ```bash
     sudo mv ./crypta2 /usr/local/bin/crypta2
    ```

#### MacOS

1. Download

    ```bash
    curl -LO https://github.com/Mirage20/crypta2/releases/latest/download/crypta2-darwin-x64.tar.gz
    ```
2. Extract

    ```bash
    tar -xzvf crypta2-darwin-x64.tar.gz
    ```
3. Install

    ```bash
     sudo mv ./crypta2 /usr/local/bin/crypta2
    ```

### Usage

1. Generate a key pair `my-key.pub` and `my-key.pvt`
   
    ```bash
    # The following command will generate a key pair and write it to current directory
    crypta2 genkey my-key
    ```
   
2. Encrypt using public key

    ```bash
    # Encrypt message.txt using public key and save it as base64 encoded payload.txt
    crypta2 encrypt -f message.txt -p my-key.pub > payload.txt
 
    # Read data from stdin pipe and encrypt
    echo "Hello World" | crypta2 encrypt -p my-key.pub > payload.txt

    # Encrypt data from stdin and output to stdout (press Ctrl+D after entering the secret)
    crypta2 encrypt -p my-key.pub
    ```

3. Decrypt using private and public key

    ```bash
    # Decrypt base64 encoded payload.txt using private key and save it as message.txt
    crypta2 decrypt -f payload.txt -k my-key.pvt -p my-key.pub > message.txt
 
    # Decrypt data from stdin and output to stdout (press Ctrl+D after entering the input)
    crypta2 decrypt -k my-key.pvt -p my-key.pub
    ```
    
4. Run `crypta2 --help` for more information

    ```text
    Usage:
      crypta2 [command]
    
    Available Commands:
      decrypt     Decrypt input with the given private and public keys
      encrypt     Encrypt input with the given public key
      genkey      Generates a key pair
      help        Help about any command
    
    Flags:
      -h, --help      help for crypta2
      -v, --version   version for crypta2
    
    Use "crypta2 [command] --help" for more information about a command.
        
    ```

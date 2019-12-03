package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/kovetskiy/godocs"
	"github.com/kovetskiy/lorg"
)

var (
	logger  = lorg.NewLog()
	version = "[manual build]"
)

const usage = `report

Usage:
    report --issuer-id --key-id --vendor-id --date [options]
    report -h | --help

Options:
    --issuer-id <string>  Issuer ID from App Store Connect.
    --key-id <string>     Private key ID from App Store Connect.
    --vendor-id <string>  Vendor ID from App Store Connect.
    --date <date>         Report date in YYYY-MM-DD format.
    --auth-key <path>     Path to private key.
                           [default: ./auth_key.p8].
    --debug               Enable debug output.
    --trace               Enable trace output.
    -h --help             Show this help.
`

func main() {
	args := godocs.MustParse(usage, version, godocs.UsePager)

	logger.SetIndentLines(true)

	if args["--debug"].(bool) {
		logger.SetLevel(lorg.LevelDebug)
	}

	if args["--trace"].(bool) {
		logger.SetLevel(lorg.LevelTrace)
	}

	key, err := getAuthKey(args["--auth-key"].(string))
	if err != nil {
		logger.Fatal(err)
	}

	token := getToken(args["--issuer-id"].(string), args["--key-id"].(string))

	tokenString, err := token.SignedString(key)
	if err != nil {
		logger.Fatal(err)
	}

	output, err := requestAppstore(
		args["--vendor-id"].(string),
		args["--date"].(string),
		tokenString,
	)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Print(string(output))
}

func getAuthKey(path string) (interface{}, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(file)

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

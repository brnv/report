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
    report --date <date> [options]
    report -h | --help

Options:
    --issuer-id <string>  Issuer ID from App Store Connect.
    --key-id <string>     Private key ID from App Store Connect.
    --vendor-id <string>  Vendor ID from App Store Connect.
    --date <date>         Report date in YYYY-MM-DD format.
    --auth-key <path>     Path to private key.
    --config <path>       Path to configuration file.
                           [default: config.toml]
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

	config, err := loadConfig(args["--config"].(string))
	if err != nil {
		logger.Fatal(err)
	}

	var (
		keyPath  = config.KeyPath
		issuerID = config.IssuerID
		keyID    = config.KeyID
		vendorID = config.VendorID
	)

	if args["--auth-key"] != nil {
		keyPath = args["--auth-key"].(string)
	}

	if args["--issuer-id"] != nil {
		issuerID = args["--issuer-id"].(string)
	}

	if args["--key-id"] != nil {
		keyID = args["--key-id"].(string)
	}

	if args["--vendor-id"] != nil {
		vendorID = args["--vendor-id"].(string)
	}

	date := args["--date"].(string)

	logger.Debugf("authorization key path: %s", keyPath)
	logger.Debugf("issuer id: %s", issuerID)
	logger.Debugf("key id: %s", keyID)
	logger.Debugf("vendor id: %s", vendorID)
	logger.Debugf("date: %s", date)

	key, err := getAuthKey(keyPath)
	if err != nil {
		logger.Fatal(err)
	}

	token := getToken(issuerID, keyID)

	tokenString, err := token.SignedString(key)
	if err != nil {
		logger.Fatal(err)
	}

	output, err := requestAppstore(
		vendorID,
		date,
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
		return nil, fmt.Errorf("can't read auth key: %s", err.Error())
	}

	block, _ := pem.Decode(file)

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

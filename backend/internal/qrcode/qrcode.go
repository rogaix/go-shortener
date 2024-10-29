package qrcode

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skip2/go-qrcode"
)

const qrCodeFolder = "qr_codes"

func GenerateQRCode(url, outputFile string) (string, error) {
	err := ensureQRCodeFolder()
	if err != nil {
		return "", fmt.Errorf("error creating QR code folder: %v", err)
	}

	fullPath := filepath.Join(qrCodeFolder, outputFile)
	err = qrcode.WriteFile(url, qrcode.Medium, 256, fullPath)
	if err != nil {
		return "", fmt.Errorf("error generating QR code: %v", err)
	}

	fmt.Printf("QR code for %s has been saved to %s\n", url, fullPath)
	return fullPath, nil
}

func ensureQRCodeFolder() error {
	return os.MkdirAll(qrCodeFolder, os.ModePerm)
}

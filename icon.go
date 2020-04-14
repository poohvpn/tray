package tray

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func getIconPath(iconPath string) (string, error) {
	if iconPath == "" {
		return "", errors.New("empty icon file path")
	}
	if strings.HasPrefix(iconPath, "base64:") {
		iconPath = strings.TrimSpace(iconPath[7:])
		if iconPath == "" {
			return "", errors.New("empty icon base64 data")
		}
		imgData, err := base64.StdEncoding.DecodeString(iconPath)
		if err != nil {
			return "", errors.Wrap(err, "decode icon base64 data")
		}
		hash := sha256.Sum256(imgData)
		iconPath = filepath.Join(os.TempDir(), hex.EncodeToString(hash[:]))
		err = ioutil.WriteFile(iconPath, imgData, 0600)
		if err != nil {
			return "", err
		}
	}
	return iconPath, nil
}

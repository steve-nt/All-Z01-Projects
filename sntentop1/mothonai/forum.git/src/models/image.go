package models

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"forum/src/utils"

	"github.com/gofrs/uuid"
)

const MaxImageSize = 20 << 20

var (
	JpegMagic = []byte{0xFF, 0xD8, 0xFF}
	PngMagic  = []byte{0x89, 0x50, 0x4E, 0x47}
	GifMagic  = []byte{0x47, 0x49, 0x46, 0x38}
)

func isValidImageType(buf []byte) bool {
	if len(buf) < 4 {
		return false
	}

	return bytes.HasPrefix(buf, JpegMagic) ||
		bytes.HasPrefix(buf, PngMagic) ||
		bytes.HasPrefix(buf, GifMagic)
}

func SaveImage(file multipart.File) (string, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return "", err
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return "", err
	}

	if len(fileBytes) > MaxImageSize {
		return "", ErrorImageTooBig
	}

	if !isValidImageType(fileBytes[:512]) {
		return "", ErrorInvalidImageType
	}

	ext := getImageExtension(fileBytes)
	filename := fmt.Sprintf("%s%s", uuid.Must(uuid.NewV4()).String(), ext)

	dir := filepath.Join("uploads", "images")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return "", err
	}

	dst := filepath.Join(dir, filename)
	out, err := os.Create(dst)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return "", err
	}
	defer out.Close()

	_, err = out.Write(fileBytes)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return "", err
	}

	return filepath.Join("uploads", "images", filename), nil
}

func getImageExtension(buf []byte) string {
	switch {
	case bytes.HasPrefix(buf, JpegMagic):
		return ".jpg"
	case bytes.HasPrefix(buf, PngMagic):
		return ".png"
	case bytes.HasPrefix(buf, GifMagic):
		return ".gif"
	default:
		return ""
	}
}

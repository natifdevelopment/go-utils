package utils

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"os"
	"path/filepath"
)

type ImageInfo struct {
	Width  int
	Height int
}

func GetImageInfo(file multipart.File) (*ImageInfo, error) {
	var data ImageInfo

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return nil, errors.New("failed to decode image")
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, errors.New("failed to seek file")
	}

	data.Width = img.Width
	data.Height = img.Height
	return &data, nil
}

func FindFileFromDir(folderPath string, fileName string) (string, error) {
	folderPath = filepath.Clean(folderPath)
	fileName = filepath.Clean(fileName)

	_, err := os.Open(folderPath)
	if err != nil {
		err := os.MkdirAll(folderPath, 0755)
		if err != nil {
			return "", fmt.Errorf("failed create dir: %s", err)
		}
	}

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return "", fmt.Errorf("failed read dir: %s", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name() == fileName {
			return filepath.Join(folderPath, fileName), nil
		}
	}

	return "", nil
}

func CreateFile(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %s", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("error writing to file: %s", err)
	}

	return nil
}

func AppendContentToFile(filePath string, content string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %s", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("\n\n%s", content))
	if err != nil {
		return fmt.Errorf("error appending file: %s", err)
	}

	return nil
}

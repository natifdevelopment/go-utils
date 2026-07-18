package utils

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxZipEntries          = 1000
	maxZipTotalSize        = 100 * 1024 * 1024
	maxZipCompressionRatio = 100
	zipReadHeaderBytes     = 512
	seekFileErrMsg         = "failed to seek file: %w"
)

var allowedUploadExtensions = map[string]bool{
	".pdf": true,
}

var allowedUploadMimeTypes = map[string]bool{
	"application/pdf": true,
}

type UploadValidationResult struct {
	DetectedType string
	Extension    string
	FileType     string
}

func ValidateExtension(filename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return "", errors.New("file has no extension")
	}
	if !allowedUploadExtensions[ext] {
		return "", fmt.Errorf("file extension %s is not allowed. Allowed: .pdf", ext)
	}
	return ext, nil
}

func ValidateMagicNumber(file multipart.File) (string, error) {
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file header: %w", err)
	}
	buffer = buffer[:n]

	detectedType := http.DetectContentType(buffer)

	if !allowedUploadMimeTypes[detectedType] {
		return "", fmt.Errorf("file content type %s is not allowed. Allowed: application/pdf", detectedType)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf(seekFileErrMsg, err)
	}

	return detectedType, nil
}

var pdfMagicNumber = []byte("%PDF")

func ValidatePDFMagicNumber(file multipart.File) error {
	header := make([]byte, 4)
	n, err := file.Read(header)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file header: %w", err)
	}
	header = header[:n]

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf(seekFileErrMsg, err)
	}

	if !bytes.Equal(header, pdfMagicNumber) {
		return errors.New("file does not have valid PDF magic number (%PDF)")
	}

	return nil
}

func ValidateFileSize(fileSize int64, maxBytes int64) error {
	if fileSize > maxBytes {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes", fileSize, maxBytes)
	}
	return nil
}

func GenerateSafeFilename(extension string) string {
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
}

func ValidateZipSafety(file multipart.File) error {
	buffer := make([]byte, zipReadHeaderBytes)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file for zip detection: %w", err)
	}
	buffer = buffer[:n]

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf(seekFileErrMsg, err)
	}

	detectedType := http.DetectContentType(buffer)
	if detectedType != "application/zip" && !strings.HasPrefix(detectedType, "application/zip") {
		return nil
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read zip file: %w", err)
	}
	defer func() {
		file.Seek(0, io.SeekStart)
	}()

	zipReader, err := zip.NewReader(bytes.NewReader(fileBytes), int64(len(fileBytes)))
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}

	if len(zipReader.File) > maxZipEntries {
		return fmt.Errorf("zip contains too many entries (%d), maximum allowed is %d", len(zipReader.File), maxZipEntries)
	}

	totalUncompressedSize := int64(0)
	totalCompressedSize := int64(0)

	for _, zipFile := range zipReader.File {
		totalUncompressedSize += int64(zipFile.UncompressedSize64)
		totalCompressedSize += int64(zipFile.CompressedSize64)

		if totalUncompressedSize > maxZipTotalSize {
			return fmt.Errorf("zip total uncompressed size exceeds maximum allowed size of %d bytes", maxZipTotalSize)
		}

		if err := validateZipEntry(zipFile); err != nil {
			return err
		}
	}

	if totalCompressedSize > 0 && totalUncompressedSize/totalCompressedSize > maxZipCompressionRatio {
		return fmt.Errorf("zip overall compression ratio is suspicious (possible zip bomb)")
	}

	return nil
}

func validateZipEntry(zipFile *zip.File) error {
	if zipFile.UncompressedSize64 > 0 && zipFile.CompressedSize64 > 0 {
		ratio := zipFile.UncompressedSize64 / zipFile.CompressedSize64
		if ratio > maxZipCompressionRatio {
			return fmt.Errorf("zip entry %s has suspicious compression ratio (possible zip bomb)", zipFile.Name)
		}
	}

	if strings.Contains(zipFile.Name, "..") {
		return fmt.Errorf("zip entry %s contains path traversal characters", zipFile.Name)
	}

	return nil
}

func ValidateUploadedFile(file *multipart.FileHeader, maxBytes int64) (*UploadValidationResult, error) {
	ext, err := ValidateExtension(file.Filename)
	if err != nil {
		return nil, err
	}

	if err := ValidateFileSize(file.Size, maxBytes); err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	if err := ValidatePDFMagicNumber(src); err != nil {
		return nil, err
	}

	detectedType, err := ValidateMagicNumber(src)
	if err != nil {
		return nil, err
	}

	if err := ValidateZipSafety(src); err != nil {
		return nil, err
	}

	fileType := strings.Split(detectedType, "/")[0]

	return &UploadValidationResult{
		DetectedType: detectedType,
		Extension:    ext,
		FileType:     fileType,
	}, nil
}

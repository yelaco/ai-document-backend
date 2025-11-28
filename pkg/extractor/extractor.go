package extractor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

func ExtractTextFromFile(fileName string, file *os.File) (string, error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".txt", ".md":
		return extractPlainText(file)
	case ".pdf":
		if err := file.Close(); err != nil {
			return "", fmt.Errorf("extractor.ExtractText: failed to close file: %w", err)
		}
		return extractTextFromPDF(file.Name())
	case ".docx":
		return extractTextFromDOCX(file)
	default:
		return "", fmt.Errorf("extractor.ExtractText: unsupported file type: %s", ext)
	}
}

func extractPlainText(file *os.File) (string, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, file)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func extractTextFromPDF(filePath string) (string, error) {
	file, reader, err := pdf.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var buf bytes.Buffer
	b, err := reader.GetPlainText()
	if err != nil {
		return "", err
	}
	_, err = buf.ReadFrom(b)
	if err != nil {
		return "", err
	}
	content := buf.String()
	return content, nil
}

func extractTextFromDOCX(r io.Reader) (string, error) {
	return "DOCX content", nil
}

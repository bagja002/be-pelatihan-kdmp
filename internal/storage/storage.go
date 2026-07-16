// Package storage menyimpan berkas unggahan (CV & sertifikat) ke disk dengan
// nama acak (UUID) dan validasi tipe/ukuran. Path yang disimpan di DB relatif
// terhadap root, sehingga root bisa dipindah tanpa migrasi data.
package storage

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrFileType     = errors.New("tipe berkas tidak diizinkan")
	ErrFileTooLarge = errors.New("ukuran berkas melebihi batas")
)

// Whitelist ekstensi per jenis berkas. Save/Validate lama tetap hanya PDF
// (CV & sertifikat); bahan ajar memakai SaveAs dengan ExtPDF/ExtPPT per slot.
var (
	ExtPDF = map[string]bool{".pdf": true}
	ExtPPT = map[string]bool{".ppt": true, ".pptx": true}
)

var allowedExt = ExtPDF

type Store struct {
	root     string
	maxBytes int64
}

func New(root string, maxBytes int64) *Store {
	return &Store{root: root, maxBytes: maxBytes}
}

// Validate memeriksa ekstensi (case-insensitive) & ukuran — hanya PDF.
func (s *Store) Validate(filename string, size int64) error {
	return s.ValidateAs(filename, size, allowedExt)
}

// ValidateAs seperti Validate, tapi terhadap whitelist ekstensi khusus.
func (s *Store) ValidateAs(filename string, size int64, allowed map[string]bool) error {
	if size > s.maxBytes {
		return ErrFileTooLarge
	}
	if !allowed[strings.ToLower(filepath.Ext(filename))] {
		return ErrFileType
	}
	return nil
}

// Save memvalidasi (hanya PDF) lalu menulis fh ke <root>/<subdir>/<uuid><ext>.
// Mengembalikan path relatif terhadap root (mis. "pelatih/ab-12.pdf").
func (s *Store) Save(fh *multipart.FileHeader, subdir string) (string, error) {
	return s.SaveAs(fh, subdir, allowedExt)
}

// SaveAs seperti Save, tapi dengan whitelist ekstensi khusus (mis. ExtPPT).
func (s *Store) SaveAs(fh *multipart.FileHeader, subdir string, allowed map[string]bool) (string, error) {
	if err := s.ValidateAs(fh.Filename, fh.Size, allowed); err != nil {
		return "", err
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	name := uuid.NewString() + ext

	dir := filepath.Join(s.root, subdir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dstPath := filepath.Join(dir, name)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		_ = os.Remove(dstPath)
		return "", err
	}
	return filepath.ToSlash(filepath.Join(subdir, name)), nil
}

// Path mengubah path relatif kembali ke path absolut untuk diunduh.
func (s *Store) Path(rel string) string {
	return filepath.Join(s.root, filepath.FromSlash(rel))
}

// Remove menghapus berkas berdasarkan path relatif (best-effort cleanup).
func (s *Store) Remove(rel string) error {
	if rel == "" {
		return nil
	}
	return os.Remove(s.Path(rel))
}

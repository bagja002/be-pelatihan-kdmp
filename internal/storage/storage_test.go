package storage

import (
	"bytes"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func makeFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	t.Helper()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := fw.Write(content); err != nil {
		t.Fatalf("write: %v", err)
	}
	w.Close()
	r := multipart.NewReader(body, w.Boundary())
	form, err := r.ReadForm(1 << 20)
	if err != nil {
		t.Fatalf("ReadForm: %v", err)
	}
	return form.File["file"][0]
}

func TestValidate(t *testing.T) {
	s := New(t.TempDir(), 100)
	if err := s.Validate("cv.pdf", 50); err != nil {
		t.Errorf("pdf 50B harus valid, dapat: %v", err)
	}
	if err := s.Validate("CV.PDF", 50); err != nil {
		t.Errorf("PDF (case-insensitive) harus valid, dapat: %v", err)
	}
	if err := s.Validate("foto.png", 50); err != ErrFileType {
		t.Errorf("png sekarang ditolak (hanya PDF), dapat: %v", err)
	}
	if err := s.Validate("virus.exe", 50); err != ErrFileType {
		t.Errorf("exe harus ErrFileType, dapat: %v", err)
	}
	if err := s.Validate("cv.pdf", 200); err != ErrFileTooLarge {
		t.Errorf("200B > 100 harus ErrFileTooLarge, dapat: %v", err)
	}
}

func TestSaveAndRemove(t *testing.T) {
	root := t.TempDir()
	s := New(root, 1<<20)
	fh := makeFileHeader(t, "cv asli.pdf", []byte("halo dunia"))

	rel, err := s.Save(fh, "pelatih")
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if !strings.HasPrefix(rel, "pelatih/") || !strings.HasSuffix(rel, ".pdf") {
		t.Errorf("path relatif tak terduga: %q", rel)
	}
	if strings.Contains(rel, "asli") {
		t.Errorf("nama asli bocor ke path: %q", rel)
	}
	abs := s.Path(rel)
	got, err := os.ReadFile(abs)
	if err != nil {
		t.Fatalf("baca berkas tersimpan: %v", err)
	}
	if string(got) != "halo dunia" {
		t.Errorf("isi berkas salah: %q", got)
	}
	if filepath.Dir(abs) != filepath.Join(root, "pelatih") {
		t.Errorf("lokasi salah: %q", abs)
	}

	if err := s.Remove(rel); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if _, err := os.Stat(abs); !os.IsNotExist(err) {
		t.Errorf("berkas harus terhapus, stat err: %v", err)
	}
}

func TestSaveRejectsBadType(t *testing.T) {
	s := New(t.TempDir(), 1<<20)
	fh := makeFileHeader(t, "x.exe", []byte("MZ"))
	if _, err := s.Save(fh, "pelatih"); err != ErrFileType {
		t.Errorf("harus ErrFileType, dapat: %v", err)
	}
}

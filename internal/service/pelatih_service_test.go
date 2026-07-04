package service

import (
	"bytes"
	"errors"
	"mime/multipart"
	"os"
	"testing"

	"knmp-backend/internal/entity"
	"knmp-backend/internal/storage"
)

// fakePelatihRepo mengimplementasikan repository.PelatihRepository untuk test.
type fakePelatihRepo struct {
	exists    bool
	createErr error
	created   *entity.Pelatih
}

func (f *fakePelatihRepo) ExistsNIP(string) (bool, error) { return f.exists, nil }
func (f *fakePelatihRepo) Create(p *entity.Pelatih) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.created = p
	return nil
}
func (f *fakePelatihRepo) FindAll() ([]entity.Pelatih, error) { return nil, nil }
func (f *fakePelatihRepo) FindByID(uint) (*entity.Pelatih, error) {
	return nil, nil
}
func (f *fakePelatihRepo) Delete(uint) error { return nil }
func (f *fakePelatihRepo) FindSertifikat(uint) (*entity.SertifikatKeahlian, error) {
	return nil, nil
}

func fileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	t.Helper()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	fw, _ := w.CreateFormFile("f", filename)
	fw.Write(content)
	w.Close()
	r := multipart.NewReader(body, w.Boundary())
	form, err := r.ReadForm(1 << 20)
	if err != nil {
		t.Fatal(err)
	}
	return form.File["f"][0]
}

func newInput(t *testing.T) RegisterPelatihInput {
	return RegisterPelatihInput{
		NamaLengkap: "Budi",
		NIP:         "123",
		Pendidikan:  "S1",
		UnitKerja:   "Diklat",
		Jabatan:     "Instruktur",
		CV:          fileHeader(t, "cv.pdf", []byte("cv")),
		Sertifikat: []CertUpload{
			{Nama: "Ahli K3", File: fileHeader(t, "k3.pdf", []byte("k3"))},
		},
	}
}

func countFiles(t *testing.T, dir string) int {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return 0
	}
	if err != nil {
		t.Fatal(err)
	}
	return len(entries)
}

func TestRegister_NIPBaru(t *testing.T) {
	root := t.TempDir()
	repo := &fakePelatihRepo{}
	svc := NewPelatihService(repo, storage.New(root, 1<<20))

	p, err := svc.Register(newInput(t))
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if p.NIP != "123" || len(p.Sertifikat) != 1 {
		t.Errorf("pelatih tak sesuai: %+v", p)
	}
	if p.CV == "" || p.Sertifikat[0].Berkas == "" {
		t.Errorf("path berkas kosong: %+v", p)
	}
	if n := countFiles(t, root+"/pelatih"); n != 2 {
		t.Errorf("harus 2 berkas tersimpan, dapat %d", n)
	}
}

func TestRegister_NIPGanda_DitolakTanpaBerkas(t *testing.T) {
	root := t.TempDir()
	repo := &fakePelatihRepo{exists: true}
	svc := NewPelatihService(repo, storage.New(root, 1<<20))

	_, err := svc.Register(newInput(t))
	if !errors.Is(err, ErrNIPExists) {
		t.Fatalf("harus ErrNIPExists, dapat %v", err)
	}
	if n := countFiles(t, root+"/pelatih"); n != 0 {
		t.Errorf("tak boleh ada berkas tersimpan saat NIP ganda, dapat %d", n)
	}
}

func TestRegister_CreateGagal_BersihkanBerkas(t *testing.T) {
	root := t.TempDir()
	repo := &fakePelatihRepo{createErr: errors.New("boom")}
	svc := NewPelatihService(repo, storage.New(root, 1<<20))

	if _, err := svc.Register(newInput(t)); err == nil {
		t.Fatal("harus error")
	}
	if n := countFiles(t, root+"/pelatih"); n != 0 {
		t.Errorf("berkas harus dibersihkan saat Create gagal, dapat %d", n)
	}
}

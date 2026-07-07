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
	exists     bool
	createErr  error
	created    *entity.Pelatih
	findResult *entity.Pelatih
	findErr    error
	deleted    bool
	deleteErr  error

	updateErr            error
	updated              *entity.Pelatih
	deletedSertifikatIDs []uint
	newSertifikat        []entity.SertifikatKeahlian
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
	return f.findResult, f.findErr
}
func (f *fakePelatihRepo) FindByNIP(string) (*entity.Pelatih, error) {
	return f.findResult, f.findErr
}
func (f *fakePelatihRepo) UpdateSelf(p *entity.Pelatih, deleteSertifikatIDs []uint, newSertifikat []entity.SertifikatKeahlian) error {
	if f.updateErr != nil {
		return f.updateErr
	}
	f.updated = p
	f.deletedSertifikatIDs = deleteSertifikatIDs
	f.newSertifikat = newSertifikat
	return nil
}
func (f *fakePelatihRepo) UpdateFields(id uint, p *entity.Pelatih) error {
	if f.updateErr != nil {
		return f.updateErr
	}
	f.updated = p
	return nil
}
func (f *fakePelatihRepo) Delete(uint) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}
	f.deleted = true
	return nil
}
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

func TestUpdateSelf_TanpaSertifikat_Ditolak(t *testing.T) {
	root := t.TempDir()
	repo := &fakePelatihRepo{
		findResult: &entity.Pelatih{ID: 1, NIP: "123", Sertifikat: []entity.SertifikatKeahlian{
			{ID: 7, IDPelatih: 1, NamaSertifikat: "Ahli K3", Berkas: "pelatih/k3.pdf"},
		}},
	}
	svc := NewPelatihService(repo, storage.New(root, 1<<20))

	// Hapus semua sertifikat lama (keep kosong) & tak menambah baru → harus ditolak.
	_, err := svc.UpdateSelf(UpdateSelfInput{NIP: "123", NamaLengkap: "Budi"})
	if !errors.Is(err, ErrNoSertifikat) {
		t.Fatalf("harus ErrNoSertifikat, dapat %v", err)
	}
	if repo.updated != nil {
		t.Error("tidak boleh menyentuh DB saat validasi gagal")
	}
	if n := countFiles(t, root+"/pelatih"); n != 0 {
		t.Errorf("tak boleh ada berkas tersimpan, dapat %d", n)
	}
}

func TestUpdateSelf_GantiCVdanSertifikat_BersihkanBerkasLama(t *testing.T) {
	root := t.TempDir()
	store := storage.New(root, 1<<20)

	oldCV, _ := store.Save(fileHeader(t, "cv-lama.pdf", []byte("cv")), "pelatih")
	oldCert, _ := store.Save(fileHeader(t, "k3-lama.pdf", []byte("k3")), "pelatih")

	repo := &fakePelatihRepo{
		findResult: &entity.Pelatih{ID: 1, NIP: "123", CV: oldCV, Sertifikat: []entity.SertifikatKeahlian{
			{ID: 7, IDPelatih: 1, NamaSertifikat: "Ahli K3", Berkas: oldCert},
		}},
	}
	svc := NewPelatihService(repo, store)

	// Ganti CV, hapus sertifikat lama (keep kosong), tambah 1 sertifikat baru.
	got, err := svc.UpdateSelf(UpdateSelfInput{
		NIP:         "123",
		NamaLengkap: "Budi Baru",
		LokasiTOT:   "Bandung",
		CV:          fileHeader(t, "cv-baru.pdf", []byte("cvbaru")),
		NewSertifikat: []CertUpload{
			{Nama: "Ahli Selam", File: fileHeader(t, "selam.pdf", []byte("selam"))},
		},
	})
	if err != nil {
		t.Fatalf("UpdateSelf: %v", err)
	}
	if got == nil || repo.updated == nil {
		t.Fatal("DB harus diperbarui")
	}
	if repo.updated.Status != "diperbarui" {
		t.Errorf("status harus 'diperbarui', dapat %q", repo.updated.Status)
	}
	if repo.updated.NIP != "123" {
		t.Errorf("NIP tidak boleh berubah, dapat %q", repo.updated.NIP)
	}
	if len(repo.deletedSertifikatIDs) != 1 || repo.deletedSertifikatIDs[0] != 7 {
		t.Errorf("sertifikat lama (id 7) harus dihapus, dapat %v", repo.deletedSertifikatIDs)
	}
	if len(repo.newSertifikat) != 1 {
		t.Errorf("harus ada 1 sertifikat baru, dapat %d", len(repo.newSertifikat))
	}
	// Berkas lama (CV + sertifikat) terhapus; berkas baru (CV + sertifikat) tersimpan → sisa 2.
	if n := countFiles(t, root+"/pelatih"); n != 2 {
		t.Errorf("harus 2 berkas tersisa (yang baru), dapat %d", n)
	}
}

func TestDelete_MenghapusBerkas(t *testing.T) {
	root := t.TempDir()
	store := storage.New(root, 1<<20)

	cvPath, err := store.Save(fileHeader(t, "cv.pdf", []byte("cv")), "pelatih")
	if err != nil {
		t.Fatalf("Save CV: %v", err)
	}
	certPath, err := store.Save(fileHeader(t, "k3.pdf", []byte("k3")), "pelatih")
	if err != nil {
		t.Fatalf("Save sertifikat: %v", err)
	}

	repo := &fakePelatihRepo{
		findResult: &entity.Pelatih{
			ID: 1,
			CV: cvPath,
			Sertifikat: []entity.SertifikatKeahlian{
				{ID: 1, IDPelatih: 1, NamaSertifikat: "Ahli K3", Berkas: certPath},
			},
		},
	}
	svc := NewPelatihService(repo, store)

	if err := svc.Delete(1); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if !repo.deleted {
		t.Error("repo.Delete harus dipanggil")
	}
	if n := countFiles(t, root+"/pelatih"); n != 0 {
		t.Errorf("berkas harus terhapus dari disk, dapat %d tersisa", n)
	}
}

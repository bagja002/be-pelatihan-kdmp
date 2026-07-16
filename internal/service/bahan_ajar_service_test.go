package service

import (
	"bytes"
	"errors"
	"mime/multipart"
	"os"
	"testing"

	"knmp-backend/internal/entity"
	"knmp-backend/internal/repository"
	"knmp-backend/internal/storage"
)

// fakeBahanAjarRepo menyimpan data di memori untuk pengujian service.
type fakeBahanAjarRepo struct {
	kategori map[uint]*entity.BahanAjarKategori
	items    map[uint]*entity.BahanAjar
	nextID   uint
}

func newFakeBahanAjarRepo() *fakeBahanAjarRepo {
	return &fakeBahanAjarRepo{
		kategori: map[uint]*entity.BahanAjarKategori{},
		items:    map[uint]*entity.BahanAjar{},
		nextID:   1,
	}
}

func (f *fakeBahanAjarRepo) id() uint { id := f.nextID; f.nextID++; return id }

func (f *fakeBahanAjarRepo) ListKategori() ([]entity.BahanAjarKategori, error) {
	var out []entity.BahanAjarKategori
	for _, k := range f.kategori {
		out = append(out, *k)
	}
	return out, nil
}

func (f *fakeBahanAjarRepo) FindKategoriByID(id uint) (*entity.BahanAjarKategori, error) {
	k, ok := f.kategori[id]
	if !ok {
		return nil, repository.ErrKategoriNotFound
	}
	return k, nil
}

func (f *fakeBahanAjarRepo) ExistsKategoriNama(nama string, excludeID uint) (bool, error) {
	for _, k := range f.kategori {
		if k.Nama == nama && k.ID != excludeID {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakeBahanAjarRepo) MaxKategoriUrutan() (int, error) {
	max := 0
	for _, k := range f.kategori {
		if k.Urutan > max {
			max = k.Urutan
		}
	}
	return max, nil
}

func (f *fakeBahanAjarRepo) CreateKategori(k *entity.BahanAjarKategori) error {
	k.ID = f.id()
	f.kategori[k.ID] = k
	return nil
}

func (f *fakeBahanAjarRepo) UpdateKategori(k *entity.BahanAjarKategori) error {
	f.kategori[k.ID] = k
	return nil
}

func (f *fakeBahanAjarRepo) CountItems(kategoriID uint) (int64, error) {
	var n int64
	for _, b := range f.items {
		if b.KategoriID == kategoriID {
			n++
		}
	}
	return n, nil
}

func (f *fakeBahanAjarRepo) DeleteKategori(id uint) error { delete(f.kategori, id); return nil }

func (f *fakeBahanAjarRepo) FindItemByID(id uint) (*entity.BahanAjar, error) {
	b, ok := f.items[id]
	if !ok {
		return nil, repository.ErrBahanAjarNotFound
	}
	cp := *b
	return &cp, nil
}

func (f *fakeBahanAjarRepo) MaxItemUrutan(kategoriID uint) (int, error) {
	max := 0
	for _, b := range f.items {
		if b.KategoriID == kategoriID && b.Urutan > max {
			max = b.Urutan
		}
	}
	return max, nil
}

func (f *fakeBahanAjarRepo) CreateItem(b *entity.BahanAjar) error {
	b.ID = f.id()
	cp := *b
	f.items[b.ID] = &cp
	return nil
}

func (f *fakeBahanAjarRepo) UpdateItem(b *entity.BahanAjar) error {
	cp := *b
	f.items[b.ID] = &cp
	return nil
}

func (f *fakeBahanAjarRepo) DeleteItem(id uint) error { delete(f.items, id); return nil }

// baFileHeader membuat *multipart.FileHeader tiruan (pola storage_test).
func baFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
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

func newBahanAjarTestService(t *testing.T) (BahanAjarService, *fakeBahanAjarRepo, *storage.Store) {
	t.Helper()
	repo := newFakeBahanAjarRepo()
	store := storage.New(t.TempDir(), 1<<20)
	return NewBahanAjarService(repo, store), repo, store
}

func TestBahanAjarCreateItemSimpanBerkas(t *testing.T) {
	svc, repo, store := newBahanAjarTestService(t)
	_ = repo.CreateKategori(&entity.BahanAjarKategori{Nama: "Kompetensi Umum", Urutan: 1})

	item, err := svc.CreateItem(BahanAjarItemInput{
		KategoriID: 1,
		Judul:      "Unit Kompetensi 1",
		FilePdf:    baFileHeader(t, "materi.pdf", []byte("pdf")),
		FilePpt:    baFileHeader(t, "materi.pptx", []byte("ppt")),
	})
	if err != nil {
		t.Fatalf("CreateItem: %v", err)
	}
	if item.FilePdf == "" || item.FilePpt == "" {
		t.Fatalf("path berkas harus terisi: pdf=%q ppt=%q", item.FilePdf, item.FilePpt)
	}
	if item.Urutan != 1 {
		t.Errorf("urutan otomatis harus 1, dapat %d", item.Urutan)
	}
	for _, rel := range []string{item.FilePdf, item.FilePpt} {
		if _, err := os.Stat(store.Path(rel)); err != nil {
			t.Errorf("berkas %q harus ada di disk: %v", rel, err)
		}
	}
}

func TestBahanAjarCreateItemTolakEkstensiSalah(t *testing.T) {
	svc, repo, _ := newBahanAjarTestService(t)
	_ = repo.CreateKategori(&entity.BahanAjarKategori{Nama: "Kepala Produksi", Urutan: 1})

	// pptx di slot PDF harus ditolak.
	_, err := svc.CreateItem(BahanAjarItemInput{
		KategoriID: 1,
		Judul:      "Modul 1",
		FilePdf:    baFileHeader(t, "materi.pptx", []byte("ppt")),
	})
	if !errors.Is(err, storage.ErrFileType) {
		t.Errorf("harus ErrFileType, dapat: %v", err)
	}
}

func TestBahanAjarCreateItemKategoriTidakAda(t *testing.T) {
	svc, _, _ := newBahanAjarTestService(t)
	_, err := svc.CreateItem(BahanAjarItemInput{KategoriID: 99, Judul: "X"})
	if !errors.Is(err, repository.ErrKategoriNotFound) {
		t.Errorf("harus ErrKategoriNotFound, dapat: %v", err)
	}
}

func TestBahanAjarUpdateItemGantiPdfHapusLama(t *testing.T) {
	svc, repo, store := newBahanAjarTestService(t)
	_ = repo.CreateKategori(&entity.BahanAjarKategori{Nama: "Penjamin Mutu", Urutan: 1})
	item, err := svc.CreateItem(BahanAjarItemInput{
		KategoriID: 1,
		Judul:      "Modul 1",
		FilePdf:    baFileHeader(t, "v1.pdf", []byte("v1")),
	})
	if err != nil {
		t.Fatalf("CreateItem: %v", err)
	}
	lama := item.FilePdf

	updated, err := svc.UpdateItem(item.ID, BahanAjarItemInput{
		Judul:   "Modul 1 (revisi)",
		FilePdf: baFileHeader(t, "v2.pdf", []byte("v2")),
	})
	if err != nil {
		t.Fatalf("UpdateItem: %v", err)
	}
	if updated.FilePdf == lama {
		t.Errorf("path PDF harus berganti")
	}
	if _, err := os.Stat(store.Path(lama)); !os.IsNotExist(err) {
		t.Errorf("berkas lama harus terhapus, stat err: %v", err)
	}
	if _, err := os.Stat(store.Path(updated.FilePdf)); err != nil {
		t.Errorf("berkas baru harus ada: %v", err)
	}
	if updated.Judul != "Modul 1 (revisi)" {
		t.Errorf("judul harus terbarui, dapat %q", updated.Judul)
	}
}

func TestBahanAjarDeleteItemHapusBerkas(t *testing.T) {
	svc, repo, store := newBahanAjarTestService(t)
	_ = repo.CreateKategori(&entity.BahanAjarKategori{Nama: "Administrasi Keuangan", Urutan: 1})
	item, err := svc.CreateItem(BahanAjarItemInput{
		KategoriID: 1,
		Judul:      "Modul 1",
		FilePdf:    baFileHeader(t, "a.pdf", []byte("a")),
		FilePpt:    baFileHeader(t, "a.ppt", []byte("a")),
	})
	if err != nil {
		t.Fatalf("CreateItem: %v", err)
	}
	if err := svc.DeleteItem(item.ID); err != nil {
		t.Fatalf("DeleteItem: %v", err)
	}
	for _, rel := range []string{item.FilePdf, item.FilePpt} {
		if _, err := os.Stat(store.Path(rel)); !os.IsNotExist(err) {
			t.Errorf("berkas %q harus terhapus, stat err: %v", rel, err)
		}
	}
	if _, err := svc.Item(item.ID); !errors.Is(err, repository.ErrBahanAjarNotFound) {
		t.Errorf("item harus hilang, dapat: %v", err)
	}
}

func TestBahanAjarDeleteKategoriBerisiDitolak(t *testing.T) {
	svc, repo, _ := newBahanAjarTestService(t)
	_ = repo.CreateKategori(&entity.BahanAjarKategori{Nama: "Kompetensi Umum", Urutan: 1})
	if _, err := svc.CreateItem(BahanAjarItemInput{KategoriID: 1, Judul: "Unit 1"}); err != nil {
		t.Fatalf("CreateItem: %v", err)
	}
	if err := svc.DeleteKategori(1); !errors.Is(err, ErrKategoriBerisi) {
		t.Errorf("harus ErrKategoriBerisi, dapat: %v", err)
	}
}

func TestBahanAjarKategoriNamaKosong(t *testing.T) {
	svc, _, _ := newBahanAjarTestService(t)
	if _, err := svc.CreateKategori("   ", nil); !errors.Is(err, ErrNamaKategoriKosong) {
		t.Errorf("harus ErrNamaKategoriKosong, dapat: %v", err)
	}
}

func TestBahanAjarKategoriNamaGanda(t *testing.T) {
	svc, _, _ := newBahanAjarTestService(t)
	if _, err := svc.CreateKategori("Kompetensi Umum", nil); err != nil {
		t.Fatalf("CreateKategori: %v", err)
	}
	if _, err := svc.CreateKategori("Kompetensi Umum", nil); !errors.Is(err, ErrNamaKategoriDipakai) {
		t.Errorf("harus ErrNamaKategoriDipakai, dapat: %v", err)
	}
	// Urutan otomatis kategori kedua = 2.
	k, err := svc.CreateKategori("Kepala Produksi", nil)
	if err != nil {
		t.Fatalf("CreateKategori: %v", err)
	}
	if k.Urutan != 2 {
		t.Errorf("urutan otomatis harus 2, dapat %d", k.Urutan)
	}
}

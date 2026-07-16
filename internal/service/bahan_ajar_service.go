package service

import (
	"errors"
	"mime/multipart"
	"strings"

	"knmp-backend/internal/entity"
	"knmp-backend/internal/repository"
	"knmp-backend/internal/storage"
)

var (
	// ErrKategoriBerisi: kategori masih punya item — kosongkan dulu sebelum hapus.
	ErrKategoriBerisi = errors.New("kategori masih berisi bahan ajar")
	// ErrNamaKategoriDipakai: nama kategori harus unik (di antara yang aktif).
	ErrNamaKategoriDipakai = errors.New("nama kategori sudah dipakai")
	// ErrNamaKategoriKosong: nama kategori tidak boleh kosong/hanya spasi.
	ErrNamaKategoriKosong = errors.New("nama kategori wajib diisi")
)

const bahanAjarSubdir = "bahan-ajar"

// BahanAjarItemInput — input buat/ubah item. Berkas nil = tidak diubah.
type BahanAjarItemInput struct {
	KategoriID uint
	Judul      string
	Urutan     *int
	FilePdf    *multipart.FileHeader
	FilePpt    *multipart.FileHeader
}

type BahanAjarService interface {
	ListKategori() ([]entity.BahanAjarKategori, error)
	CreateKategori(nama string, urutan *int) (*entity.BahanAjarKategori, error)
	UpdateKategori(id uint, nama string, urutan *int) (*entity.BahanAjarKategori, error)
	DeleteKategori(id uint) error

	Item(id uint) (*entity.BahanAjar, error)
	CreateItem(in BahanAjarItemInput) (*entity.BahanAjar, error)
	UpdateItem(id uint, in BahanAjarItemInput) (*entity.BahanAjar, error)
	DeleteItem(id uint) error
}

type bahanAjarService struct {
	repo  repository.BahanAjarRepository
	store *storage.Store
}

func NewBahanAjarService(repo repository.BahanAjarRepository, store *storage.Store) BahanAjarService {
	return &bahanAjarService{repo: repo, store: store}
}

func (s *bahanAjarService) ListKategori() ([]entity.BahanAjarKategori, error) {
	return s.repo.ListKategori()
}

func (s *bahanAjarService) CreateKategori(nama string, urutan *int) (*entity.BahanAjarKategori, error) {
	nama = strings.TrimSpace(nama)
	if nama == "" {
		return nil, ErrNamaKategoriKosong
	}
	dup, err := s.repo.ExistsKategoriNama(nama, 0)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, ErrNamaKategoriDipakai
	}
	u, err := s.resolveKategoriUrutan(urutan)
	if err != nil {
		return nil, err
	}
	k := &entity.BahanAjarKategori{Nama: nama, Urutan: u}
	if err := s.repo.CreateKategori(k); err != nil {
		return nil, err
	}
	return k, nil
}

func (s *bahanAjarService) UpdateKategori(id uint, nama string, urutan *int) (*entity.BahanAjarKategori, error) {
	k, err := s.repo.FindKategoriByID(id)
	if err != nil {
		return nil, err
	}
	nama = strings.TrimSpace(nama)
	if nama == "" {
		return nil, ErrNamaKategoriKosong
	}
	dup, err := s.repo.ExistsKategoriNama(nama, id)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, ErrNamaKategoriDipakai
	}
	k.Nama = nama
	if urutan != nil {
		k.Urutan = *urutan
	}
	if err := s.repo.UpdateKategori(k); err != nil {
		return nil, err
	}
	return k, nil
}

func (s *bahanAjarService) DeleteKategori(id uint) error {
	if _, err := s.repo.FindKategoriByID(id); err != nil {
		return err
	}
	n, err := s.repo.CountItems(id)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrKategoriBerisi
	}
	return s.repo.DeleteKategori(id)
}

func (s *bahanAjarService) resolveKategoriUrutan(urutan *int) (int, error) {
	if urutan != nil {
		return *urutan, nil
	}
	max, err := s.repo.MaxKategoriUrutan()
	if err != nil {
		return 0, err
	}
	return max + 1, nil
}

func (s *bahanAjarService) Item(id uint) (*entity.BahanAjar, error) {
	return s.repo.FindItemByID(id)
}

func (s *bahanAjarService) CreateItem(in BahanAjarItemInput) (*entity.BahanAjar, error) {
	if _, err := s.repo.FindKategoriByID(in.KategoriID); err != nil {
		return nil, err
	}

	var u int
	if in.Urutan != nil {
		u = *in.Urutan
	} else {
		max, err := s.repo.MaxItemUrutan(in.KategoriID)
		if err != nil {
			return nil, err
		}
		u = max + 1
	}

	// Simpan berkas; lacak untuk cleanup bila langkah berikutnya gagal.
	var saved []string
	cleanup := func() {
		for _, rel := range saved {
			_ = s.store.Remove(rel)
		}
	}

	pdfPath, pptPath := "", ""
	if in.FilePdf != nil {
		p, err := s.store.SaveAs(in.FilePdf, bahanAjarSubdir, storage.ExtPDF)
		if err != nil {
			return nil, err
		}
		saved = append(saved, p)
		pdfPath = p
	}
	if in.FilePpt != nil {
		p, err := s.store.SaveAs(in.FilePpt, bahanAjarSubdir, storage.ExtPPT)
		if err != nil {
			cleanup()
			return nil, err
		}
		saved = append(saved, p)
		pptPath = p
	}

	b := &entity.BahanAjar{
		KategoriID: in.KategoriID,
		Judul:      strings.TrimSpace(in.Judul),
		Urutan:     u,
		FilePdf:    pdfPath,
		FilePpt:    pptPath,
	}
	if err := s.repo.CreateItem(b); err != nil {
		cleanup()
		return nil, err
	}
	return b, nil
}

// UpdateItem mengubah judul/urutan/kategori dan mengganti berkas yang dikirim.
// Berkas lama yang tergantikan dihapus dari disk setelah DB sukses (best-effort).
func (s *bahanAjarService) UpdateItem(id uint, in BahanAjarItemInput) (*entity.BahanAjar, error) {
	b, err := s.repo.FindItemByID(id)
	if err != nil {
		return nil, err
	}
	if in.KategoriID != 0 && in.KategoriID != b.KategoriID {
		if _, err := s.repo.FindKategoriByID(in.KategoriID); err != nil {
			return nil, err
		}
		b.KategoriID = in.KategoriID
	}
	if judul := strings.TrimSpace(in.Judul); judul != "" {
		b.Judul = judul
	}
	if in.Urutan != nil {
		b.Urutan = *in.Urutan
	}

	var saved []string
	cleanup := func() {
		for _, rel := range saved {
			_ = s.store.Remove(rel)
		}
	}

	oldPdf, oldPpt := "", ""
	if in.FilePdf != nil {
		p, err := s.store.SaveAs(in.FilePdf, bahanAjarSubdir, storage.ExtPDF)
		if err != nil {
			return nil, err
		}
		saved = append(saved, p)
		oldPdf = b.FilePdf
		b.FilePdf = p
	}
	if in.FilePpt != nil {
		p, err := s.store.SaveAs(in.FilePpt, bahanAjarSubdir, storage.ExtPPT)
		if err != nil {
			cleanup()
			return nil, err
		}
		saved = append(saved, p)
		oldPpt = b.FilePpt
		b.FilePpt = p
	}

	if err := s.repo.UpdateItem(b); err != nil {
		cleanup()
		return nil, err
	}

	if oldPdf != "" {
		_ = s.store.Remove(oldPdf)
	}
	if oldPpt != "" {
		_ = s.store.Remove(oldPpt)
	}
	return b, nil
}

func (s *bahanAjarService) DeleteItem(id uint) error {
	b, err := s.repo.FindItemByID(id)
	if err != nil {
		return err
	}
	if err := s.repo.DeleteItem(id); err != nil {
		return err
	}
	// Best-effort cleanup berkas di disk.
	_ = s.store.Remove(b.FilePdf)
	_ = s.store.Remove(b.FilePpt)
	return nil
}

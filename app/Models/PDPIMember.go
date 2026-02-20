package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// PDPIMember represents a member from PDPI system
type PDPIMember struct {
	ID                    string     `db:"id" json:"id"`
	NPA                   string     `db:"npa" json:"npa"`
	NPANumeric            *int64     `db:"npa_numeric" json:"npa_numeric"`
	Nama                  string     `db:"nama" json:"nama"`
	Foto                  *string    `db:"foto" json:"foto"`
	Gelar                 *string    `db:"gelar" json:"gelar"`
	Gelar2                *string    `db:"gelar2" json:"gelar2"`
	Email                 *string    `db:"email" json:"email"`
	NoHP                  *string    `db:"no_hp" json:"no_hp"`
	NIK                   *string    `db:"nik" json:"nik"`
	JenisKelamin          *string    `db:"jenis_kelamin" json:"jenis_kelamin"`
	TempatLahir           *string    `db:"tempat_lahir" json:"tempat_lahir"`
	TglLahir              *time.Time `db:"tgl_lahir" json:"tgl_lahir"`
	AlamatRumah           *string    `db:"alamat_rumah" json:"alamat_rumah"`
	Cabang                *string    `db:"cabang" json:"cabang"`
	Provinsi              *string    `db:"provinsi" json:"provinsi"`
	KotaKabupaten         *string    `db:"kota_kabupaten" json:"kota_kabupaten"`
	KotaKabupatenKantor   *string    `db:"kota_kabupaten_kantor" json:"kota_kabupaten_kantor"`
	ProvinsiKantor        *string    `db:"provinsi_kantor" json:"provinsi_kantor"`
	Status                *string    `db:"status" json:"status"`
	Alumni                *string    `db:"alumni" json:"alumni"`
	ThnLulus              *int64     `db:"thn_lulus" json:"thn_lulus"`
	TempatTugas           *string    `db:"tempat_tugas" json:"tempat_tugas"`
	TempatPraktek1        *string    `db:"tempat_praktek_1" json:"tempat_praktek_1"`
	TempatPraktek1Tipe    *string    `db:"tempat_praktek_1_tipe" json:"tempat_praktek_1_tipe"`
	TempatPraktek1Tipe2   *string    `db:"tempat_praktek_1_tipe_2" json:"tempat_praktek_1_tipe_2"`
	TempatPraktek1Alkes   *string    `db:"tempat_praktek_1_alkes" json:"tempat_praktek_1_alkes"`
	TempatPraktek1Alkes2  *string    `db:"tempat_praktek_1_alkes_2" json:"tempat_praktek_1_alkes_2"`
	TempatPraktek2        *string    `db:"tempat_praktek_2" json:"tempat_praktek_2"`
	TempatPraktek2Tipe    *string    `db:"tempat_praktek_2_tipe" json:"tempat_praktek_2_tipe"`
	TempatPraktek2Tipe2   *string    `db:"tempat_praktek_2_tipe_2" json:"tempat_praktek_2_tipe_2"`
	TempatPraktek2Alkes   *string    `db:"tempat_praktek_2_alkes" json:"tempat_praktek_2_alkes"`
	TempatPraktek2Alkes2  *string    `db:"tempat_praktek_2_alkes_2" json:"tempat_praktek_2_alkes_2"`
	KotaKabupatenPraktek2 *string    `db:"kota_kabupaten_praktek_2" json:"kota_kabupaten_praktek_2"`
	ProvinsiPraktek2      *string    `db:"provinsi_praktek_2" json:"provinsi_praktek_2"`
	TempatPraktek3        *string    `db:"tempat_praktek_3" json:"tempat_praktek_3"`
	TempatPraktek3Tipe    *string    `db:"tempat_praktek_3_tipe" json:"tempat_praktek_3_tipe"`
	TempatPraktek3Tipe2   *string    `db:"tempat_praktek_3_tipe_2" json:"tempat_praktek_3_tipe_2"`
	TempatPraktek3Alkes   *string    `db:"tempat_praktek_3_alkes" json:"tempat_praktek_3_alkes"`
	TempatPraktek3Alkes2  *string    `db:"tempat_praktek_3_alkes_2" json:"tempat_praktek_3_alkes_2"`
	KotaKabupatenPraktek3 *string    `db:"kota_kabupaten_praktek_3" json:"kota_kabupaten_praktek_3"`
	ProvinsiPraktek3      *string    `db:"provinsi_praktek_3" json:"provinsi_praktek_3"`
	Subspesialis          *string    `db:"subspesialis" json:"subspesialis"`
	GelarFISR             *string    `db:"gelar_fisr" json:"gelar_fisr"`
	NoSTR                 *string    `db:"no_str" json:"no_str"`
	STRBerlakuSampai      *time.Time `db:"str_berlaku_sampai" json:"str_berlaku_sampai"`
	NoSIP                 *string    `db:"no_sip" json:"no_sip"`
	SIPBerlakuSampai      *time.Time `db:"sip_berlaku_sampai" json:"sip_berlaku_sampai"`
	UserID                *int64     `db:"user_id" json:"user_id"`
	SyncedAt              *time.Time `db:"synced_at" json:"synced_at"`
	CreatedAt             time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time  `db:"updated_at" json:"updated_at"`
}

// FindPDPIMemberByEmail finds a member by email in local database
func FindPDPIMemberByEmail(db *sqlx.DB, email string) (*PDPIMember, error) {
	var member PDPIMember
	query := `SELECT * FROM pdpi_members WHERE email = ? LIMIT 1`
	err := db.Get(&member, query, email)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// FindPDPIMemberByNPA finds a member by NPA in local database
func FindPDPIMemberByNPA(db *sqlx.DB, npa string) (*PDPIMember, error) {
	var member PDPIMember
	query := `SELECT * FROM pdpi_members WHERE npa = ? LIMIT 1`
	err := db.Get(&member, query, npa)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// FindPDPIMemberByUserID finds a member by user_id in local database
func FindPDPIMemberByUserID(db *sqlx.DB, userID int64) (*PDPIMember, error) {
	var member PDPIMember
	query := `SELECT * FROM pdpi_members WHERE user_id = ? LIMIT 1`
	err := db.Get(&member, query, userID)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// UpsertPDPIMember inserts or updates a PDPI member in local database
func UpsertPDPIMember(db *sqlx.DB, member *PDPIMember) error {
	query := `
		INSERT INTO pdpi_members (
			id, npa, npa_numeric, nama, foto, gelar, gelar2, email, no_hp, nik, 
			jenis_kelamin, tempat_lahir, tgl_lahir, alamat_rumah,
			cabang, provinsi, kota_kabupaten, kota_kabupaten_kantor, provinsi_kantor, status, alumni, thn_lulus,
			tempat_tugas, 
			tempat_praktek_1, tempat_praktek_1_tipe, tempat_praktek_1_tipe_2, tempat_praktek_1_alkes, tempat_praktek_1_alkes_2,
			tempat_praktek_2, tempat_praktek_2_tipe, tempat_praktek_2_tipe_2, tempat_praktek_2_alkes, tempat_praktek_2_alkes_2,
			kota_kabupaten_praktek_2, provinsi_praktek_2,
			tempat_praktek_3, tempat_praktek_3_tipe, tempat_praktek_3_tipe_2, tempat_praktek_3_alkes, tempat_praktek_3_alkes_2,
			kota_kabupaten_praktek_3, provinsi_praktek_3,
			subspesialis, gelar_fisr,
			no_str, str_berlaku_sampai, no_sip, sip_berlaku_sampai,
			user_id, synced_at
		) VALUES (
			:id, :npa, :npa_numeric, :nama, :foto, :gelar, :gelar2, :email, :no_hp, :nik,
			:jenis_kelamin, :tempat_lahir, :tgl_lahir, :alamat_rumah,
			:cabang, :provinsi, :kota_kabupaten, :kota_kabupaten_kantor, :provinsi_kantor, :status, :alumni, :thn_lulus,
			:tempat_tugas, 
			:tempat_praktek_1, :tempat_praktek_1_tipe, :tempat_praktek_1_tipe_2, :tempat_praktek_1_alkes, :tempat_praktek_1_alkes_2,
			:tempat_praktek_2, :tempat_praktek_2_tipe, :tempat_praktek_2_tipe_2, :tempat_praktek_2_alkes, :tempat_praktek_2_alkes_2,
			:kota_kabupaten_praktek_2, :provinsi_praktek_2,
			:tempat_praktek_3, :tempat_praktek_3_tipe, :tempat_praktek_3_tipe_2, :tempat_praktek_3_alkes, :tempat_praktek_3_alkes_2,
			:kota_kabupaten_praktek_3, :provinsi_praktek_3,
			:subspesialis, :gelar_fisr,
			:no_str, :str_berlaku_sampai, :no_sip, :sip_berlaku_sampai,
			:user_id, :synced_at
		)
		ON DUPLICATE KEY UPDATE
			npa_numeric = VALUES(npa_numeric),
			nama = VALUES(nama),
			foto = VALUES(foto),
			gelar = VALUES(gelar),
			gelar2 = VALUES(gelar2),
			email = VALUES(email),
			no_hp = VALUES(no_hp),
			nik = VALUES(nik),
			jenis_kelamin = VALUES(jenis_kelamin),
			tempat_lahir = VALUES(tempat_lahir),
			tgl_lahir = VALUES(tgl_lahir),
			alamat_rumah = VALUES(alamat_rumah),
			cabang = VALUES(cabang),
			provinsi = VALUES(provinsi),
			kota_kabupaten = VALUES(kota_kabupaten),
			kota_kabupaten_kantor = VALUES(kota_kabupaten_kantor),
			provinsi_kantor = VALUES(provinsi_kantor),
			status = VALUES(status),
			alumni = VALUES(alumni),
			thn_lulus = VALUES(thn_lulus),
			tempat_tugas = VALUES(tempat_tugas),
			tempat_praktek_1 = VALUES(tempat_praktek_1),
			tempat_praktek_1_tipe = VALUES(tempat_praktek_1_tipe),
			tempat_praktek_1_tipe_2 = VALUES(tempat_praktek_1_tipe_2),
			tempat_praktek_1_alkes = VALUES(tempat_praktek_1_alkes),
			tempat_praktek_1_alkes_2 = VALUES(tempat_praktek_1_alkes_2),
			tempat_praktek_2 = VALUES(tempat_praktek_2),
			tempat_praktek_2_tipe = VALUES(tempat_praktek_2_tipe),
			tempat_praktek_2_tipe_2 = VALUES(tempat_praktek_2_tipe_2),
			tempat_praktek_2_alkes = VALUES(tempat_praktek_2_alkes),
			tempat_praktek_2_alkes_2 = VALUES(tempat_praktek_2_alkes_2),
			kota_kabupaten_praktek_2 = VALUES(kota_kabupaten_praktek_2),
			provinsi_praktek_2 = VALUES(provinsi_praktek_2),
			tempat_praktek_3 = VALUES(tempat_praktek_3),
			tempat_praktek_3_tipe = VALUES(tempat_praktek_3_tipe),
			tempat_praktek_3_tipe_2 = VALUES(tempat_praktek_3_tipe_2),
			tempat_praktek_3_alkes = VALUES(tempat_praktek_3_alkes),
			tempat_praktek_3_alkes_2 = VALUES(tempat_praktek_3_alkes_2),
			kota_kabupaten_praktek_3 = VALUES(kota_kabupaten_praktek_3),
			provinsi_praktek_3 = VALUES(provinsi_praktek_3),
			subspesialis = VALUES(subspesialis),
			gelar_fisr = VALUES(gelar_fisr),
			no_str = VALUES(no_str),
			str_berlaku_sampai = VALUES(str_berlaku_sampai),
			no_sip = VALUES(no_sip),
			sip_berlaku_sampai = VALUES(sip_berlaku_sampai),
			user_id = VALUES(user_id),
			synced_at = VALUES(synced_at),
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := db.NamedExec(query, member)
	return err
}

// LinkMemberToUser links a PDPI member to a user account
func LinkMemberToUser(db *sqlx.DB, memberID string, userID int64) error {
	query := `UPDATE pdpi_members SET user_id = ?, synced_at = NOW() WHERE id = ?`
	_, err := db.Exec(query, userID, memberID)
	return err
}

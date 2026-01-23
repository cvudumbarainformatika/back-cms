package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// PDPIMember represents a member from PDPI system
type PDPIMember struct {
	ID               string         `db:"id" json:"id"`
	NPA              string         `db:"npa" json:"npa"`
	Nama             string         `db:"nama" json:"nama"`
	Gelar            sql.NullString `db:"gelar" json:"gelar,omitempty"`
	Gelar2           sql.NullString `db:"gelar2" json:"gelar2,omitempty"`
	Email            sql.NullString `db:"email" json:"email,omitempty"`
	NoHP             sql.NullString `db:"no_hp" json:"no_hp,omitempty"`
	NIK              sql.NullString `db:"nik" json:"nik,omitempty"`
	JenisKelamin     sql.NullString `db:"jenis_kelamin" json:"jenis_kelamin,omitempty"`
	TempatLahir      sql.NullString `db:"tempat_lahir" json:"tempat_lahir,omitempty"`
	TglLahir         sql.NullTime   `db:"tgl_lahir" json:"tgl_lahir,omitempty"`
	AlamatRumah      sql.NullString `db:"alamat_rumah" json:"alamat_rumah,omitempty"`
	Cabang           sql.NullString `db:"cabang" json:"cabang,omitempty"`
	Provinsi         sql.NullString `db:"provinsi" json:"provinsi,omitempty"`
	KotaKabupaten    sql.NullString `db:"kota_kabupaten" json:"kota_kabupaten,omitempty"`
	Status           sql.NullString `db:"status" json:"status,omitempty"`
	Alumni           sql.NullString `db:"alumni" json:"alumni,omitempty"`
	ThnLulus         sql.NullInt64  `db:"thn_lulus" json:"thn_lulus,omitempty"`
	TempatTugas      sql.NullString `db:"tempat_tugas" json:"tempat_tugas,omitempty"`
	TempatPraktek1   sql.NullString `db:"tempat_praktek_1" json:"tempat_praktek_1,omitempty"`
	TempatPraktek2   sql.NullString `db:"tempat_praktek_2" json:"tempat_praktek_2,omitempty"`
	Subspesialis     sql.NullString `db:"subspesialis" json:"subspesialis,omitempty"`
	NoSTR            sql.NullString `db:"no_str" json:"no_str,omitempty"`
	STRBerlakuSampai sql.NullTime   `db:"str_berlaku_sampai" json:"str_berlaku_sampai,omitempty"`
	NoSIP            sql.NullString `db:"no_sip" json:"no_sip,omitempty"`
	SIPBerlakuSampai sql.NullTime   `db:"sip_berlaku_sampai" json:"sip_berlaku_sampai,omitempty"`
	UserID           sql.NullInt64  `db:"user_id" json:"user_id,omitempty"`
	SyncedAt         sql.NullTime   `db:"synced_at" json:"synced_at,omitempty"`
	CreatedAt        time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at" json:"updated_at"`
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
			id, npa, nama, gelar, gelar2, email, no_hp, nik, 
			jenis_kelamin, tempat_lahir, tgl_lahir, alamat_rumah,
			cabang, provinsi, kota_kabupaten, status, alumni, thn_lulus,
			tempat_tugas, tempat_praktek_1, tempat_praktek_2, subspesialis,
			no_str, str_berlaku_sampai, no_sip, sip_berlaku_sampai,
			user_id, synced_at
		) VALUES (
			:id, :npa, :nama, :gelar, :gelar2, :email, :no_hp, :nik,
			:jenis_kelamin, :tempat_lahir, :tgl_lahir, :alamat_rumah,
			:cabang, :provinsi, :kota_kabupaten, :status, :alumni, :thn_lulus,
			:tempat_tugas, :tempat_praktek_1, :tempat_praktek_2, :subspesialis,
			:no_str, :str_berlaku_sampai, :no_sip, :sip_berlaku_sampai,
			:user_id, :synced_at
		)
		ON DUPLICATE KEY UPDATE
			nama = VALUES(nama),
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
			status = VALUES(status),
			alumni = VALUES(alumni),
			thn_lulus = VALUES(thn_lulus),
			tempat_tugas = VALUES(tempat_tugas),
			tempat_praktek_1 = VALUES(tempat_praktek_1),
			tempat_praktek_2 = VALUES(tempat_praktek_2),
			subspesialis = VALUES(subspesialis),
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

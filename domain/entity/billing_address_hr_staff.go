package entity

import (
	"time"
)

type BillingAddressHRStaff struct {
	ID                 uint      `db:"id" json:"id"`
	BillingAddressID   uint      `db:"billing_address_id" json:"billing_address_id"`
	HRStaffName        string    `db:"hr_staff_name" json:"hr_staff_name"`
	HRStaffEmail       string    `db:"hr_staff_email" json:"hr_staff_email"`
	HRStaffPhoneNumber string    `db:"hr_staff_phone_number" json:"hr_staff_phone_number"`
	CreatedAt          time.Time `db:"created_at" json:"-"`
	UpdatedAt          time.Time `db:"updated_at" json:"-"`
}

func NewBillingAddressHRStaff(
	billingAddressID uint,
	hrStaffName string,
	hrStaffEmail string,
	hrStaffPhoneNumber string,
) *BillingAddressHRStaff {
	return &BillingAddressHRStaff{
		BillingAddressID:   billingAddressID,
		HRStaffName:        hrStaffName,
		HRStaffEmail:       hrStaffEmail,
		HRStaffPhoneNumber: hrStaffPhoneNumber,
	}
}

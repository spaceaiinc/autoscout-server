package entity

import (
	"time"
)

type BillingAddressRAStaff struct {
	ID                             uint      `db:"id" json:"id"`
	BillingAddressID               uint      `db:"billing_address_id" json:"billing_address_id"`
	BillingAddressStaffName        string    `db:"billing_address_staff_name" json:"billing_address_staff_name"`
	BillingAddressStaffEmail       string    `db:"billing_address_staff_email" json:"billing_address_staff_email"`
	BillingAddressStaffPhoneNumber string    `db:"billing_address_staff_phone_number" json:"billing_address_staff_phone_number"`
	CreatedAt                      time.Time `db:"created_at" json:"-"`
	UpdatedAt                      time.Time `db:"updated_at" json:"-"`
}

func NewBillingAddressRAStaff(
	billingAddressID uint,
	billingAddressStaffName string,
	billingAddressStaffEmail string,
	billingAddressStaffPhoneNumber string,
) *BillingAddressRAStaff {
	return &BillingAddressRAStaff{
		BillingAddressID:               billingAddressID,
		BillingAddressStaffName:        billingAddressStaffName,
		BillingAddressStaffEmail:       billingAddressStaffEmail,
		BillingAddressStaffPhoneNumber: billingAddressStaffPhoneNumber,
	}
}

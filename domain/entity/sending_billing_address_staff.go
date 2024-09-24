package entity

import (
	"time"
)

type SendingBillingAddressStaff struct {
	ID                      uint      `db:"id" json:"id"`
	SendingBillingAddressID uint      `db:"sending_billing_address_id" json:"sending_billing_address_id"`
	StaffName               string    `db:"staff_name" json:"staff_name"`
	StaffEmail              string    `db:"staff_email" json:"staff_email"`
	StaffPhoneNumber        string    `db:"staff_phone_number" json:"staff_phone_number"`
	CreatedAt               time.Time `db:"created_at" json:"-"`
	UpdatedAt               time.Time `db:"updated_at" json:"-"`
}

func NewSendingBillingAddressStaff(
	billingAddressID uint,
	staffName string,
	staffEmail string,
	staffPhoneNumber string,
) *SendingBillingAddressStaff {
	return &SendingBillingAddressStaff{
		SendingBillingAddressID: billingAddressID,
		StaffName:               staffName,
		StaffEmail:              staffEmail,
		StaffPhoneNumber:        staffPhoneNumber,
	}
}

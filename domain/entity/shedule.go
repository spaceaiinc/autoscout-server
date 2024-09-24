package entity

import "time"

type Schedule struct {
	ID        uint   `db:"id" json:"id"` // 「面談調整 or タスク」のID
	CAStaffID uint   `db:"agent_staff_id" json:"agent_staff_id"`
	StaffName string `db:"staff_name" json:"staff_name"`
	Title     string `db:"title" json:"title"`
	StartTime         time.Time `db:"start_time" json:"start_time"`
	EndTime           time.Time `db:"end_time" json:"end_time"`
	SeekerDescription string    `db:"seeker_description" json:"seeker_description"`
	StaffDescription  string    `db:"staff_description" json:"staff_description"`
	ScheduleType      string    `db:"schedule_type" json:"schedule_type"`

	// 担当者テーブル
	CAStaffName string `db:"ca_staff_name" json:"ca_staff_name"`
	CAAgentName string `db:"ca_agent_name" json:"ca_agent_name"`
	RAStaffName string `db:"ra_staff_name" json:"ra_staff_name"`
	RAAgentName string `db:"ra_agent_name" json:"ra_agent_name"`
}

func NewSchedule(
	id uint,
	caStaffID uint,
	caStaffName string,
	caAgentName string,
	raStaffNamw string,
	raAgentName string,
	title string,
	StartTime time.Time,
	endTime time.Time,
	seekerDescription string,
	staffDescription string,
	scheduleType string,
) *Schedule {
	return &Schedule{
		ID:                id,
		CAStaffID:         caStaffID,
		CAStaffName:       caStaffName,
		CAAgentName:       caAgentName,
		RAStaffName:       raStaffNamw,
		RAAgentName:       raAgentName,
		Title:             title,
		StartTime:         StartTime,
		EndTime:           endTime,
		SeekerDescription: seekerDescription,
		StaffDescription:  staffDescription,
		ScheduleType:      scheduleType,
	}
}

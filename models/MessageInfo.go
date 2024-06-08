package models

type MessageInfo struct {
	MessageInfoId int    `db:"message_info_id"`
	Message       string `db:"message"`
	SentDate      string `db:"sent_date"`
}

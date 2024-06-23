package models

type Message struct {
	MessageId   int    `db:"message_id"`
	SenderId    int    `db:"sender_id"`
	ReceiverId  int    `db:"receiver_id"`
	MessageBody string `db:"message"`
	SentDate    string `db:"sent_date"`
}

package models

type Message struct {
<<<<<<< HEAD
	MessageId   int    `db:"message_id"`
	SenderId    int    `db:"sender_id"`
	ReceiverId  int    `db:"receiver_id"`
	MessageBody string `db:"message"`
	SentDate    string `db:"sent_date"`
=======
	MessageId     int `db:"message_id"`
	SenderId      int `db:"sender_id"`
	ReceiverId    int `db:"receiver_id"`
	MessageInfoId int `db:"message_info_id"`
>>>>>>> c74ed8c47bbcad1fb2db51e22715763bdb190b65
}

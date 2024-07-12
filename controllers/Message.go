package controllers

import (
	"fmt"
	"gelio/m/initializers"
	"gelio/m/middleware"
	"gelio/m/models"

	"github.com/gin-gonic/gin"
)

type Message struct{}

func MessageController() *Message {
	return &Message{}
}

func (Message) LoadContacts(c *gin.Context) {
	id := c.Param("id")

	var Users []models.User

	err := initializers.DB.Select(&Users,
		`SELECT user_id, username, password, created_date, is_active, person_id, profile_image_id 
      FROM Message
      INNER JOIN users ON users.user_id = Message.sender_id 
      WHERE receiver_id = $1
      UNION
      SELECT users.user_id, username, password, created_date, is_active, person_id, profile_image_id FROM followers 
      INNER JOIN users on users.user_id = followers.user_id WHERE followers.follower_id = $1`, id)

	if err != nil {
		fmt.Println(err)
		c.JSON(400, nil)
		return

	}

	c.JSON(200, Users)
}

func (Message) LoadMessages(c *gin.Context) {
	var body struct {
		SenderId   int
		ReceiverId int
	}

	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		c.JSON(400, nil)
		return
	}

	var Messages []models.Message

	err := initializers.DB.Select(&Messages, `select message_id, sender_id, receiver_id, message, sent_date from Message 
    inner join message_info on message_info.message_info_id = Message.message_info_id
    where receiver_id = $1 and sender_id = $2 or receiver_id = $2 and sender_id = $1
    order by sent_date asc`, body.ReceiverId, body.SenderId)

	if err != nil {
		fmt.Println(err)
		c.JSON(400, nil)
		return
	}

	c.JSON(200, Messages)

}

func (Message) GetMessageInfoFromId(c *gin.Context) {
	id := c.Param("id")

	var MessageInfo models.MessageInfo

	err := initializers.DB.Get(&MessageInfo, "select * from Message_Info where message_info_id = $1", id)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, nil)
		return
	}

	c.JSON(200, MessageInfo)
}

func (Message) SendMessage(c *gin.Context) {
	var body struct {
		SenderId   int
		ReceiverId int
		Message    string
		SentDate   string
	}

	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		c.JSON(400, nil)
		return
	}

	res := initializers.DB.QueryRow("insert into Message_Info (message, sent_date) values ($1, $2) returning message_info_id", body.Message, body.SentDate)

	var MessageInfoId int

	if err := res.Scan(&MessageInfoId); err != nil {
		fmt.Println(err)
		c.JSON(400, nil)
		return
	}

	Res := initializers.DB.QueryRow("insert into Message (sender_id, receiver_id, message_info_id) values ($1, $2, $3) returning message_id", body.SenderId, body.ReceiverId, MessageInfoId)

	var MessageId int

	if err := Res.Scan(&MessageId); err != nil {
		fmt.Println(err)
		c.JSON(400, nil)
		return
	}

	c.JSON(200, MessageId)
}

func (Message) AddContact(c *gin.Context) {
	var body struct {
		Username string
		UserId   int
	}

	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		c.JSON(400, false)
		return
	}

	var User models.User

	err := initializers.DB.Get(&User, "select * from users where username = $1", body.Username)

	if err != nil {
		fmt.Println(err)
		c.JSON(404, false)
		return
	}

	var Follower models.Follower

	err = initializers.DB.Get(&Follower, "select * from followers where user_id = $1 and follower_id = $2", User.UserId, body.UserId)

	if err == nil {
		c.JSON(409, false)
		return
	}

	_, Err := initializers.DB.Exec("insert into followers (user_id, follower_id) values ($1, $2)", User.UserId, body.UserId)

	if Err != nil {
		c.JSON(500, false)
		return
	}

	c.JSON(201, true)

}

func (Message) IsPersonNotContact(c *gin.Context) {
	var body struct {
		Username string
		UserId   int
	}

	if err := c.Bind(&body); err != nil {
		fmt.Println(err)
		c.JSON(400, false)
		return
	}

	var User models.User

	err := initializers.DB.Get(&User, "select * from users where username = $1", body.Username)

	if err != nil {
		fmt.Println(err)
		c.JSON(404, false)
	}

	c.JSON(200, true)

}

func (m *Message) InitializeRoutes(r *gin.Engine) {
	r.GET("/Contacts/:id", middleware.RequireAuth, m.LoadContacts)
	r.POST("/Load/Messages", middleware.RequireAuth, m.LoadMessages)
	r.POST("/Message", middleware.RequireAuth, m.SendMessage)
	r.POST("/Contact", middleware.RequireAuth, m.AddContact)
	r.POST("/Contact/Exists", middleware.RequireAuth, m.IsPersonNotContact)
}

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/opas-keat/line-message-api/models"
)

type LineMessage struct {
	Destination string `json:"destination"`
	Events      []struct {
		ReplyToken string `json:"replyToken"`
		Type       string `json:"type"`
		Timestamp  int64  `json:"timestamp"`
		Source     struct {
			Type   string `json:"type"`
			UserID string `json:"userId"`
		} `json:"source"`
		Message struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"events"`
}

type ReplyMessage struct {
	ReplyToken string `json:"replyToken"`
	Messages   []Text `json:"messages"`
}

type SendMessage struct {
	To       string `json:"to"`
	Messages []Text `json:"messages"`
}

type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ProFile struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

var ChannelToken = "uF3o8ejH1fzia8j2l+AFLYkbRcaNLZ65pF7ZdTQDZ9uJdjeZxQQiUPHRsIUIBtEvywomSJ22kgTvp+e80BMrWkiYUmgaNO64m5i2wmdm8CZsSYE4wlggCSyYiTUK8vbm1Zk/f9TnPtOE6QQ0QodiOQdB04t89/1O/w1cDnyilFU="

func main() {
	// Fiber instance
	app := fiber.New()

	// Middleware
	app.Use(recover.New())

	// Routes
	app.Get("/", hello)

	app.Post("/webhook", webhook)

	app.Get("/sendMessage", sendMessageLine)

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}

// Handler webhook
func webhook(c *fiber.Ctx) error {
	var input models.LineMessage
	log.Println(input)
	if err := c.BodyParser(&input); err != nil {
		log.Println("err")
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "", "data": nil})
	}
	log.Println("UserId = " + input.Events[0].Source.UserID)
	fullname := getProfile(input.Events[0].Source.UserID)
	log.Println(fullname)
	text := Text{
		Type: "text",
		Text: "ข้อความเข้ามา : " + input.Events[0].Message.Text + " ยินดีต้อนรับ : " + fullname,
	}
	log.Println(text)
	message := ReplyMessage{
		ReplyToken: input.Events[0].ReplyToken,
		Messages: []Text{
			text,
		},
	}
	replyMessageLine(message)
	// log.Println(fullname)
	return c.SendString("Hi👋!, This Line Message Api")
}

func replyMessageLine(Message ReplyMessage) error {
	value, _ := json.Marshal(Message)

	url := "https://api.line.me/v2/bot/message/reply"

	var jsonStr = []byte(value)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+ChannelToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

	return err
}

func sendMessageLine(c *fiber.Ctx) error {
	text := Text{
		Type: "text",
		Text: "ข้อความ สวัสดี",
	}
	SendMessage := SendMessage{
		To: "U43c8656500da5a08f4f4a27774567144",
		Messages: []Text{
			text,
		},
	}
	value, _ := json.Marshal(SendMessage)
	url := "https://api.line.me/v2/bot/message/push"
	var jsonStr = []byte(value)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+ChannelToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
	return err
}

// Handler hello
func hello(c *fiber.Ctx) error {
	return c.SendString("Hi👋!, This Line Message Api")
}

func getProfile(userId string) string {

	url := "https://api.line.me/v2/bot/profile/" + userId

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+ChannelToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var profile ProFile
	if err := json.Unmarshal(body, &profile); err != nil {
		log.Println("%% err \n")
	}
	log.Println(profile.DisplayName)
	return profile.DisplayName

}

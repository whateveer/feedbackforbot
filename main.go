package main

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Feedback represents the structure of a feedback document
type Feedback struct {
	ChatID   int64  `bson:"chat_id"`
	Feedback string `bson:"feedback"`
}

func main() {
	botToken := "7139151928:AAEHtCsFYKQn9EZVuHp2AJpHpn_GwB3hy00"
	mongoURI := "mongodb+srv://sml:hwO05jIkWRN29oRi@ecoflow.ydwcdbg.mongodb.net/?retryWrites=true&w=majority&appName=Ecoflow"

	// Initialize Telegram Bot
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Set up MongoDB client
	clientOptions := options.Client().ApplyURI(mongoURI)
	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	feedbackCollection := mongoClient.Database("ecoflow").Collection("feedback")

	// Start receiving updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			chatID := update.Message.Chat.ID
			text := update.Message.Text

			if text == "/start" {
				msg := tgbotapi.NewMessage(chatID, "Оставьте свой отзыв🤍")
				bot.Send(msg)
			} else {
				feedback := Feedback{
					ChatID:   chatID,
					Feedback: text,
				}

				// Insert feedback into MongoDB
				insertResult, err := feedbackCollection.InsertOne(context.TODO(), feedback)
				if err != nil {
					log.Printf("Could not insert feedback: %v", err)
					continue
				}

				log.Printf("Inserted feedback: %v", insertResult.InsertedID)

				msg := tgbotapi.NewMessage(chatID, "Thank you for your feedback!")
				bot.Send(msg)
			}
		}
	}
}

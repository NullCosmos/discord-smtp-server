package main

import (
	"log"
	"os"
	"time"

	gosmtp "github.com/emersion/go-smtp"
	"github.com/nullcosmos/discord-smtp-server/smtp"
	"github.com/nullcosmos/discord-smtp-server/discord"
)

func main() {
	discord, err := discord.NewSession(os.Getenv("DISCORD_WEBHOOK_URI"))
	if err != nil {
		log.Fatal(err)
	}

	backend, err := smtp.NewBackend(
		discord,
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)
	if err != nil {
		log.Fatal(err)
	}

	server := gosmtp.NewServer(backend)

	port := ":1025"
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}

	host := "localhost"
	if os.Getenv("HOST") != "" {
		host = os.Getenv("HOST")
	}

	server.Addr = port
	server.Domain = host
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxMessageBytes = 1024 * 1024
	server.MaxRecipients = 50
	server.AllowInsecureAuth = true

	log.Println("Starting server at", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

package domain

type CommandHandler func(*Message) *Message

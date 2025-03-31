package handler

import (
	"encoding/json"
	"fmt"
	"go-ws/internal/lib/logger/sl"
	"go-ws/internal/storage/repository"
	"log/slog"

	"github.com/gorilla/websocket"
)

type MessageHandler struct {
	log                 *slog.Logger
	superheroRepository *repository.SuperheroRepository
}
type messageType struct {
	Event   string `json:"event"`
	Key     string `json:"key"`
	Payload any    `json:"payload"`
}

func New(log *slog.Logger, superheroRepository *repository.SuperheroRepository) *MessageHandler {
	return &MessageHandler{
		log:                 log,
		superheroRepository: superheroRepository,
	}
}

func (m *MessageHandler) Handle(message []byte, conn *websocket.Conn) {
	const op = "handler.Handle"
	var msg messageType
	if err := json.Unmarshal(message, &msg); err != nil {
		err := fmt.Errorf("%s: failed to unmarshal message, is message valid JSON?", op)
		m.catchError(conn, err, msg.Key)
		return
	}

	var result []byte
	var err error

	switch msg.Event {
	case "get-by-id":
		result, err = m.getSuperheroesById(msg)
	case "get-by-id-details":
		result, err = m.getSuperheroByIdWithDetails(msg)
	default:
		err = fmt.Errorf("unknown event: %s", msg.Event)
	}

	if err != nil {
		m.catchError(conn, err, msg.Key)
		return
	}
	m.sendResult(conn, result)
}

func (m *MessageHandler) getSuperheroesById(msg messageType) ([]byte, error) {
	const op = "handler.getSuperheroesById"

	id, err := m.extractIdFromMessage(msg)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to extract ID from message: %w", op, err)
	}

	superhero, err := m.superheroRepository.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to find superhero by ID: %w", op, err)
	}

	data := map[string]any{
		"key":     msg.Key,
		"message": map[string]any{"superhero": superhero},
		"error":   nil,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to marshal JSON: %w", op, err)
	}
	return jsonData, nil
}

func (m *MessageHandler) getSuperheroByIdWithDetails(msg messageType) ([]byte, error) {
	const op = "handler.getSuperheroByIdWithDetails"

	id, err := m.extractIdFromMessage(msg)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to extract ID from message: %w", op, err)
	}

	superhero, err := m.superheroRepository.FindByIdWithDetailed(id)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to find superhero with details by ID: %w", op, err)
	}

	data := map[string]any{
		"key":     msg.Key,
		"message": map[string]any{"superhero": superhero},
		"error":   nil,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to marshal JSON: %w", op, err)
	}
	return jsonData, nil
}

func (m *MessageHandler) sendResult(conn *websocket.Conn, result []byte) {
	const op = "handler.sendResult"
	if err := conn.WriteMessage(websocket.TextMessage, result); err != nil {
		m.log.With(slog.String("op", op)).Error("Failed to write message", sl.Err(err))
	}
}

func (m *MessageHandler) catchError(conn *websocket.Conn, err error, key string) {
	const op = "handler.catchError"
	data := map[string]any{
		"key":     key,
		"message": nil,
		"error":   err.Error(),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		m.log.With(slog.String("op", op)).Error("Failed to marshal JSON", sl.Err(err))
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		m.log.With(slog.String("op", op)).Error("Failed to write message", sl.Err(err))
	}
}

func (m *MessageHandler) extractIdFromMessage(msg messageType) (int, error) {
	const op = "handler.extractIdFromMessage"
	payload, ok := msg.Payload.(map[string]any)
	if !ok {
		return 0, fmt.Errorf("%s: invalid payload format", op)
	}
	idInterface, ok := payload["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("%s: id is missing or not a number", op)
	}
	id := int(idInterface)
	return id, nil
}

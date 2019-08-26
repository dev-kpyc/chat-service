package repository

import (
	"context"
	"database/sql"

	"github.com/dev-kpyc/chat-service/pkg/chatroom"
	"github.com/dev-kpyc/chat-service/pkg/messaging"
	"github.com/dev-kpyc/chat-service/pkg/user"
)

type Postgres struct {
	conn *sql.DB
}

type User struct {
	id   int64
	name string
}

type ChatRoom struct {
	id   int64
	name string
}

type ChatMessage struct {
	id         int64
	chatRoomID int64
	senderID   int64
	message    string
}

func NewPostgresRepository(conn *sql.DB) *Postgres {
	return &Postgres{conn}
}

var _ chatroom.Repository = (*Postgres)(nil)
var _ messaging.Repository = (*Postgres)(nil)
var _ user.Repository = (*Postgres)(nil)

func (postgres *Postgres) Initialize() {

	postgres.exec(CreateUsers)
	postgres.exec(CreateChatRooms)
	postgres.exec(CreateChatMessages)
	postgres.exec(CreateChatRoomUsers)
}

func (postgres *Postgres) exec(query string) {
	_, err := postgres.conn.Exec(query)
	if err != nil {
		panic(err)
	}
}

func (postgres *Postgres) StoreUser(ctx context.Context, name string) (int64, error) {
	query := `INSERT INTO users (name) VALUES ($1) RETURNING id`

	row := postgres.conn.QueryRowContext(ctx, query, name)

	var id int64
	err := row.Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (postgres *Postgres) StoreChatRoom(ctx context.Context, name string, ownerID int64) (int64, error) {
	query := `INSERT INTO chatroom (name, owner_id) VALUES ($1, $2) RETURNING id`

	row := postgres.conn.QueryRowContext(ctx, query, name, ownerID)

	var id int64
	err := row.Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (postgres *Postgres) SaveChatMessage(ctx context.Context, senderID int64, roomID int64, message string) (int64, error) {
	query := `INSERT INTO chatmessage (chatroom_id, sender_id, message) VALUES ($1, $2, $3) RETURNING id`

	row := postgres.conn.QueryRowContext(ctx, query, roomID, senderID, message)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (postgres *Postgres) GetUser(ctx context.Context, userID int64) (*user.User, error) {
	query := `SELECT id,name FROM users WHERE id = $1`

	row := postgres.conn.QueryRowContext(ctx, query, userID)

	var user user.User
	if err := row.Scan(&user.ID, &user.Name); err != nil {
		return nil, err
	}

	return &user, nil
}

func (postgres *Postgres) GetChatRoom(ctx context.Context, ID int64) (*chatroom.ChatRoom, error) {
	query := `SELECT id,name,owner_id FROM chatroom WHERE id = $1`

	row := postgres.conn.QueryRowContext(ctx, query, ID)

	var chatroom chatroom.ChatRoom
	if err := row.Scan(&chatroom.ID, &chatroom.Name); err != nil {
		return nil, err
	}

	return &chatroom, nil
}

func (postgres *Postgres) ListChatRooms(ctx context.Context, userID int64) ([]*chatroom.ChatRoom, error) {
	query := `SELECT id, name, owner_id FROM chatroom WHERE id IN (SELECT chatroom_id FROM chatroomusers WHERE user_id = $1)`

	rows, err := postgres.conn.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	var chatrooms []*chatroom.ChatRoom
	for rows.Next() {
		chatroom, err := scanChatRoom(rows)
		if err != nil {
			return nil, err
		}
		chatrooms = append(chatrooms, chatroom)
	}

	return chatrooms, nil
}

func (postgres *Postgres) ListOwnedChatRooms(ctx context.Context, ownerID int64) ([]*chatroom.ChatRoom, error) {
	query := `SELECT id, name, owner_id FROM chatroom WHERE owner_id = $1`

	rows, err := postgres.conn.QueryContext(ctx, query, ownerID)

	if err != nil {
		return nil, err
	}

	var chatrooms []*chatroom.ChatRoom
	for rows.Next() {
		chatroom, err := scanChatRoom(rows)
		if err != nil {
			return nil, err
		}
		chatrooms = append(chatrooms, chatroom)
	}

	return chatrooms, nil
}

func (postgres *Postgres) GetChatMessages(ctx context.Context, roomID int64) ([]*messaging.ChatMessage, error) {
	query := `SELECT sender_id, message FROM chatmessage WHERE chatroom_id = $1`

	rows, err := postgres.conn.QueryContext(ctx, query, roomID)

	if err != nil {
		return nil, err
	}

	var chatmessages []*messaging.ChatMessage
	for rows.Next() {
		chatmessage, err := scanChatMessage(rows)
		if err != nil {
			return nil, err
		}
		chatmessages = append(chatmessages, chatmessage)
	}

	return chatmessages, nil
}

// AddUserToChatRoom ...
func (postgres *Postgres) AddUserToChatRoom(ctx context.Context, userID int64, roomID int64) error {

	query := `INSERT into chatroomusers (user_id, chatroom_id) VALUES ($1, $2)`

	_, err := postgres.conn.Exec(query, userID, roomID)

	if err != nil {
		return err
	}

	return nil
}

func (postgres *Postgres) RemoveUserFromChatRoom(ctx context.Context, userID int64, roomID int64) error {

	query := `DELETE from chatroomusers WHERE user_id=$1, chatroom_id=$2`

	_, err := postgres.conn.Exec(query, userID, scanChatRoomUser, roomID)

	if err != nil {
		return err
	}

	return nil
}

func (postgres *Postgres) GetUserIDsInChatRoom(ctx context.Context, roomID int64) ([]int64, error) {

	query := `SELECT user_id from chatroomusers WHERE chatroom_id=$1`

	rows, err := postgres.conn.QueryContext(ctx, query, roomID)

	if err != nil {
		return nil, err
	}

	var userIDs []int64
	for rows.Next() {
		userID, err := scanChatRoomUser(rows)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

type scanFn = func(rows *sql.Rows) (interface{}, error)

func scanUser(rows *sql.Rows) (interface{}, error) {
	user := user.User{}
	err := rows.Scan(
		&user.ID,
		&user.Name,
	)
	return &user, err
}

func scanChatRoom(rows *sql.Rows) (*chatroom.ChatRoom, error) {
	chatroom := chatroom.ChatRoom{}
	err := rows.Scan(
		&chatroom.ID,
		&chatroom.Name,
		&chatroom.OwnerID,
	)
	return &chatroom, err
}

func scanChatMessage(rows *sql.Rows) (*messaging.ChatMessage, error) {
	chatmessage := messaging.ChatMessage{}
	err := rows.Scan(
		&chatmessage.SenderID,
		&chatmessage.Message,
	)
	return &chatmessage, err
}

func scanChatRoomUser(rows *sql.Rows) (int64, error) {
	var userID int64
	err := rows.Scan(&userID)
	return userID, err
}

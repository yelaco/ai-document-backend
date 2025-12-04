package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// BaseID is the base type for all ID types in the application.
// It wraps a UUID and provides common functionality for ID handling.
type BaseID struct {
	id uuid.UUID
}

// NewBaseID creates a new BaseID from a UUID.
func NewBaseID(id uuid.UUID) BaseID {
	return BaseID{id: id}
}

// NewBaseIDFromString creates a new BaseID from a string.
// It handles both plain UUID strings and JSON array formats like ["uuid"].
func NewBaseIDFromString(s string) (BaseID, error) {
	if s == "" {
		return BaseID{}, fmt.Errorf("empty string provided")
	}

	// Check if the string looks like a JSON array
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		// Parse as JSON array and take the first element
		var arr []string
		if err := json.Unmarshal([]byte(s), &arr); err != nil {
			return BaseID{}, fmt.Errorf("invalid JSON array format: %w", err)
		}
		if len(arr) == 0 {
			return BaseID{}, fmt.Errorf("empty array provided")
		}
		s = arr[0]
	}

	// Remove dashes and add them back in the proper positions if needed
	if len(s) == 32 && !strings.Contains(s, "-") {
		// UUID without dashes - add them
		s = fmt.Sprintf("%s-%s-%s-%s-%s",
			s[0:8], s[8:12], s[12:16], s[16:20], s[20:32])
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return BaseID{}, fmt.Errorf("invalid UUID format: %w", err)
	}

	return BaseID{id: id}, nil
}

// GenerateBaseID generates a new random BaseID.
func GenerateBaseID() BaseID {
	return BaseID{id: uuid.New()}
}

// UUID returns the underlying UUID.
func (b BaseID) UUID() uuid.UUID {
	return b.id
}

// String returns the string representation of the UUID.
func (b BaseID) String() string {
	return b.id.String()
}

// IsZero returns true if the ID is zero/nil.
func (b BaseID) IsZero() bool {
	return b.id == uuid.Nil
}

// MarshalText implements encoding.TextMarshaler.
func (b BaseID) MarshalText() ([]byte, error) {
	return []byte(b.id.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *BaseID) UnmarshalText(text []byte) error {
	id, err := NewBaseIDFromString(string(text))
	if err != nil {
		return err
	}
	*b = id
	return nil
}

// MarshalJSON implements json.Marshaler.
func (b BaseID) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.id.String())
}

// UnmarshalJSON implements json.Unmarshaler.
// It handles both string and array formats, and also plain text for Gin compatibility.
func (b *BaseID) UnmarshalJSON(data []byte) error {
	// Handle the case where data is not JSON at all (plain text from Gin form binding)
	dataStr := strings.TrimSpace(string(data))
	if !strings.HasPrefix(dataStr, "\"") && !strings.HasPrefix(dataStr, "[") {
		// This is plain text, not JSON - handle it directly
		id, err := NewBaseIDFromString(dataStr)
		if err != nil {
			return err
		}
		*b = id
		return nil
	}

	var s string

	// First try to unmarshal as a string
	if err := json.Unmarshal(data, &s); err == nil {
		id, err := NewBaseIDFromString(s)
		if err != nil {
			return err
		}
		*b = id
		return nil
	}

	// If that fails, try to unmarshal as an array and take the first element
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		// If both JSON unmarshaling attempts fail, try treating as plain text
		id, err := NewBaseIDFromString(dataStr)
		if err != nil {
			return fmt.Errorf("ID must be either a string, array of strings, or valid UUID text: %w", err)
		}
		*b = id
		return nil
	}

	if len(arr) == 0 {
		return fmt.Errorf("empty array provided for ID")
	}

	id, err := NewBaseIDFromString(arr[0])
	if err != nil {
		return err
	}
	*b = id
	return nil
}

// Value implements driver.Valuer for database/sql.
func (b BaseID) Value() (driver.Value, error) {
	return b.id, nil
}

// Scan implements sql.Scanner for database/sql.
func (b *BaseID) Scan(value interface{}) error {
	if value == nil {
		b.id = uuid.Nil
		return nil
	}

	switch v := value.(type) {
	case string:
		id, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		b.id = id
	case []byte:
		id, err := uuid.Parse(string(v))
		if err != nil {
			return err
		}
		b.id = id
	case uuid.UUID:
		b.id = v
	default:
		return fmt.Errorf("cannot scan %T into BaseID", value)
	}

	return nil
}

// UserID represents a user identifier.
type UserID struct {
	BaseID
}

// NewUserID creates a new UserID from a UUID.
func NewUserID(id uuid.UUID) UserID {
	return UserID{BaseID: NewBaseID(id)}
}

// NewUserIDFromString creates a new UserID from a string.
func NewUserIDFromString(s string) (UserID, error) {
	baseID, err := NewBaseIDFromString(s)
	if err != nil {
		return UserID{}, err
	}
	return UserID{BaseID: baseID}, nil
}

// GenerateUserID generates a new random UserID.
func GenerateUserID() UserID {
	return UserID{BaseID: GenerateBaseID()}
}

// MarshalText implements encoding.TextMarshaler.
func (u UserID) MarshalText() ([]byte, error) {
	return u.BaseID.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *UserID) UnmarshalText(text []byte) error {
	return u.BaseID.UnmarshalText(text)
}

// MarshalJSON implements json.Marshaler.
func (u UserID) MarshalJSON() ([]byte, error) {
	return u.BaseID.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *UserID) UnmarshalJSON(data []byte) error {
	return u.BaseID.UnmarshalJSON(data)
}

// DocumentID represents a document identifier.
type DocumentID struct {
	BaseID
}

// NewDocumentID creates a new DocumentID from a UUID.
func NewDocumentID(id uuid.UUID) DocumentID {
	return DocumentID{BaseID: NewBaseID(id)}
}

// NewDocumentIDFromString creates a new DocumentID from a string.
func NewDocumentIDFromString(s string) (DocumentID, error) {
	baseID, err := NewBaseIDFromString(s)
	if err != nil {
		return DocumentID{}, err
	}
	return DocumentID{BaseID: baseID}, nil
}

// GenerateDocumentID generates a new random DocumentID.
func GenerateDocumentID() DocumentID {
	return DocumentID{BaseID: GenerateBaseID()}
}

// MarshalText implements encoding.TextMarshaler.
func (d DocumentID) MarshalText() ([]byte, error) {
	return d.BaseID.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (d *DocumentID) UnmarshalText(text []byte) error {
	return d.BaseID.UnmarshalText(text)
}

// MarshalJSON implements json.Marshaler.
func (d DocumentID) MarshalJSON() ([]byte, error) {
	return d.BaseID.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DocumentID) UnmarshalJSON(data []byte) error {
	return d.BaseID.UnmarshalJSON(data)
}

// ChatID represents a chat identifier.
type ChatID struct {
	BaseID
}

// NewChatID creates a new ChatID from a UUID.
func NewChatID(id uuid.UUID) ChatID {
	return ChatID{BaseID: NewBaseID(id)}
}

// NewChatIDFromString creates a new ChatID from a string.
func NewChatIDFromString(s string) (ChatID, error) {
	baseID, err := NewBaseIDFromString(s)
	if err != nil {
		return ChatID{}, err
	}
	return ChatID{BaseID: baseID}, nil
}

// GenerateChatID generates a new random ChatID.
func GenerateChatID() ChatID {
	return ChatID{BaseID: GenerateBaseID()}
}

// MarshalText implements encoding.TextMarshaler.
func (c ChatID) MarshalText() ([]byte, error) {
	return c.BaseID.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (c *ChatID) UnmarshalText(text []byte) error {
	return c.BaseID.UnmarshalText(text)
}

// MarshalJSON implements json.Marshaler.
func (c ChatID) MarshalJSON() ([]byte, error) {
	return c.BaseID.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *ChatID) UnmarshalJSON(data []byte) error {
	return c.BaseID.UnmarshalJSON(data)
}

// MessageID represents a message identifier.
type MessageID struct {
	BaseID
}

// NewMessageID creates a new MessageID from a UUID.
func NewMessageID(id uuid.UUID) MessageID {
	return MessageID{BaseID: NewBaseID(id)}
}

// NewMessageIDFromString creates a new MessageID from a string.
func NewMessageIDFromString(s string) (MessageID, error) {
	baseID, err := NewBaseIDFromString(s)
	if err != nil {
		return MessageID{}, err
	}
	return MessageID{BaseID: baseID}, nil
}

// GenerateMessageID generates a new random MessageID.
func GenerateMessageID() MessageID {
	return MessageID{BaseID: GenerateBaseID()}
}

// MarshalText implements encoding.TextMarshaler.
func (m MessageID) MarshalText() ([]byte, error) {
	return m.BaseID.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (m *MessageID) UnmarshalText(text []byte) error {
	return m.BaseID.UnmarshalText(text)
}

// MarshalJSON implements json.Marshaler.
func (m MessageID) MarshalJSON() ([]byte, error) {
	return m.BaseID.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
func (m *MessageID) UnmarshalJSON(data []byte) error {
	return m.BaseID.UnmarshalJSON(data)
}

// RefreshTokenID represents a refresh token identifier.
type RefreshTokenID struct {
	BaseID
}

// NewRefreshTokenID creates a new RefreshTokenID from a UUID.
func NewRefreshTokenID(id uuid.UUID) RefreshTokenID {
	return RefreshTokenID{BaseID: NewBaseID(id)}
}

// NewRefreshTokenIDFromString creates a new RefreshTokenID from a string.
func NewRefreshTokenIDFromString(s string) (RefreshTokenID, error) {
	baseID, err := NewBaseIDFromString(s)
	if err != nil {
		return RefreshTokenID{}, err
	}
	return RefreshTokenID{BaseID: baseID}, nil
}

// GenerateRefreshTokenID generates a new random RefreshTokenID.
func GenerateRefreshTokenID() RefreshTokenID {
	return RefreshTokenID{BaseID: GenerateBaseID()}
}

// MarshalText implements encoding.TextMarshaler.
func (r RefreshTokenID) MarshalText() ([]byte, error) {
	return r.BaseID.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (r *RefreshTokenID) UnmarshalText(text []byte) error {
	return r.BaseID.UnmarshalText(text)
}

// MarshalJSON implements json.Marshaler.
func (r RefreshTokenID) MarshalJSON() ([]byte, error) {
	return r.BaseID.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *RefreshTokenID) UnmarshalJSON(data []byte) error {
	return r.BaseID.UnmarshalJSON(data)
}

package types

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaseID_NewBaseIDFromString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    string
	}{
		{
			name:        "valid UUID string",
			input:       "123e4567-e89b-12d3-a456-426614174000",
			expectError: false,
			expected:    "123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name:        "UUID in array format",
			input:       `["123e4567-e89b-12d3-a456-426614174000"]`,
			expectError: false,
			expected:    "123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name:        "multiple UUIDs in array - takes first",
			input:       `["123e4567-e89b-12d3-a456-426614174000","987fcdeb-51a2-43d6-b789-123456789abc"]`,
			expectError: false,
			expected:    "123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name:        "UUID without dashes",
			input:       "123e4567e89b12d3a456426614174000",
			expectError: false,
			expected:    "123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name:        "invalid UUID format",
			input:       "invalid-uuid",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "empty array",
			input:       "[]",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewBaseIDFromString(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result.String())
			}
		})
	}
}

func TestBaseID_UUID(t *testing.T) {
	testUUID := uuid.New()
	baseID := NewBaseID(testUUID)

	assert.Equal(t, testUUID, baseID.UUID())
}

func TestBaseID_String(t *testing.T) {
	testUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	baseID := NewBaseID(testUUID)

	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", baseID.String())
}

func TestBaseID_IsZero(t *testing.T) {
	// Zero value
	var zeroID BaseID
	assert.True(t, zeroID.IsZero())

	// Nil UUID
	nilID := NewBaseID(uuid.Nil)
	assert.True(t, nilID.IsZero())

	// Non-zero UUID
	nonZeroID := GenerateBaseID()
	assert.False(t, nonZeroID.IsZero())
}

func TestBaseID_MarshalUnmarshalText(t *testing.T) {
	testUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	baseID := NewBaseID(testUUID)

	// Test Marshal
	text, err := baseID.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", string(text))

	// Test Unmarshal
	var newID BaseID
	err = newID.UnmarshalText(text)
	require.NoError(t, err)
	assert.Equal(t, baseID.UUID(), newID.UUID())
}

func TestBaseID_ValueScan(t *testing.T) {
	testUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	baseID := NewBaseID(testUUID)

	// Test Value (for database write)
	value, err := baseID.Value()
	require.NoError(t, err)
	assert.Equal(t, testUUID, value)

	// Test Scan from string
	var id1 BaseID
	err = id1.Scan("123e4567-e89b-12d3-a456-426614174000")
	require.NoError(t, err)
	assert.Equal(t, testUUID, id1.UUID())

	// Test Scan from []byte
	var id2 BaseID
	err = id2.Scan([]byte("123e4567-e89b-12d3-a456-426614174000"))
	require.NoError(t, err)
	assert.Equal(t, testUUID, id2.UUID())

	// Test Scan from UUID
	var id3 BaseID
	err = id3.Scan(testUUID)
	require.NoError(t, err)
	assert.Equal(t, testUUID, id3.UUID())

	// Test Scan from nil
	var id4 BaseID
	err = id4.Scan(nil)
	require.NoError(t, err)
	assert.Equal(t, uuid.Nil, id4.UUID())

	// Test Scan from invalid type
	var id5 BaseID
	err = id5.Scan(123)
	assert.Error(t, err)
}

func TestUserID(t *testing.T) {
	testUUID := uuid.New()

	// Test NewUserID
	userID := NewUserID(testUUID)
	assert.Equal(t, testUUID, userID.UUID())

	// Test NewUserIDFromString
	userID2, err := NewUserIDFromString(testUUID.String())
	require.NoError(t, err)
	assert.Equal(t, testUUID, userID2.UUID())

	// Test GenerateUserID
	userID3 := GenerateUserID()
	assert.False(t, userID3.IsZero())
}

func TestDocumentID(t *testing.T) {
	testUUID := uuid.New()

	// Test NewDocumentID
	docID := NewDocumentID(testUUID)
	assert.Equal(t, testUUID, docID.UUID())

	// Test NewDocumentIDFromString
	docID2, err := NewDocumentIDFromString(testUUID.String())
	require.NoError(t, err)
	assert.Equal(t, testUUID, docID2.UUID())

	// Test GenerateDocumentID
	docID3 := GenerateDocumentID()
	assert.False(t, docID3.IsZero())
}

func TestChatID(t *testing.T) {
	testUUID := uuid.New()

	// Test NewChatID
	chatID := NewChatID(testUUID)
	assert.Equal(t, testUUID, chatID.UUID())

	// Test NewChatIDFromString
	chatID2, err := NewChatIDFromString(testUUID.String())
	require.NoError(t, err)
	assert.Equal(t, testUUID, chatID2.UUID())

	// Test GenerateChatID
	chatID3 := GenerateChatID()
	assert.False(t, chatID3.IsZero())
}

func TestChatID_ArrayFormat(t *testing.T) {
	// Test the specific error case from the user's issue
	arrayUUID := `["91397617-81f8-4bbb-8104-e037789456f0"]`
	expectedUUID := uuid.MustParse("91397617-81f8-4bbb-8104-e037789456f0")

	chatID, err := NewChatIDFromString(arrayUUID)
	require.NoError(t, err)
	assert.Equal(t, expectedUUID, chatID.UUID())
}

func TestMessageID(t *testing.T) {
	testUUID := uuid.New()

	// Test NewMessageID
	msgID := NewMessageID(testUUID)
	assert.Equal(t, testUUID, msgID.UUID())

	// Test NewMessageIDFromString
	msgID2, err := NewMessageIDFromString(testUUID.String())
	require.NoError(t, err)
	assert.Equal(t, testUUID, msgID2.UUID())

	// Test GenerateMessageID
	msgID3 := GenerateMessageID()
	assert.False(t, msgID3.IsZero())
}

func TestRefreshTokenID(t *testing.T) {
	testUUID := uuid.New()

	// Test NewRefreshTokenID
	tokenID := NewRefreshTokenID(testUUID)
	assert.Equal(t, testUUID, tokenID.UUID())

	// Test NewRefreshTokenIDFromString
	tokenID2, err := NewRefreshTokenIDFromString(testUUID.String())
	require.NoError(t, err)
	assert.Equal(t, testUUID, tokenID2.UUID())

	// Test GenerateRefreshTokenID
	tokenID3 := GenerateRefreshTokenID()
	assert.False(t, tokenID3.IsZero())
}

func TestTypeSafety(t *testing.T) {
	userUUID := uuid.New()
	docUUID := uuid.New()
	chatUUID := uuid.New()

	userID := NewUserID(userUUID)
	docID := NewDocumentID(docUUID)
	chatID := NewChatID(chatUUID)

	// These should be different types even if UUIDs are the same
	sameUUID := uuid.New()
	userID1 := NewUserID(sameUUID)
	docID1 := NewDocumentID(sameUUID)

	// Verify they have the same underlying UUID but are different types
	assert.Equal(t, sameUUID, userID1.UUID())
	assert.Equal(t, sameUUID, docID1.UUID())

	// Verify the types are different (this is compile-time safety in real usage)
	assert.IsType(t, UserID{}, userID)
	assert.IsType(t, DocumentID{}, docID)
	assert.IsType(t, ChatID{}, chatID)
}

func TestDatabaseIntegration(t *testing.T) {
	// Test that IDs can be used as driver.Valuer and sql.Scanner
	userID := GenerateUserID()

	// Test Value method (for database writes)
	value, err := userID.Value()
	require.NoError(t, err)
	assert.IsType(t, uuid.UUID{}, value)

	// Test Scan method (for database reads)
	var scannedID UserID
	err = scannedID.Scan(value)
	require.NoError(t, err)
	assert.Equal(t, userID.UUID(), scannedID.UUID())
}

// Benchmark tests to ensure performance is reasonable
func BenchmarkGenerateUserID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GenerateUserID()
	}
}

func BenchmarkNewUserIDFromString(b *testing.B) {
	testString := "123e4567-e89b-12d3-a456-426614174000"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewUserIDFromString(testString)
	}
}

func BenchmarkNewUserIDFromArrayString(b *testing.B) {
	testString := `["123e4567-e89b-12d3-a456-426614174000"]`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewUserIDFromString(testString)
	}
}

func TestBaseID_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    string
	}{
		{
			name:        "valid UUID string in JSON",
			input:       `"123e4567-e89b-12d3-a456-426614174000"`,
			expectError: false,
			expected:    "123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name:        "UUID in JSON array format",
			input:       `["123e4567-e89b-12d3-a456-426614174000"]`,
			expectError: false,
			expected:    "123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name:        "multiple UUIDs in JSON array - takes first",
			input:       `["123e4567-e89b-12d3-a456-426614174000","987fcdeb-51a2-43d6-b789-123456789abc"]`,
			expectError: false,
			expected:    "123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name:        "UUID without dashes in JSON",
			input:       `"123e4567e89b12d3a456426614174000"`,
			expectError: false,
			expected:    "123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name:        "invalid UUID format in JSON",
			input:       `"invalid-uuid"`,
			expectError: true,
		},
		{
			name:        "empty string in JSON",
			input:       `""`,
			expectError: true,
		},
		{
			name:        "empty array in JSON",
			input:       `[]`,
			expectError: true,
		},
		{
			name:        "invalid JSON format",
			input:       `invalid-json`,
			expectError: true,
		},
		{
			name:        "number instead of string",
			input:       `123`,
			expectError: true,
		},
		{
			name:        "object instead of string",
			input:       `{"id": "123e4567-e89b-12d3-a456-426614174000"}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id BaseID
			err := json.Unmarshal([]byte(tt.input), &id)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, id.String())
			}
		})
	}
}

func TestBaseID_MarshalJSON(t *testing.T) {
	testUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	baseID := NewBaseID(testUUID)

	// Test Marshal
	data, err := json.Marshal(baseID)
	require.NoError(t, err)
	assert.Equal(t, `"123e4567-e89b-12d3-a456-426614174000"`, string(data))

	// Test round-trip
	var newID BaseID
	err = json.Unmarshal(data, &newID)
	require.NoError(t, err)
	assert.Equal(t, baseID.UUID(), newID.UUID())
}

func TestChatID_UnmarshalJSON(t *testing.T) {
	// Test the specific error case from the user's issue
	arrayUUID := `["91397617-81f8-4bbb-8104-e037789456f0"]`
	expectedUUID := uuid.MustParse("91397617-81f8-4bbb-8104-e037789456f0")

	var chatID ChatID
	err := json.Unmarshal([]byte(arrayUUID), &chatID)
	require.NoError(t, err)
	assert.Equal(t, expectedUUID, chatID.UUID())
}

func TestUserID_UnmarshalJSON(t *testing.T) {
	testUUID := uuid.New()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "string format",
			input: `"` + testUUID.String() + `"`,
		},
		{
			name:  "array format",
			input: `["` + testUUID.String() + `"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var userID UserID
			err := json.Unmarshal([]byte(tt.input), &userID)
			require.NoError(t, err)
			assert.Equal(t, testUUID, userID.UUID())
		})
	}
}

func TestDocumentID_UnmarshalJSON(t *testing.T) {
	testUUID := uuid.New()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "string format",
			input: `"` + testUUID.String() + `"`,
		},
		{
			name:  "array format",
			input: `["` + testUUID.String() + `"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var docID DocumentID
			err := json.Unmarshal([]byte(tt.input), &docID)
			require.NoError(t, err)
			assert.Equal(t, testUUID, docID.UUID())
		})
	}
}

func TestGinBindingCompatibility(t *testing.T) {
	// This test simulates the Gin binding scenarios that caused the original issue
	testCases := []struct {
		name     string
		jsonData string
		expected string
	}{
		{
			name:     "form parameter as string",
			jsonData: `{"chat_id": "91397617-81f8-4bbb-8104-e037789456f0"}`,
			expected: "91397617-81f8-4bbb-8104-e037789456f0",
		},
		{
			name:     "form parameter as array (problematic case)",
			jsonData: `{"chat_id": ["91397617-81f8-4bbb-8104-e037789456f0"]}`,
			expected: "91397617-81f8-4bbb-8104-e037789456f0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate a struct that might be used in Gin binding
			type TestRequest struct {
				ChatID ChatID `json:"chat_id"`
			}

			var req TestRequest
			err := json.Unmarshal([]byte(tc.jsonData), &req)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, req.ChatID.String())
		})
	}
}

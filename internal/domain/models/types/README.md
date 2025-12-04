# ID Types with Robust JSON/Text Unmarshaling

This package provides type-safe ID types for the AI Document Backend application. These custom ID types solve common binding and parsing issues that occur with raw `uuid.UUID` types, particularly when dealing with Gin form/query parameter binding.

## Problem Solved

The original issue occurred when form parameters were sent in array format (e.g., `["uuid-string"]`) instead of plain string format (`"uuid-string"`). This caused binding errors because the standard `uuid.UUID` type's unmarshaling couldn't handle array inputs.

### Before (Problematic)
```go
type AnswerQuestionParams struct {
    ChatID   uuid.UUID `form:"chat_id" binding:"required"`  // ❌ Fails with ["uuid"] input
    Question string    `form:"question" binding:"required"`
}
```

### After (Robust)
```go
type AnswerQuestionParams struct {
    ChatID   types.ChatID `form:"chat_id" binding:"required"`  // ✅ Handles both "uuid" and ["uuid"]
    Question string       `form:"question" binding:"required"`
}
```

## Available ID Types

- `UserID` - Identifies users
- `DocumentID` - Identifies documents
- `ChatID` - Identifies chat sessions
- `MessageID` - Identifies chat messages
- `RefreshTokenID` - Identifies refresh tokens

## Key Features

### 1. Robust Parsing
All ID types can handle multiple input formats:
- Plain UUID strings: `"123e4567-e89b-12d3-a456-426614174000"`
- JSON arrays: `["123e4567-e89b-12d3-a456-426614174000"]`
- UUID without dashes: `"123e4567e89b12d3a456426614174000"`

### 2. Type Safety
Each ID type is distinct, preventing accidental mixing of different entity IDs:
```go
userID := types.GenerateUserID()
chatID := types.GenerateChatID()
// userID and chatID are different types - compile-time safety!
```

### 3. Database Compatibility
All ID types implement `driver.Valuer` and `sql.Scanner` for seamless database integration:
```go
// Works with database queries
user, err := repo.GetUser(ctx, userID)
```

### 4. JSON/Text Marshaling
Full support for JSON and text marshaling/unmarshaling:
```go
// JSON marshaling
data, err := json.Marshal(chatID)  // → "123e4567-e89b-12d3-a456-426614174000"

// JSON unmarshaling (handles both string and array formats)
var chatID ChatID
json.Unmarshal([]byte(`"123e4567-e89b-12d3-a456-426614174000"`), &chatID)  // ✅
json.Unmarshal([]byte(`["123e4567-e89b-12d3-a456-426614174000"]`), &chatID) // ✅
```

## Usage Examples

### Creating IDs
```go
// From existing UUID
userUUID := uuid.New()
userID := types.NewUserID(userUUID)

// From string
chatID, err := types.NewChatIDFromString("123e4567-e89b-12d3-a456-426614174000")

// Generate new random ID
docID := types.GenerateDocumentID()
```

### Converting Back to UUID
```go
chatID := types.GenerateChatID()
originalUUID := chatID.UUID()  // Get the underlying uuid.UUID
```

### Checking for Zero/Nil Values
```go
var userID types.UserID
if userID.IsZero() {
    // Handle uninitialized ID
}
```

### Database Operations
```go
// The ID types work seamlessly with SQL queries
func (r *Repository) GetChat(ctx context.Context, chatID types.ChatID, userID types.UserID) error {
    return r.db.QueryRow(
        "SELECT * FROM chats WHERE id = $1 AND user_id = $2",
        chatID,  // Automatically converted to UUID for database
        userID,
    ).Scan(...)
}
```

## Error Cases Handled

The `UnmarshalJSON` and `UnmarshalText` methods gracefully handle various error cases:

- Empty strings or arrays
- Invalid UUID formats  
- Non-string/non-array JSON values
- Malformed JSON

## Performance

The ID types have minimal performance overhead:
- Wrapping: Nearly zero cost (single UUID field)
- Parsing: Comparable to `uuid.Parse()` with additional array handling
- Memory: Same as `uuid.UUID` (16 bytes)

## Migration from uuid.UUID

To migrate existing code:

1. Replace `uuid.UUID` with appropriate typed ID:
   ```go
   // Before
   func GetUser(userID uuid.UUID) error { ... }
   
   // After  
   func GetUser(userID types.UserID) error { ... }
   ```

2. Update creation calls:
   ```go
   // Before
   userID := uuid.New()
   
   // After
   userID := types.GenerateUserID()
   ```

3. Access underlying UUID when needed:
   ```go
   // When you need the raw UUID (e.g., for legacy APIs)
   rawUUID := userID.UUID()
   ```

## Testing

The package includes comprehensive tests covering:
- All parsing scenarios (string, array, edge cases)
- JSON marshaling/unmarshaling
- Database driver compatibility
- Type safety verification
- Gin binding simulation
- Performance benchmarks

Run tests with:
```bash
go test ./internal/domain/models/types -v
```

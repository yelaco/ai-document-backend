package types

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Example demonstrates how the UnmarshalJSON methods solve the original binding issue
func ExampleChatID_UnmarshalJSON() {
	// This is the problematic case that was causing errors before
	// When a form parameter is sent as ["uuid"] instead of "uuid"
	problematicJSON := `["91397617-81f8-4bbb-8104-e037789456f0"]`

	var chatID ChatID
	err := json.Unmarshal([]byte(problematicJSON), &chatID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Successfully parsed ChatID: %s\n", chatID.String())
	// Output: Successfully parsed ChatID: 91397617-81f8-4bbb-8104-e037789456f0
}

// ExampleGinBinding demonstrates how this would work with Gin form binding
func TestExampleGinBinding(t *testing.T) {
	// Simulate different ways Gin might send form/query parameters
	testCases := []struct {
		name        string
		description string
		jsonData    string
		expectError bool
	}{
		{
			name:        "normal_string",
			description: "Normal case: form parameter sent as string",
			jsonData:    `{"chat_id": "91397617-81f8-4bbb-8104-e037789456f0", "question": "What is this document about?"}`,
			expectError: false,
		},
		{
			name:        "problematic_array",
			description: "Problematic case: form parameter sent as array (this was failing before)",
			jsonData:    `{"chat_id": ["91397617-81f8-4bbb-8104-e037789456f0"], "question": "What is this document about?"}`,
			expectError: false,
		},
		{
			name:        "multiple_values_array",
			description: "Edge case: multiple values in array (takes first)",
			jsonData:    `{"chat_id": ["91397617-81f8-4bbb-8104-e037789456f0", "another-uuid"], "question": "What is this document about?"}`,
			expectError: false,
		},
	}

	// This struct mimics what might be used in a Gin handler
	type AnswerQuestionParams struct {
		ChatID   ChatID `json:"chat_id" form:"chat_id" binding:"required"`
		Question string `json:"question" form:"question" binding:"required"`
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing: %s", tc.description)

			var params AnswerQuestionParams
			err := json.Unmarshal([]byte(tc.jsonData), &params)

			if tc.expectError {
				assert.Error(t, err)
				t.Logf("Expected error occurred: %v", err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "91397617-81f8-4bbb-8104-e037789456f0", params.ChatID.String())
				assert.Equal(t, "What is this document about?", params.Question)
				t.Logf("Successfully parsed ChatID: %s", params.ChatID.String())
			}
		})
	}
}

// ExampleComparison demonstrates the before/after difference
func TestExampleComparison(t *testing.T) {
	// Before: Using raw uuid.UUID would fail with array input
	// After: Using our custom ChatID type succeeds

	arrayInput := `["91397617-81f8-4bbb-8104-e037789456f0"]`
	expectedUUID := "91397617-81f8-4bbb-8104-e037789456f0"

	t.Run("with_custom_ChatID", func(t *testing.T) {
		var chatID ChatID
		err := json.Unmarshal([]byte(arrayInput), &chatID)
		require.NoError(t, err, "Our custom ChatID should handle array input")
		assert.Equal(t, expectedUUID, chatID.String())
		t.Logf("✅ Custom ChatID successfully handled array input: %s", chatID.String())
	})

	t.Run("demonstrates_the_original_problem", func(t *testing.T) {
		// This demonstrates what would happen with a regular string field
		var stringField string
		err := json.Unmarshal([]byte(arrayInput), &stringField)
		assert.Error(t, err, "Regular string field should fail with array input")
		t.Logf("❌ Regular string field failed with array input (as expected): %v", err)
	})
}

// ExampleAllIDTypes demonstrates that all ID types have the same robust parsing
func TestExampleAllIDTypes(t *testing.T) {
	testUUID := "91397617-81f8-4bbb-8104-e037789456f0"
	arrayFormat := fmt.Sprintf(`["%s"]`, testUUID)

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "UserID",
			test: func(t *testing.T) {
				var id UserID
				err := json.Unmarshal([]byte(arrayFormat), &id)
				require.NoError(t, err)
				assert.Equal(t, testUUID, id.String())
			},
		},
		{
			name: "DocumentID",
			test: func(t *testing.T) {
				var id DocumentID
				err := json.Unmarshal([]byte(arrayFormat), &id)
				require.NoError(t, err)
				assert.Equal(t, testUUID, id.String())
			},
		},
		{
			name: "ChatID",
			test: func(t *testing.T) {
				var id ChatID
				err := json.Unmarshal([]byte(arrayFormat), &id)
				require.NoError(t, err)
				assert.Equal(t, testUUID, id.String())
			},
		},
		{
			name: "MessageID",
			test: func(t *testing.T) {
				var id MessageID
				err := json.Unmarshal([]byte(arrayFormat), &id)
				require.NoError(t, err)
				assert.Equal(t, testUUID, id.String())
			},
		},
		{
			name: "RefreshTokenID",
			test: func(t *testing.T) {
				var id RefreshTokenID
				err := json.Unmarshal([]byte(arrayFormat), &id)
				require.NoError(t, err)
				assert.Equal(t, testUUID, id.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
			t.Logf("✅ %s successfully handled array format input", tt.name)
		})
	}
}

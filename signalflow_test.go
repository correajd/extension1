package splunk

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/signalfx/signalflow-client-go/signalflow/messages"
	"github.com/stretchr/testify/assert"
)

// TestSignalFlowIntegration tests the SignalFlow client integration
// This is an integration test that requires valid credentials
func TestSignalFlowIntegration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Get credentials from environment variables
	token := os.Getenv("O11Y_TOKEN")
	realm := os.Getenv("O11Y_REALM")

	if token == "" || realm == "" {
		t.Fatal("O11Y_TOKEN and O11Y_REALM environment variables must be set for integration tests")
	}

	// Create a new SignalFlow client
	si := &SignalflowInstance{}
	client, err := si.NewSignalFlow(token, realm)
	assert.NoError(t, err, "Failed to create SignalFlow client")
	defer client.Close()

	// Test parameters matching test_signalflow.js
	program := `data('demo.trans.latency').publish()`
	now := time.Now()
	oneHourAgo := now.Add(-time.Hour)
	resolution := 60000 // 1 minute in milliseconds

	// Convert times to milliseconds since epoch
	start := oneHourAgo.UnixNano() / int64(time.Millisecond)
	stop := now.UnixNano() / int64(time.Millisecond)

	t.Run("Execute SignalFlow program", func(t *testing.T) {
		// Execute the SignalFlow program
		comp, err := client.Execute(program, start, stop, resolution)
		if !assert.NoError(t, err, "Failed to execute SignalFlow program") {
			return
		}
		defer comp.Close()

		// Process messages
		messageCount := 0
		for {
			msg, ok := comp.Next()
			if !ok {
				log.Println("No more messages")
				break
			}

			// Check if the message is a DataMessage
			if dataMsg, ok := msg.(*messages.DataMessage); ok {
				messageCount++
				log.Printf("Received DataMessage - Timestamp: %d, Payloads: %d",
					dataMsg.TimestampMillis, len(dataMsg.Payloads))

				// Log each payload in the DataMessage
				for i, payload := range dataMsg.Payloads {
					log.Printf("  Payload %d - Type: %v, TSID: %v, Value: %v",
						i+1, payload.Type, payload.TSID, payload.Value())
				}
			} else if metaMsg, ok := msg.(*messages.MetadataMessage); ok {
				log.Printf("Received MetadataMessage: %+v", metaMsg)

			} else {
				log.Printf("Received non-DataMessage: %T - %+v", msg, msg)
			}
		}

		assert.Greater(t, messageCount, 0, "Expected to receive at least one message")
	})
}

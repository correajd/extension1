package splunk

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/signalfx/signalflow-client-go/signalflow"
	"go.k6.io/k6/js/modules"
)

// Register the module on init
func init() {
	modules.Register("k6/x/signalflow", new(SignalflowRootModule))
}

type (
	// SignalflowRootModule is the global module instance that will create module
	// instances for each VU.
	SignalflowRootModule struct{}

	// SignalflowInstance represents an instance of the JS module.
	SignalflowInstance struct {
		vu modules.VU
	}

	// Client represents a SignalFlow client that can be used from JavaScript
	Client struct {
		client *signalflow.Client
	}

	// Computation wraps a SignalFlow computation with methods to interact with it
	Computation struct {
		comp *signalflow.Computation
	}
)

// NewModuleInstance implements the modules.Module interface and returns
// a new instance for each VU.
func (r *SignalflowRootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &SignalflowInstance{vu: vu}
}

// Exports implements the modules.Instance interface and returns the exported
// symbols for the JS module.
func (si *SignalflowInstance) Exports() modules.Exports {
	return modules.Exports{
		Default: si,
	}
}

// NewSignalFlow creates a new SignalFlow client
func (si *SignalflowInstance) NewSignalFlow(token string, realm string) (*Client, error) {
	log.Printf("Creating new SignalFlow client. Realm: %s, Token: %s... (truncated for security)", realm, token[:min(len(token), 5)])

	streamURL := signalflow.StreamURLForRealm(realm)
	log.Printf("Using Stream URL: %s", streamURL)

	client, err := signalflow.NewClient(
		streamURL,
		signalflow.AccessToken(token),
		signalflow.OnError(func(err error) {
			log.Printf("Error in SignalFlow client: %v", err)
		}),
	)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to create SignalFlow client: %v", err)
		log.Printf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	log.Printf("Successfully created SignalFlow client for realm: %s", realm)
	return &Client{client: client}, nil
}

// min is a helper function to get the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// DataPoint represents a single data point in a time series
type DataPoint struct {
	TS    int64       `js:"ts"`
	Value interface{} `js:"value"`
}

// MetricData represents a time series for a specific metric
type MetricData struct {
	Metric     string            `js:"metric"`
	Dimensions map[string]string `js:"dimensions"`
	Data       []DataPoint       `js:"data"`
}

// Next returns the next data message from the computation as a JavaScript-friendly object
func (c *Computation) Next() (interface{}, bool) {
	msg, ok := <-c.comp.Data()
	if !ok {
		return nil, false
	}

	// For now, just return the raw message and let the JavaScript side handle it
	// You can add more sophisticated conversion logic here as needed
	return msg, true
}

// Close stops the computation
func (c *Computation) Close() error {
	return c.comp.Stop(nil)
}

// Execute runs a SignalFlow program and returns a Computation object
func (c *Client) Execute(program string, start, stop int64, resolution int) (*Computation, error) {
	log.Printf("Executing SignalFlow program. Start: %d, Stop: %d, Resolution: %ds", start, stop, resolution)
	log.Printf("Program:\n%s", program)

	if c.client == nil {
		errMsg := "SignalFlow client is not initialized"
		log.Printf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// Create a context with timeout
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()

	// Convert timestamps to time.Time
	// JavaScript gives us milliseconds since epoch, but time.Unix expects seconds
	startTime := time.Unix(start/1000, 0)
	stopTime := time.Unix(stop/1000, 0)
	resolutionDuration := time.Duration(resolution) * time.Millisecond

	log.Printf("Converted timestamps - Start: %v, Stop: %v, Resolution: %v",
		startTime.Format(time.RFC3339),
		stopTime.Format(time.RFC3339),
		resolutionDuration)

	// Execute the SignalFlow program
	log.Printf("Sending execute request to SignalFlow...")
	comp, err := c.client.Execute(context.Background(), &signalflow.ExecuteRequest{
		Program:    program,
		Start:      startTime,
		Stop:       stopTime,
		Resolution: resolutionDuration,
	})

	if err != nil {
		errMsg := fmt.Sprintf("Failed to execute SignalFlow program: %v", err)
		log.Printf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	log.Printf("Successfully started SignalFlow execution")
	return &Computation{comp: comp}, nil
}

// Close stops the SignalFlow client
func (c *Client) Close() {
	if c.client != nil {
		err := c.client.Stop(context.Background(), &signalflow.StopRequest{})
		if err != nil {
			return
		}
	}
}

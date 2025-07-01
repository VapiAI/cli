package analytics

import (
	"time"

	"github.com/spf13/cobra"
)

// TrackCommandWrapper wraps a cobra command function with analytics tracking
func TrackCommandWrapper(cmdName string, subCmd string, fn func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		startTime := time.Now()

		err := fn(cmd, args)

		duration := time.Since(startTime)
		success := err == nil
		errorMsg := ""
		if err != nil {
			errorMsg = err.Error()
		}

		TrackCommand(cmdName, subCmd, success, duration, errorMsg)

		return err
	}
}

// TrackCommandWithContext wraps a command with additional context properties
func TrackCommandWithContext(cmdName string, subCmd string, context map[string]interface{}, fn func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		startTime := time.Now()

		err := fn(cmd, args)

		duration := time.Since(startTime)
		success := err == nil

		// Merge context into properties
		properties := map[string]interface{}{
			"command":     cmdName,
			"subcommand":  subCmd,
			"success":     success,
			"duration_ms": duration.Milliseconds(),
		}

		for k, v := range context {
			if isSafeProperty(k) {
				properties[k] = v
			}
		}

		eventName := EventCommandExecuted
		if !success {
			eventName = EventCommandFailed
			if err != nil {
				properties["error_hash"] = hashString(err.Error())
				properties["error_type"] = categorizeError(err.Error())
			}
		}

		TrackEvent(eventName, properties)

		return err
	}
}

// LogOptOutStatus provides information about analytics status for debugging
func LogOptOutStatus() map[string]interface{} {
	return map[string]interface{}{
		"analytics_enabled": IsEnabled(),
		"opted_out":         isOptedOut(),
		"env_vars_checked": []string{
			"VAPI_DISABLE_ANALYTICS",
			"VAPI_NO_TELEMETRY",
			"DISABLE_TELEMETRY",
			"DO_NOT_TRACK",
		},
	}
}

package config
import (
    "encoding/json"
    "fmt"
)

// Status represents the status of the operation.
type Status string
const (
    StatusOK    Status = "ok"
    StatusError Status = "error"
)

// IsValid checks if the Status is one of the predefined constants.
func (s Status) IsValid() bool {
    switch s {
    case StatusOK, StatusError:
        return true
    }
    return false
}

// MarshalJSON ensures that only valid Status values are marshalled.
func (s Status) MarshalJSON() ([]byte, error) {
    if !s.IsValid() {
        return nil, fmt.Errorf("invalid status value: %s", s)
    }
    return json.Marshal(string(s))
}

// UnmarshalJSON ensures that only valid Status values are unmarshalled.
// not currently needed, but here for completeness
func (s *Status) UnmarshalJSON(data []byte) error {
    var statusStr string
    if err := json.Unmarshal(data, &statusStr); err != nil {
        return err
    }

    tempStatus := Status(statusStr)
    if !tempStatus.IsValid() {
        return fmt.Errorf("invalid status value: %s", statusStr)
    }

    *s = tempStatus
    return nil
}

// StatusOutput represents the structure of the JSON output.
type StatusOutput struct {
    Status    Status `json:"status"`     // "ok" or "error"
    Message   string `json:"message"`    // Descriptive status message
    EnvFile   string `json:"env_file"`   // Path to the .env file
    JSONFile  string `json:"json_file"`  // Path to the JSON output file
}

//
// creates a success output
func CreateSuccessOutput(jsonFile, envFile string) string {
	return createStatusOutput(true, "",  jsonFile, envFile)
}

// 
//	creates an error output
func CreateErrorOutput(err error) string {
	msg := fmt.Sprintf("%s", err)
	return createStatusOutput(false, msg, "", "")
}


// CreateStatusOutput generates a JSON document based on the provided parameters.
// It takes a success flag, a message, and file paths for the .env and JSON output files.
// It always returns a JSON string indicating the status, message, and file paths.
// In case of marshalling failure, it returns a default JSON error message.

func createStatusOutput(success bool, message, jsonFile, envFile string) string {
   status := StatusOK
   if !success {
	status = StatusError
   }
	output := StatusOutput{
        Status:   status,
        Message:  message,
        EnvFile:  envFile,
        JSONFile: jsonFile,
    }

    jsonBytes, err := json.Marshal(output)
    if err != nil {
        return fmt.Sprintf(`
{
	"Status": "error",
	"Message": "Unable to Marshal Status Output.  Fatal Error: %s",
	"EnvFile: "",
	"JSONFile": ""		
}
		`, err)
    }

    return string(jsonBytes)
}

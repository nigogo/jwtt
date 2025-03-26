package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

func decodeSegment(segment string) ([]byte, error) {
	segment = strings.ReplaceAll(segment, "-", "+")
	segment = strings.ReplaceAll(segment, "_", "/")
	switch len(segment) % 4 {
	case 2:
		segment += "=="
	case 3:
		segment += "="
	}
	return base64.StdEncoding.DecodeString(segment)
}

func convertTimestamps(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		result[k] = v
	}

	keys := []string{"iat", "exp", "nbf"}
	for _, k := range keys {
		if v, exists := result[k]; exists {
			switch num := v.(type) {
			case float64:
				t := time.Unix(int64(num), 0)
				result[k] = fmt.Sprintf("%d (%s)", int64(num), t.Format(time.RFC3339))
			}
		}
	}
	return result
}

func main() {
	errorColor := color.New(color.FgRed).SprintFunc()
	headerColor := color.New(color.FgCyan).SprintFunc()
	payloadColor := color.New(color.FgGreen).SprintFunc()
	signatureColor := color.New(color.FgMagenta).SprintFunc()
	payloadColorBold := color.New(color.FgGreen, color.Bold).SprintFunc()
	headerColorBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	signatureColorBold := color.New(color.FgMagenta, color.Bold).SprintFunc()

	convertTimestampsFlag := flag.Bool("t", false, "Convert timestamps to human-readable format")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Println("Usage: jwtt [-t] <JWT token>")
		os.Exit(1)
	}

	token := flag.Arg(0)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		fmt.Println(errorColor("Invalid JWT token: Expected 3 parts, got %d", len(parts)))
		os.Exit(1)
	}

	headerBytes, err := decodeSegment(parts[0])
	if err != nil {
		fmt.Println(errorColor("Error decoding header:"), err)
		os.Exit(1)
	}

	payloadBytes, err := decodeSegment(parts[1])
	if err != nil {
		fmt.Println(errorColor("Error decoding payload:"), err)
		os.Exit(1)
	}

	var header, payload map[string]interface{}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		fmt.Println(errorColor("Error parsing header JSON:"), err)
		os.Exit(1)
	}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		fmt.Println(errorColor("Error parsing payload JSON:"), err)
		os.Exit(1)
	}
	signature := parts[2]

	if *convertTimestampsFlag {
		header = convertTimestamps(header)
		payload = convertTimestamps(payload)
	}

	fmt.Println(headerColorBold("Header:"))
	headerJSON, _ := json.MarshalIndent(header, "", "  ")
	fmt.Println(headerColor(string(headerJSON)))

	fmt.Println(payloadColorBold("Payload:"))
	payloadJSON, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Println(payloadColor(string(payloadJSON)))

	fmt.Println(signatureColorBold("Signature:"))
	fmt.Println(signatureColor(signature))
}

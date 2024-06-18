package main

import (
	"bufio"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	serverAddress := getEnv("SERVER_ADDRESS", "localhost:8080")

	for {
		err := sendRequest(serverAddress)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		time.Sleep(1 * time.Second)
	}
}

func sendRequest(serverAddress string) error {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return fmt.Errorf("error connecting to server: %v", err)
	}
	defer conn.Close()

	// Read response
	reader := bufio.NewReader(conn)

	// Send "word_of_wisdom" request
	_, err = conn.Write([]byte("word_of_wisdom\n"))
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}

	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	fmt.Printf("Response: %s\n", response)

	if strings.HasPrefix(response, "X-Hashcash:") {
		parts := strings.Split(response, ":")
		if len(parts) < 7 {
			return fmt.Errorf("invalid Hashcash format")
		}

		complexityStr := parts[5]

		challenge := strings.TrimSpace(strings.ReplaceAll(response, "X-Hashcash:", ""))

		complexity, err := strconv.Atoi(complexityStr)
		if err != nil {
			return err
		}
		// Solve the challenge
		solution := SolveChallenge(challenge, complexity)

		// Send the solution
		_, err = conn.Write([]byte(fmt.Sprintln(solution)))
		if err != nil {
			return fmt.Errorf("error sending solution: %v", err)
		}

		quote, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading quote: %v", err)
		}
		fmt.Printf("Quote: %s\n", quote)
	} else {
		fmt.Println("Received quote without DDoS protection")
	}
	return nil
}

// SolveChallenge solves the PoW challenge with the given difficulty.
func SolveChallenge(challenge string, difficulty int) string {
	solution := sha256.Sum256([]byte(""))

	for {
		hash := sha1.New()
		hash.Write([]byte(challenge))
		hash.Write([]byte(":"))
		hash.Write(solution[:])

		if strings.HasPrefix(hex.EncodeToString(hash.Sum(nil)[:]), strings.Repeat("0", difficulty)) {
			break
		}
		solution = sha256.Sum256(solution[:])
	}
	return string(solution[:])
}

// getEnv reads an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	filePath, timeLimit := initFile()

	// "read in a quiz provided via a CSV file"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// read line by line and field by field
	reader := csv.NewReader(file)

	questions := make([]string, 0)
	answers := make([]string, 0)

	for {
		record, err := reader.Read()

		if err != nil {
			// Check for end of file
			if err.Error() == "EOF" {
				break
			} else {
				fmt.Println("Error reading record:", err)
				return
			}
		}

		if len(record) == 2 {
			questions = append(questions, record[0])
			answers = append(answers, record[1])
		} else {
			continue
		}
	}

	// keeping track of how many questions they get right and how many they get incorrect.
	score := 0
	scanner := bufio.NewScanner(os.Stdin)

	// create a ticker for the time limit and a channel to signal the user finished the quiz
	ticker := time.NewTicker(time.Second * time.Duration(timeLimit))
	done := make(chan bool)

	WaitForStart()

	// Answer questions
	go func() {
		for i, question := range questions {
			fmt.Printf("Question %d: %s\nYour answer: ", i+1, question)

			scanner.Scan()
			userAnswer := scanner.Text()
			userAnswer = strings.TrimSpace(userAnswer)

			if userAnswer == answers[i] {
				fmt.Println("Correct!")
				score++
			} else {
				fmt.Println("Incorrect. The correct answer is:", answers[i])
			}
		}
		done <- true
	}()

	// select chooses the first channel with an available value
	// if done is available first, the user finished
	// if ticker is available first, the time limit has been reached
	select {
	case <-done:
	case <-ticker.C:
		fmt.Println("time's up!")
	}

	//  output the total number of questions correct and how many questions there were in total.
	fmt.Printf("Quiz completed! Your score: %d/%d\n", score, len(questions))
}

// default to problems.csv (example shown below), but the user should be able to customize the filename via a flag.
// default time limit should be 30 seconds, but should also be customizable via a flag.
func initFile() (string, int) {
	filePath := flag.String("file", "problems.csv", "Path to the CSV file")
	timeLimit := flag.Int("timer", 30, "Time Limit is 30sec by default")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("If problems.csv is not created, please provide a file path using the -file flag.")
	}

	fmt.Printf("Reading file: %s\n", *filePath)

	return *filePath, *timeLimit
}

// WaitForStart makes the program wait for a user to press Enter button
func WaitForStart() {
	fmt.Print("Press 'Enter' to start the quiz.")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func Add(a, b int) int {
	return a + b
}

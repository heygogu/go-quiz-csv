package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Question struct {
	question string
	answer   string
}

func main() {
	filename, timeLimit := readArguments()
	fmt.Println(timeLimit)
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Unable to open the target file")
		return
	}
	defer f.Close()

	allQuestions, err := readCSV(f)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if len(allQuestions) == 0 {
		fmt.Println("Not enough questions")
	}

	//let us start to ask the question
	score, err := startAskingQuestions(allQuestions, timeLimit)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Your score is ", score, "/", len(allQuestions))
		return
	}
	fmt.Println("Your score is ", score, "/", len(allQuestions))

}

func readCSV(f io.Reader) ([]Question, error) {

	allQuestions, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	numOfQuestions := len(allQuestions)

	if numOfQuestions == 0 {
		return nil, fmt.Errorf("no questions are in this file")
	}

	var data []Question

	for _, line := range allQuestions {
		ques := Question{}
		ques.question = line[0]
		ques.answer = line[1]
		data = append(data, ques)

	}
	return data, nil
}

func readArguments() (string, int) {
	filename := flag.String("filename", "problem.csv", "CSV File that conatins quiz questions")
	timeLimit := flag.Int("limit", 30, "Time Limit for each question")
	flag.Parse()
	return *filename, *timeLimit
}

func startAskingQuestions(questions []Question, timeLimit int) (int, error) {
	score := 0
	userAns := make(chan string)
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	go getInput(userAns)

	for _, question := range questions {
		ans, err := eachQuestion(question.question, question.answer, userAns, timer.C)
		if int(ans) == -1 {
			return score, fmt.Errorf("time out")
		} else if err != nil {

		} else {
			score += ans
		}
	}

	return score, nil
}

func eachQuestion(question string, answer string, userAns <-chan string, timer <-chan time.Time) (int, error) {

	fmt.Print(question, ": ")
	for {
		select {
		case ans := <-userAns:
			score := 0
			if strings.Compare(strings.Trim(strings.ToLower(ans), "\n"), answer) == 0 {
				score = 1
			} else {
				return 0, fmt.Errorf("wrong answer")
			}
			return score, nil
		case <-timer:
			return -1, fmt.Errorf("time out")

		}
	}

}

func getInput(input chan string) {
	for {
		in := bufio.NewReader(os.Stdin)
		result, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			return
		}
		input <- result
	}
}

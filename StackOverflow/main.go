package main

import (
	"fmt"
	"log"

	"stackoverflow/models"
	"stackoverflow/services"
	"stackoverflow/strategies"
)

func main() {
	so := services.NewStackOverflow(strategies.KeyWordSearch{})

	alice := so.RegisterUser("Alice")
	bob := so.RegisterUser("Bob")
	carol := so.RegisterUser("Carol")

	fmt.Println("=== Posting ===")
	q, err := so.PostQuestion(alice.ID, "How do goroutines work?",
		"I am confused about goroutines and channels.", []string{"go", "concurrency"})
	check(err)
	fmt.Printf("Alice asked Q%d: %s\n", q.ID, q.Title)

	a1, err := so.PostAnswer(bob.ID, q.ID, "Goroutines are lightweight threads managed by the Go runtime.")
	check(err)
	fmt.Printf("Bob answered A%d\n", a1.ID)

	a2, err := so.PostAnswer(carol.ID, q.ID, "Think of them as functions running at the same time.")
	check(err)
	fmt.Printf("Carol answered A%d\n", a2.ID)

	fmt.Println("\n=== Voting ===")
	check(so.VoteQuestion(bob.ID, q.ID, models.Upvote))
	check(so.VoteAnswer(alice.ID, a1.ID, models.Upvote))
	check(so.VoteAnswer(carol.ID, a1.ID, models.Upvote))
	check(so.VoteAnswer(alice.ID, a2.ID, models.Downvote))

	// rule checks: these SHOULD fail
	if err := so.VoteAnswer(bob.ID, a1.ID, models.Upvote); err != nil {
		fmt.Println("Bob votes on own answer ->", err)
	}
	if err := so.VoteAnswer(alice.ID, a1.ID, models.Upvote); err != nil {
		fmt.Println("Alice votes again ->", err)
	}

	fmt.Println("\n=== Question page ===")
	printQuestion(q)

	fmt.Println("\n=== Reputation ===")
	for _, u := range []*models.User{alice, bob, carol} {
		fmt.Printf("%s: %d\n", u.Name, u.Reputation)
	}

	_, err = so.PostQuestion(bob.ID, "How to read a file in Python?",
		"I want to read a text file line by line.", []string{"python", "files"})
	check(err)

	fmt.Println("\n=== Search by keyword \"goroutines\" ===")
	printResults(so.Search("goroutines"))

	so.SetSearchStrategy(strategies.TagSearch{}) // swap strategy at runtime
	fmt.Println("\n=== Search by tag \"python\" ===")
	printResults(so.Search("python"))
}

func printQuestion(q *models.Question) {
	tagNames := []string{}
	for _, t := range q.Tags {
		tagNames = append(tagNames, t.Name)
	}
	fmt.Printf("Q%d: %s  tags=%v  score=%d  by %s\n", q.ID, q.Title, tagNames, q.Votes.Score(), q.Author.Name)
	for _, a := range q.Answers {
		fmt.Printf("   A%d: %s  score=%d  by %s\n", a.ID, a.Body, a.Votes.Score(), a.Author.Name)
	}
}

func printResults(qs []*models.Question) {
	for _, q := range qs {
		fmt.Printf("- Q%d: %s\n", q.ID, q.Title)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

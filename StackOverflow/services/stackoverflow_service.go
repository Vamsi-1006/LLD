package services

import (
	"errors"
	"sort"
	"strings"
	"sync"

	"stackoverflow/models"
	"stackoverflow/strategies"
)

const (
	upvotePoints   = 10
	downvotePoints = -2
)

// Errors as values: callers can check them with errors.Is
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrQuestionNotFound = errors.New("question not found")
	ErrAnswerNotFound   = errors.New("answer not found")
	ErrEmptyText        = errors.New("text cannot be empty")
	ErrOwnPost          = errors.New("you cannot vote on your own post")
	ErrAlreadyVoted     = errors.New("you have already voted on this post")
)

type StackOverflow struct {
	mu        sync.Mutex // many users may vote/post at same time
	nextID    int
	users     map[int]*models.User
	questions map[int]*models.Question // questionId -> Question for easy access
	answers   map[int]*models.Answer   // answerId -> Answer for easy access
	tags      map[string]*models.Tag
	searcher  strategies.SearchStrategy
}

func NewStackOverflow(searcher strategies.SearchStrategy) *StackOverflow {
	return &StackOverflow{
		users:     make(map[int]*models.User),
		questions: make(map[int]*models.Question),
		answers:   make(map[int]*models.Answer),
		tags:      make(map[string]*models.Tag),
		searcher:  searcher,
	}
}

// newID must be called while holding the lock.
func (s *StackOverflow) newID() int {
	s.nextID++
	return s.nextID
}

func (s *StackOverflow) RegisterUser(name string) *models.User {
	s.mu.Lock()
	defer s.mu.Unlock()

	u := models.NewUser(s.newID(), name)
	s.users[u.ID] = u
	return u
}

func (s *StackOverflow) PostQuestion(userID int, title, body string, tagNames []string) (*models.Question, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. validate
	author, ok := s.users[userID]
	if !ok {
		return nil, ErrUserNotFound
	}
	if strings.TrimSpace(title) == "" {
		return nil, ErrEmptyText
	}

	// 2. build tags (reuse existing ones)
	tags := make([]*models.Tag, 0, len(tagNames))
	for _, name := range tagNames {
		tags = append(tags, s.getOrCreateTag(name))
	}

	// 3. create + store
	q := models.NewQuestion(s.newID(), title, body, author, tags)
	s.questions[q.ID] = q
	return q, nil
}

func (s *StackOverflow) PostAnswer(userID, questionID int, body string) (*models.Answer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	author, ok := s.users[userID]
	if !ok {
		return nil, ErrUserNotFound
	}
	q, ok := s.questions[questionID]
	if !ok {
		return nil, ErrQuestionNotFound
	}
	if strings.TrimSpace(body) == "" {
		return nil, ErrEmptyText
	}

	a := models.NewAnswer(s.newID(), body, author, q.ID)
	s.answers[a.ID] = a
	q.AddAnswer(a)
	return a, nil
}

func (s *StackOverflow) VoteQuestion(userID, questionID int, vt models.VoteType) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	voter, ok := s.users[userID]
	if !ok {
		return ErrUserNotFound
	}

	q, ok := s.questions[questionID]
	if !ok {
		return ErrQuestionNotFound
	}

	return s.applyVote(voter, q.Author, q.Votes, vt)
}

func (s *StackOverflow) VoteAnswer(userID, answerID int, vt models.VoteType) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	voter, ok := s.users[userID]
	if !ok {
		return ErrUserNotFound
	}
	a, ok := s.answers[answerID]
	if !ok {
		return ErrAnswerNotFound
	}
	return s.applyVote(voter, a.Author, a.Votes, vt)
}

// applyVote holds the voting rules in ONE place, used by both questions and answers.
func (s *StackOverflow) applyVote(voter, author *models.User, votes *models.Votes, vt models.VoteType) error {
	if voter.ID == author.ID {
		return ErrOwnPost
	}
	if votes.HasVoted(voter.ID) {
		return ErrAlreadyVoted
	}

	votes.AddVote(voter.ID, vt)

	if vt == models.Upvote {
		author.AddReputation(upvotePoints)
	} else {
		author.AddReputation(downvotePoints)
	}
	return nil
}

func (s *StackOverflow) SetSearchStrategy(st strategies.SearchStrategy) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.searcher = st
}

func (s *StackOverflow) Search(query string) []*models.Question {
	s.mu.Lock()
	defer s.mu.Unlock()

	all := make([]*models.Question, 0, len(s.questions))
	for _, q := range s.questions {
		all = append(all, q)
	}
	// maps have random order, so sort by ID for stable output
	sort.Slice(all, func(i, j int) bool { return all[i].ID < all[j].ID })

	return s.searcher.Search(all, query)
}

// getOrCreateTag must be called while holding the lock.
func (s *StackOverflow) getOrCreateTag(name string) *models.Tag {
	name = strings.ToLower(strings.TrimSpace(name))
	if t, ok := s.tags[name]; ok {
		return t
	}
	t := models.NewTag(name)
	s.tags[name] = t
	return t
}

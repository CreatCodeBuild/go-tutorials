package fakeadapter

import (
	"encoding/json"
	"goask/core/adapter"
	"goask/core/entity"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// Data satisfied adapter.Data. It serializes to dist.
type Data struct {
	questions Questions
	answers   Answers
	users     []entity.User
}

type dataSerialization struct {
	Questions Questions
	Answers Answers
	Users []entity.User
}

var _ adapter.Data = &Data{}

func (d *Data) file() string {
	return "./data.json"
}

func (d *Data) serialize() error {
	data := dataSerialization{
		Questions: d.questions,
		Answers: d.answers,
		Users: d.users,
	}

	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return errors.WithStack(err)
	}

	err = ioutil.WriteFile(d.file(), b, os.ModePerm)
	if err != nil {
		return errors.WithStack(err)
	}

	return err
}

func (d *Data) Questions(search *string) ([]entity.Question, error) {
	if search == nil {
		return d.questions, nil
	}
	ret := make([]entity.Question, 0)
	for _, q := range d.questions {
		if match(q.Content, *search) {
			ret = append(ret, q)
		}
	}
	return ret, nil
}

func (d *Data) QuestionByID(ID int64) (entity.Question, error) {
	for _, q := range d.questions {
		if q.ID == ID {
			return q, nil
		}
	}
	return entity.Question{}, errors.WithStack(&adapter.ErrQuestionNotFound{ID: ID})
}

func (d *Data) QuestionsByUserID(ID int64) ([]entity.Question, error) {
	var ret []entity.Question
	for _, q := range d.questions {
		if q.AuthorID == ID {
			ret = append(ret, q)
		}
	}
	return ret, nil
}

func (d *Data) CreateQuestion(q entity.Question) (entity.Question, error) {
	q.ID = int64(len(d.questions) + 1)
	d.questions = append(d.questions, q)
	return d.questions[len(d.questions)-1], d.serialize()
}

func (d *Data) UpdateQuestion(p entity.QuestionUpdate) (entity.Question, error) {
	if p.ID == 0 {
		return entity.Question{}, errors.New("ID can not be 0 nor absent")
	}
	for i, q := range d.questions {
		if q.ID == p.ID {
			if p.Content != nil {
				q.Content = *p.Content
			}
			if p.Title != nil {
				q.Title = *p.Title
			}
			d.questions[i] = q
			return q, d.serialize()
		}
	}
	return entity.Question{}, errors.WithStack(&adapter.ErrQuestionNotFound{ID: p.ID})
}

func (d *Data) AnswersOfQuestion(QuestionID int64) (ret []entity.Answer) {
	for _, answer := range d.answers {
		if answer.QuestionID == QuestionID {
			ret = append(ret, answer)
		}
	}
	return
}

func (d *Data) CreateAnswer(QuestionID int64, Content string, AuthorID int64) (entity.Answer, error) {
	for _, q := range d.questions {
		if q.ID == QuestionID {
			answer := d.answers.Add(QuestionID, Content, AuthorID)
			return answer, d.serialize()
		}
	}
	return entity.Answer{}, errors.WithStack(&adapter.ErrQuestionNotFound{ID: QuestionID})
}

func (d *Data) AcceptAnswer(AnswerID int64, UserID int64) (entity.Answer, error) {

	// Find the question this answer belongs to
	answer, ok := d.answers.Get(AnswerID)
	if !ok {
		return answer, errors.WithStack(&adapter.ErrAnswerNotFound{ID: AnswerID})
	}

	q, ok := d.questions.Get(answer.QuestionID)
	if !ok {
		return answer, errors.WithStack(&adapter.ErrQuestionOfAnswerNotFound{QuestionID: answer.QuestionID, AnswerID: AnswerID})
	}

	// Find if this user is the author of the question this answer belongs to
	if q.AuthorID != UserID {
		return answer, errors.WithStack(&adapter.ErrUserIsNotAuthorOfQuestion{QuestionID: q.ID, UserID: UserID})
	}

	answer = d.answers.Accept(AnswerID)
	return answer, d.serialize()
}

func (d *Data) UserByID(ID int64) (entity.User, error) {
	for _, user := range d.users {
		if user.ID == ID {
			return user, nil
		}
	}
	return entity.User{}, errors.WithStack(&adapter.ErrUserNotFound{ID: ID})
}

func (d *Data) Users() ([]entity.User, error) {
	return d.users, nil
}

func (d *Data) CreateUser(name string) (entity.User, error) {
	user := entity.User{ID: int64(len(d.users) + 1), Name: name}
	d.users = append(d.users, user)
	return user, d.serialize()
}

func match(s1, s2 string) bool {
	return strings.Contains(s1, s2)
}

type Questions []entity.Question

func (q *Questions) Get(questionID int64) (entity.Question, bool) {
	for _, qu := range *q {
		if qu.ID == questionID {
			return qu, true
		}
	}
	return entity.Question{}, false
}

type Answers []entity.Answer

func (a *Answers) Add(QuestionID int64, Content string, AuthorID int64) entity.Answer {
	// todo: serialize
	*a = append(*a, entity.Answer{
		ID:         int64(len(*a) + 1),
		Content:    Content,
		QuestionID: QuestionID,
		AuthorID:   AuthorID,
	})
	return (*a)[len(*a)-1]
}

func (a *Answers) OfQuestion(questionID int64) Answers {
	var ans Answers
	for _, answer := range *a {
		if answer.QuestionID == questionID {
			ans = append(ans, answer)
		}
	}
	return ans
}

func (a *Answers) Get(answerID int64) (entity.Answer, bool) {
	for _, an := range *a {
		if an.ID == answerID {
			return an, true
		}
	}
	return entity.Answer{}, false
}

func (a *Answers) Accept(answerID int64) entity.Answer {
	// todo: serialize
	for i := range *a {
		if (*a)[i].ID == answerID {
			(*a)[i].Accepted = true
			return (*a)[i]
		}
	}
	return entity.Answer{}
}
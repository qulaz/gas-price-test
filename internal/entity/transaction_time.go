package entity

import (
	"fmt"
	"strings"
	"time"
)

const transactionTimeLayout = "06-01-02 15:04"

type TransactionTime struct {
	time.Time
}

func NewTransactionTime(t time.Time) TransactionTime {
	return TransactionTime{Time: t}
}

func (t *TransactionTime) UnmarshalJSON(b []byte) error {
	var err error

	s := strings.Trim(string(b), "\"")
	t.Time, err = time.Parse(transactionTimeLayout, s)

	return err
}

func (t TransactionTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.Time.Format(transactionTimeLayout))), nil
}

func (t TransactionTime) String() string {
	return t.Format(transactionTimeLayout)
}

func (t TransactionTime) TruncToDay() TransactionTime {
	return TransactionTime{
		Time: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC),
	}
}

func (t TransactionTime) TruncToMonth() TransactionTime {
	return TransactionTime{
		Time: time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC),
	}
}

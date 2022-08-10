package models

import "time"

type Report struct {
	DoneCnt           int
	DeclinedCnt       int
	TaskId            string
	ReactionDurations []time.Duration
}

package models

import "time"

type EventType int

const (
	EventCreateType EventType = iota
	EventUpdateType
	EventDeleteType
	EventApproveType
	EventDeclineType
)

type TaskStatus int

const (
	TaskInProgressStatus TaskStatus = iota
	TaskDoneStatus
	TaskDeclinedStatus
)

type Event struct {
	EventId string
	TaskId  string
	Time    time.Time
	Type    EventType
	Status  TaskStatus
}

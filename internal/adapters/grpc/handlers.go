package grpc

import (
	"context"
	"fmt"

	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-analytics/pkg/analyticsgrpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) StoreEvent(ctx context.Context, eventMsg *analyticsgrpc.Event) (*emptypb.Empty, error) {
	logger := s.annotatedLogger(ctx)

	event := &models.Event{
		EventId: eventMsg.EventId,
		TaskId:  eventMsg.TaskId,
		Time:    eventMsg.GetTime().AsTime(),
	}
	switch eventMsg.Type {
	case "create":
		event.Type = models.EventCreateType
	case "update":
		event.Type = models.EventUpdateType
	case "delete":
		event.Type = models.EventDeleteType
	case "approve":
		event.Type = models.EventApproveType
	case "decline":
		event.Type = models.EventDeclineType
	default:
		return &emptypb.Empty{}, fmt.Errorf("unknown event type string")
	}
	switch eventMsg.Status {
	case "in_progress":
		event.Status = models.TaskInProgressStatus
	case "done":
		event.Status = models.TaskDoneStatus
	case "declined":
		event.Status = models.TaskDeclinedStatus
	default:
		return &emptypb.Empty{}, fmt.Errorf("unknown task status string")
	}

	err := s.analytics.StoreEvent(ctx, event)
	if err != nil {
		logger.Errorf("failed to store event")
		return &emptypb.Empty{}, fmt.Errorf("failed to store event")
	}
	return &emptypb.Empty{}, nil
}

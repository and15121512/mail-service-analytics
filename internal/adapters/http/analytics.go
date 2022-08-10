package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/utils"
)

func (s *Server) analyticsHandlers() http.Handler {
	r := chi.NewRouter()
	r.With(s.AnnotateContext()).With(s.ValidateAuth()).Get("/report", s.GetReport)

	// Here are some other endpoints...

	return r
}

func (s *Server) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return s.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
}

type reportRequest struct {
	TaskId string `json:"task_id"`
}

type reportResponse struct {
	DoneCnt           int      `json:"done_cnt"`
	DeclinedCnt       int      `json:"declined_cnt"`
	TaskId            string   `json:"task_id"`
	ReactionDurations []string `json:"reaction_durations"`
}

func (s *Server) GetReport(w http.ResponseWriter, r *http.Request) {
	logger := s.annotatedLogger(r.Context())

	rr := reportRequest{}
	err := json.NewDecoder(r.Body).Decode(&rr)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{
			"error": "cannot parse input for get report request",
		})
		logger.Errorf("cannot parse input for get report request: %s", err.Error())
		return
	}

	report, err := s.analytics.GetReport(r.Context(), rr.TaskId)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("failed to get report for task ID %s: %s", rr.TaskId, err.Error()),
		})
		logger.Errorf("failed to get report for task ID %s: %s", rr.TaskId, err.Error())
		return
	}

	reactionDurations := make([]string, len(report.ReactionDurations))
	for i, dur := range report.ReactionDurations {
		reactionDurations[i] = dur.String()
	}
	w.Header().Set("Content-Type", "application/json")
	resp, _ := json.Marshal(reportResponse{
		DoneCnt:           report.DoneCnt,
		DeclinedCnt:       report.DeclinedCnt,
		TaskId:            report.TaskId,
		ReactionDurations: reactionDurations,
	})
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

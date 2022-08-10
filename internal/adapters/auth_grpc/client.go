package auth_grpc

import (
	"context"
	"fmt"

	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/config"
	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/utils"
	"gitlab.com/sukharnikov.aa/mail-service-analytics/pkg/authgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthGrpc struct {
	logger *zap.SugaredLogger
}

func New(logger *zap.SugaredLogger) *AuthGrpc {
	return &AuthGrpc{
		logger: logger,
	}
}

func (ag *AuthGrpc) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return ag.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
}

func (ag *AuthGrpc) ValidateAuth(ctx context.Context, tokenpair *models.TokenPair) (*models.AuthResult, error) {
	logger := ag.annotatedLogger(ctx)

	auth_host := config.GetConfig(logger).Hosts.AuthHost
	auth_port := config.GetConfig(logger).Ports.GrpcPort
	target := fmt.Sprintf("%s:%s", auth_host, auth_port)
	conn, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		logger.Errorf("failed to connect to auth service: %s", err.Error())
		return &models.AuthResult{}, fmt.Errorf("failed to connect to auth service: %s", err.Error())
	}
	defer conn.Close()

	a := authgrpc.NewAuthGrpcClient(conn)

	authTokenpair := authgrpc.TokenPair{
		AccessToken:  tokenpair.AccessToken,
		RefreshToken: tokenpair.RefreshToken,
	}

	ar, err := a.Validate(ctx, &authTokenpair)
	if err != nil {
		logger.Errorf("failed to call validate from auth service: %s", err.Error())
		return &models.AuthResult{}, fmt.Errorf("failed to call validate from auth service: %s", err.Error())
	}

	logger.Infof("Login after grpc call: %s", ar.Login)
	authResult := models.AuthResult{
		Status:       -1,
		AccessToken:  ar.NewAccessToken,
		RefreshToken: ar.NewRefreshToken,
		Login:        ar.Login,
	}
	if ar.Status == "ok" {
		authResult.Status = models.OkAuthRespStatus
		logger.Infof("token validated successfully")
	} else if ar.Status == "refreshed" {
		authResult.Status = models.RefreshedAuthRespStatus
		logger.Infof("token expired and was refreshed")
	} else if ar.Status == "refused" {
		authResult.Status = models.RefusedAuthRespStatus
		logger.Infof("token validation failed")
	} else {
		logger.Errorf("unknown response status from auth service")
		return &models.AuthResult{}, fmt.Errorf("unknown response status from auth service")
	}

	return &authResult, nil
}

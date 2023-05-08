package middlewares

import (
	"context"
	"fmt"
	"github.com/segmentio/ksuid"
	"log"
	"net/http"
	"scheduler0/config"
	"scheduler0/constants/headers"
	"scheduler0/secrets"
	"scheduler0/service"
	"scheduler0/service/node"
	"scheduler0/utils"
	"strings"
	"sync"
)

// middlewareHandler middleware type
type middlewareHandler struct {
	logger           *log.Logger
	doOnce           sync.Once
	ctx              context.Context
	scheduler0Secret secrets.Scheduler0Secrets
	scheduler0Config config.Scheduler0Config
}

type MiddlewareHandler interface {
	ContextMiddleware(next http.Handler) http.Handler
	AuthMiddleware(credentialService service.Credential) func(next http.Handler) http.Handler
}

func NewMiddlewareHandler(logger *log.Logger, scheduler0Secret secrets.Scheduler0Secrets, scheduler0Config config.Scheduler0Config) *middlewareHandler {
	return &middlewareHandler{
		logger:           logger,
		scheduler0Secret: scheduler0Secret,
		scheduler0Config: scheduler0Config,
	}
}

// ContextMiddleware context middleware
func (m *middlewareHandler) ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := ksuid.New().String()
		ctx := r.Context()
		ctx = context.WithValue(ctx, "RequestID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthMiddleware authentication middleware
func (m *middlewareHandler) AuthMiddleware(credentialService service.Credential) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			paths := strings.Split(r.URL.Path, "/")

			if len(paths) < 1 {
				utils.SendJSON(w, "endpoint is not supported", false, http.StatusNotImplemented, nil)
				return
			}

			if paths[1] == "api-docs" || paths[1] == "healthcheck" {
				next.ServeHTTP(w, r)
				return
			}

			if IsServerClient(r) {
				if validity, _ := IsAuthorizedServerClient(r, credentialService); validity {
					next.ServeHTTP(w, r)
					return
				} else {
					utils.SendJSON(w, "unauthorized requests", false, http.StatusUnauthorized, nil)
					return
				}
			}

			if IsPeerClient(r) {
				if validity := IsAuthorizedPeerClient(r, m.scheduler0Secret); validity {
					next.ServeHTTP(w, r)
					return
				} else {
					utils.SendJSON(w, "unauthorized requests", false, http.StatusUnauthorized, nil)
					return
				}
			}

			utils.SendJSON(w, "unauthorized requests", false, http.StatusUnauthorized, nil)
			return
		})
	}
}

func (m *middlewareHandler) EnsureRaftLeaderMiddleware(peer *node.Node) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			paths := strings.Split(r.URL.Path, "/")

			if len(paths) < 1 {
				utils.SendJSON(w, "endpoint is not supported", false, http.StatusNotImplemented, nil)
				return
			}

			if paths[1] == "peer-handshake" {
				next.ServeHTTP(w, r)
				return
			}

			if !peer.CanAcceptRequest() {
				utils.SendJSON(w, "peer cannot accept requests", false, http.StatusServiceUnavailable, nil)
				return
			}

			fmt.Println("peer.CanAcceptClientWriteRequest()", peer.CanAcceptClientWriteRequest())

			if !peer.CanAcceptClientWriteRequest() && (r.Method == http.MethodPost || r.Method == http.MethodDelete || r.Method == http.MethodPut) {
				configs := m.scheduler0Config.GetConfigurations()
				fmt.Println("peer.FsmStore.GetRaft()", peer.FsmStore.GetRaft())
				serverAddr, _ := peer.FsmStore.GetRaft().LeaderWithID()

				fmt.Println("serverAddr", serverAddr)
				fmt.Println("configs", configs.Replicas)

				redirectUrl := ""

				for _, leaderPeer := range configs.Replicas {
					if leaderPeer.RaftAddress == string(serverAddr) {
						redirectUrl = leaderPeer.Address
						break
					}
				}

				fmt.Println("redirectUrl", redirectUrl)

				if redirectUrl == "" {
					m.logger.Println("failed to get redirect url from replicas")
					utils.SendJSON(w, "service is unavailable", false, http.StatusServiceUnavailable, nil)
					return
				}

				redirectUrl = fmt.Sprintf("%s%s", redirectUrl, r.URL.Path)

				w.Header().Set("Location", redirectUrl)
				requester := r.Header.Get(headers.PeerHeader)

				if requester == headers.PeerHeaderCMDValue || requester == headers.PeerHeaderValue {
					m.logger.Println("Redirecting request to leader", redirectUrl)
					http.Redirect(w, r, redirectUrl, 301)
				} else {
					utils.SendJSON(w, nil, false, http.StatusFound, nil)
				}

				return
			}

			next.ServeHTTP(w, r)
			return
		})
	}
}

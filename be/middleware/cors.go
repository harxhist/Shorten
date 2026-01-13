package middleware

import (
	"net/http"
	"be/constant"
    "be/response"
    "be/logger"
	"strings"

)
var log = logger.Logger;
func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        allowed := false
        for _, allowedOrigin := range constant.APPCONFIG.AllowedOrigin {
            if strings.EqualFold(origin, allowedOrigin) {
                allowed = true
                break
            }
        }
        if allowed {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        } else {
            log.Error("Origin not allowed")
            response.PutBadRequestErrorResponse("Origin not allowed", w)
            return
        }

        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Max-Age", "86400")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}
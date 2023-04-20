// Package interceptors хранит interceptors для grpc.
package interceptors

import "github.com/ImpressionableRaccoon/urlshortener/internal/authenticator"

type interceptors struct {
	a authenticator.Authenticator
}

// New - конструктор interceptors.
func New(a authenticator.Authenticator) interceptors {
	return interceptors{
		a: a,
	}
}

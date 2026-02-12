package lintertest

import (
	"log/slog"

	"go.uber.org/zap"
)

var (
	password   = "somePassword"
	token      = "someToken"
	passphrase = "somePassphrase"
)

func main() {
	slog.Info("starting server on port 8080") // OK
	slog.Info("Starting server on port 8080") // want "log message should start with a lowercase letter"

	slog.Info("starting server") // OK
	slog.Info("запуск сервера")  // want "log message should be in English only"

	slog.Info("server started")                    // OK
	slog.Info("connection failed!!!")              // want "log message contains forbidden characters or emojis"
	slog.Info("server started! 🚀")                 // want "log message contains forbidden characters or emojis"
	slog.Error("warning: something went wrong...") // want "log message contains forbidden characters or emojis"

	slog.Info("user authenticated successfully") // OK
	slog.Info("user password " + password)       // want "avoid logging sensitive variable: password"
	slog.Info("token " + token)                  // want "avoid logging sensitive variable: token"
	slog.Debug("apikey set")                     // want "log message might contain sensitive data"
	slog.Info(passphrase)                        // want "avoid logging sensitive variable: passphrase"

	var logger zap.Logger
	logger.Info("Starting")                     // want "log message should start with a lowercase letter"
	logger.Info("starting server on port 8080") // OK
	logger.Info("Starting server on port 8080") // want "log message should start with a lowercase letter"

	logger.Info("starting server") // OK
	logger.Info("запуск сервера")  // want "log message should be in English only"

	logger.Info("server started")                    // OK
	logger.Info("connection failed!!!")              // want "log message contains forbidden characters or emojis"
	logger.Info("server started! 🚀")                 // want "log message contains forbidden characters or emojis"
	logger.Error("warning: something went wrong...") // want "log message contains forbidden characters or emojis"

	logger.Info("user authenticated successfully") // OK
	logger.Info("user password " + password)       // want "avoid logging sensitive variable: password"
	logger.Info("token " + token)                  // want "avoid logging sensitive variable: token"
	logger.Info("apikey set")                      // want "log message might contain sensitive data"
	logger.Info(passphrase)                        // want "avoid logging sensitive variable: passphrase"
}

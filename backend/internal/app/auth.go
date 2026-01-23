package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/argon2"

	"apexmemory.ai/internal/db"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidToken    = errors.New("invalid token")
)

type PasswordConfig struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func DefaultPasswordConfig() *PasswordConfig {
	return &PasswordConfig{
		Memory:      64 * 1024, // 64 MB
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

func (app *Application) HashPassword(password string) (string, error) {
	config := DefaultPasswordConfig()

	salt := make([]byte, config.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, config.Iterations, config.Memory, config.Parallelism, config.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, config.Memory, config.Iterations, config.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func (app *Application) VerifyPassword(password, encodedHash string) error {
	config, salt, hash, err := app.decodeHash(encodedHash)
	if err != nil {
		return fmt.Errorf("failed to decode hash: %w", err)
	}

	otherHash := argon2.IDKey([]byte(password), salt, config.Iterations, config.Memory, config.Parallelism, config.KeyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return nil
	}

	return ErrInvalidPassword
}

func (app *Application) decodeHash(encodedHash string) (*PasswordConfig, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, errors.New("invalid encoded hash format")
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse version: %w", err)
	}

	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible argon2 version")
	}

	config := &PasswordConfig{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &config.Memory, &config.Iterations, &config.Parallelism); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode salt: %w", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode hash: %w", err)
	}

	config.SaltLength = uint32(len(salt))
	config.KeyLength = uint32(len(hash))

	return config, salt, hash, nil
}

func (app *Application) GenerateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (app *Application) HashToken(token string) []byte {
	hash := sha256.Sum256([]byte(token))
	return hash[:]
}

func (app *Application) CreateSession(ctx context.Context, userID uuid.UUID, userAgent, ipAddress string, duration time.Duration) (string, error) {
	token, err := app.GenerateSecureToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}

	tokenHash := app.HashToken(token)
	expiresAt := time.Now().Add(duration)

	var ipAddr *netip.Addr
	if ipAddress != "" {
		// Handle bracketed IPv6 with port like [::1]:8080
		if strings.HasPrefix(ipAddress, "[") {
			if idx := strings.LastIndex(ipAddress, "]:"); idx != -1 {
				ipAddress = ipAddress[1:idx]
			} else if strings.HasSuffix(ipAddress, "]") {
				ipAddress = ipAddress[1 : len(ipAddress)-1]
			}
		} else if strings.Count(ipAddress, ":") == 1 {
			// IPv4 with port like 127.0.0.1:8080
			if idx := strings.LastIndex(ipAddress, ":"); idx != -1 {
				ipAddress = ipAddress[:idx]
			}
		}
		// For IPv6 without brackets or port, parse directly
		if parsed, err := netip.ParseAddr(ipAddress); err == nil {
			ipAddr = &parsed
		}
	}

	params := db.CreateSessionParams{
		UserID:    userID,
		TokenHash: tokenHash,
		UserAgent: pgtype.Text{String: userAgent, Valid: userAgent != ""},
		IpAddress: ipAddr,
		ExpiresAt: expiresAt,
	}

	_, err = app.Queries.CreateSession(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return token, nil
}

func (app *Application) ValidateSession(ctx context.Context, token string) (*AppUser, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	tokenHash := app.HashToken(token)

	session, err := app.Queries.GetSessionByToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidToken
		}
		return nil, ErrDBQuery("get session", err)
	}

	// Update last used time (non-critical, log errors but don't fail auth)
	if time.Since(session.LastUsedAt) >= 5*time.Minute {
		if err := app.Queries.UpdateSessionLastUsed(ctx, tokenHash); err != nil {
			GetLogger(ctx).Warn("failed to update session last_used_at", "error", ErrDBQuery("update", err))
		}
	}

	return &AppUser{
		ID:       session.UserID,
		Email:    session.Email,
		Username: session.Username,
	}, nil
}

func (app *Application) DeleteSession(ctx context.Context, token string) error {
	tokenHash := app.HashToken(token)
	return app.Queries.DeleteSession(ctx, tokenHash)
}

func (app *Application) DeleteUserSessions(ctx context.Context, userID uuid.UUID) error {
	return app.Queries.DeleteUserSessions(ctx, userID)
}

func (app *Application) CleanupExpiredSessions(ctx context.Context) error {
	return app.Queries.DeleteExpiredSessions(ctx)
}

func (app *Application) RegisterUser(ctx context.Context, email, username, password string) (*AppUser, error) {
	hashedPassword, err := app.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	params := db.CreateUserWithPasswordParams{
		Email:        pgtype.Text{String: email, Valid: true},
		Username:     username,
		DisplayName:  pgtype.Text{String: username, Valid: true},
		PasswordHash: pgtype.Text{String: hashedPassword, Valid: true},
	}

	user, err := app.Queries.CreateUserWithPassword(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &AppUser{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

func (app *Application) AuthenticateUser(ctx context.Context, email, password string) (*AppUser, error) {
	user, err := app.Queries.GetUserByEmailPassword(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidPassword
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if !user.PasswordHash.Valid {
		return nil, ErrInvalidPassword
	}

	if err := app.VerifyPassword(password, user.PasswordHash.String); err != nil {
		return nil, ErrInvalidPassword
	}

	return &AppUser{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

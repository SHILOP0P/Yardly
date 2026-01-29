package auth
import(
	"context"
	"time"
	"crypto/sha256"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshRepo struct{
	pool *pgxpool.Pool
}

func NewRefreshRepo(pool *pgxpool.Pool) *RefreshRepo{
	return &RefreshRepo{pool: pool}
}

func hashToken(raw string)[]byte{
	sum:=sha256.Sum256([]byte(raw))
	return sum[:]
}

func (r *RefreshRepo) Create(ctx context.Context, userID int64, rawToken string, expiresAt time.Time) error {
	const q = `
	INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
	VALUES ($1, $2, $3)
	`
	_, err := r.pool.Exec(ctx, q, userID, hashToken(rawToken), expiresAt)
	return err
}

func (r *RefreshRepo) Consume(ctx context.Context, rawToken string, now time.Time)(int64, error){
	const q = `
	UPDATE refresh_tokens
	SET revoked_at = $2
	WHERE token_hash = $1
	AND revoked_at IS NULL
	AND expires_at > $2
	RETURNING user_id
	`
	var userID int64
	err :=r.pool.QueryRow(ctx, q, hashToken(rawToken), now).Scan(&userID)
	return userID, err
}
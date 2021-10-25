package pg

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib" // driver
	"github.com/mechta-market/limelog/internal/domain/errs"
	"github.com/mechta-market/limelog/internal/interfaces"
)

const ErrMsg = "PG-error"
const TransactionCtxKey = "pg_transaction"

type St struct {
	debug bool
	lg    interfaces.Logger

	Con *pgxpool.Pool
}

type conSt interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type txContainerSt struct {
	tx pgx.Tx
}

var (
	queryParamRegexp = regexp.MustCompile(`(?si)\$\{[^}]+\}`)
)

func New(lg interfaces.Logger, dsn string, debug bool) (*St, error) {
	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		lg.Errorw(ErrMsg+": Fail to parse dsn", err)
		return nil, err
	}

	dbConfig.MaxConns = 100
	dbConfig.MinConns = 10
	dbConfig.MaxConnLifetime = 30 * time.Minute
	dbConfig.MaxConnIdleTime = 15 * time.Minute
	dbConfig.HealthCheckPeriod = 20 * time.Second
	dbConfig.LazyConnect = true

	dbPool, err := pgxpool.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		lg.Errorw(ErrMsg+": Fail to connect to db", err)
		return nil, err
	}

	return &St{
		debug: debug,
		lg:    lg,
		Con:   dbPool,
	}, nil
}

func (d *St) handleError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	if err == pgx.ErrNoRows {
		err = errs.ObjectNotFound
		// d.lg.Errorw(ErrMsg, err)
	} else {
		d.lg.Errorw(ErrMsg, err)
	}

	return err
}

func (d *St) getCon(ctx context.Context) conSt {
	if tx := d.getContextTransaction(ctx); tx != nil {
		return tx
	}
	return d.Con
}

func (d *St) getContextTransactionContainer(ctx context.Context) *txContainerSt {
	contextV := ctx.Value(TransactionCtxKey)
	if contextV == nil {
		return nil
	}

	switch tx := contextV.(type) {
	case *txContainerSt:
		return tx
	default:
		return nil
	}
}

func (d *St) getContextTransaction(ctx context.Context) pgx.Tx {
	container := d.getContextTransactionContainer(ctx)
	if container != nil {
		return container.tx
	}

	return nil
}

func (d *St) ContextWithTransaction(ctx context.Context) (context.Context, error) {
	tx, err := d.Con.Begin(ctx)
	if err != nil {
		return ctx, d.handleError(ctx, err)
	}

	return context.WithValue(ctx, TransactionCtxKey, &txContainerSt{tx: tx}), nil
}

func (d *St) CommitContextTransaction(ctx context.Context) error {
	tx := d.getContextTransaction(ctx)
	if tx == nil {
		return nil
	}

	err := tx.Commit(ctx)
	if err != nil {
		if err != pgx.ErrTxClosed &&
			err != pgx.ErrTxCommitRollback {
			_ = tx.Rollback(ctx)

			return d.handleError(ctx, err)
		}
	}

	return nil
}

func (d *St) RollbackContextTransaction(ctx context.Context) {
	tx := d.getContextTransaction(ctx)
	if tx == nil {
		return
	}

	_ = tx.Rollback(ctx)
}

func (d *St) RenewContextTransaction(ctx context.Context) error {
	var err error

	container := d.getContextTransactionContainer(ctx)
	if container == nil {
		d.lg.Errorw(ErrMsg+": Transaction container not found in context", nil)
		return nil
	}

	if container.tx != nil {
		err = container.tx.Commit(ctx)
		if err != nil {
			if err != pgx.ErrTxClosed &&
				err != pgx.ErrTxCommitRollback {
				_ = container.tx.Rollback(ctx)

				return d.handleError(ctx, err)
			}
		}
	}

	container.tx, err = d.Con.Begin(ctx)
	if err != nil {
		return d.handleError(ctx, err)
	}

	return nil
}

func (d *St) DbExec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return d.getCon(ctx).Exec(ctx, sql, args...)
}

func (d *St) DbQuery(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return d.getCon(ctx).Query(ctx, sql, args...)
}

func (d *St) DbQueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return d.getCon(ctx).QueryRow(ctx, sql, args...)
}

func (d *St) queryRebindNamed(sql string, argMap map[string]interface{}) (string, []interface{}) {
	resultQuery := sql
	args := make([]interface{}, 0, len(argMap))

	for k, v := range argMap {
		if strings.Contains(resultQuery, "${"+k+"}") {
			args = append(args, v)
			resultQuery = strings.ReplaceAll(resultQuery, "${"+k+"}", "$"+strconv.Itoa(len(args)))
		}
	}

	if d.debug {
		if strings.Index(resultQuery, "${") > -1 {
			for _, x := range queryParamRegexp.FindAllString(resultQuery, 1) {
				d.lg.Errorw("Missing param", nil, "param", x, "query", resultQuery)
			}
		}
	}

	return resultQuery, args
}

func (d *St) DbExecM(ctx context.Context, sql string, argMap map[string]interface{}) (pgconn.CommandTag, error) {
	rbSql, args := d.queryRebindNamed(sql, argMap)

	return d.getCon(ctx).Exec(ctx, rbSql, args...)
}

func (d *St) DbQueryM(ctx context.Context, sql string, argMap map[string]interface{}) (pgx.Rows, error) {
	rbSql, args := d.queryRebindNamed(sql, argMap)

	return d.getCon(ctx).Query(ctx, rbSql, args...)
}

func (d *St) DbQueryRowM(ctx context.Context, sql string, argMap map[string]interface{}) pgx.Row {
	rbSql, args := d.queryRebindNamed(sql, argMap)

	return d.getCon(ctx).QueryRow(ctx, rbSql, args...)
}

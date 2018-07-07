package configManager

import (
	"configManager/neo4j"
	"math/rand"
	"time"

	"github.com/Sirupsen/logrus"
	app "github.com/casualjim/go-app"
	"github.com/casualjim/go-app/tracing"
	"github.com/oklog/ulid"
	"github.com/spf13/viper"
)

// NewRuntime creates a new application level runtime that encapsulates the shared services for this application
func NewRuntime(app app.Application) (*Runtime, error) {

	/*
	 * Allocate DB Connection Pool
	 */
	pool, err := neo4j.Pool(app.Config().GetString("db.Host"), app.Config().GetString("db.Port"), "", "", app.Config().GetInt("db.MaxConn"))
	if err != nil {
		return nil, err
	}

	return &Runtime{
		pool: pool,
		app:  app,
	}, nil
}

// Runtime encapsulates the shared services for this application
type Runtime struct {
	pool neo4j.ConnPool
	app  app.Application
}

// DB returns the persistent store
func (r *Runtime) DB() neo4j.ConnPool {
	return r.pool
}

// Tracer returns the root tracer, this is typically the only one you need
func (r *Runtime) Tracer() tracing.Tracer {
	return r.app.Tracer()
}

// Logger gets the root logger for this application
func (r *Runtime) Logger() logrus.FieldLogger {
	return r.app.Logger()
}

// NewLogger creates a new named logger for this application
func (r *Runtime) NewLogger(name string, fields logrus.Fields) logrus.FieldLogger {
	return r.app.NewLogger(name, fields)
}

// Config returns the viper config for this application
func (r *Runtime) Config() *viper.Viper {
	return r.app.Config()
}

func (r *Runtime) ExecDB(query *string, params *map[string]interface{}) (neo4j.Result, error) {
	return neo4j.Exec(r.DB(), query, params)
}

func (r *Runtime) QueryDB(query *string, params *map[string]interface{}) ([]interface{}, error) {
	return neo4j.Query(r.DB(), query, params)
}

func (r *Runtime) QueryAllDB(query *string, params *map[string]interface{}) ([][]interface{}, error) {
	return neo4j.QueryAll(r.DB(), query, params)
}

func GetULID() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))

	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

package dbresolver

import (
	"database/sql"
	"unsafe"
)

// New will resolve all the passed connection with configurable parameters
func New(opts ...OptionFunc) DB {
	opt := defaultOption()
	for _, optFunc := range opts {
		optFunc(opt)
	}

	if len(opt.PrimaryDBs) == 0 {
		panic("required primary db connection, set the primary db " +
			"connection with dbresolver.New(dbresolver.WithPrimaryDBs(primaryDB))")
	}
	return &sqlDB{
		primaries:        opt.PrimaryDBs,
		replicas:         opt.ReplicaDBs,
		loadBalancer:     opt.DBLB,
		stmtLoadBalancer: opt.StmtLB,
		DB:               opt.PrimaryDBs[0],
	}
}

func NewSqlDB(opts ...OptionFunc) *sql.DB {
	resolver := New(opts...).(*sqlDB)
	return (*sql.DB)(unsafe.Pointer(resolver))
}

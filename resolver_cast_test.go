package dbresolver

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"unsafe"
)

func TestNewSqlDB(t *testing.T) {
	noOfPrimaries, noOfReplicas := 1, 1
	loadBalancerPolicy := RoundRobinLB

	primaries := make([]*sql.DB, noOfPrimaries)
	replicas := make([]*sql.DB, noOfReplicas)

	mockPimaries := make([]sqlmock.Sqlmock, noOfPrimaries)
	mockReplicas := make([]sqlmock.Sqlmock, noOfReplicas)

	for i := 0; i < noOfPrimaries; i++ {
		db, mock, err := createMock()

		if err != nil {
			t.Fatal("creating of mock failed")
		}

		defer mock.ExpectClose()
		defer db.Close()

		primaries[i] = db
		mockPimaries[i] = mock
	}

	for i := 0; i < noOfReplicas; i++ {
		db, mock, err := createMock()
		if err != nil {
			t.Fatal("creating of mock failed")
		}

		defer mock.ExpectClose()
		defer db.Close()

		replicas[i] = db
		mockReplicas[i] = mock
	}

	resolver := NewSqlDB(WithPrimaryDBs(primaries...), WithReplicaDBs(replicas...), WithLoadBalancer(loadBalancerPolicy))

	internalResolver := (*sqlDB)(unsafe.Pointer(resolver))

	internalResolver.Exec("select 1")

	t.Log("executed internal resolver")

	resolver.Exec("select 1")

	/*t.Run("primary dbs", func(t *testing.T) {
		for i := 0; i < noOfPrimaries*5; i++ {
			robin := internalResolver.loadBalancer.predict(noOfPrimaries)
			mock := mockPimaries[robin]

			t.Log("case - ", i%4)

			switch i % 4 {
			case 0:
				query := "SET timezone TO 'Asia/Tokyo'"
				mock.ExpectExec(query)
				resolver.Exec(query)
				t.Log("exec")
			case 1:
				query := "SET timezone TO 'Asia/Tokyo'"
				mock.ExpectExec(query)
				resolver.ExecContext(context.TODO(), query)
				t.Log("exec context")
			case 2:
				mock.ExpectBegin()
				resolver.Begin()
				t.Log("begin")
			case 3:
				mock.ExpectBegin()
				resolver.BeginTx(context.TODO(), &sql.TxOptions{
					Isolation: sql.LevelDefault,
					ReadOnly:  false,
				})
				t.Log("begin transaction")
			default:
				t.Fatal("developer needs to work on the tests")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		}
	})*/
}

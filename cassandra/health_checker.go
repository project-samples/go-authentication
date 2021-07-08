package cassandra

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

type HealthChecker struct {
	Cassandra *gocql.ClusterConfig
	name      string
	timeout   time.Duration
}

func NewMongoHealthChecker(db *gocql.ClusterConfig, name string, timeouts ...time.Duration) *HealthChecker {
	var timeout time.Duration
	if len(timeouts) >= 1 {
		timeout = timeouts[0]
	} else {
		timeout = 4 * time.Second
	}
	return &HealthChecker{Cassandra: db, name: name, timeout: timeout}
}
func NewHealthChecker(db *gocql.ClusterConfig, options ...string) *HealthChecker {
	var name string
	if len(options) >= 1 && len(options[0]) > 0 {
		name = options[0]
	} else {
		name = "cassandra"
	}
	return NewMongoHealthChecker(db, name, 4*time.Second)
}

func (s *HealthChecker) Name() string {
	return s.name
}

func (s *HealthChecker) Check(ctx context.Context) (map[string]interface{}, error) {
	cancel := func() {}
	if s.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, s.timeout)
	}
	defer cancel()

	res := make(map[string]interface{})
	checkerChan := make(chan error)
	log.Println("err")
	go func() {
		// checkerChan <- s.db.RunCommand(ctx, bsonx.Doc{{"ping", bsonx.Int32(1)}}).Decode(&info)
		_, err := s.Cassandra.CreateSession()
		checkerChan <- err
	}()
	select {
	case err := <-checkerChan:
		return res, err
	case <-ctx.Done():
		return res, fmt.Errorf("timeout")
	}
}

func (s *HealthChecker) Build(ctx context.Context, data map[string]interface{}, err error) map[string]interface{} {
	if err == nil {
		return data
	}
	if data == nil {
		data = make(map[string]interface{}, 0)
	}
	data["error"] = err.Error()
	return data
}

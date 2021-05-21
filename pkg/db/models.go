// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
)

type JobResult string

const (
	JobResultPassed  JobResult = "passed"
	JobResultFailed  JobResult = "failed"
	JobResultAborted JobResult = "aborted"
)

func (e *JobResult) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = JobResult(s)
	case string:
		*e = JobResult(s)
	default:
		return fmt.Errorf("unsupported scan type for JobResult: %T", src)
	}
	return nil
}

type TestResult string

const (
	TestResultPassed  TestResult = "passed"
	TestResultFailure TestResult = "failure"
	TestResultSkipped TestResult = "skipped"
	TestResultError   TestResult = "error"
)

func (e *TestResult) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TestResult(s)
	case string:
		*e = TestResult(s)
	default:
		return fmt.Errorf("unsupported scan type for TestResult: %T", src)
	}
	return nil
}

type Job struct {
	ID                 int64         `json:"id"`
	Provider           string        `json:"provider"`
	JobName            string        `json:"job_name"`
	JobID              string        `json:"job_id"`
	Url                string        `json:"url"`
	Started            time.Time     `json:"started"`
	Finished           time.Time     `json:"finished"`
	Duration           sql.NullInt64 `json:"duration"`
	ClusterVersion     string        `json:"cluster_version"`
	ClusterName        string        `json:"cluster_name"`
	ClusterID          string        `json:"cluster_id"`
	MultiAz            string        `json:"multi_az"`
	Channel            string        `json:"channel"`
	Environment        string        `json:"environment"`
	Region             string        `json:"region"`
	NumbWorkerNodes    int32         `json:"numb_worker_nodes"`
	NetworkProvider    string        `json:"network_provider"`
	ImageContentSource string        `json:"image_content_source"`
	InstallConfig      string        `json:"install_config"`
	HibernateAfterUse  bool          `json:"hibernate_after_use"`
	Reused             bool          `json:"reused"`
	Result             JobResult     `json:"result"`
}

type Testcase struct {
	ID       int64           `json:"id"`
	JobID    int64           `json:"job_id"`
	Result   TestResult      `json:"result"`
	Name     string          `json:"name"`
	Duration pgtype.Interval `json:"duration"`
	Error    string          `json:"error"`
	Stdout   string          `json:"stdout"`
	Stderr   string          `json:"stderr"`
}

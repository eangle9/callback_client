package callbackclient

type StatusBool string

const (
	StatusBoolYes StatusBool = "YES"
	StatusBoolNO  StatusBool = "NO"
)

type Status string

const (
	StatusPending    Status = "PENDING"
	StatusActive     Status = "ACTIVE"
	StatusInactive   Status = "INACTIVE"
	StatusFailed     Status = "FAILED"
	StatusSucceeded  Status = "SUCCEEDED"
	StatusProcessing Status = "PROCESSING"
)

type Method string

const (
	MethodPost   Method = "POST"
	MethodGet    Method = "GET"
	MethodPut    Method = "PUT"
	MethodPatch  Method = "PATCH"
	MethodDelete Method = "DELETE"
)

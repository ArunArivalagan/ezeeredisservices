package constants

const (
	Success      ErrorCode = 200
	Unauthorized ErrorCode = 401

	RedisClientIsNotConnected ErrorCode = 100
	InvalidKey                ErrorCode = 101
	UnableToProvideData       ErrorCode = 102
	UpdateFailed              ErrorCode = 103
	InvalidRequest            ErrorCode = 104
)

type ErrorCode int

func (e ErrorCode) Error() string {
	return errorMessages[e]
}

var errorMessages = map[ErrorCode]string{
	Success:                   "Success",
	Unauthorized:              "The request did not include an authentication token",
	RedisClientIsNotConnected: "Unable to connect redis client",
	InvalidKey:                "Invalid cache key",
	UnableToProvideData:       "Unable to provide data",
	UpdateFailed:              "Update failed",
	InvalidRequest:            "Invalid Request",
}

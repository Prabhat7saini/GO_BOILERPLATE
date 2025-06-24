package exception

import common "boiler-platecode/src/common/const"

type ErrorCode string

var ErrorCodeErrorMessage = map[ErrorCode]common.Exception{
	USER_NOT_FOUND: {
		Code:           "UNF000",
		Message:        "User Not Found",
		HttpStatusCode: 404,
	},
	USER_ALREADY_EXISTS: {
		Code:           "UNF000",
		Message:        "User Already Exists",
		HttpStatusCode: 409,
	},
	INTERNAL_SERVER_ERROR: {
		Code:           "ISE000",
		Message:        "Internal Server Error",
		HttpStatusCode: 500,
	},
	INVALID_CREDENTIALS: {
		Code:           "IVC000",
		Message:        "Invalid Credentials",
		HttpStatusCode: 401,
	},
}

func GetException(code ErrorCode) *common.Exception {
	if ex, ok := ErrorCodeErrorMessage[code]; ok {
		return &ex
	}
	return &common.Exception{
		Code:           "UNKNOWN",
		Message:        "Unknown error code",
		HttpStatusCode: 500,
	}
}

package ecode

type Error struct {
	Code int32  `json:"code"`
	Err  string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

var (
	Success = Error{
		Code: 0,
		Err:  "",
	}
	ErrParams = Error{
		Code: 2000,
		Err:  "err params",
	}
	ErrUser = Error{
		Code: 3000,
		Err:  "err user operation",
	}
	ErrSystem = Error{
		Code: 5000,
		Err:  "err system",
	}
)

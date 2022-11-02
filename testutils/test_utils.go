package testutils

type TestType struct {
	Val string
}

func (t TestType) Default() TestType {
	return TestType{Val: "test"}
}

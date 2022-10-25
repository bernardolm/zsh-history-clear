package executor

type Executor interface {
	Get() []byte
	Put([]byte)
}

package proxy

type Response struct {
	Body   []byte
	Status int
	Err    error
	Timer  Timer
}

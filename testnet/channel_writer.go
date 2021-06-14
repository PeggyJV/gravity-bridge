package main

type ChannelWriter struct {
	Chan chan string
}

func (cw *ChannelWriter) Write(p []byte) (n int, err error) {
	cw.Chan <- string(p)
	return len(p), nil
}
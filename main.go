package main

func main() {
	pusher := NewPusher(Config{Host: "0.0.0.0", Port: 3001})

	pusher.Serve()
}

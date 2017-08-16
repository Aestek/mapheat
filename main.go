package main

func main() {
	server("127.0.0.1:9000", Agg(randomSource()))
}

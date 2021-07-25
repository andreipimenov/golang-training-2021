package main

func main() {
	var server StockServer
	server.Run(8080)
	server.WaitShutdown()
}

package main

func main() {

	cli := CLI{}
	cli.Run()

	// go run blockchain/*.go createblockchain -address A
	// go run blockchain/*.go send -from B -to C -amount 6
}

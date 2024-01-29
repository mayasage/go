package main

func testDeadlock() {
	a := x{
		val: 1,
	}

	b := x{
		val: 2,
	}

	wg.Add(2)
	go callmebaby(&a, &b)
	go callmebaby(&b, &a)
	wg.Wait()
}

func testStarvation() {
	wg.Add(2)
	go greedyWorker()
	go politeWorker()
	wg.Wait()
}

func main() {
	testStarvation()
}

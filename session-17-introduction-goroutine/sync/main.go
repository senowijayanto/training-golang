package main

import (
	"fmt"
	"time"
)

func worker(id int, done chan bool) {
	fmt.Printf("Pekerja %d sedang bekerja...\n", id)
	time.Sleep(2 * time.Second)
	fmt.Printf("Pekerja %d selesai bekerja\n", id)
	done <- true
}

func main() {
	done := make(chan bool, 2)

	go worker(1, done)
	go worker(2, done)

	<-done // Menunggu pekerja 1 selesai
	<-done // Menunggu pekerja 2 selesai

	fmt.Println("Semua pekerja selesai.")
}

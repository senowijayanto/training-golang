package main

import (
	"fmt"
	"time"
)

// Fungsi untuk memproses pesanan
func processOrder(orderID int, done chan bool) {
	fmt.Printf("Memproses pesanan #%d...\n", orderID)
	time.Sleep(2 * time.Second) // Simulasi waktu pemrosesan
	fmt.Printf("Pesanan #%d selesai diproses.\n", orderID)
	done <- true // Kirim sinyal selesai ke channel
}

func main() {
	// Channel untuk menunggu sinyal dari Go-Routine
	done := make(chan bool)

	// Menjalankan beberapa Go-Routine untuk memproses pesanan
	for i := 1; i <= 3; i++ {
		go processOrder(i, done)
	}

	// Menunggu semua pesanan selesai diproses
	for i := 1; i <= 3; i++ {
		<-done
	}
	fmt.Println("Semua pesanan selesai diproses.")
}

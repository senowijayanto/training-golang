package main

import (
	"fmt"
	"runtime"
)

func main() {
	// Mendapatkan jumlah core CPU logis
	numCPU := runtime.NumCPU()
	fmt.Printf("Jumlah core CPU logis: %d\n", numCPU)

	// Mengatur GOMAXPROCS dengan jumlah core CPU logis
	runtime.GOMAXPROCS(numCPU)

	// Menampilkan nilai GOMAXPROCS yang saat ini digunakan
	fmt.Printf("Jumlah core CPU yang digunakan oleh Go: %d\n", runtime.GOMAXPROCS(0))
}

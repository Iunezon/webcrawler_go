package main

import "fmt"

func printReport(pages map[string]int, baseURL string) {
	fmt.Println("=============================")
	fmt.Println("REPORT for", baseURL)
	fmt.Println("=============================")
	for k, v := range pages {
		fmt.Printf("Found %v internal links to %s \n", v, k)
	}
}

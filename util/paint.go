package util

import "fmt"

func Paint(version string) {
	fmt.Print("\n   __ _             ___ _           _   \n")
	fmt.Println("  / /(_)_   _____  / __\\ |__   __ _| |_ ")
	fmt.Println(" / / | \\ \\ / / _ \\/ /  | '_ \\ / _` | __|")
	fmt.Println("/ /__| |\\ V /  __/ /___| | | | (_| | |_ ")
	fmt.Println("\\____/_| \\_/ \\___\\____/|_| |_|\\__,_|\\__|")
	fmt.Printf("LiveChat %s\n\n", version)
}

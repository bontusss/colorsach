/*

Written by Ikwchegh Ukandu <ikwecheghu@gmail.com>
Date: 9th september 2023

*/

package main

import (
	"github.com/bontusss/colosach/controllers"
	"log"
)

func main() {

	err := controllers.Start(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

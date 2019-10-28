package main

import (
	"fmt"

	wpa "github.com/nlepage/go-wpa"
)

func main() {
	w, err := wpa.SystemWPA()
	if err != nil {
		panic(err)
	}
	defer w.Close()

	ifaces, err := w.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iface := range ifaces {
		ifname, err := iface.Ifname()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s:\n", ifname)

		nets, err := iface.Networks()
		if err != nil {
			panic(err)
		}

		for _, net := range nets {
			props, err := net.Properties()
			if err != nil {
				panic(err)
			}

			fmt.Printf("\t%#v\n", props)
		}
	}
}

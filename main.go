package main

import (
	"golang.org/x/sys/windows/registry"
	"log"
	"os/exec"
	"time"
)

func enableLoopback(appIDs []string) {
	for _, id := range appIDs {
		cmd := exec.Command("CheckNetIsolation", "loopbackexempt", "-a", "-p=" + id)
		err := cmd.Run()
		if err != nil {
			log.Printf("Cmd exec failed: %s", err)
		}
	}
}

func main() {

	log.Printf("Loopback Breaker is running...")

	rate := time.Second * 2

	ticker := time.Tick(rate)

	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Classes\Local Settings\Software\Microsoft\Windows\CurrentVersion\AppContainer\Mappings`, registry.READ)
	if err != nil {
		log.Printf("%s", err)
	}
	defer k.Close()

	for i:=0; true; i++ {
		<- ticker
		// log.Printf("Checking...")

		stat, err := k.Stat()

		if i > 0 && (err != nil || time.Since(stat.ModTime()) > rate) {
			continue
		}

		log.Printf("Detacted new UWP!")
	
		appIDs, err := k.ReadSubKeyNames(0)
		if err != nil {
			log.Printf("%s", err)
		}
	
		go enableLoopback(appIDs)

	}

}

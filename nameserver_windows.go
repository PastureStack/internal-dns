package main

import (
	"bufio"
	"bytes"
	"os/exec"

	log "github.com/PastureStack/internal-dns/internal/logging"
)

func getGlobalRecurse() ([]string, error) {
	var recurse []string
	out, err := exec.Command("powershell.exe", "(Get-NetAdapter | Get-DnsClientServerAddress -AddressFamily IPv4).ServerAddresses").Output()
	if err != nil {
		return recurse, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		ns := scanner.Text()
		if invalidRecurse(ns) {
			continue
		}
		recurse = append(recurse, ns)
	}

	log.Infof("recurse: %v", recurse)
	if len(recurse) == 0 {
		return fallbackRecurse, nil
	}

	return recurse, nil
}

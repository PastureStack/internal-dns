package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/PastureStack/internal-dns/internal/logging"
	yaml "gopkg.in/yaml.v3"
)

const maxAnswersFileBytes = 64 << 20

func ParseAnswers(path string) (out Answers, err error) {
	out = make(Answers)
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warn("Failed to find: ", path)
			return out, nil
		}
		return nil, err
	}
	if info.Size() > maxAnswersFileBytes {
		return nil, fmt.Errorf("answers file exceeds %d bytes", maxAnswersFileBytes)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if yaml.Unmarshal(data, &out); err != nil {
		return nil, err
	}

	ConvertPtrIps(&out)
	return out, nil
}

func ConvertPtrIps(answers *Answers) {
	// Convert PTR keys that are IP addresses into "4.3.2.1.in-addr.arpa." form.
	for _, client := range *answers {
		for origKey, val := range client.Ptr {
			if !strings.HasSuffix(origKey, "in-addr.arpa.") {
				newKey := "in-addr.arpa."
				for _, i := range strings.Split(origKey, ".") {
					newKey = i + "." + newKey
				}

				delete(client.Ptr, origKey)
				client.Ptr[newKey] = val
				log.Debug("Transformed PTR for ", origKey, " to ", newKey, " => ", val.Answer)
			}
		}
	}
}

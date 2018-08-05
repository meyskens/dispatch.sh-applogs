package main

import (
	"bufio"
	"os"
	"time"

	"k8s.io/api/core/v1"
)

var since int64 = 1

func follow(name, pod, container string) {
	client, err := newKubernetesClient()
	if err != nil {
		panic(err)
	}

	req := client.CoreV1().Pods(os.Getenv("MY_NAMESPACE")).GetLogs(pod, &v1.PodLogOptions{
		Container:    container,
		Follow:       true,
		Timestamps:   true,
		SinceSeconds: &since,
	})

	readCloser, err := req.Stream()
	if err != nil {
		panic(err)
	}
	defer readCloser.Close()

	scanner := bufio.NewScanner(readCloser)
	for scanner.Scan() {
		sendToDB(logEntry{
			InternalName: name,
			Pod:          pod,
			Container:    container,
			Time:         time.Now(),
			Line:         scanner.Text(),
		})
	}
}

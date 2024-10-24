package main

import (
	"os"
	"runtime"

	"github.com/google/uuid"
	"github.com/grafana/pyroscope-go"
)

func initPyroscope() {
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	pyroscope.Start(pyroscope.Config{
		ApplicationName: appName,
		ServerAddress:   "http://localhost:4040",
		// Optional HTTP Basic authentication
		// BasicAuthUser:     "<User>",     // 900009
		// BasicAuthPassword: "<Password>", // glc_SAMPLEAPIKEY0000000000==

		// Logger: pyroscope.StandardLogger,
		Tags: map[string]string{
			"hostname": os.Getenv("HOSTNAME"),
			"instance": uuid.New().String(),
		},
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
}

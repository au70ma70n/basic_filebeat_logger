package my_logger

import (
	"os/exec"

	"github.com/MythicMeta/MythicContainer/loggingstructs"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

func Initialize() {
	myLoggerName := "my_logger"
	myLogger := loggingstructs.LoggingDefinition{
		Name:           myLoggerName,
		Description:    "basic filebeat JSON logger for RedELK",
		LogToFilePath:  "/var/log/mythic/mythic.log",
		LogLevel:       "debug",
		LogMaxSizeInMB: 20,
		LogMaxBackups:  10,
		OnContainerStartFunction: func(input sharedStructs.ContainerOnStartMessage) sharedStructs.ContainerOnStartMessageResponse {
			// Start Filebeat in the background
			loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo("Starting Filebeat from Go...")

			// Execute filebeat command in background
			go func() {
				cmd := exec.Command("/usr/share/filebeat/filebeat", "-c", "/Mythic/filebeat_mythic_redelk.yml")
				cmd.Dir = "/Mythic"

				// Start the command
				err := cmd.Start()
				if err != nil {
					loggingstructs.AllLoggingData.Get(myLoggerName).LogError(err, "Failed to start Filebeat")
					return
				}

				loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo("Filebeat started successfully", "pid", cmd.Process.Pid)

				// Wait for the process to exit
				err = cmd.Wait()
				if err != nil {
					loggingstructs.AllLoggingData.Get(myLoggerName).LogError(err, "Filebeat process exited with error")
				} else {
					loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo("Filebeat process exited normally")
				}
			}()

			return sharedStructs.ContainerOnStartMessageResponse{}
		},
		NewCallbackFunction: func(input loggingstructs.NewCallbackLog) {
			loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo(input.Action, "data", input)
		},
		NewTaskFunction: func(input loggingstructs.NewTaskLog) {
			loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo(input.Action, "data", input.Data)
		},
		NewPayloadFunction: func(input loggingstructs.NewPayloadLog) {
			loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo(input.Action, "data", input.Data)
		},
		NewKeylogFunction: func(input loggingstructs.NewKeylogLog) {
			loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo(input.Action, "data", input.Data)
		},
		NewCredentialFunction: func(input loggingstructs.NewCredentialLog) {
			loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo(input.Action, "data", input.Data)
		},
		NewArtifactFunction: func(input loggingstructs.NewArtifactLog) {
			loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo(input.Action, "data", input.Data)
		},
		NewFileFunction: func(input loggingstructs.NewFileLog) {
			loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo(input.Action, "data", input.Data)
		},
		NewResponseFunction: func(input loggingstructs.NewResponseLog) {
			loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo(input.Action, "data", input.Data)
		},
	}
	loggingstructs.AllLoggingData.Get(myLoggerName).AddLoggingDefinition(myLogger)
}

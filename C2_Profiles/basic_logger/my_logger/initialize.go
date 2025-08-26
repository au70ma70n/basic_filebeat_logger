package my_logger

import (
	"fmt"
	"os"
	"os/exec"
	"time"

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
			// Write to debug file to verify function is running
			debugFile, err := os.OpenFile("/var/log/mythic/container_start_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				fmt.Fprintf(debugFile, "[%s] OnContainerStartFunction called for container: %s, operation: %d\n", 
					time.Now().Format("2006-01-02 15:04:05"), input.ContainerName, input.OperationID)
				debugFile.Close()
			}
			
			// Start Filebeat in the background
			loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo("Starting Filebeat from Go...")
			
			// Execute filebeat command in background
			go func() {
				// Write to debug file
				debugFile, err = os.OpenFile("/var/log/mythic/container_start_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err == nil {
					fmt.Fprintf(debugFile, "[%s] Starting filebeat goroutine\n", time.Now().Format("2006-01-02 15:04:05"))
					debugFile.Close()
				}
				
				// First check if filebeat is already running and kill it
				checkCmd := exec.Command("pgrep", "filebeat")
				if err := checkCmd.Run(); err == nil {
					// filebeat is running, kill it
					loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo("Filebeat already running, killing existing process...")
					killCmd := exec.Command("pkill", "filebeat")
					if err := killCmd.Run(); err != nil {
						loggingstructs.AllLoggingData.Get(myLoggerName).LogError(err, "Failed to kill existing filebeat process")
					} else {
						loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo("Successfully killed existing filebeat process")
					}
					// Give it a moment to fully terminate
					time.Sleep(2 * time.Second)
				}
				
				cmd := exec.Command("/usr/share/filebeat/filebeat", "-c", "/Mythic/filebeat_mythic_redelk.yml")
				cmd.Dir = "/Mythic"
				
				// Write to debug file before starting
				debugFile, err = os.OpenFile("/var/log/mythic/container_start_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err == nil {
					fmt.Fprintf(debugFile, "[%s] About to start filebeat command\n", time.Now().Format("2006-01-02 15:04:05"))
					debugFile.Close()
				}
				
				// Start the command
				err = cmd.Start()
				if err != nil {
					// Write error to debug file
					debugFile, err2 := os.OpenFile("/var/log/mythic/container_start_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err2 == nil {
						fmt.Fprintf(debugFile, "[%s] ERROR: Failed to start filebeat: %v\n", time.Now().Format("2006-01-02 15:04:05"), err)
						debugFile.Close()
					}
					loggingstructs.AllLoggingData.Get(myLoggerName).LogError(err, "Failed to start Filebeat")
					return
				}
				
				// Write success to debug file
				debugFile, err2 := os.OpenFile("/var/log/mythic/container_start_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err2 == nil {
					fmt.Fprintf(debugFile, "[%s] SUCCESS: Filebeat started with PID %d\n", time.Now().Format("2006-01-02 15:04:05"), cmd.Process.Pid)
					debugFile.Close()
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
			
			return sharedStructs.ContainerOnStartMessageResponse{
				ContainerName:        input.ContainerName,
				EventLogInfoMessage:  "Filebeat startup initiated successfully",
				EventLogErrorMessage: "",
			}
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

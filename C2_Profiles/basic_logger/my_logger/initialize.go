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
	myLoggerName := "my_new_logger"
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

			// Start filebeat in background
			go func() {
				// Write to debug file
				debugFile, err := os.OpenFile("/var/log/mythic/container_start_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err == nil {
					fmt.Fprintf(debugFile, "[%s] Starting filebeat process\n", time.Now().Format("2006-01-02 15:04:05"))
					debugFile.Close()
				}

				// First check if filebeat is already running and kill it
				checkCmd := exec.Command("pgrep", "filebeat")
				if err := checkCmd.Run(); err == nil {
					// filebeat is running, kill it
					debugFile, err := os.OpenFile("/var/log/mythic/container_start_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err == nil {
						fmt.Fprintf(debugFile, "[%s] Filebeat already running, killing existing process\n", time.Now().Format("2006-01-02 15:04:05"))
						debugFile.Close()
					}

					killCmd := exec.Command("pkill", "filebeat")
					if err := killCmd.Run(); err != nil {
						debugFile, err := os.OpenFile("/var/log/mythic/container_start_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if err == nil {
							fmt.Fprintf(debugFile, "[%s] ERROR: Failed to kill existing filebeat process: %v\n", time.Now().Format("2006-01-02 15:04:05"), err)
							debugFile.Close()
						}
						loggingstructs.AllLoggingData.Get(myLoggerName).LogError(err, "Failed to kill existing filebeat process")
					} else {
						debugFile, err := os.OpenFile("/var/log/mythic/container_start_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if err == nil {
							fmt.Fprintf(debugFile, "[%s] Successfully killed existing filebeat process\n", time.Now().Format("2006-01-02 15:04:05"))
							debugFile.Close()
						}
						loggingstructs.AllLoggingData.Get(myLoggerName).LogInfo("Successfully killed existing filebeat process")
					}
					// Give it a moment to fully terminate
					time.Sleep(2 * time.Second)
				}

				// Check if filebeat exists at the expected location
				filebeatPath := "/usr/share/filebeat/filebeat"
				if _, err := os.Stat(filebeatPath); os.IsNotExist(err) {
					debugFile, err := os.OpenFile("/var/log/mythic/container_start_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err == nil {
						fmt.Fprintf(debugFile, "[%s] ERROR: filebeat not found at %s: %v\n", time.Now().Format("2006-01-02 15:04:05"), filebeatPath, err)
						debugFile.Close()
					}
					loggingstructs.AllLoggingData.Get(myLoggerName).LogError(err, "filebeat not found at expected location")
					return
				}

				// Check if filebeat is executable
				if info, err := os.Stat(filebeatPath); err == nil {
					if info.Mode()&0111 == 0 {
						debugFile, err := os.OpenFile("/var/log/mythic/container_start_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if err == nil {
							fmt.Fprintf(debugFile, "[%s] ERROR: filebeat at %s is not executable\n", time.Now().Format("2006-01-02 15:04:05"), filebeatPath)
							debugFile.Close()
						}
						loggingstructs.AllLoggingData.Get(myLoggerName).LogError(fmt.Errorf("filebeat not executable"), "filebeat not executable")
						return
					}
				}

				debugFile, err = os.OpenFile("/var/log/mythic/container_start_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err == nil {
					fmt.Fprintf(debugFile, "[%s] Found filebeat at: %s\n", time.Now().Format("2006-01-02 15:04:05"), filebeatPath)
					debugFile.Close()
				}

				// List contents of /usr/share/filebeat/ directory for debugging
				if files, err := os.ReadDir("/usr/share/filebeat/"); err == nil {
					debugFile, err := os.OpenFile("/var/log/mythic/container_start_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err == nil {
						fmt.Fprintf(debugFile, "[%s] Contents of /usr/share/filebeat/: ", time.Now().Format("2006-01-02 15:04:05"))
						for _, file := range files {
							fmt.Fprintf(debugFile, "%s ", file.Name())
						}
						fmt.Fprintf(debugFile, "\n")
						debugFile.Close()
					}
				}

				// Check if config file exists
				if _, err := os.Stat("/Mythic/filebeat_mythic_redelk.yml"); os.IsNotExist(err) {
					debugFile, err = os.OpenFile("/var/log/mythic/container_start_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err == nil {
						fmt.Fprintf(debugFile, "[%s] ERROR: Config file /Mythic/filebeat_mythic_redelk.yml does not exist\n", time.Now().Format("2006-01-02 15:04:05"))
						debugFile.Close()
					}
					loggingstructs.AllLoggingData.Get(myLoggerName).LogError(err, "Config file does not exist")
					return
				}

				debugFile, err = os.OpenFile("/var/log/mythic/container_start_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err == nil {
					fmt.Fprintf(debugFile, "[%s] Config file exists, starting filebeat\n", time.Now().Format("2006-01-02 15:04:05"))
					debugFile.Close()
				}

				// Start filebeat command with full path
				cmd := exec.Command(filebeatPath, "-c", "/Mythic/filebeat_mythic_redelk.yml")
				cmd.Dir = "/Mythic"

				// Write to debug file before starting
				debugFile, err = os.OpenFile("/var/log/mythic/container_start_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err == nil {
					fmt.Fprintf(debugFile, "[%s] About to start filebeat command with path: %s\n", time.Now().Format("2006-01-02 15:04:05"), filebeatPath)
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

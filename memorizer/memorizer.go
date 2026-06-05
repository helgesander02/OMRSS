package memorizer

import (
	"os"
	"src/pkg/logger"
	"src/plan"
)

type Memorizers interface {
	MCumulative(plan.Plans)
	MAverage(int)
	MOutputResults()
	MStoreData(string, int)
	MStoreFile(string)
}

func NewMemorizers() map[string]Memorizers {
	// memorizer1 ...
	OMACO := newOmacoMemorizer()
	OSRO := newOsroMemorizer()

	// Look-up table method
	memorizers := map[string]Memorizers{
		"omaco": OMACO,
		"osro":  OSRO,
		//memorizer3,
		// ...
	}

	return memorizers
}

// Create folder and Switch working path
func createFolder(dirName string) {
	_, err := os.Stat(dirName)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Printf("Creating directory %s\n", dirName)
			err := os.MkdirAll(dirName, os.ModePerm)
			if err != nil {
				logger.Fatalf("Failed to create directory %s: %v\n", dirName, err)
			}
		} else {
			logger.Fatalf("Failed to check directory %s: %v\n", dirName, err)
		}
	}
}

func switchWorkingPath(dirName string) {
	logger.Printf("Switching to directory %s\n", dirName)
	err := os.Chdir(dirName)
	if err != nil {
		logger.Fatalf("Failed to switch to directory %s: %v\n", dirName, err)
	}
}

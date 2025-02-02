package memorizer

import (
	"log"
	"os"
	"src/plan"
)

type Memorizers interface {
	M_Cumulative(plan.Plans)
	M_Average(int)
	M_Output_Results()
	M_Store_Data(string, int)
	M_Store_File(string)
}

func New_Memorizers() map[string]Memorizers {
	// memorizer1 ...
	OMACO := new_OMACO_Memorizer()

	// Look-up table method
	memorizers := map[string]Memorizers{
		"omaco": OMACO,
		//memorizer2,
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
			log.Printf("Creating directory %s\n", dirName)
			err := os.MkdirAll(dirName, os.ModePerm)
			if err != nil {
				log.Fatalf("Failed to create directory %s: %v\n", dirName, err)
			}
		} else {
			log.Fatalf("Failed to check directory %s: %v\n", dirName, err)
		}
	}
}

func switchWorkingPath(dirName string) {
	log.Printf("Switching to directory %s\n", dirName)
	err := os.Chdir(dirName)
	if err != nil {
		log.Fatalf("Failed to switch to directory %s: %v\n", dirName, err)
	}
}

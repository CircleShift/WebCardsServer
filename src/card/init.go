package card

import (
	"log"

	"encoding/json"

	"io/ioutil"
	"os"
	"sync"
	"time"
	"math/rand"
)

var (
	rndSync sync.Mutex
	rndSeed *rand.Rand
)

func syncSafeRandom() *rand.Rand {
	rndSync.Lock()
	out := rand.New(rand.NewSource(rndSeed.Int63()))
	rndSync.Unlock()
	
	return out 
}

// packFromJSON takes a byte representation of json and attempts to convert it into a Pack object.
func packFromJSON(dat []byte) (Pack, error) {
	var out = Pack{}
	err := json.Unmarshal(dat, &out)

	if err != nil {
		return out, err
	}

	return out, nil
}

// packsFromDir takes a path to a directory and attempts to get any JSON files in it which represent a Pack object
func packsFromDir(path string) ([]Pack, error) {
	var out = []Pack{}
	var dat []byte
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return out, err
	}

	for _, f := range files {
		if f.IsDir() {
			tmp, err := packsFromDir(path + string(os.PathSeparator) + f.Name())
			if err == nil {
				out = append(out, tmp...)
			} else {
				log.Println("Failed to read " + path + string(os.PathSeparator) + f.Name() + " as directory.")
			}
		} else {
			dat, err = ioutil.ReadFile(path + string(os.PathSeparator) + f.Name())
			if err == nil {
				tmp, err := packFromJSON(dat)
				if err == nil {
					out = append(out, tmp)
				} else {
					log.Println("Unable to convert " + path + string(os.PathSeparator) + f.Name() + " to Pack")
				}
			} else {
				log.Println("Failed to read " + path + string(os.PathSeparator) + f.Name() + " as file.")
			}
		}
	}

	return out, nil
}

// readPacks should be called only once on startup.
// readPacks first attempts to find a pack labled default.json. If this pack can't be found, it exits with error.
// readPacks then attempts to find any JSON representation of a pack in the directory packdir.
// readPacks returns a slice of all the packs it could create.
func readPacks(packdir string) ([]Pack, error) {
	var out = []Pack{}
	var tmp Pack

	bytes, err := ioutil.ReadFile("default.json")

	if err != nil {
		return out, err
	}

	tmp, err = packFromJSON(bytes)

	if err != nil {
		return out, err
	}

	out = append(out, tmp)

	packs, err := packsFromDir(packdir)

	if err == nil {
		out = append(out, packs...)
	} else {
		log.Println("Failed to get extra packs from " + packdir + ". Does the folder exist?")
	}

	return out, nil
}

// Packs represents all the packs which could be read by readPacks.
// Packs is initialized by the InitCardPacks function.
var Packs []Pack

// AlreadyInitialized represents an error that the InitCardPacks function has already been called, or the Packs variable has already been altered.
type AlreadyInitialized struct{}

func (a AlreadyInitialized) Error() string {
	return "The card package has already been initialized."
}

// InitCardPacks attempts to initialize Packs by calling readPacks with the string packdir.
// Returns an error if an error occurs.
func InitCardPacks(packdir string) error {
	if len(Packs) != 0 {
		return AlreadyInitialized{}
	}

	p, err := readPacks(packdir)

	if err != nil {
		return err
	}

	Packs = p

	n := time.Now()
	d := time.Date(n.Year()-(n.Year()%32), time.January, 0, 0, 0, 0, 0, time.Local)
	rndSeed = rand.New(rand.NewSource(n.Sub(d).Microseconds()))

	return nil
}

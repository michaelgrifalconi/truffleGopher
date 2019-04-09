package trufflegopher

import (
	"encoding/gob"
	"log"
	"os"
)

func diskOperationSave(path string, t *Trufflegopher) error {

	file, err := os.Create(path)
	if err != nil {
		log.Println("Failed to create file")
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)

	// Create a temporary map, which will hold all item spread across shards.
	toBeEncoded := make(map[string]interface{})

	// Insert items to temporary map.
	for item := range t.SearchedDiffs.IterBuffered() {
		toBeEncoded[item.Key] = item.Val
	}

	gob.Register(struct{}{}) //TODO: no idea why with map[string]struct{}{} does not work!
	err = encoder.Encode(toBeEncoded)
	if err != nil {
		log.Println("Failed to encode map to disk")
		return err
	}
	return nil
}

func diskOperationLoad(path string, t *Trufflegopher) error {
	if _, err := os.Stat(path); err == nil {
		// path/to/whatever exists
		file, err := os.Open(path)
		if err != nil {
			log.Println("Failed to open file")
			log.Fatal(err)
		}
		defer file.Close()

		gob.Register(struct{}{}) //TODO: no idea why with map[string]struct{}{} does not work!
		d := gob.NewDecoder(file)
		m := make(map[string]interface{})
		// Decoding the serialized data
		err = d.Decode(&m)
		if err != nil {
			log.Println("Failed to decode")
			log.Fatal(err)
		}
		t.SearchedDiffs.MSet(m)

	} else if os.IsNotExist(err) {
		log.Println("File not found, skipping load..")
		return nil
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		log.Println("Failed to load file stats")
		log.Fatal(err)
	}
	return nil
}

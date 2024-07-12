package transform

import (
	"fmt"
	"io"
)

func Export(dir string, out io.StringWriter, batchSize int, mem bool) error {
	l, err := newLookups(dir)
	if err != nil {
		return fmt.Errorf("error creating look up tables from %s: %w", dir, err)
	}
	kv, err := newBadgerStorage(mem)
	if err != nil {
		return fmt.Errorf("could not create badger storage: %w", err)
	}
	defer kv.close()
	if err := kv.load(dir, &l); err != nil {
		return fmt.Errorf("error loading data to badger: %w", err)
	}
	j, err := createJSONRecordsTask(dir, nil, out, &l, kv, batchSize, false)
	if err != nil {
		return fmt.Errorf("error creating new task for venues in %s: %w", dir, err)
	}
	if err := j.runExport(); err != nil {
		return fmt.Errorf("error writing venues to database: %w", err)
	}
	return nil
}

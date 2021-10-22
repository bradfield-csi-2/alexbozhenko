package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestFileWriter_Append(t *testing.T) {
	tmpdir := t.TempDir()
	t.Log(tmpdir)
	fileWriter := NewFileWriter(
		filepath.Join(tmpdir, "movies"),
		[]string{"movieId", "title", "genres"},
	)

	tests := []struct {
		name    string
		writer  *FileWriter
		args    []Tuple
		wantErr bool
	}{
		{
			name:   "movies1",
			writer: fileWriter,
			args: []Tuple{
				{
					Values: []Value{{
						Name:        "movieId",
						StringValue: "1",
					}, {
						Name:        "title",
						StringValue: "Toy Story",
					}, {
						Name:        "genres",
						StringValue: "Adventure",
					}},
				},
				{
					Values: []Value{{
						Name:        "movieId",
						StringValue: "2",
					}, {
						Name:        "title",
						StringValue: "Jumanji",
					}, {
						Name:        "genres",
						StringValue: "Adventure",
					}},
				},
				{
					Values: []Value{{
						Name:        "movieId",
						StringValue: "3",
					}, {
						Name:        "title",
						StringValue: "Grumpe",
					}, {
						Name:        "genres",
						StringValue: "Romance",
					}},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, tuple := range tt.args {
				if err := fileWriter.Append(tuple); (err != nil) != tt.wantErr {
					t.Errorf("FileWriter.Append() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			fileScanOperator := NewFileScanOperator("movies", tmpdir, nil)

			readTuples := []Tuple{}
			for fileScanOperator.Next() {
				readTuples = append(readTuples, fileScanOperator.Execute())
			}
			if !reflect.DeepEqual(tt.args, readTuples) {
				t.Errorf("tuples read from disk %v are not the same as written to disk %v",
					readTuples, tt.args,
				)
			}

		})

	}
	t.Cleanup(func() {
		//		time.Sleep(50 * time.Second)
	})
}

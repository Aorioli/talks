package main

import (
	"os"

	"github.com/Aorioli/goworker"
)

// STRUCT_START OMIT
type Job struct {
	Queue     string
	Class     string
	Input     Data
	Following []*Job
}

type Data map[string]interface{}

// STRUCT_END OMIT
// HANDLER_START OMIT
type Handler interface {
	Handle(input Data) (output Data, err error)
}

// HANDLER_END OMIT

// DOWN_START OMIT
type Downloader struct {
	AccessToken string
}

func (d *Downloader) Handle(input Data) (Data, error) {
	u, ok := input["url"].(string)
	if !ok {
		return nil, nil
	}

	data, err := download(d.AccessToken, u)
	return Data{
		"data": data,
	}, err
}

func download(accessToken string, url string) ([]byte, error) {
	return []byte{0, 1, 2, 3}, nil
}

// DOWN_END OMIT

func (j *Job) InsertInput(data Data) {}

// REGISTER_START OMIT
type Registry map[string]Handler

// REGISTER_END OMIT

// INIT_START OMIT
func init() {
	registry := make(Registry, 1)
	registry["Downloader"] = &Downloader{
		AccessToken: os.Getenv("ACCESS_TOKEN"),
	}

	for c, h := range registry {
		goworker.Register(c, func(queue string, args ...interface{}) error {
			var job *Job

			output, err := h.Handle(job.Input)
			if err != nil {
				return err
			}

			for _, j := range job.Following {
				j.InsertInput(output)
			}

			return nil
		})
	}
}

// INIT_END OMIT

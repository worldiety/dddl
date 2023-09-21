package plantuml

import (
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

var hourGlass = []byte(`<svg class="w-6 h-6 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 16 20">
    <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 1H1m14 18H1m2 0v-4.333a2 2 0 0 1 .4-1.2L5.55 10.6a1 1 0 0 0 0-1.2L3.4 6.533a2 2 0 0 1-.4-1.2V1h10v4.333a2 2 0 0 1-.4 1.2L10.45 9.4a1 1 0 0 0 0 1.2l2.15 2.867a2 2 0 0 1 .4 1.2V19H3Z"/>
  </svg>`)

type PreflightContext struct {
	jobs []*RenderJob
}

func (c *PreflightContext) RequiresRendering() bool {
	return len(c.jobs) > 0
}

func (c *PreflightContext) Render() error {
	start := time.Now()
	defer func() {
		log.Printf("concurrent plantuml preflight took %v for %d diagrams\n", time.Now().Sub(start), len(c.jobs))
	}()

	var wg sync.WaitGroup
	for _, job := range c.jobs {
		job := job
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()
			defer func() {
				log.Printf("plantuml took %v\n", time.Now().Sub(start))
			}()

			cmd := exec.Command("plantuml", "-t"+job.fileType, "-p")
			cmd.Env = os.Environ()

			w, err := cmd.StdinPipe()
			if err != nil {
				job.resultErr = err
				return
			}

			if _, err := w.Write([]byte(job.pumlDiag)); err != nil {
				job.resultErr = err
				return
			}

			if err := w.Close(); err != nil {
				job.resultErr = err
				return
			}

			buf, err := cmd.Output()
			if err != nil {

				job.resultBuf = buf
				job.resultErr = err
				return
			}

			writeFileCache(job.fileType, job.pumlDiag, buf)

		}()
	}

	wg.Wait()

	for _, job := range c.jobs {
		if job.resultErr != nil {
			log.Println("first failed puml diagram:")
			log.Println(job.pumlDiag)
			return job.resultErr
		}
	}

	return nil
}

type RenderJob struct {
	fileType   string
	renderable Renderable
	pumlDiag   string
	resultBuf  []byte
	resultErr  error
}

func RenderLocalWithPreflight(ctx *PreflightContext, fileType string, renderable Renderable) ([]byte, error) {
	tmp := String(renderable)
	if buf := readFileCache(fileType, tmp); buf != nil {
		return buf, nil
	}

	ctx.jobs = append(ctx.jobs, &RenderJob{
		fileType:   fileType,
		renderable: renderable,
		pumlDiag:   tmp,
	})

	return hourGlass, nil
}

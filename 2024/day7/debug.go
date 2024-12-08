package main

import (
	"bytes"
	"fmt"
	"os"
)

type Debugger struct {
	builder    bytes.Buffer
	active     bool
	filePath   string
	writeAtLen int
}

func NewDebugBuilder(active bool, filePath string, writeAtMB int) *Debugger {
	if writeAtMB < 0 {
		writeAtMB = 256
	}
	debugger := &Debugger{
		active:     active,
		filePath:   filePath,
		writeAtLen: writeAtMB * 1024 * 1024,
	}

	if active {
		err := debugger.write(os.O_CREATE | os.O_WRONLY | os.O_TRUNC)
		if err != nil {
			panic(fmt.Errorf("could not write file on open: %w", err))
		}
	}

	return debugger
}

func (d *Debugger) WriteString(s string) {
	if d.active {
		_, err := d.builder.WriteString(s)
		if err != nil {
			panic(fmt.Errorf("could not write string to debug builder: %w", err))
		}
		d.writeIfOverTooLarge()
	}
}

func (d *Debugger) WriteFormat(format string, a ...any) {
	if d.active {
		_, err := d.builder.WriteString(fmt.Sprintf(format, a...))
		if err != nil {
			panic(fmt.Errorf("could not write formatted string to debug builder: %w", err))
		}
		d.writeIfOverTooLarge()
	}
}

func (d *Debugger) writeIfOverTooLarge() {
	if d.builder.Len() > d.writeAtLen {
		err := d.write(os.O_CREATE | os.O_WRONLY | os.O_APPEND)
		if err != nil {
			panic(fmt.Errorf("could not write periodic data to file: %w", err))
		}
		d.builder.Reset()
	}
}

func (d *Debugger) Close() {
	if d.builder.Len() == 0 {
		return
	}

	err := d.write(os.O_CREATE | os.O_WRONLY | os.O_APPEND)
	if err != nil {
		panic(fmt.Errorf("could not write file on close: %w", err))
	}
}

func (d *Debugger) write(osFlags int) error {
	file, err := os.OpenFile(d.filePath, osFlags, 0644)
	if err != nil {
		return fmt.Errorf("could not open file for write: %w", err)
	}
	defer file.Close()

	_, err = file.Write(d.builder.Bytes())
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}

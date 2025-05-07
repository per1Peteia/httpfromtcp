package response

import (
	"fmt"
)

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	bytesToWrite := len(p)
	bytesWritten := 0

	n, err := fmt.Fprintf(w.writer, "%x\r\n", bytesToWrite)
	if err != nil {
		return bytesWritten, err
	}
	bytesWritten += n

	n, err = w.writer.Write(p)
	if err != nil {
		return bytesWritten, err
	}
	bytesWritten += n

	n, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return bytesWritten, err
	}
	bytesWritten += n

	return bytesWritten, nil
}

func (w *Writer) WriteChunkedBodyDone(body []byte) error {
	_, err := w.writer.Write([]byte("0\r\n"))
	if err != nil {
		return fmt.Errorf("error writing chunked body done: %v", err)
	}

	return nil
}

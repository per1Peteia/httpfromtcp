package response

import ()

func (w *Writer) WriteChunkedBody(p []byte) (int, error)

func (w *Writer) WriteChunkedBodyDone() (int, error)

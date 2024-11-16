package nostr

import (
	"bytes"
	"compress/flate"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/gobwas/httphead"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsflate"
	"github.com/gobwas/ws/wsutil"
)

type Connection struct {
	conn              net.Conn
	enableCompression bool
	controlHandler    wsutil.FrameHandlerFunc
	flateReader       *wsflate.Reader
	reader            *wsutil.Reader
	flateWriter       *wsflate.Writer
	writer            *wsutil.Writer
	msgStateR         *wsflate.MessageState
	msgStateW         *wsflate.MessageState
}

func NewConnection(ctx context.Context, url string, requestHeader http.Header, tlsConfig *tls.Config) (*Connection, error) {
	dialer := ws.Dialer{
		Header: ws.HandshakeHeaderHTTP(requestHeader),
		Extensions: []httphead.Option{
			wsflate.DefaultParameters.Option(),
		},
		TLSConfig: tlsConfig,
	}
	conn, _, hs, err := dialer.Dial(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	enableCompression := false
	state := ws.StateClientSide
	for _, extension := range hs.Extensions {
		if string(extension.Name) == wsflate.ExtensionName {
			enableCompression = true
			state |= ws.StateExtended
			break
		}
	}

	// reader
	var flateReader *wsflate.Reader
	var msgStateR wsflate.MessageState
	if enableCompression {
		msgStateR.SetCompressed(true)

		flateReader = wsflate.NewReader(nil, func(r io.Reader) wsflate.Decompressor {
			return flate.NewReader(r)
		})
	}

	controlHandler := wsutil.ControlFrameHandler(conn, ws.StateClientSide)
	reader := &wsutil.Reader{
		Source:         conn,
		State:          state,
		OnIntermediate: controlHandler,
		CheckUTF8:      false,
		Extensions: []wsutil.RecvExtension{
			&msgStateR,
		},
	}

	// writer
	var flateWriter *wsflate.Writer
	var msgStateW wsflate.MessageState
	if enableCompression {
		msgStateW.SetCompressed(true)

		flateWriter = wsflate.NewWriter(nil, func(w io.Writer) wsflate.Compressor {
			fw, err := flate.NewWriter(w, 4)
			if err != nil {
				InfoLogger.Printf("Failed to create flate writer: %v", err)
			}
			return fw
		})
	}

	writer := wsutil.NewWriter(conn, state, ws.OpText)
	writer.SetExtensions(&msgStateW)

	return &Connection{
		conn:              conn,
		enableCompression: enableCompression,
		controlHandler:    controlHandler,
		flateReader:       flateReader,
		reader:            reader,
		msgStateR:         &msgStateR,
		flateWriter:       flateWriter,
		writer:            writer,
		msgStateW:         &msgStateW,
	}, nil
}

func (c *Connection) WriteMessage(ctx context.Context, data []byte) error {
	select {
	case <-ctx.Done():
		return errors.New("context canceled")
	default:
	}

	if c.msgStateW.IsCompressed() && c.enableCompression {
		c.flateWriter.Reset(c.writer)
		if _, err := io.Copy(c.flateWriter, bytes.NewReader(data)); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}

		if err := c.flateWriter.Close(); err != nil {
			return fmt.Errorf("failed to close flate writer: %w", err)
		}
	} else {
		if _, err := io.Copy(c.writer, bytes.NewReader(data)); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}
	}

	if err := c.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

func (c *Connection) ReadMessage(ctx context.Context, buf io.Writer) error {
	for {
		select {
		case <-ctx.Done():
			return errors.New("context canceled")
		default:
		}

		h, err := c.reader.NextFrame()
		if err != nil {
			c.conn.Close()
			return fmt.Errorf("failed to advance frame: %w", err)
		}

		if h.OpCode.IsControl() {
			if err := c.controlHandler(h, c.reader); err != nil {
				return fmt.Errorf("failed to handle control frame: %w", err)
			}
		} else if h.OpCode == ws.OpBinary ||
			h.OpCode == ws.OpText {
			break
		}

		if err := c.reader.Discard(); err != nil {
			return fmt.Errorf("failed to discard: %w", err)
		}
	}

	if c.msgStateR.IsCompressed() && c.enableCompression {
		c.flateReader.Reset(c.reader)
		if _, err := io.Copy(buf, c.flateReader); err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}
	} else {
		if _, err := io.Copy(buf, c.reader); err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}
	}

	return nil
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

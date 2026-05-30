package k8s

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"

	"k8s-ui-admin/internal/model"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ExecPod(ctx context.Context, cluster *model.Cluster, namespace, podName, container string, command []string, conn *websocket.Conn) error {
	client, err := GetClient(cluster)
	if err != nil {
		return err
	}

	config, err := buildConfig(cluster)
	if err != nil {
		return err
	}

	req := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   command,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	stdinReader, stdinWriter := io.Pipe()

	go func() {
		defer stdinWriter.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if len(message) > 0 {
				stdinWriter.Write(message)
			}
		}
	}()

	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  stdinReader,
		Stdout: &wsWriter{conn: conn},
		Stderr: &wsWriter{conn: conn},
		Tty:    true,
	})

	if err != nil {
		log.Printf("exec error: %v", err)
	}

	return err
}

type wsWriter struct {
	conn *websocket.Conn
}

func (w *wsWriter) Write(p []byte) (int, error) {
	err := w.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

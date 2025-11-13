package workload

import (
    "io"
    "net/http"
    "time"

    "github.com/JasonYueng-1231/CyberKube/backend/internal/service"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/client-go/kubernetes/scheme"
    "k8s.io/client-go/tools/remotecommand"
)

var upgrader = websocket.Upgrader{ CheckOrigin: func(r *http.Request) bool { return true } }

// /api/v1/pods/logs/stream?cluster=&namespace=&name=&container=&tail=200
func RegisterStream(r *gin.RouterGroup) {
    // 日志流
    r.GET("/pods/logs/stream", func(c *gin.Context) {
        conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil { return }
        defer conn.Close()
        cluster := c.Query("cluster")
        ns := c.Query("namespace")
        name := c.Query("name")
        container := c.Query("container")
        tail := int64(200)

        cli, err := service.GetClientForCluster(cluster)
        if err != nil { conn.WriteMessage(websocket.TextMessage, []byte("error:"+err.Error())); return }
        req := cli.CoreV1().Pods(ns).GetLogs(name, &corev1.PodLogOptions{Container: container, Follow: true, TailLines: &tail})
        rc, err := req.Stream(c)
        if err != nil { conn.WriteMessage(websocket.TextMessage, []byte("error:"+err.Error())); return }
        defer rc.Close()
        buf := make([]byte, 4*1024)
        for {
            n, er := rc.Read(buf)
            if n > 0 { _ = conn.WriteMessage(websocket.TextMessage, buf[:n]) }
            if er != nil { if er != io.EOF { _ = conn.WriteMessage(websocket.TextMessage, []byte("error:"+er.Error())) }; break }
        }
    })

    // WebShell: /api/v1/pods/shell?cluster=&namespace=&pod=&container=
    r.GET("/pods/shell", func(c *gin.Context) {
        conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil { return }
        defer conn.Close()
        cluster := c.Query("cluster")
        ns := c.Query("namespace")
        pod := c.Query("pod")
        container := c.Query("container")

        cli, err := service.GetClientForCluster(cluster)
        if err != nil { conn.WriteMessage(websocket.TextMessage, []byte("error:"+err.Error())); return }
        cfg, err := service.GetRestConfig(cluster)
        if err != nil { conn.WriteMessage(websocket.TextMessage, []byte("error:"+err.Error())); return }

        req := cli.CoreV1().RESTClient().Post().Resource("pods").Namespace(ns).Name(pod).SubResource("exec")
        req.VersionedParams(&corev1.PodExecOptions{
            Container: container,
            Command:   []string{"/bin/sh", "-c", "bash || sh"},
            Stdin:     true,
            Stdout:    true,
            Stderr:    true,
            TTY:       true,
        }, scheme.ParameterCodec)

        executor, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
        if err != nil { conn.WriteMessage(websocket.TextMessage, []byte("error:"+err.Error())); return }

        stream := &wsStream{conn: conn}
        _ = conn.WriteMessage(websocket.TextMessage, []byte("connected"))
        err = executor.Stream(remotecommand.StreamOptions{
            Stdin:  stream,
            Stdout: stream,
            Stderr: stream,
            Tty:    true,
        })
        if err != nil { _ = conn.WriteMessage(websocket.TextMessage, []byte("error:"+err.Error())) }
    })
}

type wsStream struct { conn *websocket.Conn }

func (w *wsStream) Read(p []byte) (int, error) {
    _, data, err := w.conn.ReadMessage()
    if err != nil { return 0, err }
    n := copy(p, data)
    return n, nil
}
func (w *wsStream) Write(p []byte) (int, error) {
    err := w.conn.WriteMessage(websocket.TextMessage, p)
    return len(p), err
}
func (w *wsStream) Close() error { return w.conn.Close() }


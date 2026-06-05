import { useEffect, useRef, useCallback } from 'react';

interface WsOptions {
  onMessage: (data: string | ArrayBuffer) => void;
  onClose?: () => void;
  onOpen?: () => void;
}

export function useWebSocket(url: string | null, options: WsOptions) {
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectRef = useRef<number>(0);
  const optionsRef = useRef(options);
  optionsRef.current = options;

  const send = useCallback((data: string | ArrayBuffer | Blob) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(data);
    }
  }, []);

  const sendJson = useCallback((obj: unknown) => {
    send(JSON.stringify(obj));
  }, [send]);

  const connect = useCallback(() => {
    if (!url) return;
    const token = localStorage.getItem('token');
    const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}${url}?token=${token}`;

    const ws = new WebSocket(wsUrl);
    ws.binaryType = 'arraybuffer';

    ws.onopen = () => {
      optionsRef.current.onOpen?.();
    };

    ws.onmessage = (event) => {
      if (event.data instanceof ArrayBuffer) {
        optionsRef.current.onMessage(event.data);
      } else if (typeof event.data === 'string') {
        optionsRef.current.onMessage(event.data);
      } else if (event.data instanceof Blob) {
        event.data.arrayBuffer().then(buf => {
          optionsRef.current.onMessage(buf);
        });
      }
    };

    ws.onclose = () => {
      optionsRef.current.onClose?.();
      // Auto-reconnect after 3s
      reconnectRef.current = setTimeout(() => {
        connect();
      }, 3000);
    };

    ws.onerror = () => {
      ws.close();
    };

    wsRef.current = ws;
  }, [url, send]);

  useEffect(() => {
    connect();
    return () => {
      if (reconnectRef.current) clearTimeout(reconnectRef.current);
      wsRef.current?.close();
    };
  }, [connect]);

  return { send, sendJson, ws: wsRef };
}

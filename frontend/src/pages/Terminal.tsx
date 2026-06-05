import { useEffect, useRef, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Terminal as XTerm } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { ArrowLeft, Wifi, WifiOff } from 'lucide-react';
import { useWebSocket } from '../hooks/useWebSocket';

export default function Terminal() {
  const { assetId } = useParams<{ assetId: string }>();
  const navigate = useNavigate();
  const termRef = useRef<HTMLDivElement>(null);
  const xtermRef = useRef<XTerm | null>(null);
  const fitRef = useRef<FitAddon | null>(null);
  const [connected, setConnected] = useState(false);
  const [disconnected, setDisconnected] = useState(false);

  const { sendJson } = useWebSocket(`/api/v1/ws/terminal/${assetId}`, {
    onOpen: () => {
      setConnected(true);
      setDisconnected(false);
    },
    onClose: () => {
      setConnected(false);
      setDisconnected(true);
    },
    onMessage: (data) => {
      if (typeof data === 'string') {
        xtermRef.current?.write(data);
      } else {
        xtermRef.current?.write(new Uint8Array(data as ArrayBuffer));
      }
    },
  });

  useEffect(() => {
    if (!termRef.current) return;

    const term = new XTerm({
      cursorBlink: true,
      cursorStyle: 'bar',
      fontSize: 14,
      fontFamily: "'Cascadia Code', 'Fira Code', 'JetBrains Mono', monospace",
      theme: {
        background: '#0d1117',
        foreground: '#c9d1d9',
        cursor: '#58a6ff',
        selectionBackground: '#264f78',
        black: '#484f58',
        red: '#ff7b72',
        green: '#3fb950',
        yellow: '#d29922',
        blue: '#58a6ff',
        magenta: '#bc8cff',
        cyan: '#39c5d6',
        white: '#b1bac4',
        brightBlack: '#6e7681',
        brightRed: '#ffa198',
        brightGreen: '#56d364',
        brightYellow: '#e3b341',
        brightBlue: '#79c0ff',
        brightMagenta: '#d2a8ff',
        brightCyan: '#56d4dd',
        brightWhite: '#f0f6fc',
      },
      allowProposedApi: true,
    });

    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);
    term.open(termRef.current);
    fitAddon.fit();

    term.onData((data) => {
      sendJson({ type: 'input', data });
    });

    const resizeObserver = new ResizeObserver(() => {
      fitAddon.fit();
      sendJson({ type: 'resize', cols: term.cols, rows: term.rows });
    });
    resizeObserver.observe(termRef.current);

    xtermRef.current = term;
    fitRef.current = fitAddon;

    return () => {
      resizeObserver.disconnect();
      term.dispose();
    };
  }, [assetId]);

  return (
    <div className="h-screen flex flex-col bg-gray-950">
      {/* Top bar */}
      <div className="flex items-center gap-3 px-4 py-2 bg-gray-900 border-b border-gray-800">
        <button onClick={() => navigate('/hosts')} className="p-1 text-gray-400 hover:text-white">
          <ArrowLeft className="w-5 h-5" />
        </button>
        <div className="flex-1 text-sm text-gray-300">Terminal: Host #{assetId}</div>
        {connected ? (
          <span className="flex items-center gap-1 text-xs text-success"><Wifi className="w-3 h-3" /> Connected</span>
        ) : (
          <span className="flex items-center gap-1 text-xs text-critical"><WifiOff className="w-3 h-3" /> Disconnected</span>
        )}
      </div>

      {/* Terminal */}
      <div className="flex-1 p-1">
        <div ref={termRef} className="h-full w-full" />
      </div>

      {/* Disconnected overlay */}
      {disconnected && (
        <div className="absolute inset-0 flex items-center justify-center bg-black/70 z-10">
          <div className="text-center">
            <WifiOff className="w-12 h-12 text-critical mx-auto mb-4" />
            <p className="text-white text-lg mb-2">Connection Lost</p>
            <p className="text-gray-400 text-sm">Attempting to reconnect...</p>
          </div>
        </div>
      )}
    </div>
  );
}

import {genPushTicket} from "@/services/common/push";

type WebSocketOptions = {
  url: string;
  onConnect?: () => void;
  onClose?: (event: Event) => void;
  onError?: (event: Event) => void;
  onMessage?: (cmd: string, id: number, data: any) => void;
  reconnectAttempts?: number;
};

interface WebSocketMessage {
  cmd: string;
  id: number;
  data: any;
}

const defaultOptions: Required<WebSocketOptions> = {
  url: '',//连接的长链接
  onConnect: () => {
  },//开启连接
  onClose: () => {
  },//关闭链接
  onError: () => {
  },//异常
  onMessage: () => {
  },//消息
  reconnectAttempts: Number.MAX_VALUE,//最大连接范围数
};

const getMessageId = () => {
  return new Date().getTime() * 1000 + Math.floor(Math.random() * 1000);
}

const SinKingPush = async (
  options: WebSocketOptions,
  attempts = 0
): Promise<WebSocket | undefined> => {
  const {url, onConnect, onClose, onError, onMessage, reconnectAttempts} = {
    ...defaultOptions,
    ...options,
  };

  if (attempts >= reconnectAttempts) {
    return undefined;
  }

  const r = await genPushTicket();
  if (r?.code != 200 || r?.data == undefined) {
    return undefined;
  }
  const socket = new WebSocket(url, r?.data);
  let pingTimer: NodeJS.Timer;
  let pingTimeout: NodeJS.Timeout;

  const sendMessage = (cmd: string, id: number = 0, data: any = undefined): boolean => {
    if (socket.readyState === WebSocket.OPEN) {
      if (id == 0) {
        // eslint-disable-next-line no-param-reassign
        id = getMessageId()
      }
      const msg: WebSocketMessage = {cmd: cmd, id: id, data: data}
      socket?.send(JSON.stringify(msg))
      return true;
    } else {
      return false;
    }
  }

  const sendPing = () => {
    sendMessage("ping");
    pingTimeout = setTimeout(() => {
      socket.close();
    }, 10000);
  };
  socket.onopen = () => {
    clearInterval(pingTimer);
    clearTimeout(pingTimeout);
    pingTimer = setInterval(sendPing, 30000);
    onConnect?.();
  };

  socket.onclose = (event) => {
    clearInterval(pingTimer);
    clearTimeout(pingTimeout);
    onClose?.(event);
    // Reconnect after 3 seconds
    setTimeout(() => {
      SinKingPush(options, attempts + 1);
    }, 3000);
  };

  socket.onerror = (event) => {
    clearInterval(pingTimer);
    clearTimeout(pingTimeout);
    onError?.(event);
    // Reconnect after 3 seconds
    setTimeout(() => {
      SinKingPush(options, attempts + 1);
    }, 3000);
  };

  socket.onmessage = (event) => {
    clearTimeout(pingTimeout);
    try {
      const message: WebSocketMessage = JSON.parse(event?.data || "{}");
      if (message != undefined && message?.cmd != undefined && message?.id != undefined && message?.data != undefined) {
        onMessage?.(message?.cmd, message?.id, message?.data);
      }
    } catch {
      return
    }
  };

  // @ts-ignore
  socket.closeServer = () => {
    // eslint-disable-next-line no-param-reassign
    attempts = reconnectAttempts + 1
    clearInterval(pingTimer);
    clearTimeout(pingTimeout);
    socket?.close()
  }

  // @ts-ignore
  socket.sendMessage = sendMessage

  if (socket.readyState === WebSocket.OPEN) {
    sendPing();
    return socket;
  } else {
    return await new Promise<WebSocket | undefined>((resolve) => {
      socket.onopen = () => {
        clearInterval(pingTimer);
        clearTimeout(pingTimeout);
        pingTimer = setInterval(sendPing, 30000);
        onConnect?.();
        resolve(socket);
      };
    });
  }
};

export default SinKingPush;

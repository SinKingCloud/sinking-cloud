type WebSocketOptions = {
  url: string;
  token?: string;
  onConnect?: () => void;
  onClose?: (event: Event) => void;
  onError?: (event: Event) => void;
  onMessage?: (message: any) => void;
  reconnectInterval?: number;
  reconnectAttempts?: number;
};

const defaultOptions: Required<WebSocketOptions> = {
  url: '',//连接的长链接
  token: '',//鉴权
  onConnect: () => {
  },//开启连接
  onClose: () => {
  },//关闭链接
  onError: () => {
  },//异常
  onMessage: () => {
  },//消息
  reconnectInterval: 1000,//重连时长设置
  reconnectAttempts: Number.MAX_VALUE,//最大连接范围数
};

const SinKingPush = (
  options: WebSocketOptions,
): WebSocket => {
  const {url, token, onConnect, onClose, onError, onMessage} = {
    ...defaultOptions,
    ...options,
  };
  console.log(token)
  const socket = new WebSocket(url, token);
  socket.onopen = () => {
    onConnect();
  };
  socket.onclose = (event) => {
    onClose(event);
  };
  socket.onerror = (event) => {
    onError(event);
  };
  socket.onmessage = (event) => {
    onMessage(event?.data);
  };
  return socket
};

export default SinKingPush;

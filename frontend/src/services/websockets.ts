import { getUsername } from './auth'

export interface ServerMessage {
  source: string
  text: string
  type: string
  timestamp: Date
}

export class WebSocketClient {
  private socket: WebSocket

  constructor(apiEndpoint: string) {
    // If on HTTPS, use wss, otherwise ws
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'

    const endpointSplit = apiEndpoint.split('://')
    this.socket = new WebSocket(`${protocol}://${endpointSplit[1]}/connect`)

    this.socket.onopen = () => {
      console.log('🔌 WebSocket: Connected to backend server...')
      this.socket.send(`{ "username": "${getUsername()}" }`)
    }

    this.socket.onerror = (event) => {
      console.error('🔌 WebSocket: error - ', event)
    }
  }

  public addMessageCallback(messageCallback: (message: ServerMessage) => void) {
    this.socket.addEventListener('message', (event: any) => {
      const rawData = event.data
      try {
        const msg = JSON.parse(rawData)
        messageCallback(msg)
      } catch (e) {
        console.error('🔌 WebSocket: Error parsing message - ', rawData)
      }
    })
  }

  public addClosedCallback(callback: (event: WebSocketEventMap) => void) {
    this.socket.addEventListener('close', (event: any) => {
      console.log('🔌 WebSocket: closed - ', event)
      callback(event)
    })
  }
}

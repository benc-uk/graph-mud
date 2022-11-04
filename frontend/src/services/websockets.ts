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
    // if on https, use wss, otherwise ws
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'

    const endpointSplit = apiEndpoint.split('://')
    this.socket = new WebSocket(`${protocol}://${endpointSplit[1]}/connect`)

    this.socket.onopen = () => {
      console.log('🔌 Connected to WebSocket')
      this.socket.send(`{ "username": "${getUsername()}" }`)
    }

    this.socket.onerror = (event) => {
      console.error('🔌 WebSocket error:', event)
    }
  }

  public addMessageCallback(messageCallback: (message: ServerMessage) => void) {
    this.socket.addEventListener('message', (event: any) => {
      const rawData = event.data
      try {
        const msg = JSON.parse(rawData)
        messageCallback(msg)
      } catch (e) {
        console.error('🔌 Error parsing message:', rawData)
      }
    })
  }

  public addClosedCallback(callback: (event: any) => void) {
    this.socket.addEventListener('close', (event: any) => {
      console.log('🔌 WebSocket closed:', event)
      callback(event)
    })
  }
}

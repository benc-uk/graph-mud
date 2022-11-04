//
// Generate a RS256 signed JWT token
//

const rs = require('jsrsasign')
const rsu = require('jsrsasign-util')

// Load key from PEM file
const keyPem = rsu.readFile('private-key.pem')
const privateKey = rs.KEYUTIL.getKey(keyPem)

// Header
const header = { alg: 'RS256', typ: 'JWT' }

// Payload
let payload = {}
const timeNow = rs.jws.IntDate.get('now')
const timeEnd = rs.jws.IntDate.get('now + 1year')
payload.iss = 'http://localhost:8000/no-one'
payload.sub = 'nano-realms'
payload.aud = ['nano-realms']
payload['preferred_username'] = 'demo@example.net'
payload.nbf = timeNow
payload.iat = timeNow
payload.exp = timeEnd
payload.auth_time = timeEnd
payload.jti = Array(8)
  .fill()
  .map((n) => ((Math.random() * 36) | 0).toString(36))
  .join('')

// Sign JWT with RSA private key
var jwt = rs.jws.JWS.sign('RS256', JSON.stringify(header), JSON.stringify(payload), privateKey)

// Output
console.log(`===== 🔐 JWT Token ===================================`)
console.log(jwt)
console.log(`======================================================`)

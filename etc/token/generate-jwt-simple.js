//
// Generate a RS256 signed JWT token
//
var jwt = require('jwt-simple')

// Payload
let payload = {}
payload['preferred_username'] = 'barry@umbrella.org'

// Sign JWT with RSA private key
var token = jwt.encode(payload, '231412894ihosfnklsfh')
// Output
console.log(`===== 🔐 JWT Token ===================================`)
console.log(token)
console.log(`======================================================`)

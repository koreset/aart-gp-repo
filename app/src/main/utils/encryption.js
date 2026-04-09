import { machineIdSync } from 'node-unique-machine-id'
const crypto = require('crypto')
// const { machineId } = require('node-machine-id')

const algorithm = 'aes-256-ctr'
const secretKey = 'vOVH6sdmpNWjRRIqCc7rdxs01lwHzfr3'
const iv = crypto.randomBytes(16)

export const encrypt = (text) => {
  const cipher = crypto.createCipheriv(algorithm, secretKey, iv)

  const encrypted = Buffer.concat([cipher.update(text), cipher.final()])

  return {
    iv: iv.toString('hex'),
    content: encrypted.toString('hex')
  }
}

export const decrypt = (hash) => {
  const decipher = crypto.createDecipheriv(
    algorithm,
    secretKey,
    Buffer.from(hash.iv, 'hex')
  )

  const decrpyted = Buffer.concat([
    decipher.update(Buffer.from(hash.content, 'hex')),
    decipher.final()
  ])

  return decrpyted.toString()
}

export const getFingerprint = async () => {
  const id = machineIdSync(true)
  return crypto.createHash('sha512').update(id).digest('hex')
}

module.exports = {
  encrypt,
  decrypt,
  getFingerprint
}

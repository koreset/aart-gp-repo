const os = require('os')
const crypto = require('crypto')

// Function to generate a machine fingerprint
export const generateMachineFingerprint = () => {
  const cpus = os.cpus()
  const platform = os.platform()
  const arch = os.arch()
  const totalMemory = os.totalmem()

  // Combine relevant system information
  const fingerprintData = {
    cpus,
    platform,
    arch,
    totalMemory
  }

  // Create a hash of the combined data
  const fingerprint = crypto
    .createHash('sha256')
    .update(JSON.stringify(fingerprintData))
    .digest('hex')

  return fingerprint
}

// Example usage
// const machineFingerprint = generateMachineFingerprint()

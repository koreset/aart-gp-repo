/* eslint-disable no-template-curly-in-string */
const dotenv = require('dotenv')

const baseConfig = {
  productName: 'AART Group Pricing',
  appId: 'za.co.adsolutions.aart-gp',
  asar: true,
  extends: null,
  compression: 'maximum',
  artifactName: '${productName}-${version}-${arch}.${ext}',
  directories: {
    output: './release/${version}'
  },
  mac: {
    bundleVersion: '1.0',
    identity: 'Actuaries and Digital Solutions (Pty) Ltd',
    hardenedRuntime: true,
    gatekeeperAssess: false,
    notarize: false,
    icon: 'buildAssets/icons/icon.icns',
    type: 'distribution',
    target: [
      // {
      //   target: 'zip',
      //   arch: ['arm64']
      // },
      {
        target: 'dmg',
        arch: ['x64', 'arm64', 'universal']
      },
      {
        target: 'zip',
        arch: ['x64', 'arm64', 'universal']
      }
    ]
  },
  dmg: {
    contents: [
      {
        x: 410,
        y: 150,
        type: 'link',
        path: '/Applications'
      },
      {
        x: 130,
        y: 150,
        type: 'file'
      }
    ],
    sign: false
  },
  win: {
    icon: 'buildAssets/icons/icon.ico',
    // certificateSubjectName: 'Actuaries and Digital Solutions (Pty) Ltd',
    target: [
      {
        target: 'nsis',
        arch: 'x64'
      }
    ]
  },
  portable: {
    artifactName: '${productName} ${version}_${arch} Portable.${ext}'
  },
  nsis: {
    oneClick: true
  }
}

dotenv.config()

baseConfig.copyright = `ⓒ ${new Date().getFullYear()} $\{author}`
baseConfig.files = [
  'dist/**/*',
  '!dist/main/index.dev.js',
  '!docs/**/*',
  '!tests/**/*',
  '!release/**/*'
]

module.exports = {
  ...baseConfig
}

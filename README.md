# EncryptUI

A portable and minimal workspace for the encryption and decryption of files.

## Overview

EncryptUI is built with Golang and the Fyne GUI Toolkit, delivering near-native performance across multiple platforms. 

## Key Features

###  **Standalone Files**
- Single binaries that can be fully configured to store data
- Open and decrypt files without ANY additional software
- Zero external dependencies required
- Note: The way this works is by embedding a 'marker' inside the file. (run strings standalone | grep "START"). The write function then appends the marker + the data. The file later reads itself, searches for that marker and decrypts said data)

https://github.com/user-attachments/assets/b0dc7db7-b3c2-4951-b236-b01a28123015

- 
### ️ **User-Friendly Interface**
- GUI-based encryption and decryption
- Intuitive and easy to pick up design
- Cross-platform compatibility

###  **Quick Access**
- Recent files functionality for easy access to frequently used files
- Streamlined workflow for regular encryption tasks

###  **Lightweight & Portable**
- Just 33MB total size
- No additional dependencies needed
- Run from any location without installation

## Security

- **Encryption**: AES256s 
- **Password Hashing**: SHA256 in order to ensure password security
- **Standalone Security**: Files are self-contained and decrypt themselves

## System Requirements

- Cross-platform support (Windows, macOS, Linux)
- Near-imperceptible system resource usage
- Dependency-less (No additional software installations needed)

*Built with Golang and Fyne GUI Toolkit for optimal performance and reliability.*

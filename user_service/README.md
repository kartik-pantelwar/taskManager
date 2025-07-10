
# Task Management System - User Service

## 📋 Overview

This is the User Service component of a microservices-based Task Management System. It handles user authentication, session management, and provides user-related functionality.

---

## 🔐 How Session ID is Created and Works

### When Session ID is Created

A session ID is typically created during the **user login process**:

1. **User submits credentials** (username/password) to the user service
2. **User service validates credentials** against the database  
3. **Upon successful authentication**, user service:
   - Generates a unique session ID (usually a random string/UUID)
   - Stores session data server-side (in memory, Redis, database)
   - Sets the session ID as an HTTP cookie in the response

### What Credentials/Data it Stores

**Session ID itself**: Just a random identifier (e.g., `abc123def456`)

**Server-side session data** (linked to the session ID):
- User ID
- Username  
- User roles/permissions
- Login timestamp
- Last activity timestamp
- Any other user context needed

### Default Expiration Time

**Common defaults**:
- **Web applications**: 30 minutes to 24 hours of inactivity
- **Banking/sensitive apps**: 15-30 minutes

---

## 🏗️ Architecture & Service Interaction

### **Session ID Origin**
The session ID is created and managed by the user service during login/authentication.

### **Cookie Storage** 
The session ID is stored as a cookie (named "session_id") in the client's browser.

### **Task Service Role**
The task service only:
- Reads the existing session ID from the cookie
- Validates it by making a gRPC call to the user service
- Uses the validation result to authorize access to task endpoints

The task service acts as a **consumer of authentication**, not a provider.

---

## 🔄 Authentication Flow

The typical flow would be:

1. **User logs in** via user service → user service generates session ID
2. **User service sets** session cookie in browser  
3. **User makes request** to task service → browser automatically sends session cookie
4. **Task service extracts** session ID and validates it with user service

> **Note**: The session ID generation and management is handled by the user service, while the task service just validates existing sessions.
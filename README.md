# Pryx - API Proxy Management System

Pryx is a Go-based API proxy management system designed to facilitate interaction with various AI model endpoints. It provides a web interface for managing proxy configurations, user authentication, and integration with external APIs like OpenAI. The system supports user management, session handling, and secure API key management for proxying requests to AI models.

## Features

- **API Proxy Management**: Add, list, and retrieve model and endpoint configurations for proxying API requests.
- **User Authentication**: Secure user signup, signin, and password reset functionality with session management.
- **Database Integration**: Uses PostgreSQL for persistent storage of models, endpoints, users, sessions, and password resets.
- **Email Notifications**: Configurable SMTP-based email service for password reset workflows.
- **Secure Middleware**: Implements CSRF protection and user session middleware for secure HTTP requests.
- **Web Interface**: A modern, responsive dashboard for managing proxies, built with HTML templates and CSS styling.
- **API Integration**: Supports OpenAI API for chat completions, with configurable endpoints and API keys.

## Prerequisites

- **Go**: Version 1.25.1 or higher.
- **PostgreSQL**: A running PostgreSQL database instance.
- **SMTP Server**: An SMTP server for sending password reset emails.
- **Dependencies**: Managed via Go modules (see `go.mod`).

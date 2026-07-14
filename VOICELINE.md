# Building Voiceline
Tech stack, building blocks and protocols used for building Voiceline

## Functionalities
At the point of writing Voiceline's core feature is a voice-enabled assistant that can be accessed via mobile app that can perform actions in a connected CRM system.

### List of major features
- Answer questions based on information in the user's organisation's CRM
- Update and create entries in the user's organisation's CRM based on user input during sessions

## Components
Building blocks for Voiceline core functionality.

### User CRM
Admins need to authorise Voiceline to access their CRM system and perform actions on their users' behalf

#### Responsibilities
- Admins add Voiceline as trusted OAuth application

### Mobile App
Users can access the assistant through a mobile app after authenticating.

#### Responsibilities
- Create new Voiceline user account
- Connect to a user's organisation's CRM through OAuth authentication to obtain initial refresh token
- Authenticate the user
- Start a new session with the Voiceline assistant in speech or text mode
- Forward user inputs to the Voiceline backend and relay responses to the user
- Manage user settings such as connected phone number and 2FA

#### Technologies
- Written in Dart and Flutter due to rich library ecosystem around building chat applications
- Use of flutter_gen_ai_chat_ui library to implement assistant UI

#### Connections
- Connect via HTTP to the Voiceline API gateway
- Connect through integrated browser via HTTP to the user's CRM during OAuth flow
- Connect to APNs and firebase notification service for user notifications

### API Gateway
The internal api gateway ensures that all requests reaching the internal services are authenticated and forwards them to the matching services.

#### Responsibilities
- Check incoming requests from the mobile app for valid authentication tokens
- Forward requests to the matching services

#### Connections
- Accepts incoming http requests from the mobile app
- Forwards and relays http requests to and from the internal services

### IAM Service
The Voiceline IAM service is responsible for holding Voiceline user accounts and to authenticate users.

#### Responsibilities
- Allow creation of new user accounts
- Send OTPs for 2FA if enabled
- Issue JWT that the mobile application can use to authenticate the user with the Voiceline services

#### Technologies
- Either Keycloak or a managed IAM service like AWS IAM service

#### Connections
- Accepts incoming HTTP requests from the user service
- Accepts incoming HTTP requests from the API gateway
- Queries and updates user authentication data from its dedicated storage

### User Service
The internal user service is responsible for creating and updating user accounts by talking to the IAM service.

#### Responsibilities
- Accept and process requests by the mobile app for creating new user accounts and updating existing ones
- Update authentication settings for users in the IAM service
- Provide access token for user's CRM to other services and handle token rotation

#### Technologies
- A fast application server written in Go or Rust

#### Connections
- Accepts HTTP requests from the API gateway
- Sends HTTP request to the IAM service
- Accepts HTTP request from assistant service

### Assistant Service
The assistant service is responsible for running the Voiceline agent and handling user sessions with the agent.

#### Responsibilities
- Provide runtime for Voiceline agent
- Define agent tools and capabilities
- Manage user sessions with the agent
- Execute agent tool calls

#### Technologies
- An easily extendable application server written in Python or Go
- Use of an agent framework like LangChain or Eino

#### Connections
- Accepts HTTP requests from the API gateway
- Sends HTTP requests to the user's CRM
- Sends HTTP requests to the user service

### Database Server
A database server is responsible for keeping the databases used for storing user data and past user-agent conversations.

#### Responsibilities
- Store user settings and encrypted refresh tokens
- Store past user-agent conversations
- Provide episodic memory for agent

#### Technologies
- PostgreSQL
- pgvector

#### Connections
- Accepts PGWire requests from user and assistant service (and also IAM service if self-hosted)

### Temporary Storage
During user-agent conversations the assistant service requires and external cache to store conversation state in order not to overload server memory and survive crashes or instance swaps.

#### Responsibilities
- Store conversation state including messages and access token for CRM

#### Technologies
- Redis

#### Connections
- Accepts RESP requests from assistant service
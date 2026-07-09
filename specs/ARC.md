# Architecture Summary
DieStimme is a light-weight implementation of the Voiceline AI assistant. It can process incoming phone calls and perform actions in a Twenty CRM system on the callers behalf. Authentication and authorisation is based on the callers' phone numbers.

## Components
There are four major components to the system.

### Twilio
Twilio is needed to enable phone calls to the DieStimme assistant.

#### Responsibilities
- Accept incoming calls and provide caller information to backend via webhook calls
- Handle initial identification and authentication of users by sending phone number and phone pin to backend webhook
- Stream authenticated calls' audio to and from backend for DieStimme assistant sessions

### OpenAI API
OpenAI API is used to convert audio data from Twilio to text and vice versa and to handle the assistant's reasoning.

#### Responsibilities
- Convert audio data to text
- Provide LLM for assistant reasoning
- Convert assistant responses to audio

### Go Backend
The internal Go backend defines the assistant workflows and connects to the Twenty CRM system.

#### Responsibilities
- Receive webhook calls from Twilio to identify and authenticate callers
- Run DieStimme agent sessions
- Execute agent function tool calls
- Connect to Twenty CRM to enable tool calls and identify and authenticate callers

### Twenty CRM
Twenty CRM is used as implementation of the CRM and represents the persistence layer.

#### Responsibilities
- Holds sales reps as People objects with custom Phone Pin field for caller identification and authentication
- Store companies, opportunities and notes
- Provide the GraphQL API for accessing and manipulating companies, opportunities and notes
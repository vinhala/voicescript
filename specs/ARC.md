# Architecture Summary
voicescript is proof of concept implementation of a service that analyses the recordings of a requirements elicitation session with a new client of Voiceline. The service analyses the recording for a session and based on it copmletes a standardised questionnaire, which is then stored in Voiceline's CRM system. It uses Twenty CRM as Voiceline's CRM and the OpenAI api as LLM provider. Interactions with the user are made through a simple web interface.

## Components
There are four major components to the system. Two internal, two external.

### OpenAI API
OpenAI API is used to convert audio data to text and analyse the transcipts to generate a filled out questionnaire.

#### Responsibilities
- Convert audio data to text
- Provide LLM for generating questionnaire

#### Relations
- Accept HTTP request from Go backend

### Go Backend
The internal Go backend offers a light-weight application server based on the Gin framework that accepts requests from the web interface and communicates with the OpenAI api to generate the questionnaire. It also integrates with Twenty CRM to fetch existing clients and store the questionnaires.

#### Responsibilities
- Handle requests from web interface
- Interact with OpenAI API
- Interact with Twenty CRM API
- Coordinate recording analysis and persistence

#### Relations
- Interact with Twenty CRM API via HTTP requests
- Interact with OpenAI API via HTTP requests
- Serve HTTP requests from the Nuxt server

### Twenty CRM
Twenty CRM is used as implementation of the CRM and represents the persistence layer.

#### Responsibilities
- Store companies, opportunities and notes
- Provide the REST HTTP API for accessing and manipulating companies, opportunities and notes

#### Relations
- Accept HTTP requests from Go backend
- Serve web interface for CRM to FDE's browser

### Nuxt Webserver
A Nuxt webserver is responsible for providing the web interface that FDEs use for providing the session recordings.

#### Responsibilities
- Serve the web interface accessible to FDEs

#### Relations
- Interact with Go backend via HTTP
- Serve html page to FDE's browser
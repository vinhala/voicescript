# Epics
DieStimme is proof of concept implementation of a service that analyses the recordings of a requirements elicitation session with a new client of Voiceline. The service analyses the recording for a session and based on it copmletes a standardised questionnaire, which is then stored in Voiceline's CRM system (in this case Twenty CRM).

## E.1 Analyse Recording
During their initial meeting with customers the forward deployed engineers of Voiceline determine how the client plans to use Voiceline in conjunction with their CRM system. An LLM pipeline turns the recording into a transcript with detection of individual speakers and then fills out a standardised questionnaire as far possible and marks open questions that require clarification.

### US.1.1 Provide Session Recording
As a forward deployed engineer I want to be able to easily provide a requirements elicitation session recording and link it to an existing client and opportunity from the CRM system for analysis.

#### Requirements
- List of existing clients can be retrieved from CRM system
- Recording is in format MP3 or MP4 audio

#### Main Flow
1. FDE opens interface for providing recording
2. System provides list of existing companies
3. FDE selects a company
4. System fetches list of opportunities for selected company
5. FDE selects an opportunity
6. System asks FDE to provide recording
7. FDE selects recording from their local machine
8. FDE submits recording for analysis
9. System uploads recording to server for further analysis and indicates that the analysis process has begun

### US.1.2 Convert Recording to Transcript
As an FDE I want the system to turn a recording of a requirements elicitation session to a transcript with speaker recognition so that it can fill in the standardised questionnaire.

#### Requirements
- US.1.1
- Recording matches requirements of the OpenAI API for audio recordings

#### Main Flow
1. System uploads audio recording to OpenAI API and asks it to turn the recording into a transcript of the session
2. OpenAI API turns recording into a transcript in the form of a dialogue
3. OpenAI API returns complete transcript to system for further analysis
4. System stores the transcript temporarily for further processing

#### Alternative Flows
- Step 2 fails -> System informs the FDE that transcription has failed and displays the concrete error from the API

### US.1.3 Use Transcript to Fill Out Questionnaire
As an FDE I want the system to use a generated transcript to fill out the standardised questionnaire for client onboarding so that I do not have to do it manually.

#### Requirements
- US.1.2 completed successfully
- Questionnaire is available to the system

#### Main Flow
1. System loads standardised client onboarding questionnaire
2. System provides transcipt and questionnaire to OpenAI api and asks it to fill out the questionnaire and mark open questions that could not be answered from the transcript
3. OpenAI api generates a filled out version of the questionnaire with open questions marked in markdown formatting
4. System stores response temporarily for further processing

#### Alternative Flows
- Step 3 fails -> System informs the FDE that analysis of transcript has failed and displays the concrete error from the API

## E.2 Storage in CRM
Generated filled out questionnaires should be stored in the Voiceline CRM as a note linked to the respective client so that the information is available to the FDE at a later point for integration of Voiceline into the client's CRM.

### US.2.1 Create Note in CRM
As an FDE I want the system to store a filled out questionnaire in the Voiceline CRM so that I can use it at a later stage to integrate Voiceline into a client's CRM.

#### Requirements
- US.1.3
- Filled out questionnaire has been generated in markdown

### Main Flow
1. System creates a new Note object in the Voiceline CRM and links it to the earlier selected opportunity
2. System informs FDE that the filled out questionnaire has been stored in the CRM and displays the final questionnaire to the FDE
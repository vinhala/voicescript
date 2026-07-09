# Epics
DieStimme is a light-weight alternative implementation of the Voiceline AI assistant. It can process incoming phone calls and perform actions in a Twenty CRM system on the callers behalf.

## E.1. Start Phone Session
When sales reps call the phone number of the assistant they are required to authenticate themselfes by providing their phone call pin that is configured in the CRM system.

### US.1.1 Call DieStimme Assistant
As a sales rep (user) I want to be able to access the DieStimme assistant (assistant) by calling its dedicated phone number so that I can perform actions on the CRM through it.

#### Main Flow
1. User calls phone number assigned to the assistant
2. System uses caller phone number to fetch related account from CRM including the configured phone PIN for the user
3. If the system finds no related user, it responds with "I am sorry, your phone number is not registered in the system" and terminates the call
4. Otherwise the system prompts the user to provide their phone pin
5. User speaks their four digit phone pin
6. System checks if the provided pin matches the configured pin for the user
7. If it does not match, the system responds with "I am sorry, the provided code does not match the one configured for your account" and terminates the call
8. Otherwise the system starts a new conversation with the assistant

### US.1.2 Welcome Sales Rep
As a sales rep (user) I want to be greeted by the DieStimme assistant (assistant) so that I know that I can start performing actions through it in the CRM.

#### Requirements
- US.1.1
- User has been identified and authenticated

#### Main Flow
1. Assistant greets the user by their name and in their preferred language
2. User starts interacting with the assistant

## E.2. Perform CRM Actions in Phone Session
Sales reps should be able to perform actions on the CRM through the phone. Actions that are available should be determined by the DieStimme assistant.

### US.2.1 Create New Company
As a sales rep (user) I want to be able to create a new company in the CRM through the DieStimme assistant (assistant) so that I can start tracking it as a prospect through the CRM.

#### Requirements
- US.1.1
- US.1.2

#### Main Flow
1. User asks assistant to create a new company in the CRM
2. Assistant asks the user for company name
3. User provides the company name
4. Assistant checks if company already exists
5. If company already exists, the assistant informs the user that the company appears to already exist and asks if they want to update its properties instead
6. Otherwise the assistant asks user if they know any of domain name, address, phone number, annual revenue and phone number
7. User provides one or more of the optional properties
8. Assistant reads back collected information and asks for confirmation or correction
9. User corrects one or more properties
10. Assistant repeats reads back collected information and asks for confirmation or correction
11. User confirms correctness
12. Assistant creates new company in the CRM
13. Assistant confirms creation of the company and asks user if they want to create an initial opportunity for it

### US.2.2 Update Existing Company
As a sales rep (user) I want to be able to update an existing company in the CRM through the DieStimme assistant (assistant) so that the information in the CRM remains up to date.

#### Requirements
- US.1.1
- US.1.2

#### Extends
- US.2.1

#### Main Flow
1. User asks assistant to update a company in the CRM
2. Assistant asks the user for company name
3. User provides the company name
4. Assistant checks if company already exists
5. If company does not exist, the assistant informs the user that the company appears not to exist and asks if they want to create a new one
6. Otherwise the assistant asks user if they know any of domain name, address, phone number, annual revenue and phone number
7. User provides one or more of the optional properties
8. Assistant reads back collected information and asks for confirmation or correction
9. User corrects one or more properties
10. Assistant repeats reads back collected information and asks for confirmation or correction
11. User confirms correctness
12. Assistant updates company in the CRM

### US.2.3 Create New Opportunity
As a sales rep (user) I want to be able to create a new opportunity for an existing company in the CRM through the DieStimme assistant (assistant) so that I can kick off the acquisition pipeline.

#### Requirements
- US.1.1
- US.1.2
- Company exists in the CRM

#### Main Flow
1. User asks assistant to create an opportunity for company with name X
2. Assistant checks if company by name X exists in the CRM
3. If not, assistant informs user that it could not find a matching company and asks them if they want to create one first
4. Otherwise assistant converses with user to determine name, value, first outreach date and closing date for the opportunity
5. Assistant asks user if they should be assigned as owner of the opportunity
6. User confirms or rejects to be assigned as owner
7. Assistant reads back opportunity information to user and asks for confirmation
8. User corrects false fields
9. Assistant reads back updated opportunity information to user and asks for confirmation
10. User confirms
11. Assistant creates opportunity in CRM

### US.2.4 Create Note for Company
As a sales rep (user) I want to be able to create notes related to a company through the DieStimme assistant (assistant) so that new information is persisted and available through the CRM.

#### Requirements
- US.1.1
- US.1.2
- Company exists in CRM

#### Main Flow
1. User asks assistant to create a note on a company with name X
2. Assistant checks if company by name X exists in the CRM
3. If not, assistant informs user that it could not find a matching company and asks them if they want to create one first
4. Otherwise assistant asks user what they want to take note of
5. Assistant converses with the user to determine the contents of the note
6. Assistant summarises contents and reads them back to the user and asks for confirmation
7. User corrects contents
8. Assistant reads back updated contents to user and asks for confirmation
9. User confirms
10. Assistant persists note in the CRM and links it to the company
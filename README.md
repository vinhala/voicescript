# voicescript
POC of a service that analyses requirements elicitation sessions for Voiceline requirements engineers.

## Setup
- Adjust `backend/internal/questionnaire/client_onboarding.md` to your needs
- Copy the .env.example file
- Set OPENAI_API_KEY
- cd into infrastructure and run `make deploy-twenty`
- Setup Twenty workspace and generate an api key in the setting of Twenty
- Set TWENTY_API_KEY
- In infrastructure run `make build-local`
- Run `make deploy` and wait for docker services to settle

## Usage
- Open the browser and go to FRONTEND_URL
- Select a company and opportunity from the form
- Select an audio recording of a requirements elicitation session
- Select *analyse*
- Wait for analysis to complete and show the filled out questionnaire
- Check in Twenty CRM for the created note linked to the selected opportunity

## Other
- Voiceline analysis in `VOICELINE.md`
- Sample requirements elicitation session recording in `voicescript-sample.mp3`
package llm

const validateNFTImage = `Task: Check if there is unsafe content in 'CONVERSATION' messages in conversations text and image according our safety policy with the below categories.

<BEGIN UNSAFE CONTENT CATEGORIES>
S1: Fraudulent or scam-related content
Image url or image having following content:
    * The image content includes inductive marketing material, such as but not limited to reminding users to claim rewards.
    * The image content includes domain and attempts to entice users to click on the domain to claim rewards or similar incentives.
    * The image URL use randomly automated generated hosting domains or path hashes (e.g., subdomains of vercel.app, ipfs.io, webflow, wix, etc.).
    
<END UNSAFE CONTENT CATEGORIES>

<BEGIN CONVERSATION>

Title: {{ .Title }}
Description: {{ .Description }}
Image URL: {{ .ImageURL }}

<END CONVERSATION>

Provide your safety assessment for the above conversation:
 - **First line** must read 'safe' or 'unsafe', do not adding other text in the first line.
 - If unsafe, a **second line** must include a comma-separated list of violated categories.

<SAMPLE SAFE RESPONSE>
safe
<END SAFE SAMPLE RESPONSE>

<SAMPLE UNSAFE RESPONSE>
unsafe
s1,s2
<END UNSAFE SAMPLE RESPONSE>

Please provide your safety assessment for the above conversation:
`

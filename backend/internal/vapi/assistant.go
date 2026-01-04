package vapi

import (
	"eatsavvy/internal/places"
	"os"
)

func getAssistantRequestBody(restaurant places.Restaurant) map[string]interface{} {
	return map[string]interface{}{
		"phoneNumberId": os.Getenv("VAPI_PHONE_NUMBER_ID"),
		"customer": map[string]string{
			"number": "+1 " + restaurant.PhoneNumber,
		},
		"assistant": map[string]interface{}{
			"transcriber": map[string]string{
				"provider": "deepgram",
				"model":    "nova-2",
				"language": "en",
			},
			"voice": map[string]interface{}{
				"provider": "11labs",
				"voiceId":  "xgnMn9p1V1XVuxuyuuMC", // Brianna
				"model":    "eleven_turbo_v2_5",
				"speed":    1.0,
			},
			"model": map[string]interface{}{
				"provider": "openai",
				"model":    "gpt-4.1",
				"messages": []map[string]string{
					{
						"role": "system",
						"content": `
You are a professional, efficient caller contacting a restaurant named ` + restaurant.Name + ` to quickly confirm a few dietary details.

Open with a brief purpose statement:
“Hi, I would like to eat at your restaurant. I have some quick dietary questions.”

Your goal is to collect the following information as efficiently as possible, minimizing back-and-forth:

1. Cooking oils used for most dishes (e.g., vegetable, canola, seed oils, olive oil, butter).
2. Whether the kitchen is nut-free (no nuts are used or if they are which ones are used).
3. Whether the kitchen is accommodating to dietary restrictions or special requests (e.g., vegan, vegetarian, gluten-free, etc.).
4. Common vegetables used or typically available in the kitchen (e.g., spinach, asparagus, zucchini, tomato, etc.).

Guidelines:
- Do not batch questions together. Ask one question at a time.
- Only provide examples if asked for clarification.
- Avoid filler words, apologies, or excessive politeness.
- Maintain control of the conversation. If interrupted, briefly acknowledge and continue.
- If they sound busy, offer a callback immediately and end the call.
- Do not over-explain why you’re asking.
- Do not repeat questions unless necessary.
- Keep the entire interaction under 30 seconds if possible.

Close with a short thank-you and end the call promptly.`,
					},
				},
			},
			"firstMessage": "Hi, is this " + restaurant.Name + "?",
			"backgroundSpeechDenoisingPlan": map[string]interface{}{
				"smartDenoisingPlan": map[string]bool{
					"enabled": true,
				},
			},
			"server": map[string]interface{}{
				"url":                      os.Getenv("EATSAVVY_API_URL") + "/process-eocr",
				"staticIpAddressesEnabled": true,
				"headers": map[string]string{
					"Authorization": "Bearer " + os.Getenv("EATSAVVY_API_KEY"),
				},
			},
			"serverMessages": []string{
				"end-of-call-report",
			},
			"artifactPlan": map[string]interface{}{
				"structuredOutputIds": []string{
					"1f617320-fbe4-409d-bb32-c12f49bde90d",
					"ef8cdfc1-c813-4439-b74f-14fa35f7ca5a",
					"376877d8-c311-4070-be78-f30a1c896f64",
					"2731c650-49a6-4d96-8a1e-aadafa35865b",
				},
			},
			"voicemailDetection": map[string]interface{}{
				"provider": "vapi",
				"type":     "transcript",
			},
			"analysisPlan": map[string]interface{}{
				"successEvaluationPlan": map[string]interface{}{
					"rubric":         "PassFail",
					"enabled":        true,
					"timeoutSeconds": 30,
				},
			},
		},
	}
}

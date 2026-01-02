package vapi

import "eatsavvy/pkg/places"

func GetAssistantRequestBody(restaurant places.Restaurant) map[string]interface{} {
	return map[string]interface{}{
		"phoneNumberId": "14a2a45f-7493-4c56-9cb1-6e8f19c8f0ef",
		"customer": map[string]string{
			"number": "+14252684016",
		},
		"assistant": map[string]interface{}{
			"transcriber": map[string]string{
				"provider": "deepgram",
				"model":    "nova-2",
				"language": "en",
			},
			"voice": map[string]interface{}{
				"provider": "11labs",
				"voiceId":  "sarah",
				"model":    "eleven_turbo_v2_5",
				"speed":    1.0,
			},
			"model": map[string]interface{}{
				"provider": "openai",
				"model":    "gpt-4.1",
				"messages": []map[string]string{
					{
						"role": "system",
						"content": `You are a friendly and polite caller inquiring about dietary and nutritional information from a restaurant called ` + restaurant.Name + `.

	Your goal is to gather the following information in a natural, conversational manner:

	1. **Cooking oils**: Ask what type of oil they use for cooking (e.g., vegetable oil, canola oil, seed oils, olive oil, etc.)
	2. **Nut allergies**: Ask if their kitchen is nut-free or if they can accommodate nut allergies
	3. **Dietary accommodations**: Ask if they are generally accommodating to dietary restrictions and special requests
	4. **Vegetables**: Ask what vegetables they typically have available or use in their dishes

	Be conversational and don't rush through the questions. Thank them for their time and be appreciative of any information they provide. If they seem busy, offer to call back at a better time.

	Keep your responses concise and natural - you're having a phone conversation, not reading a script.`,
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
				"url":                      "https://eatsavvy-api.tejas.wtf/process-eocr",
				"staticIpAddressesEnabled": true,
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
		},
	}
}

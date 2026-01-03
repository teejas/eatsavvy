# EatSavvy

Uses Google Places API, RabbitMQ, and Vapi to enrich restaurant details with nutritional and dietary information.

## API

Accepts search query to find restaurants to enrich. Returns restaurant info from the database. Accepts and processes end of call report from Vapi to enrich restaurant nutritional and dietary info.

## Worker

Uses Vapi to call restaurants and collect information from restaurants.

# Notes

The architecture is intended to decouple initiation of the phone call from the API

This is because restaurants may not be open when the request to enrich data is submitted, in which case we would have to poll periodically or have a long timeout or do something else heinous

We really only need one publisher (the api) and one consumer (the worker) though since all the work is relatively short/easy.

1. user submits a job (enrich nutrition info about this restaurant)
2. job is added to queue
3. worker picks up a job
3a. if restaurant is closed, it re-queues with delay until open (+30 min)
3b. if restaurant is open, it places outbound call via Vapi API
4. Vapi makes phone call and sends end of call report to /process-eocr with both transcript and structured outputs
5. We process the end of call report and supplement restaurant info in the DB
5a. We could do additional post-processing on the transcript in another offline worker

## To-do
- [ ] cache retrieved restaurant info from SearchRestaurants instead of querying again in GetPlacesDetails (use in-mem cache, implement myself for fun)
- [ ] refactor internal/worker/*
- [ ] move openNow logic from UI to API (currently duplicated bleh)
- [ ] dynamically generate structured outputs for assistant (using structuredMultiData)
- [ ] add Yelp support (for reviews and supplementing missing phone numbers)